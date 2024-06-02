package nodes

import (
	"context"
	"crypto"
	"encoding/hex"
	"net"
	"sync"
	"time"

	proto "github.com/zacksfF/gRPC-P2P-UTXO-Blocker/Proto"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const blockTime = 5 * time.Second

type Mempool struct {
	lock sync.RWMutex
	txx  map[string]*proto.Transaction
}

func NewMempool() *Mempool {
	return &Mempool{
		txx: make(map[string]*proto.Transaction),
	}
}

func (m *Mempool) Clear() []*proto.Transaction {
	m.lock.Lock()
	defer m.lock.Unlock()

	txx := make([]*proto.Transaction, len(m.txx))
	it := 0
	for k, v := range m.txx {
		delete(m.txx, k)
		txx[it] = v
		it++
	}

	return txx
}

func (m *Mempool) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return len(m.txx)
}

func (m *Mempool) Has(tx *proto.Transaction) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	hash := hex.EncodeToString(types.HashTransaction(tx))
	_, ok := m.txx[hash]

	return ok
}

func (m *Mempool) Add(tx *proto.Transaction) bool {
	if m.Has(tx) {
		return false
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	hash := hex.EncodeToString(types.HashTransaction(tx))
	m.txx[hash] = tx

	return true
}

type ServerConfig struct {
	Version    string
	ListenAddr string
	PrivateKey *crypto.PrivateKey
}

type Node struct {
	ServerConfig
	logger   *zap.SugaredLogger
	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.Version
	mempool  *Mempool
	proto.UnimplementedNodeServer
}

func NewNode(cfg ServerConfig) *Node {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()
	return &Node{
		peers:        make(map[proto.NodeClient]*proto.Version),
		logger:       logger.Sugar(),
		mempool:      NewMempool(),
		ServerConfig: cfg,
	}
}

func (n *Node) Start(listenAddr string, bootstrapNodes []string) error {
	n.ListenAddr = listenAddr
	var (
		opts       []grpc.ServerOption
		grpcServer = grpc.NewServer(opts...)
	)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	proto.RegisterNodeServer(grpcServer, n)

	n.logger.Infow("Starting node...", "on", n.ListenAddr)

	// bootstrap network with the known nodes
	if len(bootstrapNodes) > 0 {
		go func() {
			err := n.bootstrapNetwork(bootstrapNodes)
			if err != nil {
				n.logger.Errorw("Bootstrap error", "error", err)
			}
		}()
	}

	if n.PrivateKey != nil {
		go n.validatorLoop()
	}

	return grpcServer.Serve(ln)
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	client, err := makeNodeClient(v.ListenAddr)
	if err != nil {
		return nil, err
	}

	n.addPeer(client, v)

	return n.getVersion(), nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	p, _ := peer.FromContext(ctx)
	hash := hex.EncodeToString(types.HashTransaction(tx))

	if n.mempool.Add(tx) {
		n.logger.Debugw("Received transaction", "from", p.Addr, "hash", hash, "we", n.ListenAddr)

		go func() {
			if err := n.broadcast(tx); err != nil {
				n.logger.Errorw("Broadcast error", "error", err)
			}
		}()
	}

	return &proto.Ack{}, nil
}

func (n *Node) broadcast(msg any) error {
	for p := range n.peers {
		switch v := msg.(type) {
		case *proto.Transaction:
			_, err := p.HandleTransaction(context.Background(), v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *Node) bootstrapNetwork(bootstrapNodes []string) error {
	for _, node := range bootstrapNodes {
		if !n.canConnectWith(node) {
			continue
		}
		n.logger.Debugw("Dialing remote nodes...", "ourNode", n.ListenAddr, "remoteNode", node)

		client, v, err := n.dialRemoteNode(node)
		if err != nil {
			return err
		}

		n.addPeer(client, v)
	}

	return nil
}

func (n *Node) validatorLoop() {
	//n.logger.Infow("Starting validator loop...", "publicKey", n.PrivateKey.Public(), "blockTime", blockTime)
	ticker := time.NewTicker(blockTime)

	for {
		<-ticker.C

		txx := n.mempool.Clear()

		n.logger.Debugw("Creating new block...", "lenTx", len(txx))
	}
}

func (n *Node) addPeer(client proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.peers[client] = v

	if len(v.PeerList) > 0 {
		go func() {
			go func() {
				err := n.bootstrapNetwork(v.PeerList)
				if err != nil {
					n.logger.Errorw("Bootstrap error", "error", err)
				}
			}()
		}()
	}

	n.logger.Debugw("New peer successfully connected.",
		"ourNode", n.ListenAddr,
		"remoteNode", v.ListenAddr,
		"height", v.Height)
}

func (n *Node) deletePeer(client proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	delete(n.peers, client)
}

func (n *Node) dialRemoteNode(addr string) (proto.NodeClient, *proto.Version, error) {
	client, err := makeNodeClient(addr)
	if err != nil {
		return nil, nil, err
	}

	v, err := client.Handshake(context.Background(), n.getVersion())
	if err != nil {
		return nil, nil, err
	}

	return client, v, nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "0.0.1",
		Height:     0,
		ListenAddr: n.ListenAddr,
		PeerList:   n.getPeerList(),
	}
}

func (n *Node) canConnectWith(addr string) bool {
	if n.ListenAddr == addr {
		return false
	}

	connectedPeers := n.getPeerList()
	for _, connectedAddr := range connectedPeers {
		if addr == connectedAddr {
			return false
		}
	}

	return true
}

func (n *Node) getPeerList() []string {
	n.peerLock.RLock()
	defer n.peerLock.RUnlock()

	var peers []string
	for _, v := range n.peers {
		peers = append(peers, v.ListenAddr)
	}

	return peers
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	conn, err := grpc.Dial(listenAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(conn), nil
}
