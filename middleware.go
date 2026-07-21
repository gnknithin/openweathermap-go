package openweathermap

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// retryableTransport wraps a standard http.RoundTripper to inject automatic
// exponential backoff and retry handling when encountering HTTP 429 or 5xx status codes.
type retryableTransport struct {
	next       http.RoundTripper
	maxRetries int
	baseDelay  time.Duration
}

// RoundTrip intercepts the standard request lifecycle to evaluate errors and auto-retry if viable.
func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= t.maxRetries; i++ {
		// 🏆 Fixed: Check context error type accurately before attempting execution
		if req.Context().Err() != nil {
			return nil, req.Context().Err()
		}

		var reqClone *http.Request
		if req.Body != nil {
			bodyBytes, cloneErr := io.ReadAll(req.Body)
			if cloneErr != nil {
				return nil, cloneErr
			}
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			reqClone = req.Clone(req.Context())
			reqClone.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		} else {
			reqClone = req.Clone(req.Context())
		}

		resp, err = t.next.RoundTrip(reqClone)

		if err == nil && resp.StatusCode != http.StatusTooManyRequests && resp.StatusCode < 500 {
			return resp, nil
		}

		if i == t.maxRetries {
			return resp, err
		}

		if resp != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}

		// 🏆 Modern Go 1.20+ thread-safe randomized exponential backoff with jitter
		jitter := time.Duration(rand.Intn(50)) * time.Millisecond
		backoffDelay := (t.baseDelay * (1 << uint(i))) + jitter

		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(backoffDelay):
		}
	}

	return resp, err
}
