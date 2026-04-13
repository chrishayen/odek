# Requirement: "a pipeline-first string transformation toolkit"

Build a reusable transformation pipeline and apply it to strings.

std: (all units exist)

str_pipe
  str_pipe.new
    @ () -> pipe_state
    + returns an empty pipeline
    # construction
  str_pipe.then_trim
    @ (state: pipe_state) -> pipe_state
    + appends a step that strips leading and trailing whitespace
    # step
  str_pipe.then_lower
    @ (state: pipe_state) -> pipe_state
    + appends a step that lowercases all characters
    # step
  str_pipe.then_replace
    @ (state: pipe_state, from: string, to: string) -> pipe_state
    + appends a step that replaces every occurrence of from with to
    # step
  str_pipe.apply
    @ (state: pipe_state, input: string) -> string
    + runs the pipeline over input and returns the result
    + returns input unchanged when the pipeline is empty
    # apply
