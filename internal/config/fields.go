package config

import (
	"fmt"
	"slices"

	"github.com/go-playground/validator/v10"
)

type FieldType string

const (
	FieldTypeString   FieldType = "string"
	FieldTypeBool     FieldType = "bool"
	FieldTypeInt      FieldType = "int"
	FieldTypeFloat    FieldType = "float"
	FieldTypeDuration FieldType = "duration"
)

// Field defines a single configuration field with all metadata
type Field struct {
	Name         string
	Group        string
	Type         FieldType
	Default      any
	Description  string
	Docstring    string
	Hidden       bool
	Shorthand    string
	ValidValues  []any
	ValidateTag  string
	ValidateFunc func(any) error
	Example      string
	Deprecated   string
}

func (f *Field) Validate(value any) error {
	if value == nil {
		return nil
	}

	if f.ValidValues != nil {
		if slices.Contains(f.ValidValues, value) {
			return nil
		}

		return fmt.Errorf("%s: valid values: %v", f.Name, f.ValidValues)
	}

	if f.ValidateTag != "" {
		if err := validator.New().Var(value, f.ValidateTag); err != nil {
			return fmt.Errorf("%s: %w", f.Name, err)
		}
	}

	if f.ValidateFunc != nil {
		if err := f.ValidateFunc(value); err != nil {
			return fmt.Errorf("%s: %w", f.Name, err)
		}
	}

	return nil
}

type FieldCollection []*Field

type FieldMap map[string]*Field
type FieldGroup map[string]FieldCollection

var Fields = FieldCollection{}

func (fc *FieldCollection) Add(f ...*Field) {
	for _, field := range f {
		for _, existingField := range *fc {
			if existingField.Name == field.Name {
				*existingField = *field
				break
			}
		}
		*fc = append(*fc, field)
	}
}

func (fc *FieldCollection) Group() FieldGroup {
	groups := make(map[string]FieldCollection)
	for _, f := range *fc {
		groups[f.Group] = append(groups[f.Group], f)
	}
	return groups
}

func (fc *FieldCollection) Map() FieldMap {
	fields := make(FieldMap)
	for _, f := range *fc {
		fields[f.Name] = f
	}
	return fields
}

// region Global Flag Fields

var FieldFlagConfig = &Field{
	Name:         "config",
	Type:         "string",
	Shorthand:    "c",
	Description:  "Config file path",
	Docstring:    `The configuration file should of one of the extensions: yaml, yml, json, toml, hcl, env`,
	ValidateTag:  "filepath",
	ValidateFunc: validateConfigFile,
}

var FieldFlagVerbose = &Field{
	Name:        "verbose",
	Type:        "bool",
	Shorthand:   "v",
	Description: "Display more verbose output in console output.",
	Docstring:   "",
}
