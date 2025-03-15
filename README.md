# Layer Edge Light Node

## Introduction

Layer Edge Light Node is a client that connects to the Layer Edge network to verify Merkle trees by collecting random samples from available trees and verifying their integrity. The light node performs Zero-Knowledge proof verification operations through a local ZK prover service, and submits verified proofs to the network.

Key features:
- Automatically discovers available Merkle trees from the network
- Collects random samples from trees for verification
- Generates and verifies Zero-Knowledge proofs
- Implements intelligent sleep mechanisms to avoid redundant work on unchanged trees
- Submits verified proofs to earn rewards

## Building and Running

### Prerequisites

- Go 1.18 or higher
- Access to a Layer Edge gRPC endpoint
- Local ZK prover service running on port 3001

### Build Instructions

1. Clone the repository:
   ```bash
   git clone https://github.com/Layer-Edge/light-node.git
   cd light-node
   ```

2. Run the RISC0 Merkle service
    ```bash
    cd risc0-merkle-service
    cargo build && cargo run
2. Build the binary:
   ```bash
   go build
   ```

3. Run the node:
   ```bash
   ./light-node
   ```

## Configuration

The light node uses several configuration parameters that you might need to customize:

### GRPC Connection
In `clients/cosmos_client.go`, modify:
```go
const (
    grpcURL      = "34.31.74.109:9090"     // Replace with your gRPC endpoint
    contractAddr = "cosmos1ufs3tlq4umljk0qfe8k5ya0x6hpavn897u2cnf9k0en9jr7qarqqt56709"  // Replace with your contract address
)
```

### ZK Prover Service
In `node/prover.go`, modify the ZK prover URL if needed:
```go
resp, err := clients.PostRequest[ZKProverPayload, ZKProverResponse]("http://127.0.0.1:3001/process", ...)
```

If you do so, you need to modify the URL in `risc0-merkle-service/host/src/main.rs` as well

```rust
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv::dotenv().ok();
    println!("Starting server on port 3001...");
    
    HttpServer::new(|| {
        App::new()
            .route("/process", web::post().to(process))
    })
    .bind("127.0.0.1:3001")?
    .run()
    .await
}
```

### Worker Configuration
In `main.go`, you can adjust the worker polling interval:
```go
time.Sleep(5 * time.Second)  // Adjust as needed
```

## Setting Your Wallet Address and Signature

To receive rewards for submitting verified proofs, you need to update the wallet address and signature in the `CollectSampleAndVerify` function located in `node/prover.go`:

```go
if receipt != nil {
    walletAddress := "YOUR_WALLET_ADDRESS"  // Replace with your wallet address
    signature := "YOUR_SIGNATURE"           // Replace with your signature
    
    err = SubmitVerifiedProof(walletAddress, signature, *proof, *receipt)
    if err != nil {
        log.Printf("Failed to submit verified proof: %v", err)
    } else {
        log.Printf("Successfully submitted verified proof for tree %s", treeId)
    }
}
```

### Generating a Signature

The signature should be created based on your wallet's private key. Specific instructions for generating a valid signature will depend on your wallet type and the Layer Edge signature requirements. Please refer to the Layer Edge documentation for detailed instructions on signature generation.

## Advanced Configuration

### Tree Sleep Mechanism

The light node implements a sleep mechanism to avoid redundant work on trees that haven't changed. You can adjust the sleep parameters in `node/prover.go`:

```go
if state.ConsecutiveSame >= 3 {  // Number of consecutive checks before sleep
    sleepDuration := 5 * time.Minute  // Sleep duration
    state.SleepUntil = time.Now().Add(sleepDuration)
}
```

## Logging and Monitoring

The light node provides detailed logging about its operations. You can monitor the log output to track:
- Tree discovery
- Proof generation and verification
- Submission of verified proofs
- Sleep state of trees

## Troubleshooting

If you encounter issues:

1. Check your gRPC connection to the Layer Edge network
2. Ensure the ZK prover service is running and accessible
3. Verify your wallet address and signature format
4. Check logs for specific error messages

## License

This project is licensed under the MIT License - see the LICENSE file for details.