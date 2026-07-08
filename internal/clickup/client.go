package clickup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"clickup-mcp/internal/config"
)

type apiVersion string

const (
	apiV2 apiVersion = "v2"
	apiV3 apiVersion = "v3"
)

// Client is a typed, retrying HTTP client for the ClickUp REST API.
type Client struct {
	httpClient  *http.Client
	baseURLv2   string
	baseURLv3   string
	token       string
	defaultTeam string
	maxRetries  int
	logger      *slog.Logger
}

// NewClient builds a Client from a loaded config.Config.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient:  &http.Client{Timeout: cfg.HTTPTimeout},
		baseURLv2:   cfg.BaseURLv2,
		baseURLv3:   cfg.BaseURLv3,
		token:       cfg.APIToken,
		defaultTeam: cfg.TeamID,
		maxRetries:  cfg.MaxRetries,
		// stderr only: stdout is the stdio MCP JSON-RPC channel.
		logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
}

// DefaultTeamID returns the workspace ID to fall back to when a tool call
// doesn't supply one explicitly.
func (c *Client) DefaultTeamID() string { return c.defaultTeam }

type requestParams struct {
	Method     string
	APIVersion apiVersion // "" defaults to v2
	Path       string     // leading slash, path params already substituted
	Query      url.Values
	Body       any
}

// do executes a JSON request/response round trip with retry/backoff, and
// decodes a successful response body into out (skipped if out is nil).
func (c *Client) do(ctx context.Context, p requestParams, out any) error {
	base := c.baseURLv2
	if p.APIVersion == apiV3 {
		base = c.baseURLv3
	}
	fullURL := base + p.Path
	if len(p.Query) > 0 {
		fullURL += "?" + p.Query.Encode()
	}

	var bodyBytes []byte
	if p.Body != nil {
		b, err := json.Marshal(p.Body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyBytes = b
	}

	return c.doRaw(ctx, func(ctx context.Context) (*http.Request, error) {
		var reader io.Reader
		if bodyBytes != nil {
			reader = bytes.NewReader(bodyBytes)
		}
		req, err := http.NewRequestWithContext(ctx, p.Method, fullURL, reader)
		if err != nil {
			return nil, err
		}
		if bodyBytes != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		return req, nil
	}, out)
}

// doRaw executes a request/response round trip with retry/backoff for a
// request built fresh on every attempt by newRequest (so non-replayable
// bodies, e.g. multipart uploads, are rebuilt per attempt). Auth/Accept
// headers are injected on every attempt regardless of what newRequest sets.
func (c *Client) doRaw(ctx context.Context, newRequest func(ctx context.Context) (*http.Request, error), out any) error {
	for attempt := 0; ; attempt++ {
		req, err := newRequest(ctx)
		if err != nil {
			return fmt.Errorf("build request: %w", err)
		}
		req.Header.Set("Authorization", c.token)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if attempt >= c.maxRetries {
				return fmt.Errorf("clickup request failed: %w", err)
			}
			if !c.sleep(ctx, exponentialBackoff(attempt)) {
				return ctx.Err()
			}
			continue
		}

		retry, delay, callErr := c.handleResponse(attempt, resp, out)
		if !retry || attempt >= c.maxRetries {
			return callErr
		}
		c.logger.Warn("retrying clickup request", "attempt", attempt+1, "delay", delay, "error", callErr)
		if !c.sleep(ctx, delay) {
			return ctx.Err()
		}
	}
}

func (c *Client) handleResponse(attempt int, resp *http.Response, out any) (retry bool, delay time.Duration, err error) {
	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return false, 0, fmt.Errorf("read response body: %w", readErr)
	}

	switch {
	case resp.StatusCode == http.StatusTooManyRequests:
		d := computeBackoff(attempt, resp.Header.Get("Retry-After"), resp.Header.Get("X-RateLimit-Reset"), time.Now())
		return true, d, decodeAPIError(resp.StatusCode, body)
	case resp.StatusCode >= 500:
		return true, exponentialBackoff(attempt), decodeAPIError(resp.StatusCode, body)
	case resp.StatusCode >= 400:
		return false, 0, decodeAPIError(resp.StatusCode, body)
	case resp.StatusCode == http.StatusNoContent || out == nil || len(body) == 0:
		return false, 0, nil
	default:
		if err := json.Unmarshal(body, out); err != nil {
			return false, 0, fmt.Errorf("decode response body: %w", err)
		}
		return false, 0, nil
	}
}

func (c *Client) sleep(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return true
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

// computeBackoff picks a retry delay: Retry-After header, then
// X-RateLimit-Reset header, then exponential backoff as a last resort. It's a
// pure function of its inputs so it's testable without real sleeps.
func computeBackoff(attempt int, retryAfterHeader, rateLimitResetHeader string, now time.Time) time.Duration {
	if retryAfterHeader != "" {
		if secs, err := strconv.Atoi(strings.TrimSpace(retryAfterHeader)); err == nil && secs >= 0 {
			return time.Duration(secs) * time.Second
		}
	}
	if rateLimitResetHeader != "" {
		if ts, err := strconv.ParseInt(strings.TrimSpace(rateLimitResetHeader), 10, 64); err == nil {
			if d := time.Unix(ts, 0).Sub(now); d > 0 {
				return d
			}
		}
	}
	return exponentialBackoff(attempt)
}

func exponentialBackoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	const (
		base     = 500 * time.Millisecond
		maxDelay = 30 * time.Second
	)
	delay := base * time.Duration(math.Pow(2, float64(attempt)))
	if delay > maxDelay {
		delay = maxDelay
	}
	jitter := time.Duration(rand.Int63n(int64(delay)/2 + 1))
	return delay + jitter
}
