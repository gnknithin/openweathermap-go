package openweathermap_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestGetHistoricalWeather_Success(t *testing.T) {
	// 1. Define a mock JSON matching the One Call 3.0 timemachine payload format
	mockJSON := `{
		"lat": 40.7128,
		"lon": -74.0060,
		"timezone": "America/New_York",
		"timezone_offset": -14400,
		"data": [
			{
				"dt": 1691215600,
				"temp": 294.15,
				"feels_like": 294.50,
				"pressure": 1012,
				"humidity": 72,
				"dew_point": 288.75,
				"clouds": 40,
				"wind_speed": 4.6,
				"wind_deg": 180,
				"weather": [
					{
						"id": 802,
						"main": "Clouds",
						"description": "scattered clouds",
						"icon": "03d"
					}
				]
			}
		]
	}`

	// 2. Setup the local mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data/3.0/onecall/timemachine" {
			t.Errorf("expected path '/data/3.0/onecall/timemachine', got '%s'", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("lat") != "40.712800" || query.Get("lon") != "-74.006000" {
			t.Errorf("unexpected spatial coordinates mapped to query parameters")
		}
		if query.Get("dt") != "1691234567" {
			t.Errorf("expected historical timestamp '1691234567', got '%s'", query.Get("dt"))
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

	// 4. Execute the historical method using a test Unix timestamp
	response, err := client.GetHistoricalWeather(ctx, 40.7128, -74.0060, 1691234567)

	// 5. Run structural validations
	if err != nil {
		t.Fatalf("expected no connection errors, got: %v", err)
	}

	if response.Timezone != "America/New_York" || response.TimezoneOffset != -14400 {
		t.Errorf("timezone historical metadata parsed incorrectly")
	}

	if len(response.Data) != 1 {
		t.Fatalf("expected 1 record block inside data history timeline, got %d", len(response.Data))
	}

	historyItem := response.Data[0]
	if historyItem.Temp != 294.15 || historyItem.Humidity != 72 {
		t.Errorf("failed to map base metrics inside historical payload item")
	}

	if len(historyItem.Weather) == 0 || historyItem.Weather[0].Main != "Clouds" {
		t.Errorf("reused weather nested block slice failed structural alignment checks")
	}
}
