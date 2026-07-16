package openweathermap_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gnknithin/openweathermap-go"
)

func TestClientRetry_OnTransientFailure(t *testing.T) {
	var requestCounter int32

	// 1. Setup a dynamic mock server that simulates structural failures recovery
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCounter, 1)

		// On the first request, simulate an enterprise rate limit drop
		if count == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"cod":429, "message": "Rate limit exceeded"}`))
			return
		}

		// On the second request, recover automatically
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"coord":{"lon":10.99,"lat":44.34},"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}],"base":"stations","main":{"temp":289.5},"visibility":10000,"dt":1661870592,"id":3163858,"name":"Zocca","cod":200}`))
	}))
	defer server.Close()

	// 2. Initialize the Client pointing at our dynamic retry testing server
	client := openweathermap.NewClient(
		"test-api-key",
		openweathermap.WithBaseURL(server.URL),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 3. Trigger an evaluation execution call
	res, err := client.GetCurrentWeather(ctx, 44.34, 10.99)
	if err != nil {
		t.Fatalf("expected client to successfully transparently recover, got error: %v", err)
	}

	// 4. Run metrics verification assertions
	if res.Name != "Zocca" {
		t.Errorf("expected final successful response to match payload values, got name %s", res.Name)
	}

	finalCount := atomic.LoadInt32(&requestCounter)
	if finalCount != 2 {
		t.Errorf("expected exactly 2 network calls to execute via retry loop, got %d", finalCount)
	}
}
