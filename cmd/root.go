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

// Global flag variables that are bound to the root command
var (
	FlagConfig  string // Path to the configuration file
	FlagVerbose bool   // Enable verbose output
)

// rootCmd represents the base command when called without any subcommands.
// It serves as the entry point for the CLI application.
var rootCmd = &cobra.Command{
	Use:     consts.AppName,
	Short:   "Go application with structured configuration example",
	Long:    `A demonstration of production-ready configuration management in Go using Viper and Cobra.`,
	Version: consts.AppVersion,
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
	// Initialize configuration before any command is executed
	cobra.OnInitialize(initConfig)

	// Set up persistent flags that are available to all commands
	setupPersistentFlags()
}

// setupPersistentFlags configures the global flags that are available to all commands.
// These flags are bound to Viper for automatic configuration integration.
func setupPersistentFlags() {
	// Determine default config file location
	configDir, _ := os.UserConfigDir()
	defaultConfig := path.Join(configDir, consts.AppName, "config.yaml")

	// Config file flag
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

	// Verbose flag
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
// This function is called by Cobra before any command execution.
func initConfig() {
	if err := config.Init(); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Display configuration information in verbose mode
	if FlagVerbose {
		printConfigInfo()
	}
}

// printConfigInfo displays information about the current configuration
// when verbose mode is enabled.
func printConfigInfo() {
	cfgFile := config.ReadFieldString(config.FieldFlagConfig)
	if cfgFile != "" {
		fmt.Printf("# Using config file: %s\n", cfgFile)
	} else {
		fmt.Printf("# No config file found, using default values\n")
	}
}
