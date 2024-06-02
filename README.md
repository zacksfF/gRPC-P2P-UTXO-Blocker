## GrpcP2PUtxoBlockchain: Secure UTXO Blockchain with Go & gRPC

This project implements a secure, peer-to-peer (P2P) blockchain using the Unspent Transaction Output (UTXO) model in Golang with gRPC for communication.

### Project Description

This blockchain utilizes the UTXO model, where transactions spend outputs from previous transactions, promoting efficient verification. Secure communication between nodes is facilitated through gRPC for reliable data exchange.

### Requirements

* **Programming Language:** ``go1.22.2``
* **Dependencies:**
    * Protobuf compiler ``(protoc)``
    * gRPC libraries ([https://grpc.io/](https://grpc.io/))
    * Additional libraries for cryptography and data structures (to be specified in installation) ......


### Functional Requirements

* **Cryptography:** Implement private and public key cryptography for secure transactions.
* **Data Structures:** Utilize protobuf for defining message formats and the Merkle Tree for efficient block verification.
* **UTXO Model:** Manage transactions using the UTXO model for efficient spending.
* **P2P Communication:** Establish communication between nodes using gRPC for reliable data exchange.
* **Peer Discovery:** Implement a custom gossip protocol for peer discovery.
* **Blockchain Management:**
    * Add and validate new blocks.
    * Manage the transaction mempool.
    * Create and validate transactions.
    * Store UTXO data efficiently.
* **Circuit Breaking & Rate Limiting:** Implement mechanisms to prevent overloading and ensure system stability.

### Features

* Secure UTXO blockchain implementation with Golang and gRPC.
* Efficient transaction verification with the Merkle Tree.
* Scalable peer-to-peer network with custom gossip protocol.
* Robust transaction management with creation, validation, and UTXO storage.
* Integrated circuit breaking and rate limiting for system stability.

### Installation

1. **Install Golang:** Download and install Golang from the official website ([https://go.dev/](https://go.dev/)). Ensure you set up the environment variables (``GOPATH``, ``GOROOT``).
2. **Install Protobuf compiler:** Follow the installation instructions for `protoc` based on your operating system ([https://protobuf.dev/](https://protobuf.dev/)).
3. **Install gRPC libraries:** Use `go get` to download the necessary gRPC libraries:

```
go get -u google.golang.org/grpc
```

4. **Install additional dependencies:** Specific libraries for cryptography and data structures might be required. Refer to the project code for details and installation instructions.
5. **Build the project:** Navigate to the project directory and run:

```
go run main.go
```
## Example Some Usage
1. 
```
TestPrivateKeySign
    === RUN   TestGeneratePrivateKey
    --- PASS: TestGeneratePrivateKey (0.00s)
PASS
ok      github.com/zacksfF/gRPC-P2P-UTXO-Blocker/encrypted      1.085s

output : Addrees Len = 2c22b31027a2683deeec8f5d3c1bdd8a0a31b952f774684db6f23b8feaa1
```
2. 
```
TestVerifyBlock:
Running tool: /usr/local/bin/go test -timeout 30s -run ^TestVerifyBlock$ github.com/zacksfF/gRPC-P2P-UTXO-Blocker/types

=== RUN   TestVerifyBlock
--- PASS: TestVerifyBlock (0.00s)
PASS
ok      github.com/zacksfF/gRPC-P2P-UTXO-Blocker/types  1.156s

```

## Demo 
I don't Know Why this Demo push like this but you can see this in better quality https://www.linkedin.com/feed/update/urn:li:activity:7203072833566949376/
![Screen Recording 2024-06-02 at 17 31 49](https://github.com/zacksfF/gRPC-P2P-UTXO-Blocker/assets/129240583/08436c7b-9009-4b07-8af3-92626451063e)

## Project Thoughts

**Security**: Implementing a secure blockchain is crucial. Ensure that your cryptographic functions (hashing, signing, verifying) are correctly implemented and up-to-date with the latest standards.
gRPC for Communication: Using gRPC for communication between nodes is an excellent choice as it allows for efficient, language-agnostic communication with strong type-checking.

**Concurrency in Go:** Go is well-suited for handling concurrent operations, which is vital in a blockchain for handling multiple transactions and blocks simultaneously.

**UTXO Model:** The UTXO model is a great choice for simplicity and security. It is also the model used by Bitcoin, making it a tried and tested approach.

### Potential Problems to Solve

**Scalability:**
Problem: As the number of transactions and blocks increases, the blockchain can become large and cumbersome to manage.
Solution: Implement techniques such as pruning, sharding, or layer 2 solutions (e.g., payment channels) to improve scalability.

**Consensus Mechanism:**
Problem: Ensuring consensus across nodes in a decentralized manner can be challenging.
Solution: Evaluate and implement a robust consensus mechanism (e.g., Proof of Work, Proof of Stake, Practical Byzantine Fault Tolerance) that suits your blockchain's requirements.

**Network Latency and Partitioning:**
Problem: Network delays and partitioning can lead to inconsistencies and forks.
Solution: Implement strategies to handle network partitions and ensure quick resolution of forks.

**Transaction Throughput:**
Problem: High transaction throughput is essential for usability, but blockchains often struggle with this.
Solution: Optimize the block size and interval, and consider batch processing or off-chain solutions to increase throughput.

**Security Against Attacks:**
Problem: Blockchain systems can be vulnerable to various attacks (e.g., 51% attack, Sybil attack, double-spending).
Solution: Enhance security measures, such as increasing the difficulty of the proof-of-work algorithm, using robust peer authentication methods, and implementing consensus rules that mitigate these risks.

**Privacy:**
Problem: Transaction data on blockchains is typically public, which can be a privacy concern.
Solution: Implement privacy-enhancing technologies such as zero-knowledge proofs, ring signatures, or confidential transactions.

## Building from source
Environment requirement: Go +1.22

Compile:

```
git clone https://github.com/zacksfF/gRPC-P2P-UTXO-Blocker.git
cd gRPC-P2P-UTXO-Blocker
code .
```

## Enjoy
