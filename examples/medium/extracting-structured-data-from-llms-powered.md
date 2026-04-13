# Requirement: "a library for extracting structured data from language model responses using schemas"

Given a schema and a model's free-form response, coerce it into a validated record, retrying with a repair prompt when needed.

std
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses any JSON value (object, array, scalar) into a generic value
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

extractor
  extractor.define_schema
    @ (fields: list[tuple[string, string]]) -> schema
    + returns a schema from a list of (field_name, type_name) pairs
    ? type_name is one of "string", "int", "float", "bool"
    # schema
  extractor.build_prompt
    @ (schema: schema, instruction: string) -> string
    + returns a prompt that embeds the schema and instructs the model to reply in JSON
    # prompting
    -> std.json.encode_object
  extractor.parse_response
    @ (schema: schema, response: string) -> result[map[string, string], string]
    + returns the extracted record when the response parses and every field matches its declared type
    - returns error when the response contains no JSON object
    - returns error when a required field is missing
    - returns error when a field's value does not match its declared type
    # parsing
    -> std.json.parse_value
  extractor.build_repair_prompt
    @ (original: string, error: string) -> string
    + returns a follow-up prompt asking the model to fix its previous response, given the error
    # repair
  extractor.extract
    @ (schema: schema, instruction: string, call_model: fn(string) -> string, max_attempts: i32) -> result[map[string, string], string]
    + runs the prompt, parses the response, and retries up to max_attempts using the repair prompt on failure
    - returns the last parse error when attempts are exhausted
    # orchestration
