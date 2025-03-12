package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type ZKRequestClient struct {
	client  *resty.Client
	NodeUrl string
}

func (zkrc *ZKRequestClient) Init(node_url string) {
	zkrc.client = resty.New()
	zkrc.client.BaseURL = node_url
}

func PostRequest[T any, R any](zkrc *ZKRequestClient, requestData T) (*R, error) {
	// Set default headers, timeout
	zkrc.client.
		SetTimeout(time.Second*10).
		SetHeader("Authorization", "Bearer your-token-here")

	// Make request
	resp, err := zkrc.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestData).
		Post(fmt.Sprintf("%s/process", zkrc.NodeUrl))

	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	// Check status code
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s",
			resp.StatusCode(), string(resp.Body()))
	}

	// Parse response into generic type
	var response R
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &response, nil
}
