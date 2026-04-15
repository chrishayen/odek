# Requirement: "a form and JSON parsing and validation library with multipart and file support"

Parses request bodies in multiple encodings into a typed form, then runs declared validators.

std
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses arbitrary JSON
      - returns error on malformed input
      # serialization
  std.url
    std.url.parse_query
      fn (raw: string) -> map[string, string]
      + parses an application/x-www-form-urlencoded body into a map
      # url
  std.multipart
    std.multipart.parse
      fn (body: bytes, boundary: string) -> result[list[multipart_part], string]
      + parses a multipart body into parts with headers and content
      - returns error on malformed boundary
      # multipart

forms
  forms.new_schema
    fn () -> form_schema
    + creates an empty validation schema
    # construction
  forms.field
    fn (schema: form_schema, name: string, type_name: string, required: bool) -> form_schema
    + declares a field with a type keyword like "string", "int", "float", "bool"
    # schema
  forms.rule
    fn (schema: form_schema, field: string, rule: validation_rule) -> form_schema
    + attaches a validation rule to a field (min, max, pattern, etc.)
    # schema
  forms.file_field
    fn (schema: form_schema, name: string, max_bytes: i64, required: bool) -> form_schema
    + declares a file upload field with a maximum size
    # schema
  forms.parse_body
    fn (schema: form_schema, content_type: string, body: bytes) -> result[parsed_form, string]
    + dispatches on content type to parse JSON, urlencoded, or multipart bodies
    - returns error on unsupported content type
    # parsing
    -> std.json.parse
    -> std.url.parse_query
    -> std.multipart.parse
  forms.validate
    fn (schema: form_schema, form: parsed_form) -> result[parsed_form, list[validation_error]]
    + runs all declared rules and returns the form or a list of errors
    - returns errors for missing required fields
    - returns errors for files exceeding max_bytes
    # validation
  forms.get_string
    fn (form: parsed_form, name: string) -> result[string, string]
    + returns a string field value
    - returns error when the field is missing or of another type
    # access
  forms.get_file
    fn (form: parsed_form, name: string) -> result[uploaded_file, string]
    + returns an uploaded file and its metadata
    - returns error when the field is missing
    # access
