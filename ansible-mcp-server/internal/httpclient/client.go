package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"
)

// Client wraps go-retryablehttp with rate limiting and standard headers.
type Client struct {
	inner   *retryablehttp.Client
	global  *rate.Limiter
	session *rate.Limiter
	version string
}

// New creates a new HTTP client with the given settings.
func New(timeout time.Duration, maxRetries int, globalLimit, sessionLimit, version string) (*Client, error) {
	rc := retryablehttp.NewClient()
	rc.RetryMax = maxRetries
	rc.HTTPClient = &http.Client{Timeout: timeout}
	rc.RetryWaitMin = 1 * time.Second
	rc.RetryWaitMax = 30 * time.Second
	rc.CheckRetry = checkRetry
	rc.Backoff = backoff
	rc.Logger = nil // suppress default retry logging

	globalLimiter, err := parseLimiter(globalLimit)
	if err != nil {
		return nil, fmt.Errorf("invalid global rate limit: %w", err)
	}

	sessionLimiter, err := parseLimiter(sessionLimit)
	if err != nil {
		return nil, fmt.Errorf("invalid session rate limit: %w", err)
	}

	return &Client{
		inner:   rc,
		global:  globalLimiter,
		session: sessionLimiter,
		version: version,
	}, nil
}

// Get performs a GET request with rate limiting and retry logic.
func (c *Client) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	if err := c.global.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}
	if err := c.session.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	req.Header.Set("User-Agent", "ansible-mcp-server/"+c.version)
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.inner.Do(req)
}

func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if err != nil {
		return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	}
	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true, nil
	}
	return false, nil
}

func backoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	if resp != nil && resp.StatusCode == 429 {
		if reset := resp.Header.Get("x-ratelimit-reset"); reset != "" {
			if ts, err := strconv.ParseInt(reset, 10, 64); err == nil {
				wait := time.Until(time.Unix(ts, 0))
				if wait > 0 && wait <= max {
					return wait
				}
			}
		}
	}
	return retryablehttp.LinearJitterBackoff(min, max, attemptNum, resp)
}

func parseLimiter(spec string) (*rate.Limiter, error) {
	parts := strings.SplitN(spec, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected rate:burst format, got %q", spec)
	}
	r, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid rate %q: %w", parts[0], err)
	}
	b, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid burst %q: %w", parts[1], err)
	}
	return rate.NewLimiter(rate.Limit(r), b), nil
}
