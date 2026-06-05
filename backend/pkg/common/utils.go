package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

func ReadSettings[T any](path string, data *T) error {
	ext := filepath.Ext(path)

	switch ext {
	case ".yml", ".yaml":
		return ReadYAMLFile(path, data)
	case ".json":
		return ReadJSONFile(path, data)
	default:
		return fmt.Errorf("unknown file extension %s", ext)
	}
}

func ReadJSONFile[T any](path string, data *T) error {
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

func ReadYAMLFile[T any](path string, data *T) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(data); err != nil {
		return err
	}

	return nil
}
