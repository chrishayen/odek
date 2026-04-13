# Requirement: "a library for validating configuration files in multiple formats"

Detects the format from the extension or content, parses it, and reports syntax errors with file, line, and message.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads an entire file into a string
      - returns error when the file cannot be opened
      # filesystem
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns all file paths under root
      - returns error when root is not a directory
      # filesystem
  std.json
    std.json.parse
      @ (raw: string) -> result[void, parse_error]
      + succeeds on valid JSON
      - returns a parse error with line and column
      # parsing
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[void, parse_error]
      + succeeds on valid YAML
      - returns a parse error with line and column
      # parsing
  std.toml
    std.toml.parse
      @ (raw: string) -> result[void, parse_error]
      + succeeds on valid TOML
      - returns a parse error with line and column
      # parsing

config_validator
  config_validator.detect_format
    @ (path: string, content: string) -> config_format
    + returns json, yaml, toml, or unknown based on extension
    + falls back to content sniffing when the extension is ambiguous
    # detection
  config_validator.validate_file
    @ (path: string) -> result[void, validation_error]
    + reads and validates a file using its detected format
    - returns a validation error with path, line, column, and message on failure
    - returns error when the format cannot be detected
    # validation
    -> std.fs.read_all
    -> std.json.parse
    -> std.yaml.parse
    -> std.toml.parse
  config_validator.validate_string
    @ (content: string, format: config_format) -> result[void, validation_error]
    + validates content against an explicit format
    - returns a validation error when parsing fails
    # validation
  config_validator.validate_tree
    @ (root: string) -> list[validation_error]
    + walks a directory and returns errors for every file that fails validation
    # validation
    -> std.fs.walk
