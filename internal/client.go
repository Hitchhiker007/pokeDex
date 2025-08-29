package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// client is the HTTP client for the PokeAPI
type Client struct {
	httpClient http.Client
}

// newClient creates a new client with a timeout
func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

// get fetches the URL and returns the response body
func (c *Client) Get(url string) ([]byte, error) {
	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with status %d: %s", res.StatusCode, body)
	}

	return body, nil
}
