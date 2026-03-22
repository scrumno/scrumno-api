package middleware

import (
	"net/http"

	"github.com/scrumno/scrumno-api/internal/api/utils"
)

// RequireOrganizationID проверяет наличие параметра organizationId в query
// и возвращает 400, если он не передан. Подходит для всех эндпоинтов iiko,
// которым нужна организация.
func RequireOrganizationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("organizationId") == "" {
			utils.JSONResponse(w, map[string]any{
				"isSuccess": false,
				"error":     "organizationId query parameter is required",
			}, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
