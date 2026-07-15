package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GeocodeCity looks up coordinates for a given city name.
// Notice it returns a slice (`[]GeocodeLocation`) instead of a single object pointer.
func (c *Client) GeocodeCity(ctx context.Context, cityName string, limit int) ([]GeocodeLocation, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/geo/1.0/direct", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse geocoding base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("q", cityName)
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create geocoding request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute geocoding request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close() // Satisfies the errcheck linter rule instantly!
	}()

	// 🏆 Integrated central error tracking
	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var result []GeocodeLocation
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	return result, nil
}
