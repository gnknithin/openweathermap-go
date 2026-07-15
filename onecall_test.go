package openweathermap_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestGetOneCallWeather_SuccessWithOptions(t *testing.T) {
	// 1. Setup a rich mock JSON containing the core parent fields, current metrics, and hourly slices
	mockJSON := `{
		"lat": 33.44,
		"lon": -94.04,
		"timezone": "America/Chicago",
		"timezone_offset": -18000,
		"current": {
			"dt": 1684929490,
			"temp": 292.55,
			"feels_like": 292.87,
			"pressure": 1014,
			"humidity": 65,
			"weather": [{"id": 801, "main": "Clouds", "description": "few clouds", "icon": "02d"}]
		},
		"hourly": [
			{
				"dt": 1684926000,
				"temp": 292.01,
				"feels_like": 292.33,
				"pressure": 1014,
				"humidity": 68,
				"weather": [{"id": 801, "main": "Clouds", "description": "few clouds", "icon": "02d"}]
			}
		]
	}`

	// 2. Setup the local mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data/3.0/onecall" {
			t.Errorf("expected path '/data/3.0/onecall', got '%s'", r.URL.Path)
		}

		// Verify coordinates and auth parameter mappings
		query := r.URL.Query()
		if query.Get("lat") != "33.440000" || query.Get("lon") != "-94.040000" {
			t.Errorf("unexpected coordinate arguments mapped to URL query")
		}

		// Validate that our optional configurations successfully mapped to query params
		if query.Get("units") != "metric" {
			t.Errorf("expected units parameter 'metric', got '%s'", query.Get("units"))
		}
		if query.Get("exclude") != "minutely,daily" {
			t.Errorf("expected format execution for exclusions, got '%s'", query.Get("exclude"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockJSON))
	}))
	defer server.Close()

	// 3. Initialize Client pointing at the mock server
	client := openweathermap.NewClient(
		"test-api-key",
		openweathermap.WithBaseURL(server.URL),
	)

	// 4. Construct our advanced configuration payload
	opts := &openweathermap.OneCallOptions{
		Units:   "metric",
		Exclude: []openweathermap.ExcludeBlock{openweathermap.ExcludeMinutely, openweathermap.ExcludeDaily},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 5. Execute method
	response, err := client.GetOneCallWeather(ctx, 33.44, -94.04, opts)

	// 6. Assertions
	if err != nil {
		t.Fatalf("expected no execution errors, got: %v", err)
	}

	if response.Timezone != "America/Chicago" {
		t.Errorf("expected timezone payload string match failed")
	}

	if response.Current.Temp != 292.55 {
		t.Errorf("failed to map base current metrics cleanly")
	}

	if len(response.Hourly) != 1 || response.Hourly[0].Temp != 292.01 {
		t.Errorf("failed to accurately parse dynamic nested slices within the response array")
	}
}
