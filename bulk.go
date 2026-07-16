package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// BulkDownloadRequest defines the target parameter set for historical bulk archives.
type BulkDownloadRequest struct {
	LocationID string `json:"location_id"`
	Type       string `json:"type"`
}

// StreamBulkRecords executes a massive bulk data download and streams the records sequentially.
// This prevents high RAM consumption by decoding elements one by one out of the network stream socket.
func (c *Client) StreamBulkRecords(ctx context.Context, reqPayload BulkDownloadRequest, callback func(jsonRecord []byte) error) error {
	endpoint, err := url.Parse(fmt.Sprintf("%s/data/3.0/bulk/download", c.baseURL))
	if err != nil {
		return fmt.Errorf("failed to parse bulk data base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("location_id", reqPayload.LocationID)
	query.Set("type", reqPayload.Type)
	query.Set("appid", c.apiKey)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create bulk stream request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute bulk stream request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := c.checkResponse(resp); err != nil {
		return err
	}

	// 🏆 Initialize a Tokenized JSON Stream Decoder to extract elements one by one
	decoder := json.NewDecoder(resp.Body)

	// Read the open array token bracket '['
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to parse opening bulk array token: %w", err)
	}

	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("expected JSON array matrix start, got: %v", token)
	}

	// Loop dynamically through the active elements sequence inside the array stream
	for decoder.More() {
		var rawMessage json.RawMessage
		if err := decoder.Decode(&rawMessage); err != nil {
			return fmt.Errorf("failed to parse streaming bulk record token element: %w", err)
		}

		// Dispatch chunk directly to host application memory space via callback hook execution
		if err := callback(rawMessage); err != nil {
			return fmt.Errorf("application callback returned execution interruption error: %w", err)
		}
	}

	// Read the close array token bracket ']'
	_, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to parse closing bulk array token: %w", err)
	}

	return nil
}
