package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/shared/jwt"
)

type contextKey string

const ContextKeyClaims contextKey = "jwt_claims"

var (
	ErrAuthorizationHeaderRequired      = errors.New("Обязателен заголовок Authorization")
	ErrInvalidAuthorizationHeaderFormat = errors.New("Неверный формат заголовка Authorization")
	ErrTokenRequired                    = errors.New("Обязателен токен")
)

type AuthMiddleware struct {
	jwtManager *jwt.Manager
}

func NewAuthMiddleware(jwtManager *jwt.Manager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")

		token, err := validateHeaderParam(auth)
		if err != nil {
			utils.JSONResponse(w, ErrorResponse{
				IsSuccess: false,
				Error:     err.Error(),
			}, http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtManager.ValidateAccessToken(token)
		if err != nil {
			utils.JSONResponse(w, ErrorResponse{
				IsSuccess: false,
				Error:     err.Error(),
			}, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyClaims, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClaimsFromRequest(r *http.Request) *jwt.Claims {
	return claimsFromContext(r.Context())
}

func validateHeaderParam(auth string) (string, error) {
	if auth == "" {
		return "", ErrAuthorizationHeaderRequired
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrInvalidAuthorizationHeaderFormat
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrTokenRequired
	}

	return token, nil
}

func claimsFromContext(ctx context.Context) *jwt.Claims {
	v := ctx.Value(ContextKeyClaims)
	if v == nil {
		return nil
	}
	c, _ := v.(*jwt.Claims)
	return c
}

type ErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}
