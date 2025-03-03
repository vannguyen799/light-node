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
	ProofRequest *string  `json:"proof_request"` // Using pointer since it can be null
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

// Verifier handles proof verification operations
type Verifier struct {
	restClient   *clients.RestClient
	cosmosClient *clients.CosmosQueryClient
	cfg          *config.Config
}

// NewVerifier creates a new verifier with the given configuration
func NewVerifier(cfg *config.Config) *Verifier {
	return &Verifier{
		restClient:   clients.NewRestClient(cfg),
		cosmosClient: clients.NewCosmosQueryClient(cfg),
		cfg:          cfg,
	}
}

func (v *Verifier) proveProof(data []string, proofRequest string) (*Proof, error) {
	payload := ZKProverPayload{
		Operation:    "prove",
		Data:         data,
		ProofRequest: &proofRequest,
		Proof:        nil,
	}
	
	var response ZKProverResponse
	err := v.restClient.PostRequest(
		v.cfg.ZKProver.URL+"/process", 
		payload,
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf("proof verification error: %v", err)
	}
	
	log.Printf("verification done: %v", response)
	return response.Proof, nil
}

func (v *Verifier) verifyProofs(data []string, proof Proof) (*string, error) {
	payload := ZKProverPayload{
		Operation:    "verify",
		Data:         data,
		ProofRequest: nil,
		Proof:        &proof,
	}
	
	var response ZKProverResponse
	err := v.restClient.PostRequest(
		v.cfg.ZKProver.URL+"/process", 
		payload,
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf("proof verification error: %v", err)
	}
	
	log.Printf("verification done: %v\n", response)
	return &response.Receipt, nil
}

func (v *Verifier) CollectSampleAndVerify() error {
	err := v.cosmosClient.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize cosmos query client: %v", err)
	}
	defer v.cosmosClient.Close()

	treeIds, err := v.cosmosClient.ListMerkleTreeIds()
	if err != nil {
		return fmt.Errorf("failed to fetch tree ids: %v", err)
	}
	log.Println("Received tree ids: ", treeIds)
	
	if len(treeIds) == 0 {
		return fmt.Errorf("no tree ids found")
	}

	tree, err := v.cosmosClient.GetMerkleTreeData(treeIds[0])
	if err != nil {
		return fmt.Errorf("failed to fetch tree data: %v", err)
	}

	sample := utils.RandomElement[string](tree.Leaves)

	proof, err := v.proveProof(tree.Leaves, sample)
	if err != nil {
		return fmt.Errorf("failed to prove sample: %v", err)
	}

	receipt, err := v.verifyProofs(tree.Leaves, *proof)
	if err != nil {
		return fmt.Errorf("failed to verify sample: %v", err)
	}

	log.Printf("Sample Data %v verified with receipt %v\n", sample, receipt)
	return nil
}