# Requirement: "a library that generates http input/output handling boilerplate"

Given a description of a handler's request shape and response shape, produce the decoding, validation, and encoding code as a single source string.

std: (all units exist)

httpgen
  httpgen.describe_field
    @ (name: string, type_name: string, source: string, required: bool) -> field_spec
    + builds a field spec; source is "path", "query", "header", or "body"
    - returns an error spec when source is outside the supported set
    # specification
  httpgen.describe_handler
    @ (name: string, method: string, path: string, inputs: list[field_spec], output_type: string) -> handler_spec
    + builds a handler spec with method, path, inputs, and the response type name
    # specification
  httpgen.render_decoder
    @ (h: handler_spec) -> string
    + returns the source text that reads each input field from its source and checks required-ness
    # rendering
  httpgen.render_encoder
    @ (h: handler_spec) -> string
    + returns the source text that serializes the response type to the http response body
    # rendering
  httpgen.render_handler
    @ (h: handler_spec) -> string
    + returns the combined decoder plus dispatch plus encoder as a single function body
    # rendering
