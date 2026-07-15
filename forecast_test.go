package openweathermap_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestGet5DayForecast_Success(t *testing.T) {
	// 1. Define a rich mock JSON containing the core parent fields, list timeline array, and city block
	mockJSON := `{
		"cod": "200",
		"message": 0,
		"cnt": 1,
		"list": [
			{
				"dt": 1661871600,
				"main": {
					"temp": 296.76,
					"feels_like": 296.98,
					"temp_min": 296.76,
					"temp_max": 297.87,
					"pressure": 1015,
					"humidity": 69
				},
				"weather": [
					{
						"id": 500,
						"main": "Rain",
						"description": "light rain",
						"icon": "10d"
					}
				],
				"visibility": 10000,
				"pop": 0.32,
				"dt_txt": "2022-08-30 15:00:00"
			}
		],
		"city": {
			"id": 2643743,
			"name": "London",
			"coord": {
				"lat": 51.5073,
				"lon": -0.1277
			},
			"country": "GB",
			"population": 1000000,
			"timezone": 3600,
			"sunrise": 1661834187,
			"sunset": 1661883418
		}
	}`

	// 2. Setup the local mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data/2.5/forecast" {
			t.Errorf("expected path '/data/2.5/forecast', got '%s'", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("lat") != "51.507300" || query.Get("lon") != "-0.127700" {
			t.Errorf("unexpected coordinate parameters mapped to URL query")
		}
		if query.Get("units") != "imperial" || query.Get("lang") != "en" {
			t.Errorf("failed to accurately forward optional units and language parameters")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockJSON))
	}))
	defer server.Close()

	// 3. Initialize Client pointing at our mock server URL
	client := openweathermap.NewClient(
		"test-api-key",
		openweathermap.WithBaseURL(server.URL),
	)

	// 4. Construct configuration payload
	opts := &openweathermap.ForecastOptions{
		Units: "imperial",
		Lang:  "en",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 5. Execute method
	response, err := client.Get5DayForecast(ctx, 51.5073, -0.1277, opts)

	// 6. Run validations across our composed data structure fields
	if err != nil {
		t.Fatalf("expected no execution errors, got: %v", err)
	}

	if response.City.Name != "London" || response.City.Country != "GB" {
		t.Errorf("failed to properly unpack nested city metadata attributes")
	}

	if len(response.List) != 1 {
		t.Fatalf("expected 1 item inside timeline timeline block, got %d", len(response.List))
	}

	if response.List[0].Main.Temp != 296.76 {
		t.Errorf("reused main structure values failed to extract accurately")
	}

	if len(response.List[0].Weather) == 0 || response.List[0].Weather[0].Main != "Rain" {
		t.Errorf("reused secondary nested struct arrays failed to map cleanly")
	}
}
