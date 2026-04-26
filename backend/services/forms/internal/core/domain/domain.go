package domain

import "time"

// Package declaration for the current time function. Allows for easier testing by enabling the injection of a
// mock time function.
var Now = time.Now
