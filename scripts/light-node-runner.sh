#!/bin/bash
# start-light-node.sh
# This script builds and runs the light node client

# Configuration
LIGHT_NODE_DIR="."
LIGHT_NODE_BIN="light-node"
ZK_PROVER_PID_FILE="zk_prover_pid.txt"

# Color codes
BLUE='\033[0;34m'
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Log function
log() {
  echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
  echo -e "${RED}[$(date '+%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
  exit 1
}

success() {
  echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}"
}

# Check if ZK prover PID file exists
if [ ! -f "$ZK_PROVER_PID_FILE" ]; then
  error "ZK prover PID file not found. Please run start-risc0-service.sh first."
fi

# Read ZK prover PID
ZK_PROVER_PID=$(cat "$ZK_PROVER_PID_FILE")

# Check if ZK prover is running
if ! ps -p "$ZK_PROVER_PID" > /dev/null; then
  error "ZK prover is not running. Please start it using start-risc0-service.sh."
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
  error "Go is not installed. Please install Go before running this script."
fi

# Build light node
log "Building light node..."
cd "$LIGHT_NODE_DIR" || error "Failed to navigate to light node directory"
go build -o "$LIGHT_NODE_BIN" || error "Failed to build light node"
success "Light node built successfully"

# Run light node
success "Starting light node..."
success "Light node logs will appear below:"
echo -e "${GREEN}----------------------------------------${NC}"
./"$LIGHT_NODE_BIN"

exit 0