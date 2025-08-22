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
	"github.com/spf13/pflag"
)

// FlagShowHidden controls whether hidden fields are displayed in output
var FlagShowHidden bool

// configCmd represents the config command group
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Example: `confapp config list
confapp config describe
confapp config set --log.level debug`,
}

// configListCmd lists configuration values
var configListCmd = &cobra.Command{
	Use:               "list [prefix] ...",
	Short:             "List configuration values",
	RunE:              listConfig,
	ValidArgsFunction: generateFieldCompletions,
	Example: `confapp config list
confapp config list log
confapp config list proxy
confapp config list --hidden`,
}

// configDescribeCmd describes configuration parameters in detail
var configDescribeCmd = &cobra.Command{
	Use:               "describe [prefix] ...",
	Short:             "Describe configuration parameters",
	RunE:              describeConfig,
	ValidArgsFunction: generateFieldCompletions,
	Example: `confapp config describe
confapp config describe log
confapp config describe update
confapp config describe --hidden`,
}

// configSetCmd sets configuration values
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  `For more information about the configuration values, use the "describe" command.`,
	Args:  cobra.NoArgs,
	RunE:  setConfig,
	Example: `confapp config set --log.level info
confapp config set --log.level debug --log.output /var/log/app.log
confapp config set --update.auto true --update.period 1h
confapp config set --proxy.http http://proxy:8080`,
	ValidArgsFunction: generateSetCompletions,
}

func init() {
	// Add commands to the command tree
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configDescribeCmd)
	configCmd.AddCommand(configSetCmd)

	// Add flags for showing hidden fields
	configDescribeCmd.Flags().BoolVarP(&FlagShowHidden, "hidden", "", false, "Show hidden fields")
	configListCmd.Flags().BoolVarP(&FlagShowHidden, "hidden", "", false, "Show hidden fields")

	// Set up flags for all configuration fields
	setupConfigFlags()
}

// setupConfigFlags creates CLI flags for all configuration fields.
// Each field gets a corresponding flag that can be used with the 'set' command.
func setupConfigFlags() {
	flags := configSetCmd.Flags()

	for _, field := range config.Fields {
		if field.Hidden {
			continue
		}

		switch field.Type {
		case config.FieldTypeString:
			setupStringFlag(flags, field)
		case config.FieldTypeBool:
			setupBoolFlag(flags, field)
		case config.FieldTypeInt:
			setupIntFlag(flags, field)
		case config.FieldTypeFloat:
			setupFloatFlag(flags, field)
		case config.FieldTypeDuration:
			setupDurationFlag(flags, field)
		default:
			log.Panicf("Unsupported field type: %s\n", field.Type)
		}
	}
}

// setupStringFlag creates a string flag for a configuration field
func setupStringFlag(flags *pflag.FlagSet, field *config.Field) {
	defaultVal := ""
	if field.Default != nil {
		defaultVal = field.Default.(string)
	}
	flags.StringP(field.Name, field.Shorthand, defaultVal, field.Description)

	// Add completion for valid values if specified
	if len(field.ValidValues) > 0 {
		configSetCmd.RegisterFlagCompletionFunc(field.Name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return generateValidValueCompletions(field.ValidValues, toComplete)
		})
	}
}

// setupBoolFlag creates a bool flag for a configuration field
func setupBoolFlag(flags *pflag.FlagSet, field *config.Field) {
	// Use string flag to allow unsetting the value
	defaultVal := ""
	if field.Default != nil {
		defaultVal = strconv.FormatBool(field.Default.(bool))
	}
	flags.StringP(field.Name, field.Shorthand, defaultVal, field.Description)
}

// setupIntFlag creates an int flag for a configuration field
func setupIntFlag(flags *pflag.FlagSet, field *config.Field) {
	defaultVal := 0
	if field.Default != nil {
		defaultVal = field.Default.(int)
	}
	flags.IntP(field.Name, field.Shorthand, defaultVal, field.Description)
}

// setupFloatFlag creates a float64 flag for a configuration field
func setupFloatFlag(flags *pflag.FlagSet, field *config.Field) {
	defaultVal := 0.0
	if field.Default != nil {
		defaultVal = field.Default.(float64)
	}
	flags.Float64P(field.Name, field.Shorthand, defaultVal, field.Description)
}

// setupDurationFlag creates a duration flag for a configuration field
func setupDurationFlag(flags *pflag.FlagSet, field *config.Field) {
	defaultVal := time.Duration(0)
	if field.Default != nil {
		defaultVal = field.Default.(time.Duration)
	}
	flags.DurationP(field.Name, field.Shorthand, defaultVal, field.Description)
}

// generateFieldCompletions provides shell completion for field names
func generateFieldCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var completions []string
	for key := range config.Fields.Map() {
		if len(toComplete) == 0 || (len(key) >= len(toComplete) && key[:len(toComplete)] == toComplete) {
			completions = append(completions, key)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// generateSetCompletions provides shell completion for the set command
func generateSetCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
}

// generateValidValueCompletions provides shell completion for field values
func generateValidValueCompletions(validValues []any, toComplete string) ([]string, cobra.ShellCompDirective) {
	completions := make([]string, 0, len(validValues))
	for _, value := range validValues {
		if strings.HasPrefix(value.(string), toComplete) {
			completions = append(completions, value.(string))
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// listConfig displays configuration values in a simple key=value format
func listConfig(cmd *cobra.Command, args []string) error {
	selectedFields := selectFieldsByPrefix(args)

	if len(selectedFields) == 0 {
		return fmt.Errorf("no configuration fields found")
	}

	// Calculate maximum field name length for alignment
	maxNameLen := calculateMaxFieldNameLength(selectedFields)

	// Display each field
	for _, field := range selectedFields {
		if field.Hidden && !FlagShowHidden {
			continue
		}
		fmt.Printf("%-*s = %v\n", maxNameLen+1, field.Name, config.ReadField(field))
	}

	return nil
}

// describeConfig displays detailed information about configuration fields
func describeConfig(cmd *cobra.Command, args []string) error {
	selectedFields := selectFieldsByPrefix(args)

	if len(selectedFields) == 0 {
		return fmt.Errorf("no configuration fields found")
	}

	w, cleanup := Pager()
	defer cleanup()

	// Display fields grouped by their group name
	for group, fields := range selectedFields.GroupIter() {
		fmt.Fprintf(w, "\n%s\n", group)

		for _, field := range fields {
			if field.Hidden && !FlagShowHidden {
				continue
			}
			writeFieldDescription(w, field)
		}
		fmt.Fprintf(w, "\n")
	}

	return nil
}

// selectFieldsByPrefix filters fields based on the provided prefixes
func selectFieldsByPrefix(prefixes []string) config.FieldCollection {
	if len(prefixes) == 0 {
		return config.Fields
	}

	selectedFields := config.FieldCollection{}
	for _, prefix := range prefixes {
		for _, field := range config.Fields {
			if strings.HasPrefix(field.Name, prefix) {
				selectedFields = append(selectedFields, field)
			}
		}
	}
	return selectedFields
}

// calculateMaxFieldNameLength determines the longest field name for alignment
func calculateMaxFieldNameLength(fields config.FieldCollection) int {
	maxLen := 0
	for _, field := range fields {
		if len(field.Name) > maxLen {
			maxLen = len(field.Name)
		}
	}
	return maxLen
}

// writeFieldDescription writes detailed information about a single field
func writeFieldDescription(w interface{ Write([]byte) (int, error) }, field *config.Field) {
	// Field name and deprecation status
	if field.Deprecated != "" {
		fmt.Fprintf(w, "\n  %s (deprecated: %s)\n", field.Name, field.Deprecated)
	} else {
		fmt.Fprintf(w, "\n  %s\n", field.Name)
	}

	// Basic field information
	fmt.Fprintf(w, "    %s\n", field.Description)
	fmt.Fprintf(w, "    Type: %s\n", field.Type)

	// Current value (if different from default)
	val := config.ReadField(field)
	if val != nil && val != field.Default {
		fmt.Fprintf(w, "    Value: %v\n", val)
	}

	// Default value
	if field.Default != nil {
		fmt.Fprintf(w, "    Default: %v\n", field.Default)
	}

	// Valid values
	if field.ValidValues != nil {
		fmt.Fprintf(w, "    Valid values: %v\n", field.ValidValues)
	}

	// Validation rules
	if field.ValidateTag != "" {
		fmt.Fprintf(w, "    Validation: %s\n", field.ValidateTag)
	}

	// Example value
	if field.Example != "" {
		fmt.Fprintf(w, "    Example: %s\n", field.Example)
	}

	// Detailed documentation
	if field.Docstring != "" {
		fmt.Fprintf(w, "    Doc:\n")
		lines := strings.SplitSeq(field.Docstring, "\n")
		for line := range lines {
			fmt.Fprintf(w, "      %s\n", line)
		}
	}
}

// setConfig sets configuration values based on provided flags
func setConfig(cmd *cobra.Command, args []string) error {
	// Collect and process all set flags
	fieldCount := 0
	for _, field := range config.Fields {
		if cmd.Flags().Changed(field.Name) {
			if err := processFieldFlag(cmd, field); err != nil {
				return err
			}
			fieldCount++
		}
	}

	if fieldCount == 0 {
		return cmd.Help()
	}

	// Save the configuration to file
	if err := config.Save(); err != nil {
		fmt.Printf("Failed to save configuration: %v\n", err)
		os.Exit(1)
	}

	return nil
}

// processFieldFlag processes a single field flag and sets its value
func processFieldFlag(cmd *cobra.Command, field *config.Field) error {
	flag := cmd.Flags().Lookup(field.Name)
	if flag == nil {
		return fmt.Errorf("flag not found: %s", field.Name)
	}

	if FlagVerbose {
		fmt.Printf("# Setting: %s: %s\n", field.Name, flag.Value.String())
	}

	if err := config.WriteField(field, flag.Value.String()); err != nil {
		return fmt.Errorf("%s: %w", field.Name, err)
	}

	return nil
}
