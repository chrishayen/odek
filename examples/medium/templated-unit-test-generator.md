# Requirement: "a templated unit test generator"

Like a plain test skeleton generator, but the output shape is controlled by a user-supplied template string with named placeholders.

std: (all units exist)

tpl_testgen
  tpl_testgen.describe_function
    @ (name: string, params: list[param], returns: list[string]) -> function_spec
    + builds a function spec with a name, parameter list, and return types
    # specification
  tpl_testgen.parse_template
    @ (raw: string) -> result[template, string]
    + parses a template containing {{name}}, {{params}}, {{returns}}, and {{case}} placeholders
    - returns error on an unterminated placeholder
    - returns error on an unknown placeholder name
    # template_parsing
  tpl_testgen.bindings
    @ (f: function_spec, case_name: string) -> map[string, string]
    + returns the placeholder-to-text map for rendering a single case
    # binding
  tpl_testgen.render
    @ (tpl: template, bindings: map[string, string]) -> string
    + returns the template expanded with the given bindings
    + unknown placeholders in the bindings map are ignored
    # rendering
  tpl_testgen.render_file
    @ (tpl: template, f: function_spec, cases: list[string]) -> string
    + concatenates one render per case and returns the full test file text
    # rendering
