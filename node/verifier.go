package node

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/Layer-Edge/light-node/clients"
	"github.com/Layer-Edge/light-node/utils"
	"github.com/go-resty/resty/v2"
)

type VerifyPayload struct {
	Operation    string      `json:"operation"`
	Data         []string    `json:"data"`
	ProofRequest interface{} `json:"proof_request"` // Using interface{} since it can be null
	Proof        *Proof      `json:"proof"`
}

type Proof struct {
	LeafValue string          `json:"leaf_value"`
	ProofPath [][]interface{} `json:"proof_path"` // Each element is [string, bool]
}

func postRequest(url string, data VerifyPayload) error {
	// Import at top of file: "github.com/go-resty/resty/v2"
	client := resty.New()

	// Set default headers, timeout
	client.
		SetTimeout(time.Second*10).
		SetHeader("Authorization", "Bearer your-token-here")

	// Make request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(url)

	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	// Check status code
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s",
			resp.StatusCode(), string(resp.Body()))
	}

	fmt.Printf("Response: %s\n", string(resp.Body()))
	return nil
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

	sampleSize := rand.Intn(len(tree.Leaves)-1) + 1

	sample := utils.RandomSample(tree.Leaves, sampleSize)

	err = postRequest("http://34.57.83.179:3000/process", VerifyPayload{
		Operation:    "verify",
		Data:         sample,
		ProofRequest: nil,
		Proof:        nil, // WIP to hash mapping from contract response
	})
	if err != nil {
		log.Fatalf("failed to verify sample: %v", err)
	}
}
