package middleware

import (
	"context"
	"net/http"
	"strings"
)

type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
	ForceRefreshToken(ctx context.Context) (string, error)
}

type AuthRetryTransport struct {
	base          http.RoundTripper
	tokenProvider TokenProvider
}

func NewAuthRetryTransport(base http.RoundTripper, tokenProvider TokenProvider) *AuthRetryTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &AuthRetryTransport{
		base:          base,
		tokenProvider: tokenProvider,
	}
}

func (t *AuthRetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Для получения access_token не добавляем bearer и не делаем ретрай,
	// иначе можно попасть в рекурсию.
	if isAccessTokenRequest(req) {
		return t.base.RoundTrip(req)
	}

	token, err := t.tokenProvider.GetToken(req.Context())
	if err != nil {
		return nil, err
	}

	firstReq, canRetry, err := cloneRequestForAttempt(req, token, false)
	if err != nil {
		return nil, err
	}

	resp, err := t.base.RoundTrip(firstReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusUnauthorized || !canRetry {
		return resp, nil
	}

	_ = resp.Body.Close()

	newToken, err := t.tokenProvider.ForceRefreshToken(req.Context())
	if err != nil {
		return nil, err
	}

	retryReq, _, err := cloneRequestForAttempt(req, newToken, true)
	if err != nil {
		return nil, err
	}
	return t.base.RoundTrip(retryReq)
}

func cloneRequestForAttempt(req *http.Request, token string, forceGetBody bool) (*http.Request, bool, error) {
	cloned := req.Clone(req.Context())

	canRetry := req.Body == nil || req.GetBody != nil
	if req.Body != nil {
		if forceGetBody || req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, false, err
			}
			cloned.Body = body
		} else {
			cloned.Body = req.Body
		}
	}

	cloned.Header = req.Header.Clone()
	cloned.Header.Set("Authorization", "Bearer "+token)
	return cloned, canRetry, nil
}

func isAccessTokenRequest(req *http.Request) bool {
	if req == nil || req.URL == nil {
		return false
	}
	return strings.HasSuffix(req.URL.Path, "/api/1/access_token")
}
