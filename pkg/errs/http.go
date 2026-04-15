package errs

import (
	"encoding/json"
	"errors"
	"net/http"
)

// WriteHTTPError inspects the error against sentinel values and writes the
// appropriate HTTP status code and JSON error body. Falls back to 500 for
// unknown errors.
func WriteHTTPError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var status int
	switch {
	case errors.Is(err, ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, ErrConflict):
		status = http.StatusConflict
	default:
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
