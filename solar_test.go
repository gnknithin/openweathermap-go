package openweathermap_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestGetSolarIrradianceForecast_Success(t *testing.T) {
	// 1. Define a mock JSON matching the exact OpenWeatherMap 3.0 solar data footprint
	mockJSON := `{
		"lat": 37.7749,
		"lon": -122.4194,
		"list": [
			{
				"dt": 1783080000,
				"ghi": 850.45,
				"dni": 920.12,
				"dhi": 110.34,
				"clear": 865.20
			}
		]
	}`

	// 2. Setup the local mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data/3.0/solar/forecast" {
			t.Errorf("expected path '/data/3.0/solar/forecast', got '%s'", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("lat") != "37.774900" || query.Get("lon") != "-122.419400" {
			t.Errorf("unexpected spatial coordinates passed to query parameters")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockJSON))
	}))
	defer server.Close()

	// 3. Initialize the Client pointing at our mock server
	client := openweathermap.NewClient(
		"test-api-key",
		openweathermap.WithBaseURL(server.URL),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 4. Execute the solar endpoint method
	response, err := client.GetSolarIrradianceForecast(ctx, 37.7749, -122.4194)

	// 5. Run explicit field validations
	if err != nil {
		t.Fatalf("expected no connection errors, got: %v", err)
	}

	if response.Lat != 37.7749 || response.Lon != -122.4194 {
		t.Errorf("root alignment coordinates mapped incorrectly")
	}

	if len(response.List) != 1 {
		t.Fatalf("expected 1 item inside timeline telemetry block, got %d", len(response.List))
	}

	item := response.List[0]
	if item.GHI != 850.45 || item.DNI != 920.12 || item.DHI != 110.34 {
		t.Errorf("failed to extract raw scientific radiation metrics accurately")
	}

	if item.Clear != 865.20 {
		t.Errorf("failed to extract clear sky radiation metric accurately")
	}
}
