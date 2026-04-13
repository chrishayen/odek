# Requirement: "a library that generates web form interfaces from typed function signatures"

Takes a function description with typed parameters and produces html form markup plus a request decoder.

std
  std.html
    std.html.escape
      @ (s: string) -> string
      + escapes &, <, >, ", and ' for safe html text
      # html
  std.url
    std.url.parse_form
      @ (raw: string) -> map[string,string]
      + parses application/x-www-form-urlencoded into a map
      # url

uiform
  uiform.declare_param
    @ (name: string, kind: string, required: bool) -> param_spec
    + builds a parameter spec with a primitive kind
    ? kind is one of "string", "int", "float", "bool"
    # declaration
  uiform.declare_function
    @ (name: string, params: list[param_spec]) -> function_spec
    + builds a function spec from a name and ordered parameters
    # declaration
  uiform.render_form
    @ (fn: function_spec, action_url: string) -> string
    + returns an html form whose inputs correspond to each parameter
    + each input has name, label, type, and required attributes
    # rendering
    -> std.html.escape
  uiform.decode_submission
    @ (fn: function_spec, body: string) -> result[map[string,string], string]
    + parses a submitted form body and returns only declared parameters
    - returns error when a required parameter is missing
    - returns error when a value fails its declared kind
    # decoding
    -> std.url.parse_form
  uiform.render_result_page
    @ (fn: function_spec, result_text: string) -> string
    + returns an html page showing the rendered form and the most recent result
    # rendering
    -> std.html.escape
