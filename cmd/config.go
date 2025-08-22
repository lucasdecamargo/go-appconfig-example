package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lucasdecamargo/go-appconfig-example/internal/config"
	"github.com/spf13/cobra"
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
	Use:   "list [prefix] ...",
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
	Use:   "describe [prefix] ...",
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

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  `For more information about the configuration values, use the "describe" command.`,
	Args:  cobra.NoArgs,
	RunE:  setConfig,
	Example: `config set --log.level info
config set --log.level info --log.output /var/log/app.log --update.auto true
`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			completions := make([]string, 0, len(config.Fields))
			for _, field := range config.Fields {
				if field.Hidden {
					continue
				}
				completions = append(completions, "--"+field.Name)
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		}

		return nil, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configDescribeCmd)
	configCmd.AddCommand(configSetCmd)

	configDescribeCmd.Flags().BoolVarP(&FlagShowHidden, "hidden", "", false, "Show hidden fields")
	configListCmd.Flags().BoolVarP(&FlagShowHidden, "hidden", "", false, "Show hidden fields")

	flags := configSetCmd.Flags()

	for _, field := range config.Fields {
		if field.Hidden {
			continue
		}

		switch field.Type {
		case config.FieldTypeString:
			default_ := ""
			if field.Default != nil {
				default_ = field.Default.(string)
			}
			flags.StringP(field.Name, field.Shorthand, default_, field.Description)
			// Add completion function for valid values
			if len(field.ValidValues) > 0 {
				configSetCmd.RegisterFlagCompletionFunc(field.Name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
					completions := make([]string, 0, len(field.ValidValues))
					for _, value := range field.ValidValues {
						if strings.HasPrefix(value.(string), toComplete) {
							completions = append(completions, value.(string))
						}
					}
					return completions, cobra.ShellCompDirectiveNoFileComp
				})
			}
		case config.FieldTypeBool:
			// Use string flag to allow unsetting the value.
			default_ := ""
			if field.Default != nil {
				default_ = strconv.FormatBool(field.Default.(bool))
			}
			flags.StringP(field.Name, field.Shorthand, default_, field.Description)
		case config.FieldTypeInt:
			default_ := 0
			if field.Default != nil {
				default_ = field.Default.(int)
			}
			flags.IntP(field.Name, field.Shorthand, default_, field.Description)
		case config.FieldTypeFloat:
			default_ := 0.0
			if field.Default != nil {
				default_ = field.Default.(float64)
			}
			flags.Float64P(field.Name, field.Shorthand, default_, field.Description)
		case config.FieldTypeDuration:
			default_ := time.Duration(0)
			if field.Default != nil {
				default_ = field.Default.(time.Duration)
			}
			flags.DurationP(field.Name, field.Shorthand, default_, field.Description)
		default:
			log.Panicf("Unsupported field type: %s\n", field.Type)
		}
	}
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

	for group, fields := range selectedFields.GroupIter() {
		fmt.Fprintf(w, "\n%s\n", group)

		for _, field := range fields {
			if field.Hidden && !FlagShowHidden {
				continue
			}

			if field.Deprecated != "" {
				fmt.Fprintf(w, "\n  %s (deprecated: %s)\n", field.Name, field.Deprecated)
			} else {
				fmt.Fprintf(w, "\n  %s\n", field.Name)
			}

			fmt.Fprintf(w, "    %s\n", field.Description)
			fmt.Fprintf(w, "    Type: %s\n", field.Type)

			val := config.ReadField(field)

			if val != nil && val != field.Default {
				fmt.Fprintf(w, "    Value: %v\n", val)
			}

			if field.Default != nil {
				fmt.Fprintf(w, "    Default: %v\n", field.Default)
			}

			if field.ValidValues != nil {
				fmt.Fprintf(w, "    Valid values: %v\n", field.ValidValues)
			}

			if field.ValidateTag != "" {
				fmt.Fprintf(w, "    Validation: %s\n", field.ValidateTag)
			}

			if field.Example != "" {
				fmt.Fprintf(w, "    Example: %s\n", field.Example)
			}

			if field.Docstring != "" {
				fmt.Fprintf(w, "    Doc:\n")
				// Print field.Docstring with 2-space indentation, handling newlines
				lines := strings.SplitSeq(field.Docstring, "\n")
				for line := range lines {
					fmt.Fprintf(w, "      %s\n", line)
				}
			}
		}

		fmt.Fprintf(w, "\n")
	}

	return nil
}

func setConfig(cmd *cobra.Command, args []string) error {
	// Collect only the Flags that have been set.
	fieldCount := 0
	for _, field := range config.Fields {
		if cmd.Flags().Changed(field.Name) {
			flag := cmd.Flags().Lookup(field.Name)
			if flag == nil {
				continue
			}

			if FlagVerbose {
				fmt.Printf("# Setting: %s: %s\n", field.Name, flag.Value.String())
			}

			fieldCount++

			if err := config.WriteField(field, flag.Value.String()); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}
		}
	}

	if fieldCount == 0 {
		return cmd.Help()
	}

	if err := config.Save(); err != nil {
		fmt.Printf("save: %v\n", err)
		os.Exit(1)
	}

	return nil
}
