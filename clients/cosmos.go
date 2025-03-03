package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Layer-Edge/light-node/config"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MerkleTree struct {
	Root     string   `json:"root"`
	Leaves   []string `json:"leaves"`
	Metadata string   `json:"metadata"`
}

type QueryGetTree struct {
	GetMerkleTree struct {
		ID string `json:"id"`
	} `json:"get_merkle_tree"`
}

type QueryListTreeIDs struct {
	ListMerkleTreeIds struct {
	} `json:"list_merkle_tree_ids"`
}

type CosmosQueryClient struct {
	conn        *grpc.ClientConn
	queryClient wasmtypes.QueryClient
	cfg         *config.Config
}

func NewCosmosQueryClient(cfg *config.Config) *CosmosQueryClient {
	return &CosmosQueryClient{
		cfg: cfg,
	}
}

func (cqc *CosmosQueryClient) Init() error {
	// Connect to gRPC client
	conn, err := grpc.NewClient(cqc.cfg.GRPC.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC: %v", err)
	}

	cqc.conn = conn
	cqc.queryClient = wasmtypes.NewQueryClient(conn)
	return nil
}

func (cqc *CosmosQueryClient) Close() {
	if cqc.conn != nil {
		cqc.conn.Close()
	}
}

func (cqc *CosmosQueryClient) GetMerkleTreeData(id string) (*MerkleTree, error) {
	query := QueryGetTree{}
	query.GetMerkleTree.ID = id

	queryBytes, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Failed to marshal query: %v", err)
	}

	res, err := cqc.queryClient.SmartContractState(
		context.Background(),
		&wasmtypes.QuerySmartContractStateRequest{
			Address:   cqc.cfg.Cosmos.ContractAddr,
			QueryData: queryBytes,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query contract: %v", err)
	}

	// Parse response JSON into struct
	var tree MerkleTree
	err = json.Unmarshal(res.Data, &tree)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tree data: %v", err)
	}
	log.Println("Received tree: ", tree)
	return &tree, nil
}

func (cqc *CosmosQueryClient) ListMerkleTreeIds() ([]string, error) {
	query := QueryListTreeIDs{}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Failed to marshal query: %v", err)
	}

	res, err := cqc.queryClient.SmartContractState(
		context.Background(),
		&wasmtypes.QuerySmartContractStateRequest{
			Address:   cqc.cfg.Cosmos.ContractAddr,
			QueryData: queryBytes,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query contract: %v", err)
	}

	// Parse response JSON into struct
	var treeIds []string
	err = json.Unmarshal(res.Data, &treeIds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tree data: %v", err)
	}

	return treeIds, nil
}