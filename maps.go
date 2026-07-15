package openweathermap

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetMapTile fetches a raw binary image tile stream for a specific layer and coordinate grid.
// Crucial Design Note: The caller MUST call `.Close()` on the returned io.ReadCloser
// when they are done reading the stream to prevent severe resource leaks.
func (c *Client) GetMapTile(ctx context.Context, layer MapLayer, z, x, y int) (io.ReadCloser, error) {
	// 1. Target the standard map path structure: /map/2.0/{layer}/{z}/{x}/{y}
	endpoint, err := url.Parse(fmt.Sprintf("%s/map/2.0/%s/%d/%d/%d", c.baseURL, layer, z, x, y))
	if err != nil {
		return nil, fmt.Errorf("failed to parse map tile URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	// 2. Set up the context-bound HTTP stream request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create map tile request: %w", err)
	}

	// 3. Dispatch execution
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute map tile request: %w", err)
	}

	// 4. Handle boundary status failure checks perfectly
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close() // Safely kill body allocation loop right away on errors
		return nil, fmt.Errorf("unexpected status code from map tile API: %d", resp.StatusCode)
	}

	// 5. Pass the raw body pointer back out to the developer completely intact!
	return resp.Body, nil
}
