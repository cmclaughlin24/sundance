package common

import (
	"encoding/json"
	"net/http"
	"os"
)

func ReadJsonFile[T any](path string, data *T) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err = decoder.Decode(data); err != nil {
		return err
	}

	return nil
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
