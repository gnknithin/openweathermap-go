package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetHistoricalWeather fetches historical weather snapshots for given coordinates at a precise Unix timestamp.
func (c *Client) GetHistoricalWeather(ctx context.Context, lat, lon float64, timestamp int64) (*HistoricalWeatherResponse, error) {
	// 1. Target the standard One Call 3.0 time machine endpoint path
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/3.0/onecall/timemachine", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse historical base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("dt", fmt.Sprintf("%d", timestamp))
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	// 2. Setup context-bound request structure
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create historical request: %w", err)
	}

	// 3. Dispatch execution loop via client http pool
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute historical request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 4. Centralized Error Management Pipeline Check
	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	// 5. Decode validated stream payload values
	var result HistoricalWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode historical response: %w", err)
	}

	return &result, nil
}
