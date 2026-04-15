# Requirement: "a tag-based configuration parser that loads values from multiple providers into typed fields"

Given a schema of fields with type tags and source hints, look up values across a chain of providers (env vars, file, defaults) and produce a populated config map. Parsing primitive types goes through a thin std utility.

std
  std.env
    std.env.get
      fn (key: string) -> optional[string]
      + returns the environment variable value when set
      # env
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the full contents of the file
      - returns error when the file does not exist
      # filesystem
  std.strconv
    std.strconv.parse_i64
      fn (s: string) -> result[i64, string]
      + parses a decimal integer
      - returns error on non-digit input
      # parsing
    std.strconv.parse_bool
      fn (s: string) -> result[bool, string]
      + parses true/false/1/0/yes/no
      - returns error on other input
      # parsing
    std.strconv.parse_f64
      fn (s: string) -> result[f64, string]
      + parses a decimal floating-point number
      - returns error on non-numeric input
      # parsing

tagconf
  tagconf.schema_new
    fn () -> schema_state
    + creates an empty schema
    # construction
  tagconf.schema_field
    fn (schema: schema_state, name: string, type_tag: string, env_key: string, default_value: string) -> schema_state
    + declares a field with its type tag, env key, and default
    ? type_tag is one of "string", "i64", "f64", "bool"
    # schema
  tagconf.load_file_provider
    fn (path: string) -> result[map[string, string], string]
    + reads a simple key=value file into a provider map
    - returns error when the file is not readable
    # provider
    -> std.fs.read_all
  tagconf.env_provider
    fn (schema: schema_state) -> map[string, string]
    + looks up each field's env_key and collects present values
    # provider
    -> std.env.get
  tagconf.resolve
    fn (schema: schema_state, providers: list[map[string, string]]) -> result[map[string, typed_value], string]
    + for each field picks the first provider that supplies a value, falling back to the default, and parses per type tag
    - returns error when a required field has no value and no default
    - returns error when a value fails to parse under its type tag
    # resolution
    -> std.strconv.parse_i64
    -> std.strconv.parse_f64
    -> std.strconv.parse_bool
  tagconf.get_i64
    fn (config: map[string, typed_value], name: string) -> result[i64, string]
    + returns the integer field
    - returns error when the field is missing or not an integer
    # access
  tagconf.get_string
    fn (config: map[string, typed_value], name: string) -> result[string, string]
    + returns the string field
    - returns error when the field is missing or not a string
    # access
