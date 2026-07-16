package openweathermap_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestStreamBulkRecords_Success(t *testing.T) {
	// 1. Simulate a multi-element JSON array data payload streaming over a connection socket
	mockBulkStreamJSON := `[
		{"record_id": 1, "temp": 298.15},
		{"record_id": 2, "temp": 299.50},
		{"record_id": 3, "temp": 300.12}
	]`

	// 2. Setup the streaming mock HTTP server environment
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data/3.0/bulk/download" {
			t.Errorf("expected path '/data/3.0/bulk/download', got '%s'", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("location_id") != "loc_usa_99" || query.Get("type") != "hourly" {
			t.Errorf("unexpected query parameters passed to bulk stream engine")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBulkStreamJSON))
	}))
	defer server.Close()

	// 3. Initialize the Client pointing at our streaming server
	client := openweathermap.NewClient(
		"test-api-key",
		openweathermap.WithBaseURL(server.URL),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	reqPayload := openweathermap.BulkDownloadRequest{
		LocationID: "loc_usa_99",
		Type:       "hourly",
	}

	var extractedRecords [][]byte

	// 4. Trigger the sequential low-RAM tokenized array execution loop method
	err := client.StreamBulkRecords(ctx, reqPayload, func(record []byte) error {
		extractedRecords = append(extractedRecords, record)
		return nil
	})

	// 5. Run structural evaluations and assertions
	if err != nil {
		t.Fatalf("expected clean bulk array stream processing, got: %v", err)
	}

	if len(extractedRecords) != 3 {
		t.Fatalf("expected exactly 3 elements to stream out of loop, got %d", len(extractedRecords))
	}

	// Verify structural values inside parsed slice items
	var firstRecord struct {
		RecordID int     `json:"record_id"`
		Temp     float64 `json:"temp"`
	}
	if err := json.Unmarshal(extractedRecords[0], &firstRecord); err != nil {
		t.Fatalf("failed to decode extracted record chunk: %v", err)
	}

	if firstRecord.RecordID != 1 || firstRecord.Temp != 298.15 {
		t.Errorf("stream chunk values mapped incorrectly inside host space app loop")
	}
}
