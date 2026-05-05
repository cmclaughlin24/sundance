package domain

import (
	"time"

	"github.com/google/uuid"
)

// Package declaration for the current time function. Allows for easier testing by enabling the injection of a
// mock time function.
var Now = time.Now

// Creates a new random UUIDV7 and returns it as a string or panics.
func NewID() string {
	id, err := uuid.NewV7()

	if err != nil {
		panic(err)
	}

	return id.String()
}
