package openweathermap_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestGeocodeCity_Success(t *testing.T) {
	// 1. Define our expected mock JSON array payload
	mockJSON := `[
		{
			"name": "London",
			"lat": 51.5073219,
			"lon": -0.1276474,
			"country": "GB",
			"state": "England"
		}
	]`

	// 2. Setup the local mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/geo/1.0/direct" {
			t.Errorf("expected path '/geo/1.0/direct', got '%s'", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("q") != "London" {
			t.Errorf("expected query param q to be 'London', got '%s'", query.Get("q"))
		}
		if query.Get("limit") != "1" {
			t.Errorf("expected limit to be '1', got '%s'", query.Get("limit"))
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

	// 4. Run the method
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	locations, err := client.GeocodeCity(ctx, "London", 1)

	// 5. Assertions
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(locations) != 1 {
		t.Fatalf("expected exactly 1 location returned, got %d", len(locations))
	}

	if locations[0].Name != "London" {
		t.Errorf("expected name 'London', got '%s'", locations[0].Name)
	}

	if locations[0].Latitude != 51.5073219 {
		t.Errorf("expected lat 51.5073219, got %f", locations[0].Latitude)
	}
}
