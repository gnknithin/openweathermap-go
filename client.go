package openweathermap

import (
	"encoding/json"
	"fmt"
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

// checkResponse inspects the HTTP response status code. If it's a success status code,
// it returns nil. If it's an error status code, it reads and returns a populated APIError.
func (c *Client) checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	// Draft a safe fallback error structure in case the response body is empty or unparseable
	apiErr := &APIError{
		Code:    resp.StatusCode,
		Message: fmt.Sprintf("unexpected response status: %s", resp.Status),
	}

	// Try to decode the explicit JSON error payload returned by OpenWeatherMap
	if resp.Body != nil {
		var openWeatherErr struct {
			Cod     interface{} `json:"cod"` // OpenWeatherMap sometimes sends cod as a string or int
			Message string      `json:"message"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&openWeatherErr); err == nil {
			if openWeatherErr.Message != "" {
				apiErr.Message = openWeatherErr.Message
			}
			// Safe type-assertion case mapping if cod comes back unexpectedly
			if intCod, ok := openWeatherErr.Cod.(float64); ok {
				apiErr.Code = int(intCod)
			} else if strCod, ok := openWeatherErr.Cod.(string); ok {
				var parsedCod int
				if _, fmtErr := fmt.Sscanf(strCod, "%d", &parsedCod); fmtErr == nil {
					apiErr.Code = parsedCod
				}
			}
		}
	}

	return apiErr
}
