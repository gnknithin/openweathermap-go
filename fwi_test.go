package openweathermap_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestGetFireWeatherIndex_Success(t *testing.T) {
	// 1. Define a mock JSON payload matching the OpenWeatherMap FWI output footprint
	mockJSON := `{
		"lat": -23.5505,
		"lon": -46.6333,
		"list": [
			{
				"dt": 1783080000,
				"fwi": 18.75
			}
		]
	}`

	// 2. Setup the local mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data/3.0/fwi" {
			t.Errorf("expected path '/data/3.0/fwi', got '%s'", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("lat") != "-23.550500" || query.Get("lon") != "-46.633300" {
			t.Errorf("unexpected spatial coordinates mapped to query parameters")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockJSON))
	}))
	defer server.Close()

	// 3. Initialize the Client pointing at our mock server URL
	client := openweathermap.NewClient(
		"test-api-key",
		openweathermap.WithBaseURL(server.URL),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 4. Execute the FWI method
	response, err := client.GetFireWeatherIndex(ctx, -23.5505, -46.6333)

	// 5. Run explicit field validation assertions
	if err != nil {
		t.Fatalf("expected no connection errors, got: %v", err)
	}

	if response.Lat != -23.5505 || response.Lon != -46.6333 {
		t.Errorf("root alignment coordinates mapped incorrectly")
	}

	if len(response.List) != 1 {
		t.Fatalf("expected 1 item inside FWI telemetry block, got %d", len(response.List))
	}

	item := response.List[0]
	if item.FWI != 18.75 || item.Time != 1783080000 {
		t.Errorf("failed to extract raw scientific FWI values accurately")
	}
}
