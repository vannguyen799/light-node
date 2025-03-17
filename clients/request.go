package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Layer-Edge/light-node/utils"
	"github.com/go-resty/resty/v2"
)

// Default timeout in seconds if environment variable is not set
const DEFAULT_TIMEOUT = 100

func PostRequest[T any, R any](url string, requestData T) (*R, error) {
	client := resty.New()

	// Get timeout from environment variable or use default
	timeout := DEFAULT_TIMEOUT
	envTimeout := utils.GetEnv("API_REQUEST_TIMEOUT", "100")
	if t, err := strconv.Atoi(envTimeout); err == nil {
		timeout = t
	}

	// Set default headers, timeout
	client.
		SetTimeout(time.Second*time.Duration(timeout)).
		SetHeader("Authorization", "Bearer your-token-here")

	// Make request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestData).
		Post(url)

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
