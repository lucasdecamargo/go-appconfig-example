# Go Application Configuration Template

A production-ready template for managing configuration parameters in Go applications using **Viper** and **Cobra**. This template demonstrates a clean, maintainable approach to configuration management with a single source of truth for all configuration metadata.

## 🎯 Key Features

- **Single Source of Truth**: Configuration fields are defined once with all metadata (validation, documentation, defaults)
- **Type Safety**: Strongly typed configuration with validation
- **CLI Integration**: Seamless integration with Cobra for command-line interfaces
- **Multiple Formats**: Support for YAML, JSON, TOML, HCL, and ENV files
- **Environment Variables**: Automatic binding with configurable prefix
- **Shell Completion**: Auto-completion for configuration field names and values
- **Documentation**: Built-in help and documentation generation
- **Validation**: Multiple validation strategies (tags, custom functions, valid values)

## 🏗️ Architecture Overview

### Core Concept: Field-Driven Configuration

The central idea is to define configuration fields as structured data that contains everything needed to:
- Validate values
- Generate CLI flags
- Provide documentation
- Set defaults
- Handle serialization

```go
type Field struct {
    Name         string            // Configuration key
    Group        string            // Logical grouping
    Type         FieldType         // Data type
    Default      any               // Default value
    Description  string            // CLI help text
    Docstring    string            // Detailed documentation
    ValidValues  []any             // Allowed values
    ValidateTag  string            // Go validator tag
    ValidateFunc func(any) error   // Custom validation
    // ... more metadata
}
```

### Repository Structure

```
go-appconfig-example/
├── cmd/                    # CLI command implementations
│   ├── root.go            # Root command and global flags
│   ├── config.go          # Configuration management commands
│   └── pager.go           # Output pagination utility
├── internal/
│   ├── config/            # Configuration management core
│   │   ├── fields.go      # Field definition and collection
│   │   ├── fields_app.go  # Application-specific fields
│   │   ├── fields_network.go # Network-specific fields
│   │   ├── config.go      # Viper integration
│   │   ├── defaults.go    # Default value conversion
│   │   └── validators.go  # Custom validation functions
│   └── consts/            # Application constants
│       └── consts.go      # App name, version, etc.
├── main.go                # Application entry point
└── go.mod                 # Go module definition
```

## 🚀 Getting Started

### Prerequisites

- Go 1.25 or later
- Git

### Installation

```bash
git clone https://github.com/lucasdecamargo/go-appconfig-example.git
cd go-appconfig-example
go mod download
```

### Building

```bash
# Basic build
go build -o confapp

# Build with custom defaults
go build -ldflags "-X github.com/lucasdecamargo/go-appconfig-example/internal/config.DefaultAppEnvironment=prod" -o confapp
```

### Usage

The example application `confapp` provides several commands for configuration management:

```bash
# List all configuration values
./confapp config list

# List specific configuration groups
./confapp config list log
./confapp config list proxy

# Describe configuration fields (with pager)
./confapp config describe

# Describe specific fields
./confapp config describe log.level update

# Set configuration values
./confapp config set --log.level debug
./confapp config set --log.level info --log.output /var/log/app.log

# Show hidden fields
./confapp config list --hidden
./confapp config describe --hidden

# Verbose output
./confapp --verbose config list
```

## 📋 Configuration Fields

### Application Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `environment` | string | `dev` | Application environment (hidden) |
| `log.level` | string | `info` | Logging level (debug, info, warn, error) |
| `log.output` | string | `nil` | Log output file path |
| `log.format` | string | `text` | Log format (json, text) |
| `update.unstable` | bool | `false` | Enable unstable updates |
| `update.auto` | bool | `false` | Enable auto-updates |
| `update.period` | duration | `15m` | Update check period |

### Network Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `proxy.all` | string | `nil` | Proxy for all traffic |
| `proxy.http` | string | `nil` | HTTP proxy |
| `proxy.https` | string | `nil` | HTTPS proxy |

## 🔧 Configuration Sources

The application reads configuration from multiple sources in order of precedence:

1. **Command-line flags** (highest priority)
2. **Environment variables** (with `CONFAPP_` prefix)
3. **Configuration file** (YAML, JSON, TOML, HCL, ENV)
4. **Default values** (lowest priority)

### Environment Variables

Environment variables are automatically bound with the `CONFAPP_` prefix. Dots in field names are replaced with underscores:

```bash
export CONFAPP_LOG_LEVEL=debug
export CONFAPP_UPDATE_AUTO=true
export CONFAPP_PROXY_HTTP=http://proxy:8080
```

### Configuration File

By default, the application looks for a configuration file at:
- `~/.config/confapp/config.yaml` (Unix/macOS)
- `%APPDATA%\confapp\config.yaml` (Windows)

You can specify a custom location with the `--config` flag:

```bash
./confapp --config /path/to/config.yaml config list
```

Example configuration file (`config.yaml`):
```yaml
log:
  level: debug
  output: /var/log/app.log
  format: json

update:
  auto: true
  period: 1h

proxy:
  http: http://proxy:8080
```

## 🛠️ Extending the Template

### Adding New Configuration Fields

1. **Define the field** in the appropriate `fields_*.go` file:

```go
// In internal/config/fields_app.go
var FieldAppDatabase = &Field{
    Name:        "database.url",
    Group:       GroupApplication,
    Type:        FieldTypeString,
    Default:     defaultString("postgres://localhost:5432/app"),
    Description: "Database connection URL",
    ValidateTag: "url",
    Example:     "postgres://user:pass@host:5432/db",
}
```

2. **Register the field** in the `init()` function:

```go
func init() {
    Fields.Add(
        // ... existing fields
        FieldAppDatabase,
    )
}
```

3. **Add validation** if needed in `validators.go`:

```go
func validateDatabaseURL(v any) error {
    // Custom validation logic
    return nil
}
```

### Adding New Field Types

1. **Define the type** in `fields.go`:

```go
const (
    // ... existing types
    FieldTypeURL FieldType = "url"
)
```

2. **Add default conversion** in `defaults.go`:

```go
func defaultURL(val string) any {
    if val == "" {
        return nil
    }
    // URL parsing logic
    return val
}
```

3. **Add CLI flag setup** in `cmd/config.go`:

```go
case config.FieldTypeURL:
    setupURLFlag(flags, field)
```

### Adding New Commands

1. **Create the command** in `cmd/`:

```go
var customCmd = &cobra.Command{
    Use:   "custom",
    Short: "Custom command description",
    RunE:  runCustom,
}
```

2. **Register the command** in the appropriate `init()` function:

```go
func init() {
    rootCmd.AddCommand(customCmd)
}
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/config -v
```

## 📚 Best Practices

### Field Definition

- **Use descriptive names**: `log.level` instead of `loglevel`
- **Group related fields**: Use consistent group names
- **Provide examples**: Help users understand expected values
- **Document thoroughly**: Use both `Description` and `Docstring`
- **Validate inputs**: Use validation tags and custom functions

### Configuration Management

- **Single source of truth**: Define fields once, use everywhere
- **Type safety**: Use strongly typed fields
- **Sensible defaults**: Provide good default values
- **Environment awareness**: Use environment-specific defaults
- **Validation**: Validate at multiple levels

### CLI Design

- **Consistent naming**: Use kebab-case for flags
- **Helpful descriptions**: Write clear, concise descriptions
- **Shell completion**: Enable auto-completion for better UX
- **Error handling**: Provide clear error messages
- **Verbose mode**: Include debug information when needed

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Update documentation
6. Submit a pull request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Viper](https://github.com/spf13/viper) - Configuration solution for Go applications
- [Cobra](https://github.com/spf13/cobra) - Framework for creating powerful modern CLI applications
- [Go Validator](https://github.com/go-playground/validator) - Struct validation library

---

**Note**: This template is designed to be a starting point for production Go applications. Adapt it to your specific needs and requirements.
