package openweathermap_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestGetMapTile_Success(t *testing.T) {
	// 1. Generate a mock binary payload (representing fake raw PNG file data)
	mockImageBytes := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	// 2. Setup the mock HTTP streaming server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/map/2.0/clouds_new/6/20/35"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("appid") != "test-api-key" {
			t.Errorf("missing or invalid appid auth token")
		}

		// Set binary headers and write raw stream bytes
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockImageBytes)
	}))
	defer server.Close()

	// 3. Initialize Client pointing at our mock server URL
	client := openweathermap.NewClient(
		"test-api-key",
		openweathermap.WithBaseURL(server.URL),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 4. Execute the streaming map method
	stream, err := client.GetMapTile(ctx, openweathermap.LayerClouds, 6, 20, 35)
	if err != nil {
		t.Fatalf("expected no connection errors, got: %v", err)
	}
	defer func() {
		_ = stream.Close() // 🏆 Crucial step: Freeing memory buffers after reading
	}()

	// 5. Read all streaming bytes into memory for assertion check
	resultBytes, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("failed to read complete data stream: %v", err)
	}

	// 6. Assertions
	if len(resultBytes) != len(mockImageBytes) {
		t.Errorf("byte mismatch: expected size %d, got %d", len(mockImageBytes), len(resultBytes))
	}

	for i, b := range resultBytes {
		if b != mockImageBytes[i] {
			t.Errorf("data corruption detected at index %d", i)
		}
	}
}
