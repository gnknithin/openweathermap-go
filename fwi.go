package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetFireWeatherIndex fetches current and forecast Fire Weather Index telemetry calculations for given coordinates.
func (c *Client) GetFireWeatherIndex(ctx context.Context, lat, lon float64) (*FireWeatherIndexResponse, error) {
	// 1. Target the advanced 3.0 structural environmental path mapping
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/3.0/fwi", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse fwi base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	// 2. Set up the context-bound request structure
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create fwi request: %w", err)
	}

	// 3. Dispatch request via HTTP client connection pool
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute fwi request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 4. Centralized Enterprise Error Management Pipeline Check
	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	// 5. Stream decode validated JSON bytes directly into the target structural types
	var result FireWeatherIndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode fwi response: %w", err)
	}

	return &result, nil
}
