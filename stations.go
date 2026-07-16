package openweathermap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// RegisterStation creates a new weather station on the OpenWeatherMap infrastructure (POST).
func (c *Client) RegisterStation(ctx context.Context, station StationRegisterRequest) (*StationResponse, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/3.0/stations", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse stations base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	// Marshal the Go struct into raw JSON bytes for the HTTP request body payload
	bodyBytes, err := json.Marshal(station)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal station registration request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create station registration request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute station registration request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var result StationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode station registration response: %w", err)
	}

	return &result, nil
}

// GetStationByID retrieves metadata for a specifically registered weather station by its unique server ID (GET).
func (c *Client) GetStationByID(ctx context.Context, stationID string) (*StationResponse, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/3.0/stations/%s", c.baseURL, stationID))
	if err != nil {
		return nil, fmt.Errorf("failed to parse station retrieval URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create station retrieval request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute station retrieval request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var result StationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode station retrieval response: %w", err)
	}

	return &result, nil
}

// DeleteStation removes a weather station permanently from the OpenWeatherMap infrastructure (DELETE).
func (c *Client) DeleteStation(ctx context.Context, stationID string) error {
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/3.0/stations/%s", c.baseURL, stationID))
	if err != nil {
		return fmt.Errorf("failed to parse station deletion URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create station deletion request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute station deletion request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return c.checkResponse(resp)
}
