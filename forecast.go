package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ForecastOptions contains optional configuration for the 5-Day Forecast API request.
type ForecastOptions struct {
	Units string
	Lang  string
}

// Get5DayForecast fetches the 5-day / 3-hour weather forecast for the given coordinates.
func (c *Client) Get5DayForecast(ctx context.Context, lat, lon float64, opts *ForecastOptions) (*ForecastResponse, error) {
	// 1. Target the standard 2.5 forecast path
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/2.5/forecast", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse forecast base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("appid", c.apiKey)

	// 2. Process optional query parameters if provided
	if opts != nil {
		if opts.Units != "" {
			query.Set("units", opts.Units)
		}
		if opts.Lang != "" {
			query.Set("lang", opts.Lang)
		}
	}
	endpoint.RawQuery = query.Encode()

	// 3. Create the context-bound request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create forecast request: %w", err)
	}

	// 4. Dispatch the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute forecast request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close() // Keep our linter smiling
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from forecast API: %d", resp.StatusCode)
	}

	// 5. Stream parse directly into the composed struct
	var result ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	return &result, nil
}
