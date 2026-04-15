# Requirement: "a fast JSON schema validator driven by generated validation code"

The schema is compiled once into a validator program; runtime validation then walks a value and interprets the compiled instructions.

std
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses arbitrary JSON into a tagged value tree
      - returns error on malformed input
      # serialization
    std.json.type_of
      fn (value: json_value) -> string
      + returns one of "null", "bool", "number", "string", "array", "object"
      # serialization
    std.json.get_field
      fn (value: json_value, name: string) -> optional[json_value]
      + returns the field of an object by name
      # serialization
    std.json.array_items
      fn (value: json_value) -> list[json_value]
      + returns the elements of a JSON array
      # serialization
  std.regex
    std.regex.compile
      fn (pattern: string) -> result[regex_state, string]
      + compiles a regex pattern
      - returns error on invalid syntax
      # regex
    std.regex.matches
      fn (re: regex_state, input: string) -> bool
      + reports whether the input matches
      # regex

schema_validator
  schema_validator.compile
    fn (schema_raw: string) -> result[validator_program, string]
    + compiles a JSON Schema document into an executable validator
    - returns error when the schema itself is malformed
    # compilation
    -> std.json.parse_value
  schema_validator.compile_type_check
    fn (expected_type: string) -> validator_instruction
    + emits an instruction that asserts the node type
    # compilation
  schema_validator.compile_properties
    fn (fields: map[string, validator_program]) -> validator_instruction
    + emits an instruction that validates each named field against its subschema
    # compilation
  schema_validator.compile_items
    fn (item_program: validator_program) -> validator_instruction
    + emits an instruction that validates every array element
    # compilation
  schema_validator.compile_pattern
    fn (pattern: string) -> result[validator_instruction, string]
    + emits an instruction that matches a string against a regex
    - returns error when the pattern is invalid
    # compilation
    -> std.regex.compile
  schema_validator.compile_required
    fn (names: list[string]) -> validator_instruction
    + emits an instruction that asserts the listed fields are present
    # compilation
  schema_validator.execute
    fn (program: validator_program, raw: string) -> result[void, list[string]]
    + parses the input and runs the compiled program against it
    - returns the list of validation failures when any instruction rejects
    # execution
    -> std.json.parse_value
  schema_validator.step
    fn (instruction: validator_instruction, node: json_value) -> list[string]
    + executes a single instruction, returning the errors it produced
    # execution
    -> std.json.type_of
    -> std.json.get_field
    -> std.json.array_items
    -> std.regex.matches
  schema_validator.validate
    fn (program: validator_program, raw: string) -> bool
    + returns true when execute produces no errors
    # execution
  schema_validator.first_error
    fn (program: validator_program, raw: string) -> optional[string]
    + returns the first validation failure, if any
    # execution
