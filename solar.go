package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetSolarIrradianceForecast fetches predicted solar radiation indexes (GHI, DNI, DHI) for a location.
func (c *Client) GetSolarIrradianceForecast(ctx context.Context, lat, lon float64) (*SolarIrradianceResponse, error) {
	// 1. Target the modern 3.0 advanced scientific solar endpoint layout path
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/3.0/solar/forecast", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse solar base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	// 2. Set up context-bound request architecture
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create solar request: %w", err)
	}

	// 3. Dispatch execution loop via client http pool
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute solar request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 4. Centralized Enterprise Error Management Pipeline Check
	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	// 5. Decode validated stream payload values
	var result SolarIrradianceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode solar response: %w", err)
	}

	return &result, nil
}
