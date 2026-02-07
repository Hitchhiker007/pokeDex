package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	pokecache "github.com/Hitchhiker007/pokeDex/internal"
)

// client is a wrapper around the standard http client
// adds caching support using a simple in-memory cache
type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
	// below is for my TokenBucket rate limting
	currentTokens float64
	tokenCapacity int
	refillRate    float64
	lastRefill    time.Time
	mu            sync.Mutex
}

// NewClient creates a new pokeapi client with a timeout for requests
// initialize an in-memory cache with a  expiration time of 5 minutes
func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(5 * time.Minute),
		// here rate limting is set on the new client creation
		tokenCapacity: 5,
		refillRate:    0.5, // 5 requests per 10 secs -> 5 / 10 = 0.5
		currentTokens: 5,   // start full
		lastRefill:    time.Now(),
	}
}

// get fetches the  URL from the PokeAPI, first check the cache, if the response is cached
// return cached value
// if not perform a http get request, read the response body, caches it, and then return it
func (c *Client) Get(url string) ([]byte, error) {

	// Check cache first, no tokens consumed
	if cached, found := c.cache.Get(url); found {
		return cached, nil
	}

	// lock mutex to protect shared state tokens and lastRefill
	// multiple goroutines, concurrent requests
	c.mu.Lock()
	defer c.mu.Unlock()

	// how much time has passed since the last refill of tokens
	elapsedTime := time.Now().Sub(c.lastRefill)
	// refill bucket by -> elapsed seconds * refillRate (only done when a request happens)
	refill := elapsedTime.Seconds() * c.refillRate
	c.currentTokens = c.currentTokens + refill
	// ensure tokens do no exceed the bucket capacity
	if c.currentTokens > float64(c.tokenCapacity) {
		c.currentTokens = float64(c.tokenCapacity)
	}
	// if token exists consume, else throw err
	if c.currentTokens >= 1 {
		c.currentTokens -= 1
	} else {
		return nil, fmt.Errorf("Rate Limiting Has Been Triggered Please Wait and Try again")
	}
	// update timestamp
	c.lastRefill = time.Now()

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
