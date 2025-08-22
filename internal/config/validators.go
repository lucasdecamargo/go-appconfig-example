package config

import (
	"fmt"
	"path"
	"slices"
	"strings"
	"time"
)

// validConfigFileExts defines the supported configuration file extensions
var validConfigFileExts = []string{"yaml", "yml", "json", "toml", "hcl", "env"}

// validateConfigFile validates that a config file path has a supported extension.
// Returns an error if the file extension is not supported, or nil if valid.
func validateConfigFile(v any) error {
	val, ok := v.(string)
	if !ok {
		return fmt.Errorf("config file path must be a string")
	}

	if val == "" {
		return nil // empty value is allowed
	}

	// Extract file extension and validate
	cfgFileExt := strings.ToLower(path.Ext(val))
	if cfgFileExt == "" {
		return fmt.Errorf("config file must have an extension")
	}

	// Remove the leading dot from extension
	cfgFileExt = cfgFileExt[1:]

	if !slices.Contains(validConfigFileExts, cfgFileExt) {
		return fmt.Errorf("unsupported file extension: %s (supported: %v)", cfgFileExt, validConfigFileExts)
	}

	return nil
}

// validateDuration validates that a value can be interpreted as a duration.
// Accepts numeric values (interpreted as seconds) or duration strings.
// Returns an error if the value cannot be parsed as a duration.
func validateDuration(v any) error {
	switch v := v.(type) {
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		// Numeric values are valid (interpreted as seconds)
		return nil
	case float32, float64:
		// Float values are valid (interpreted as seconds)
		return nil
	case string:
		if v == "" {
			return nil // empty value is allowed
		}
		_, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("invalid duration format: %s (examples: 1h30m, 15m, 10s)", v)
		}
		return nil
	default:
		return fmt.Errorf("duration must be a string or numeric value")
	}
}
