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

func Init() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(consts.ConfigEnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for _, field := range Fields {
		if field.Default != nil {
			viper.SetDefault(field.Name, field.Default)
		}
	}

	cfgFile := viper.GetString(FieldFlagConfig.Name)
	if err := FieldFlagConfig.Validate(cfgFile); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if cfgFile != "" {
		cfgFileDir := path.Dir(cfgFile)
		cfgFileBase := path.Base(cfgFile)
		cfgFileExt := path.Ext(cfgFile)
		cfgFileName := strings.TrimSuffix(cfgFileBase, cfgFileExt)

		viper.SetConfigName(cfgFileName)
		viper.SetConfigType(cfgFileExt[1:])
		viper.AddConfigPath(cfgFileDir)

		if err := viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return fmt.Errorf("config read: %w", err)
		}
	}

	return nil
}

func ReadField(f *Field) any {
	return viper.Get(f.Name)
}

func ReadFieldString(f *Field) string {
	return viper.GetString(f.Name)
}

func ReadFieldBool(f *Field) bool {
	return viper.GetBool(f.Name)
}

func ReadFieldInt(f *Field) int {
	return viper.GetInt(f.Name)
}

func ReadFieldDuration(f *Field) time.Duration {
	return viper.GetDuration(f.Name)
}

func WriteField(f *Field, value any) error {
	if err := f.Validate(value); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	viper.Set(f.Name, value)

	return nil
}

func Save() error {
	cfgFile := viper.GetString(FieldFlagConfig.Name)
	if cfgFile == "" {
		return fmt.Errorf("config file not set")
	}

	if err := viper.WriteConfigAs(cfgFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// If the config file is not found, create it.
			if err := os.MkdirAll(path.Dir(cfgFile), 0755); err != nil {
				return fmt.Errorf("create config file directory: %w", err)
			}

			if f, err := os.Create(cfgFile); err != nil {
				return fmt.Errorf("create config file: %s: %w", cfgFile, err)
			} else {
				f.Close()
			}

			if err := viper.WriteConfigAs(cfgFile); err != nil {
				return fmt.Errorf("write config: %w", err)
			}

			return nil
		}

		return fmt.Errorf("write config: %s: %w", cfgFile, err)
	}

	return nil
}
