package config

import (
	"log"
	"strconv"
	"time"
)

// defaultString converts a string to a default value, returning nil if empty.
// This allows build-time definition of defaults using -ldflags.
func defaultString(val string) any {
	if val == "" {
		return nil
	}
	return val
}

// defaultBool converts a string to a boolean default value.
// Valid values are "true", "false", or empty string (returns nil).
// Panics on invalid values to catch configuration errors early.
func defaultBool(val string) any {
	switch val {
	case "":
		return nil
	case "true":
		return true
	case "false":
		return false
	default:
		log.Fatalf("Invalid default bool value: %s (must be 'true', 'false', or empty)", val)
		return false // unreachable, but satisfies compiler
	}
}

// defaultInt converts a string to an integer default value.
// Returns nil if the string is empty, panics on invalid integers.
func defaultInt(val string) any {
	if val == "" {
		return nil
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Invalid default int value: %s (%v)", val, err)
	}

	return i
}

// defaultDuration converts a string to a duration default value.
// Returns nil if the string is empty, panics on invalid duration strings.
// Supports Go duration format (e.g., "1h30m", "15m", "10s").
func defaultDuration(val string) any {
	if val == "" {
		return nil
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("Invalid default duration value: %s (%v)", val, err)
	}

	return d
}
