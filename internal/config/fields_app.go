package config

// GroupApplication is the logical group name for application-related configuration fields
const GroupApplication = "Application"

// Default values for application configuration fields.
// These are defined as string variables to allow build-time definition with -ldflags.
// For example: go build -ldflags "-X github.com/lucasdecamargo/go-appconfig-example/internal/config.DefaultAppEnvironment=prod"
var (
	DefaultAppEnvironment    = "dev"   // Default application environment
	DefaultAppLogLevel       = "info"  // Default logging level
	DefaultAppLogOutput      = ""      // Default log output (empty = stdout)
	DefaultAppLogFormat      = "text"  // Default log format
	DefaultAppUpdateUnstable = "false" // Default unstable update setting
	DefaultAppUpdateAuto     = "false" // Default auto-update setting
	DefaultAppUpdatePeriod   = "15m"   // Default update check period
)

func init() {
	// Register all application configuration fields
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

// FieldAppEnvironment defines the application environment setting
var FieldAppEnvironment = &Field{
	Name:        "environment",
	Group:       GroupApplication,
	Type:        FieldTypeString,
	Default:     defaultString(DefaultAppEnvironment),
	Description: "The environment in which the application runs.",
	Example:     "prod, test, dev",
	Hidden:      true, // This field is not shown in the config list
}

// Logging configuration fields

// FieldAppLogLevel defines the application logging level
var FieldAppLogLevel = &Field{
	Name:        "log.level",
	Group:       GroupApplication,
	Type:        FieldTypeString,
	Default:     defaultString(DefaultAppLogLevel),
	Description: "The log level to use for the application.",
	ValidValues: []any{"debug", "info", "warn", "error"},
}

// FieldAppLogOutput defines the log output destination
var FieldAppLogOutput = &Field{
	Name:        "log.output",
	Group:       GroupApplication,
	Type:        FieldTypeString,
	Default:     defaultString(DefaultAppLogOutput),
	Description: "The output file to use for the application logs, if set.",
	ValidateTag: "filepath",
	Example:     "/var/log/app.log",
}

// FieldAppLogFormat defines the log output format
var FieldAppLogFormat = &Field{
	Name:        "log.format",
	Group:       GroupApplication,
	Type:        FieldTypeString,
	Default:     defaultString(DefaultAppLogFormat),
	Description: "The format to use for the application log file, if set.",
	ValidValues: []any{"json", "text"},
}

// Update configuration fields

// FieldAppUpdateUnstable controls whether to receive unstable version updates
var FieldAppUpdateUnstable = &Field{
	Name:        "update.unstable",
	Group:       GroupApplication,
	Type:        FieldTypeBool,
	Default:     defaultBool(DefaultAppUpdateUnstable),
	Description: "Receive updates for unstable versions.",
}

// FieldAppUpdateAuto controls automatic application updates
var FieldAppUpdateAuto = &Field{
	Name:        "update.auto",
	Group:       GroupApplication,
	Type:        FieldTypeBool,
	Default:     defaultBool(DefaultAppUpdateAuto),
	Description: "Automatically update the application when a new version is available.",
}

// FieldAppUpdatePeriod defines how often to check for updates
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
