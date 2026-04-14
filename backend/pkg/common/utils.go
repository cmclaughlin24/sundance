package common

import (
	"encoding/json"
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
