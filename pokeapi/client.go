package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"time"

	pokecache "github.com/Hitchhiker007/pokeDex/internal"
)

// client is a wrapper around the standard http client
// adds caching support using a simple in-memory cache
type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

// NewClient creates a new pokeapi client with a timeout for requests
// initialize an in-memory cache with a  expiration time of 5 minutes
func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(5 * time.Minute),
	}
}

// get fetches the  URL from the PokeAPI, first check the cache, if the response is cached
// return cached value
// if not perform a http get request, read the response body, caches it, and then return it
func (c *Client) Get(url string) ([]byte, error) {
	// Check cache first
	if cached, found := c.cache.Get(url); found {
		return cached, nil
	}

	// perform request
	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// check response for failure
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with status %d: %s", res.StatusCode, body)
	}

	// Add response to cache
	c.cache.Add(url, body)
	// return raw body
	return body, nil
}
