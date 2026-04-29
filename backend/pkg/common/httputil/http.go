package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
)

var (
	ErrDecodeJSON = errors.New("failed to parse json request")
)

type APIResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type APIErrorResponse struct {
	Message    string `json:"message"`
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

// Reads the JSON payload from the http.Request and decodes it into the provided data structure of type T.
func ReadJSONPayload[T any](r *http.Request, data T) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("%w: %w", ErrDecodeJSON, err)
	}

	return nil
}

// Reads the JSON payload from the http.Request, decodes it into the provided data structure of type T, and
// validates the decoded data.
func ReadValidateJSONPayload[T any](r *http.Request, data T) error {
	if err := ReadJSONPayload(r, data); err != nil {
		return err
	}

	return validate.ValidateStruct(data)
}

// Encodes the provided data as JSON and writes it to the http.ResponseWriter with the specified status code and
// optional headers.
func SendJSONResponse(w http.ResponseWriter, statusCode int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	for _, header := range headers {
		for key, value := range header {
			w.Header()[key] = value
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(out)

	return err
}

// Sends an error response with the appropriate HTTP status code based on the type of error provided.
func SendErrorResponse(w http.ResponseWriter, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrDecodeJSON) || errors.Is(err, common.ErrInvalidID) || validate.IsValidationErr(err):
		return SendJSONResponse(w, http.StatusBadRequest, APIErrorResponse{
			Message:    "Bad Request",
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	case errors.Is(err, common.ErrUnauthorized):
		return SendJSONResponse(w, http.StatusUnauthorized, APIErrorResponse{
			Message:    "Unauthorized",
			Error:      err.Error(),
			StatusCode: http.StatusUnauthorized,
		})
	case errors.Is(err, common.ErrNotFound):
		return SendJSONResponse(w, http.StatusNotFound, APIErrorResponse{
			Message:    "Not Found",
			Error:      err.Error(),
			StatusCode: http.StatusNotFound,
		})
	case errors.Is(err, common.ErrExists):
		return SendJSONResponse(w, http.StatusConflict, APIErrorResponse{
			Message:    "Conflict",
			Error:      err.Error(),
			StatusCode: http.StatusConflict,
		})
	default:
		return SendJSONResponse(w, http.StatusInternalServerError, APIErrorResponse{
			Message:    "Internal Server Error",
			Error:      "An unexpected error occurred. Please contact support if the issue persists.",
			StatusCode: http.StatusInternalServerError,
		})
	}
}

// Reads the JSON response from the http.Response and decodes it into the provided data structure of type T or
// returns an error if the response status code indicates a failure or if the JSON decoding fails.
func DecodeJSONResponse[T any](resp *http.Response, data T) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
