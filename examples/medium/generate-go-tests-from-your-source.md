# Requirement: "a unit test generator that reads source code"

Given a description of functions in a source file, produce a skeleton test for each one with placeholders for inputs and expected outputs.

std: (all units exist)

testgen
  testgen.describe_function
    @ (name: string, params: list[param], returns: list[string]) -> function_spec
    + builds a function spec with a name, parameter list, and return types
    # specification
  testgen.plan_cases
    @ (f: function_spec) -> list[test_case]
    + returns one happy-path case plus one case per parameter exercising a zero/empty value
    ? plans do not invent error cases for functions that do not return an error-like type
    # planning
  testgen.case_name
    @ (c: test_case) -> string
    + returns a descriptive snake_case name for the test case
    # naming
  testgen.render_case
    @ (f: function_spec, c: test_case) -> string
    + returns the source text of a single subtest body with placeholders for inputs and expected values
    # rendering
  testgen.render_file
    @ (f: function_spec) -> string
    + returns a complete test file covering every case planned for the function
    # rendering
