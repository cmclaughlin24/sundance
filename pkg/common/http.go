package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ApiResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type ApiErrorResponse struct {
	Message    string `json:"message"`
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

// Reads the JSON payload from the `http.Request` and decodes it into the provided data structure of type `T`.
func ReadJsonPayload[T any](r *http.Request, data T) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

// Encodes the provided data as JSON and writes it to the `http.ResponseWriter` with the specified status code and
// optional headers.
func SendJsonResponse(w http.ResponseWriter, statusCode int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(out)

	return nil
}

// Sends an error response with the appropriate HTTP status code based on the type of error provided.
func SendErrorResponse(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, ErrNotFound):
		SendJsonResponse(w, http.StatusNotFound, ApiErrorResponse{
			Message:    "Not Found",
			Error:      err.Error(),
			StatusCode: http.StatusNotFound,
		})
	case errors.Is(err, ErrExists):
		SendJsonResponse(w, http.StatusConflict, ApiErrorResponse{
			Message:    "Conflict",
			Error:      err.Error(),
			StatusCode: http.StatusConflict,
		})
	default:
		SendJsonResponse(w, http.StatusInternalServerError, ApiErrorResponse{
			Message:    "Internal Server Error",
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}
}
