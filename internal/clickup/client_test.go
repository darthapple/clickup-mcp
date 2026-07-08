package clickup

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"clickup-mcp/internal/config"
)

func testClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := NewClient(&config.Config{
		APIToken:    "pk_test_token",
		TeamID:      "123",
		BaseURLv2:   srv.URL,
		BaseURLv3:   srv.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  2,
	})
	return c, srv
}

func TestDoSendsRawAuthHeader(t *testing.T) {
	var gotAuth string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	if err := c.do(context.Background(), requestParams{Method: http.MethodGet, Path: "/x"}, &map[string]any{}); err != nil {
		t.Fatalf("do: %v", err)
	}
	if gotAuth != "pk_test_token" {
		t.Errorf("Authorization header = %q, want raw token (no Bearer prefix)", gotAuth)
	}
}

func TestDoDecodesSuccessBody(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"42","name":"hello"}`))
	})

	var out map[string]any
	if err := c.do(context.Background(), requestParams{Method: http.MethodGet, Path: "/x"}, &out); err != nil {
		t.Fatalf("do: %v", err)
	}
	if out["id"] != "42" || out["name"] != "hello" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestDoDecodesAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"err":"Task not found","ECODE":"TASK_001"}`))
	})

	err := c.do(context.Background(), requestParams{Method: http.MethodGet, Path: "/x"}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not an *APIError: %v", err)
	}
	if apiErr.StatusCode != http.StatusNotFound || apiErr.ECode != "TASK_001" || apiErr.Err != "Task not found" {
		t.Errorf("decoded APIError = %+v", apiErr)
	}
	if !IsNotFound(err) {
		t.Error("IsNotFound(err) = false, want true")
	}
	if IsRateLimited(err) {
		t.Error("IsRateLimited(err) = true, want false")
	}
}

func TestDoRetriesOn429ThenSucceeds(t *testing.T) {
	var attempts int32
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&attempts, 1) == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"err":"rate limited","ECODE":"RATE_001"}`))
			return
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	var out map[string]any
	if err := c.do(context.Background(), requestParams{Method: http.MethodGet, Path: "/x"}, &out); err != nil {
		t.Fatalf("do: %v", err)
	}
	if atomic.LoadInt32(&attempts) != 2 {
		t.Errorf("attempts = %d, want 2", attempts)
	}
	if out["ok"] != true {
		t.Errorf("decoded = %+v", out)
	}
}

func TestDoGivesUpAfterMaxRetries(t *testing.T) {
	var attempts int32
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"err":"boom","ECODE":"SRV_001"}`))
	})

	err := c.do(context.Background(), requestParams{Method: http.MethodGet, Path: "/x"}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// maxRetries=2 means 1 initial attempt + 2 retries = 3 total.
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Errorf("attempts = %d, want 3", got)
	}
}

func TestDoContextCancellationAbortsQuickly(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"err":"boom","ECODE":"SRV_001"}`))
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := c.do(ctx, requestParams{Method: http.MethodGet, Path: "/x"}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("cancellation took too long: %v", elapsed)
	}
}

func TestArrayQueryParamEncoding(t *testing.T) {
	var gotQuery url.Values
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{}`))
	})

	q := url.Values{}
	addArrayParam(q, "statuses[]", []string{"open", "in progress"})
	if err := c.do(context.Background(), requestParams{Method: http.MethodGet, Path: "/x", Query: q}, &map[string]any{}); err != nil {
		t.Fatalf("do: %v", err)
	}
	got := gotQuery["statuses[]"]
	if len(got) != 2 || got[0] != "open" || got[1] != "in progress" {
		t.Errorf("statuses[] query = %v", got)
	}
}

func TestComputeBackoff(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("retry-after header wins", func(t *testing.T) {
		d := computeBackoff(0, "5", "", now)
		if d != 5*time.Second {
			t.Errorf("got %v, want 5s", d)
		}
	})

	t.Run("rate limit reset header used when no retry-after", func(t *testing.T) {
		reset := now.Add(10 * time.Second).Unix()
		d := computeBackoff(0, "", strconv.FormatInt(reset, 10), now)
		if d <= 9*time.Second || d > 10*time.Second {
			t.Errorf("got %v, want ~10s", d)
		}
	})

	t.Run("falls back to exponential backoff", func(t *testing.T) {
		d := computeBackoff(0, "", "", now)
		if d <= 0 || d > 30*time.Second {
			t.Errorf("got %v, want a positive bounded backoff", d)
		}
	})
}
