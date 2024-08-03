package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpStatus int

// ParseJson decodes the request body into the payload.
func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

// WriteJson writes the payload as a JSON response.
func WriteJson(w http.ResponseWriter, status HttpStatus, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	return json.NewEncoder(w).Encode(payload)
}

// WriteError writes a HTTP error as a JSON response.
func WriteError(w http.ResponseWriter, status HttpStatus, err error) {
	WriteJson(w, status, map[string]string{"error": err.Error()})
}
