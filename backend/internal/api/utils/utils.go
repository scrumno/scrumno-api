package utils

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Ошибка кодирвоать JSON:", "error", err.Error())
		return
	}
}

func DecodeJSONBody(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return err
	}

	return nil
}

func CheckRequiredFieldsInBody(body io.ReadCloser, requiredFields []string) error {
	if body == nil {
		return errors.New("body is empty")
	}

	var bodyMap map[string]interface{}
	if err := json.NewDecoder(body).Decode(&bodyMap); err != nil {
		return errors.New("invalid body")
	}

	for _, field := range requiredFields {
		if _, ok := bodyMap[field]; !ok {
			return errors.New("field " + field + " is required")
		}
	}

	return nil
}
