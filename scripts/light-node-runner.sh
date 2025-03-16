#!/bin/bash
# Light Node Run Script
# This script starts the ZK prover service and launches the light node client

# Configuration
ZK_PROVER_DIR="risc0-merkle-service"
LIGHT_NODE_DIR="."
LIGHT_NODE_BIN="light-node"

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Log function
log() {
  echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
  echo -e "${RED}[$(date '+%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
  exit 1
}

# Check if zk-prover directory exists
if [ ! -d "$ZK_PROVER_DIR" ]; then
  error "ZK prover directory not found. Please run this script from the project root."
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
  error "Go is not installed. Please install Go before running this script."
fi

# Start ZK prover service
log "Starting ZK prover service..."
cd "$ZK_PROVER_DIR" || error "Failed to navigate to ZK prover directory"
cargo run > zk-prover.log 2>&1 &
ZK_PROVER_PID=$!
cd - > /dev/null || error "Failed to return to original directory"

# Wait for ZK prover to start
log "Waiting for ZK prover to initialize (5 seconds)..."
sleep 5

# Check if ZK prover is running
if ! ps -p $ZK_PROVER_PID > /dev/null; then
  error "ZK prover failed to start. Check zk-prover.log for details."
fi

log "ZK prover started successfully with PID: $ZK_PROVER_PID"

# Build light node
log "Building light node..."
cd "$LIGHT_NODE_DIR" || error "Failed to navigate to light node directory"
go build -o "$LIGHT_NODE_BIN" || error "Failed to build light node"
log "Light node built successfully"

# Run light node
log "Starting light node..."
log "Light node logs will appear below:"
echo -e "${GREEN}----------------------------------------${NC}"
./"$LIGHT_NODE_BIN"

# Handle exit - cleanup
cleanup() {
  log "Shutting down services..."
  kill $ZK_PROVER_PID 2>/dev/null
  log "ZK prover service stopped"
  log "Cleanup complete"
}

# Set trap to catch script termination
trap cleanup EXIT

# Exit
exit 0