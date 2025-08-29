package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"time"

	pokecache "github.com/Hitchhiker007/pokeDex/internal"
)

// client is the HTTP client for the PokeAPI
type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

// newClient creates a new client with a timeout
func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(5 * time.Minute),
	}
}

// get fetches the URL and returns the response body
func (c *Client) Get(url string) ([]byte, error) {
	// Check cache first
	if cached, found := c.cache.Get(url); found {
		return cached, nil
	}

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

	// Add response to cache
	c.cache.Add(url, body)

	return body, nil
}
