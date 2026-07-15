package openweathermap_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestGetCurrentWeather_Success(t *testing.T) {
	// 1. Define our expected mock JSON response payload
	mockJSON := `{
		"coord": {"lon": 10.99, "lat": 44.34},
		"weather": [{"id": 804, "main": "Clouds", "description": "overcast clouds", "icon": "04d"}],
		"main": {"temp": 289.92, "feels_like": 289.53, "temp_min": 287.88, "temp_max": 290.44, "pressure": 1017, "humidity": 77},
		"name": "Zocca",
		"cod": 200
	}`

	// 2. Set up the local mock HTTP server in memory
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that our request contains the proper parameters
		if r.URL.Path != "/data/2.5/weather" {
			t.Errorf("expected path '/data/2.5/weather', got '%s'", r.URL.Path)
		}

		// Verify query parameters are passed down properly
		query := r.URL.Query()
		if query.Get("lat") != "44.340000" || query.Get("lon") != "10.990000" {
			t.Errorf("unexpected lat/lon parameters received")
		}
		if query.Get("appid") != "test-api-key" {
			t.Errorf("expected appid 'test-api-key', got '%s'", query.Get("appid"))
		}

		// Respond with HTTP 200 OK and our mock weather data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockJSON))
	}))
	defer server.Close() // Ensure the local port is closed when the test ends

	// 3. Initialize our SDK Client pointing at our mock server URL instead of production
	client := openweathermap.NewClient(
		"test-api-key",
		openweathermap.WithBaseURL(server.URL),
	)

	// 4. Execute the method we are testing
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response, err := client.GetCurrentWeather(ctx, 44.34, 10.99)

	// 5. Run assertions on our output data
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if response.Name != "Zocca" {
		t.Errorf("expected location name 'Zocca', got '%s'", response.Name)
	}

	if response.Main.Temp != 289.92 {
		t.Errorf("expected temperature 289.92, got %f", response.Main.Temp)
	}

	if len(response.Weather) == 0 || response.Weather[0].Main != "Clouds" {
		t.Errorf("failed to accurately parse nested slice weather data")
	}
}

func TestGetCurrentWeather_FailureStatus(t *testing.T) {
	// 1. Set up a mock server that explicitly returns an HTTP 401 Unauthorized error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"cod":401, "message": "Invalid API key"}`))
	}))
	defer server.Close()

	// 2. Initialize our client pointing to the failing mock server
	client := openweathermap.NewClient(
		"bad-api-key",
		openweathermap.WithBaseURL(server.URL),
	)

	// 3. Execute the method
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response, err := client.GetCurrentWeather(ctx, 44.34, 10.99)

	// 4. Assertions: Expect an error and NO response structure
	if err == nil {
		t.Fatalf("expected an error due to 401 status code, got nil")
	}

	if response != nil {
		t.Errorf("expected response to be nil on failure, got %v", response)
	}

	// 5. 🏆 Gold-Standard Type Assertion Verification
	// We check if the error is exactly the custom *APIError type we defined
	apiErr, ok := err.(*openweathermap.APIError)
	if !ok {
		t.Fatalf("expected returned error to be type *openweathermap.APIError, got %T", err)
	}

	if apiErr.Code != 401 {
		t.Errorf("expected APIError.Code to be 401, got %d", apiErr.Code)
	}

	expectedMsg := "Invalid API key"
	if apiErr.Message != expectedMsg {
		t.Errorf("expected APIError.Message to be '%s', got '%s'", expectedMsg, apiErr.Message)
	}
}
