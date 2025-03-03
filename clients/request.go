package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Layer-Edge/light-node/config"
	"github.com/go-resty/resty/v2"
)

// RestClient encapsulates the API client with configuration
type RestClient struct {
	client *resty.Client
	cfg    *config.Config
}

// NewRestClient creates a new REST client with the given configuration
func NewRestClient(cfg *config.Config) *RestClient {
	client := resty.New()
	
	// Set default headers, timeout from config
	client.
		SetTimeout(time.Second * time.Duration(cfg.API.Timeout)).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", cfg.API.AuthToken))
		
	return &RestClient{
		client: client,
		cfg:    cfg,
	}
}

// PostRequest makes a POST request and unmarshals the response into the provided response object
func (rc *RestClient) PostRequest(url string, requestData interface{}, response interface{}) error {
	// Make request
	resp, err := rc.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestData).
		Post(url)

	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	// Check status code
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s",
			resp.StatusCode(), string(resp.Body()))
	}

	// Parse response into provided response object
	if err := json.Unmarshal(resp.Body(), response); err != nil {
		return fmt.Errorf("error parsing response: %v", err)
	}

	return nil
}