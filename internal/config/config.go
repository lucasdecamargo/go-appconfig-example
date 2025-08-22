package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/lucasdecamargo/go-appconfig-example/internal/consts"
	"github.com/spf13/viper"
)

// Init initializes the configuration system by setting up Viper with environment
// variables, defaults, and optionally reading from a configuration file.
// This function must be called before any configuration values are accessed.
func Init() error {
	// Enable automatic environment variable binding
	viper.AutomaticEnv()
	viper.SetEnvPrefix(consts.ConfigEnvPrefix)
	// Replace dots with underscores in environment variable names
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values for all defined fields
	for _, field := range Fields {
		if field.Default != nil {
			viper.SetDefault(field.Name, field.Default)
		}
	}

	// Validate and process config file if specified
	cfgFile := viper.GetString(FieldFlagConfig.Name)
	if err := FieldFlagConfig.Validate(cfgFile); err != nil {
		return fmt.Errorf("config file validation failed: %w", err)
	}

	if cfgFile != "" {
		if err := loadConfigFile(cfgFile); err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
	}

	return nil
}

// loadConfigFile loads configuration from the specified file path.
// The file extension determines the format (yaml, json, toml, etc.).
func loadConfigFile(cfgFile string) error {
	cfgFileDir := path.Dir(cfgFile)
	cfgFileBase := path.Base(cfgFile)
	cfgFileExt := path.Ext(cfgFile)
	cfgFileName := strings.TrimSuffix(cfgFileBase, cfgFileExt)

	viper.SetConfigName(cfgFileName)
	viper.SetConfigType(cfgFileExt[1:])
	viper.AddConfigPath(cfgFileDir)

	// Attempt to read the config file, but don't fail if it doesn't exist
	if err := viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
}

// ReadField retrieves the current value of a configuration field.
// Returns the value as an interface{} type.
func ReadField(f *Field) any {
	return viper.Get(f.Name)
}

// ReadFieldString retrieves the current value of a configuration field as a string.
func ReadFieldString(f *Field) string {
	return viper.GetString(f.Name)
}

// ReadFieldBool retrieves the current value of a configuration field as a boolean.
func ReadFieldBool(f *Field) bool {
	return viper.GetBool(f.Name)
}

// ReadFieldInt retrieves the current value of a configuration field as an integer.
func ReadFieldInt(f *Field) int {
	return viper.GetInt(f.Name)
}

// ReadFieldDuration retrieves the current value of a configuration field as a duration.
func ReadFieldDuration(f *Field) time.Duration {
	return viper.GetDuration(f.Name)
}

// WriteField sets a configuration field value after validating it.
// The value is validated using the field's validation rules before being set.
func WriteField(f *Field, value any) error {
	if err := f.Validate(value); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	viper.Set(f.Name, value)
	return nil
}

// Save writes the current configuration to the specified config file.
// Creates the directory structure if it doesn't exist.
func Save() error {
	cfgFile := viper.GetString(FieldFlagConfig.Name)
	if cfgFile == "" {
		return fmt.Errorf("no config file specified")
	}

	// Try to write the config file
	if err := viper.WriteConfigAs(cfgFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Create directory structure and file if they don't exist
			return createAndWriteConfigFile(cfgFile)
		}
		return fmt.Errorf("failed to write config file %s: %w", cfgFile, err)
	}

	return nil
}

// createAndWriteConfigFile creates the directory structure and config file,
// then writes the current configuration to it.
func createAndWriteConfigFile(cfgFile string) error {
	// Create the directory structure
	if err := os.MkdirAll(path.Dir(cfgFile), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create an empty config file
	if f, err := os.Create(cfgFile); err != nil {
		return fmt.Errorf("failed to create config file %s: %w", cfgFile, err)
	} else {
		f.Close()
	}

	// Write the configuration to the file
	if err := viper.WriteConfigAs(cfgFile); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	return nil
}
