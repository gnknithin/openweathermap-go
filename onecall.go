package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ExcludeBlock defines a custom type for specifying blocks to omit from the One Call payload.
type ExcludeBlock string

const (
	ExcludeCurrent  ExcludeBlock = "current"
	ExcludeMinutely ExcludeBlock = "minutely"
	ExcludeHourly   ExcludeBlock = "hourly"
	ExcludeDaily    ExcludeBlock = "daily"
	ExcludeAlerts   ExcludeBlock = "alerts"
)

// OneCallOptions contains optional configuration for the One Call API request.
type OneCallOptions struct {
	Units   string
	Exclude []ExcludeBlock
}

// GetOneCallWeather fetches comprehensive weather data using the One Call 4.0 API.
func (c *Client) GetOneCallWeather(ctx context.Context, lat, lon float64, opts *OneCallOptions) (*OneCallResponse, error) {
	// 1. Target the modern 3.0/4.0 endpoint path
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/3.0/onecall", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse onecall base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("appid", c.apiKey)

	// 2. Dynamically process optional settings if provided by the caller
	if opts != nil {
		if opts.Units != "" {
			query.Set("units", opts.Units)
		}
		if len(opts.Exclude) > 0 {
			var exclusions []string
			for _, b := range opts.Exclude {
				exclusions = append(exclusions, string(b))
			}
			query.Set("exclude", strings.Join(exclusions, ","))
		}
	}
	endpoint.RawQuery = query.Encode()

	// 3. Prepare and dispatch the HTTP context-bound channel
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create onecall request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute onecall request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from onecall API: %d", resp.StatusCode)
	}

	var result OneCallResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode onecall response: %w", err)
	}

	return &result, nil
}
