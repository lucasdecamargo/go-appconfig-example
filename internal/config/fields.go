package config

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FieldType represents the data type of a configuration field
type FieldType string

const (
	FieldTypeString   FieldType = "string"
	FieldTypeBool     FieldType = "bool"
	FieldTypeInt      FieldType = "int"
	FieldTypeFloat    FieldType = "float"
	FieldTypeDuration FieldType = "duration"
)

// Field defines a single configuration field with all metadata.
// This struct serves as the single source of truth for configuration parameters,
// containing everything needed to define, validate, and document a config field.
type Field struct {
	Name         string          // The configuration key name (e.g., "log.level")
	Group        string          // Logical grouping for organization (e.g., "Application")
	Type         FieldType       // Data type of the field
	Default      any             // Default value if not specified
	Description  string          // Short description for CLI help
	Docstring    string          // Detailed documentation for pager output
	Hidden       bool            // Whether to hide from normal listing
	Shorthand    string          // Short flag name (e.g., "v" for verbose)
	ValidValues  []any           // Allowed values for validation
	ValidateTag  string          // Go validator tag for validation
	ValidateFunc func(any) error // Custom validation function
	Example      string          // Example value for documentation
	Deprecated   string          // Deprecation message if field is deprecated
}

// Validate performs validation on a field value using the configured validation rules.
// Returns nil if validation passes, or an error describing the validation failure.
func (f *Field) Validate(value any) error {
	if value == nil {
		return nil
	}

	// Check against valid values if specified
	if f.ValidValues != nil {
		if slices.Contains(f.ValidValues, value) {
			return nil
		}
		return fmt.Errorf("%s: valid values: %v", f.Name, f.ValidValues)
	}

	// Apply validator tag if specified
	if f.ValidateTag != "" {
		if err := validator.New().Var(value, f.ValidateTag); err != nil {
			return fmt.Errorf("%s: %w", f.Name, err)
		}
	}

	// Apply custom validation function if specified
	if f.ValidateFunc != nil {
		if err := f.ValidateFunc(value); err != nil {
			return fmt.Errorf("%s: %w", f.Name, err)
		}
	}

	return nil
}

// FieldCollection represents a collection of configuration fields
type FieldCollection []*Field

// FieldMap provides O(1) lookup by field name
type FieldMap map[string]*Field

// FieldGroup organizes fields by their group name
type FieldGroup map[string]FieldCollection

// Fields is the global collection of all configuration fields
var Fields = FieldCollection{}

// Add adds fields to the collection, replacing existing fields with the same name
func (fc *FieldCollection) Add(fields ...*Field) {
	for _, field := range fields {
		// Check if field already exists and replace it
		for i, existingField := range *fc {
			if existingField.Name == field.Name {
				(*fc)[i] = field
				goto nextField
			}
		}
		// Field doesn't exist, append it
		*fc = append(*fc, field)
	nextField:
	}
}

// Group organizes fields by their group name
func (fc *FieldCollection) Group() FieldGroup {
	groups := make(map[string]FieldCollection)
	for _, f := range *fc {
		groups[f.Group] = append(groups[f.Group], f)
	}
	return groups
}

// GroupIter returns an iterator that yields groups and their fields in sorted order.
// This is useful for consistent output ordering in CLI commands.
func (fc *FieldCollection) GroupIter() iter.Seq2[string, FieldCollection] {
	groups := fc.Group()
	groupNames := make([]string, 0, len(groups))
	for group := range groups {
		groupNames = append(groupNames, group)
	}
	slices.Sort(groupNames)

	return func(yield func(string, FieldCollection) bool) {
		for _, group := range groupNames {
			fields := groups[group]
			// Sort fields within each group for consistent output
			slices.SortFunc(fields, func(a, b *Field) int {
				return strings.Compare(a.Name, b.Name)
			})
			if !yield(group, fields) {
				break
			}
		}
	}
}

// Map returns a map for O(1) field lookup by name
func (fc *FieldCollection) Map() FieldMap {
	fields := make(FieldMap)
	for _, f := range *fc {
		fields[f.Name] = f
	}
	return fields
}

// Global flag fields that are available to all commands

// FieldFlagConfig defines the config file flag
var FieldFlagConfig = &Field{
	Name:         "config",
	Type:         "string",
	Shorthand:    "c",
	Description:  "Config file path",
	Docstring:    `The configuration file should of one of the extensions: yaml, yml, json, toml, hcl, env`,
	ValidateTag:  "filepath",
	ValidateFunc: validateConfigFile,
}

// FieldFlagVerbose defines the verbose flag
var FieldFlagVerbose = &Field{
	Name:        "verbose",
	Type:        "bool",
	Shorthand:   "v",
	Description: "Display more verbose output in console output.",
	Docstring:   "",
}
