package openweathermap_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestGetAirPollution_Success(t *testing.T) {
	// 1. Define a rich mock JSON mirroring the array-style coordinate payload layout
	mockJSON := `{
		"coord": [50.0, 50.0],
		"list": [
			{
				"dt": 1601868000,
				"main": {
					"aqi": 1
				},
				"components": {
					"co": 201.94,
					"no": 0.01,
					"no2": 0.23,
					"o3": 69.38,
					"so2": 0.34,
					"pm2_5": 0.5,
					"pm10": 0.54,
					"nh3": 0.12
				}
			}
		]
	}`

	// 2. Setup the local mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data/2.5/air_pollution" {
			t.Errorf("expected path '/data/2.5/air_pollution', got '%s'", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("lat") != "50.000000" || query.Get("lon") != "50.000000" {
			t.Errorf("unexpected coordinate parameters forwarded to query strings")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockJSON))
	}))
	defer server.Close()

	// 3. Initialize Client pointing at the mock server URL
	client := openweathermap.NewClient(
		"test-api-key",
		openweathermap.WithBaseURL(server.URL),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 4. Execute method
	response, err := client.GetAirPollution(ctx, 50.0, 50.0)

	// 5. Run validations across our array-bound types definitions
	if err != nil {
		t.Fatalf("expected no execution errors, got: %v", err)
	}

	if len(response.Coord) != 2 || response.Coord[0] != 50.0 || response.Coord[1] != 50.0 {
		t.Errorf("failed to accurately parse array-based root coordinates layout")
	}

	if len(response.List) != 1 {
		t.Fatalf("expected 1 element inside timeline metrics block, got %d", len(response.List))
	}

	if response.List[0].Main.AQI != 1 {
		t.Errorf("expected AQI metric value '1', got %d", response.List[0].Main.AQI)
	}

	if response.List[0].Components.PM25 != 0.5 || response.List[0].Components.O3 != 69.38 {
		t.Errorf("failed to capture micro-concentration chemical pollutant components accurately")
	}
}
