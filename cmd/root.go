package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/lucasdecamargo/go-appconfig-example/internal/config"
	"github.com/lucasdecamargo/go-appconfig-example/internal/consts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	FlagConfig  string
	FlagVerbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   consts.AppName,
	Short: "Go application with structured configuration example",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// initConfig is called before any command is executed.
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// UserConfigDir returns the default root directory to use for user-specific configuration data.
	// Users should create their own application-specific subdirectory within this one and use that.
	configDir, _ := os.UserConfigDir()
	defaultConfig := path.Join(configDir, consts.AppName, "config.yaml")

	// If the default config file is in the user's home directory, use "~" instead.
	homeDir, _ := os.UserHomeDir()
	if len(homeDir) > 0 && len(defaultConfig) > len(homeDir) && defaultConfig[:len(homeDir)] == homeDir {
		defaultConfig = "~" + defaultConfig[len(homeDir):]
	}

	rootCmd.PersistentFlags().StringVarP(
		&FlagConfig,
		config.FieldFlagConfig.Name,
		config.FieldFlagConfig.Shorthand,
		defaultConfig,
		config.FieldFlagConfig.Description,
	)
	viper.BindPFlag(
		config.FieldFlagConfig.Name,
		rootCmd.PersistentFlags().Lookup(config.FieldFlagConfig.Name),
	)

	rootCmd.PersistentFlags().BoolVarP(
		&FlagVerbose,
		config.FieldFlagVerbose.Name,
		config.FieldFlagVerbose.Shorthand,
		false,
		config.FieldFlagVerbose.Description,
	)
	viper.BindPFlag(
		config.FieldFlagVerbose.Name,
		rootCmd.PersistentFlags().Lookup(config.FieldFlagVerbose.Name),
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if err := config.Init(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if FlagVerbose {
		cfgFile := config.ReadFieldString(config.FieldFlagConfig)
		if cfgFile != "" {
			fmt.Printf("# Using config file: %s\n", cfgFile)
		} else {
			fmt.Printf("# No config file found, using default values\n")
		}
	}
}
