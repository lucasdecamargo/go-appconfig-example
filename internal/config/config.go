package config

import (
	"errors"
	"fmt"
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
		return err
	}

	viper.Set(f.Name, value)

	return nil
}

func Save() error {
	return viper.WriteConfig()
}
