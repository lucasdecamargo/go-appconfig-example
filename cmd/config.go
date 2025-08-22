package cmd

import (
	"fmt"
	"strings"

	"github.com/lucasdecamargo/go-appconfig-example/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	FlagShowHidden bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
}

var configListCmd = &cobra.Command{
	Use:   "list [prefix]",
	Short: "List configuration values",
	RunE:  listConfig,
	// args: the arguments that the user has already typed before the one currently being completed.
	// toComplete: the current word fragment the user is typing and wants completion for.
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		for key := range config.Fields.Map() {
			if len(toComplete) == 0 || len(key) >= len(toComplete) && key[:len(toComplete)] == toComplete {
				completions = append(completions, key)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
}

var configDescribeCmd = &cobra.Command{
	Use:   "describe [prefix]",
	Short: "Describe configuration parameters",
	RunE:  describeConfig,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		for key := range config.Fields.Map() {
			if len(toComplete) == 0 || len(key) >= len(toComplete) && key[:len(toComplete)] == toComplete {
				completions = append(completions, key)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configDescribeCmd)

	configDescribeCmd.Flags().BoolVarP(&FlagShowHidden, "hidden", "", false, "Show hidden fields")
	viper.BindPFlag("flag.show-hidden", configDescribeCmd.Flags().Lookup("hidden"))
	configListCmd.Flags().BoolVarP(&FlagShowHidden, "hidden", "", false, "Show hidden fields")
	viper.BindPFlag("flag.show-hidden", configListCmd.Flags().Lookup("hidden"))
}

func listConfig(cmd *cobra.Command, args []string) error {
	selectedFields := config.FieldCollection{}

	if len(args) > 0 {
		for _, prefix := range args {
			for _, field := range config.Fields {
				if strings.HasPrefix(field.Name, prefix) {
					selectedFields = append(selectedFields, field)
				}
			}
		}
	} else {
		selectedFields = config.Fields
	}

	maxNameLen := 0
	for _, field := range selectedFields {
		if len(field.Name) > maxNameLen {
			maxNameLen = len(field.Name)
		}
	}

	for _, field := range selectedFields {
		if field.Hidden && !FlagShowHidden {
			continue
		}

		fmt.Printf("%-*s = %v\n", maxNameLen+1, field.Name, config.ReadField(field))
	}

	return nil
}

func describeConfig(cmd *cobra.Command, args []string) error {
	selectedFields := config.FieldCollection{}

	if len(args) > 0 {
		for _, prefix := range args {
			for _, field := range config.Fields {
				if strings.HasPrefix(field.Name, prefix) {
					selectedFields = append(selectedFields, field)
				}
			}
		}
	} else {
		selectedFields = config.Fields
	}

	w, cleanup := Pager()
	defer cleanup()

	for _, field := range selectedFields {
		if field.Hidden && !FlagShowHidden {
			continue
		}

		if field.Deprecated != "" {
			fmt.Fprintf(w, "\n%s (deprecated: %s)\n", field.Name, field.Deprecated)
		} else {
			fmt.Fprintf(w, "\n%s\n", field.Name)
		}

		fmt.Fprintf(w, "  %s\n", field.Description)
		fmt.Fprintf(w, "  Type: %s\n", field.Type)

		val := config.ReadField(field)

		if val != nil && val != field.Default {
			fmt.Fprintf(w, "  Value: %v\n", val)
		}

		if field.Default != nil {
			fmt.Fprintf(w, "  Default: %v\n", field.Default)
		}

		if field.ValidValues != nil {
			fmt.Fprintf(w, "  Valid values: %v\n", field.ValidValues)
		}

		if FlagVerbose {
			if field.ValidateTag != "" {
				fmt.Fprintf(w, "  Validation: %s\n", field.ValidateTag)
			}

			if field.Example != "" {
				fmt.Fprintf(w, "  Example: %s\n", field.Example)
			}

			if field.Docstring != "" {
				fmt.Fprintf(w, "  Doc:\n")
				// Print field.Docstring with 2-space indentation, handling newlines
				lines := strings.SplitSeq(field.Docstring, "\n")
				for line := range lines {
					fmt.Fprintf(w, "    %s\n", line)
				}
			}
		}
	}

	return nil
}
