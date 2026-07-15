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
		return nil, fmt.Errorf("failed to parse geocode base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("q", cityName)
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create geocode request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute geocode request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close() // Satisfies the errcheck linter rule instantly!
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from geocoding API: %d", resp.StatusCode)
	}

	var result []GeocodeLocation
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode geocode response: %w", err)
	}

	return result, nil
}
