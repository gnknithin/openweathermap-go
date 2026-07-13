package openweathermap

import (
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.openweathermap.org"

// Client is the core SDK structure used to communicate with OpenWeatherMap APIs.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Option defines a function type used to configure the Client.
type Option func(*Client)

// NewClient initializes and returns a new OpenWeatherMap SDK client.
func NewClient(apiKey string, opts ...Option) *Client {
	// 1. Establish enterprise defaults
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Prevents production goroutine leaks
		},
	}

	// 2. Evaluate any functional options passed by the caller
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithHTTPClient allows overriding the default internal HTTP client.
func WithHTTPClient(customClient *http.Client) Option {
	return func(c *Client) {
		if customClient != nil {
			c.httpClient = customClient
		}
	}
}

// WithBaseURL allows overriding the base URL (useful for mocking/testing).
func WithBaseURL(url string) Option {
	return func(c *Client) {
		if url != "" {
			c.baseURL = url
		}
	}
}
