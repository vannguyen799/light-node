package node

import (
	"fmt"
	"log"

	"github.com/Layer-Edge/light-node/clients"
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

func proveProof(data []string, proof_request string) (*Proof, error) {
	resp, err := clients.PostRequest[ZKProverPayload, ZKProverResponse]("http://34.57.83.179:3000/process", ZKProverPayload{
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

func verifyProofs(data []string, proof Proof) (*string, error) {
	resp, err := clients.PostRequest[ZKProverPayload, ZKProverResponse]("http://34.57.83.179:3000/process", ZKProverPayload{
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

func CollectSampleAndVerify() {
	cosmosQueryClient := clients.CosmosQueryClient{}
	err := cosmosQueryClient.Init()
	if err != nil {
		log.Fatalf("failed to initialize cosmos query client: %v", err)
	}
	defer cosmosQueryClient.Close()

	treeIds, err := cosmosQueryClient.ListMerkleTreeIds()
	if err != nil {
		log.Fatalf("failed to fetch tree ids: %v", err)
	}

	tree, err := cosmosQueryClient.GetMerkleTreeData(treeIds[len(treeIds)-1])
	if err != nil {
		log.Fatalf("failed to fetch tree data: %v", err)
	}

	sample := utils.RandomElement[string](tree.Leaves)

	proof, err := proveProof(tree.Leaves, sample)
	if err != nil {
		log.Fatalf("failed to prove sample: %v", err)
	}

	receipt, err := verifyProofs(tree.Leaves, *proof)
	if err != nil {
		log.Fatalf("failed to verify sample: %v", err)
	}

	log.Printf("Sample Data %v verified with receipt %v\n", sample, receipt)
}
