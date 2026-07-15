package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetAirPollution fetches current air pollution metrics for given coordinates.
func (c *Client) GetAirPollution(ctx context.Context, lat, lon float64) (*AirPollutionResponse, error) {
	// 1. Target the standard 2.5 air pollution endpoint path
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/2.5/air_pollution", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse air pollution base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	// 2. Create context-bound network request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create air pollution request: %w", err)
	}

	// 3. Dispatch execution
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute air pollution request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 🏆 Integrated central error tracking
	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var result AirPollutionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode air pollution response: %w", err)
	}

	return &result, nil
}
