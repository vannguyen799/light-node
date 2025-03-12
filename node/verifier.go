package node

import (
	"fmt"
	"log"

	"github.com/Layer-Edge/light-node/clients"
	"github.com/Layer-Edge/light-node/config"
	"github.com/Layer-Edge/light-node/utils"
)

type ZKProverPayload struct {
	Operation    string   `json:"operation"`
	Data         []string `json:"data"`
	ProofRequest *string  `json:"proof_request"` // Using interface{} since it can be null
	Proof        *Proof   `json:"proof"`
}

type ZKProverResponse struct {
	Root          string             `json:"root"`
	Proof         *Proof             `json:"proof"`
	Verified      bool               `json:"verified"`
	Visualization *TreeVisualization `json:"visualization"`
	Receipt       string             `json:"receipt"`
}

type TreeVisualization struct {
	LevelSizes        []uint64   `json:"level_sizes"`
	TreeStructure     [][]string `json:"tree_structure"`
	DataToHashMapping [][]string `json:"data_to_hash_mapping"`
}

type Proof struct {
	LeafValue string          `json:"leaf_value"`
	ProofPath [][]interface{} `json:"proof_path"` // Each element is [string, bool]
}

func proveProof(zkrc *clients.ZKRequestClient, data []string, proof_request string) (*Proof, error) {
	resp, err := clients.PostRequest[ZKProverPayload, ZKProverResponse](zkrc, ZKProverPayload{
		Operation:    "prove",
		Data:         data,
		ProofRequest: &proof_request,
		Proof:        nil,
	})
	if err != nil {
		return nil, fmt.Errorf("proof verification error: %v", err)
	}
	log.Printf("verification done: %v", resp)
	return resp.Proof, nil
}

func verifyProofs(zkrc *clients.ZKRequestClient, data []string, proof Proof) (*string, error) {
	resp, err := clients.PostRequest[ZKProverPayload, ZKProverResponse](zkrc, ZKProverPayload{
		Operation:    "verify",
		Data:         data,
		ProofRequest: nil,
		Proof:        &proof,
	})
	if err != nil {
		return nil, fmt.Errorf("proof verification error: %v", err)
	}
	log.Printf("verification done: %v\n", resp)
	return &resp.Receipt, nil
}

func CollectSampleAndVerify(cfg *config.Config) {
	cosmosQueryClient := clients.CosmosQueryClient{}
	err := cosmosQueryClient.Init(cfg.LayerEdgeCosmos.GrpcEndpoint, cfg.LayerEdgeCosmos.MerkleTreeContractAddress)
	if err != nil {
		log.Fatalf("failed to initialize cosmos query client: %v", err)
	}
	defer cosmosQueryClient.Close()

	zkRequestClient := clients.ZKRequestClient{}
	zkRequestClient.Init(cfg.ZKComputationNode.NODE_URL)

	treeIds, err := cosmosQueryClient.ListMerkleTreeIds()
	if err != nil {
		log.Fatalf("failed to fetch tree ids: %v", err)
	}

	tree, err := cosmosQueryClient.GetMerkleTreeData(treeIds[len(treeIds)-1])
	if err != nil {
		log.Fatalf("failed to fetch tree data: %v", err)
	}

	sample := utils.RandomElement(tree.Leaves)

	proof, err := proveProof(&zkRequestClient, tree.Leaves, sample)
	if err != nil {
		log.Fatalf("failed to prove sample: %v", err)
	}

	receipt, err := verifyProofs(&zkRequestClient, tree.Leaves, *proof)
	if err != nil {
		log.Fatalf("failed to verify sample: %v", err)
	}

	log.Printf("Sample Data %v verified with receipt %v\n", sample, receipt)
}
