package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// TokenRefreshFunc should return a fresh iiko access token.
// It is called only when the original request was answered with HTTP 401.
type TokenRefreshFunc func(ctx context.Context) (string, error)

// AuthRefreshRoundTripper retries the request once after HTTP 401 by:
// 1) calling refreshFunc to get a new token
// 2) updating Authorization header
// 3) replaying request body (if any)
//
// It is local to the given http.Client, so it won't affect other parts of the backend.
type AuthRefreshRoundTripper struct {
	base         http.RoundTripper
	refreshFunc  TokenRefreshFunc
	tokenPtr     *string
	bearerPrefix string

	// guard prevents concurrent refresh storms
	mu sync.Mutex
}

const retryHeader = "X-IIKO-Token-Refreshed"

func NewAuthRefreshRoundTripper(base http.RoundTripper, tokenPtr *string, refreshFunc TokenRefreshFunc) *AuthRefreshRoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &AuthRefreshRoundTripper{
		base:         base,
		refreshFunc:  refreshFunc,
		tokenPtr:     tokenPtr,
		bearerPrefix: "Bearer ",
	}
}

func (t *AuthRefreshRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t == nil || t.base == nil {
		return nil, fmt.Errorf("auth roundtripper is not configured")
	}
	if t.refreshFunc == nil {
		return t.base.RoundTrip(req)
	}
	if req.Header.Get(retryHeader) == "1" {
		// retry already attempted for this request chain
		return t.base.RoundTrip(req)
	}

	bodyBytes, err := readAndRestoreRequestBody(req)
	if err != nil {
		return nil, err
	}

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}

	// If refresh fails, we must return the original response (including resp.Body).
	newToken, refreshErr := t.refreshTokenOnce(req.Context())
	if refreshErr != nil || newToken == "" {
		return resp, nil
	}

	// We will retry, so close the original body to avoid leaks.
	// Caller will receive the response from the retry below.
	_ = drainAndClose(resp.Body)

	req2 := req.Clone(req.Context())
	req2.Header = req.Header.Clone()
	req2.Header.Set(retryHeader, "1")
	req2.Header.Set("Authorization", t.bearerPrefix+newToken)
	if bodyBytes != nil {
		req2.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	return t.base.RoundTrip(req2)
}

func (t *AuthRefreshRoundTripper) refreshTokenOnce(ctx context.Context) (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	token, err := t.refreshFunc(ctx)
	if err != nil {
		return "", err
	}
	if t.tokenPtr != nil {
		*t.tokenPtr = token
	}
	return token, nil
}

func readAndRestoreRequestBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}

	b, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	// restore so base.RoundTrip can read it
	req.Body = io.NopCloser(bytes.NewReader(b))
	return b, nil
}

func drainAndClose(rc io.ReadCloser) error {
	if rc == nil {
		return nil
	}
	_, _ = io.Copy(io.Discard, rc)
	return rc.Close()
}
