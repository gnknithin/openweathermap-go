package openweathermap_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestRegisterStation_Success(t *testing.T) {
	mockResponseJSON := `{
		"id": "60a1b2c3d4e5f60001234567",
		"external_id": "SF_STATION_001",
		"name": "San Francisco Central Station",
		"latitude": 37.7749,
		"longitude": -122.4194,
		"altitude": 15.5,
		"rank": 1,
		"created_at": "2026-07-16T00:00:00Z"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected Method 'POST', got '%s'", r.Method)
		}
		if r.URL.Path != "/data/3.0/stations" {
			t.Errorf("expected path '/data/3.0/stations', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(mockResponseJSON))
	}))
	defer server.Close()

	client := openweathermap.NewClient("test-api-key", openweathermap.WithBaseURL(server.URL))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	reqPayload := openweathermap.StationRegisterRequest{
		ExternalID: "SF_STATION_001",
		Name:       "San Francisco Central Station",
		Latitude:   37.7749,
		Longitude:  -122.4194,
		Altitude:   15.5,
	}

	res, err := client.RegisterStation(ctx, reqPayload)
	if err != nil {
		t.Fatalf("expected no errors during registration, got: %v", err)
	}

	if res.ID != "60a1b2c3d4e5f60001234567" || res.Name != "San Francisco Central Station" {
		t.Errorf("station confirmation payload fields mapped incorrectly")
	}
}

func TestGetStationByID_Success(t *testing.T) {
	mockResponseJSON := `{
		"id": "station_abc123",
		"external_id": "TEST_02",
		"name": "Validation Node"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/data/3.0/stations/station_abc123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponseJSON))
	}))
	defer server.Close()

	client := openweathermap.NewClient("test-api-key", openweathermap.WithBaseURL(server.URL))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := client.GetStationByID(ctx, "station_abc123")
	if err != nil {
		t.Fatalf("expected clean station extraction, got: %v", err)
	}

	if res.ExternalID != "TEST_02" || res.Name != "Validation Node" {
		t.Errorf("extracted station details failed field verification assertions")
	}
}

func TestDeleteStation_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected Method 'DELETE', got '%s'", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := openweathermap.NewClient("test-api-key", openweathermap.WithBaseURL(server.URL))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.DeleteStation(ctx, "station_abc123")
	if err != nil {
		t.Fatalf("expected successful deletion, got error: %v", err)
	}
}
