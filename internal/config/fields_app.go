package config

const GroupApplication = "Application"

var (
	// Defaults are set as string variables to allow build-time definition with -ldflags
	DefaultAppEnvironment    = "dev"
	DefaultAppLogLevel       = "info"
	DefaultAppLogOutput      = ""
	DefaultAppLogFormat      = "text"
	DefaultAppUpdateUnstable = "false"
	DefaultAppUpdateAuto     = "false"
	DefaultAppUpdatePeriod   = "15m"
)

func init() {
	Fields.Add(
		FieldAppEnvironment,
		FieldAppLogLevel,
		FieldAppLogOutput,
		FieldAppLogFormat,
		FieldAppUpdateUnstable,
		FieldAppUpdateAuto,
		FieldAppUpdatePeriod,
	)
}

var FieldAppEnvironment = &Field{
	Name:        "environment",
	Group:       GroupApplication,
	Type:        FieldTypeString,
	Default:     defaultString(DefaultAppEnvironment),
	Description: "The environment in which the application runs.",
	Example:     "prod, test, dev",
	Hidden:      true, // This field is not shown in the config list
}

// region Logging

var FieldAppLogLevel = &Field{
	Name:        "log.level",
	Group:       GroupApplication,
	Type:        FieldTypeString,
	Default:     defaultString(DefaultAppLogLevel),
	Description: "The log level to use for the application.",
	ValidValues: []any{"debug", "info", "warn", "error"},
}

var FieldAppLogOutput = &Field{
	Name:        "log.output",
	Group:       GroupApplication,
	Type:        FieldTypeString,
	Default:     defaultString(DefaultAppLogOutput),
	Description: "The output file to use for the application logs, if set.",
	ValidateTag: "dirpath",
	Example:     "/var/log/app.log",
}

var FieldAppLogFormat = &Field{
	Name:        "log.format",
	Group:       GroupApplication,
	Type:        FieldTypeString,
	Default:     defaultString(DefaultAppLogFormat),
	Description: "The format to use for the application log file, if set.",
	ValidValues: []any{"json", "text"},
}

// region Updates

var FieldAppUpdateUnstable = &Field{
	Name:        "update.unstable",
	Group:       GroupApplication,
	Type:        FieldTypeBool,
	Default:     defaultBool(DefaultAppUpdateUnstable),
	Description: "Receive updates for unstable versions.",
}

var FieldAppUpdateAuto = &Field{
	Name:        "update.auto",
	Group:       GroupApplication,
	Type:        FieldTypeBool,
	Default:     defaultBool(DefaultAppUpdateAuto),
	Description: "Automatically update the application when a new version is available.",
}

var FieldAppUpdatePeriod = &Field{
	Name:         "update.period",
	Group:        GroupApplication,
	Type:         FieldTypeDuration,
	Default:      defaultDuration(DefaultAppUpdatePeriod),
	Description:  "The period to check for updates, if enabled.",
	Docstring:    `The period can be a number of seconds, or a valid duration string.`,
	ValidateFunc: validateDuration,
	Example:      "1h, 15m, 10 (seconds)",
}
