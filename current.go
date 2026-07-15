package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetCurrentWeather fetches the current weather for the given latitude and longitude.
func (c *Client) GetCurrentWeather(ctx context.Context, lat, lon float64) (*CurrentWeatherResponse, error) {
	// 1. Construct the URL securely using Go's standard url package
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/2.5/weather", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse current weather base URL: %w", err)
	}

	// 2. Build out the query parameters
	query := endpoint.Query()
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	// 3. Create the HTTP request bound to our lifecycle Context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create current weather request: %w", err)
	}

	// 4. Execute the network request using our configured client
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute current weather request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}() // Ensures connection drains and closes to avoid leaks

	// 5. Call our new centralized error verification loop
	if err := c.checkResponse(resp); err != nil {
		return nil, err // Returns our native type-safe APIError directly to the caller!
	}

	// 6. Decode the validated successful response payload stream
	var result CurrentWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode current weather response: %w", err)
	}

	return &result, nil
}
