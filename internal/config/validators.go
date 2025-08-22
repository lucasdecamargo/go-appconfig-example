package config

import (
	"fmt"
	"path"
	"slices"
	"strings"
	"time"
)

var validConfigFileExts = []string{"yaml", "yml", "json", "toml", "hcl", "env"}

func validateConfigFile(v any) error {
	val, ok := v.(string)
	if !ok {
		return fmt.Errorf("config file must be a string")
	}

	if val == "" {
		return nil // empty value
	}

	cfgFileExt := strings.ToLower(path.Ext(val)[1:])

	if !slices.Contains(validConfigFileExts, cfgFileExt) {
		return fmt.Errorf("valid extensions: %v", validConfigFileExts)
	}

	return nil
}

func validateDuration(v any) error {
	switch v := v.(type) {
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return nil
	case float32, float64:
		return nil
	case string:
		if v == "" {
			return nil // empty value
		}
		_, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("duration must be a valid duration string")
		}
		return nil
	default:
		return fmt.Errorf("duration must be a string or a number")
	}
}
