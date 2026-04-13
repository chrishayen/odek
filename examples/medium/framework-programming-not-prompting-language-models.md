# Requirement: "a framework for composing language-model programs from typed modules instead of free-form prompts"

A module takes a typed input signature and produces a typed output; a compiler links modules and optimizes their instructions from examples. The model call is a std primitive.

std
  std.llm
    std.llm.complete
      @ (prompt: string) -> result[string, string]
      + returns the completion text for the prompt
      - returns error when the backend rejects the request
      # llm

lm_programs
  lm_programs.new_signature
    @ (inputs: list[string], outputs: list[string]) -> signature
    + returns a signature describing named inputs and outputs
    # construction
  lm_programs.new_module
    @ (sig: signature, instruction: string) -> module_state
    + returns a module that invokes the model under the given instruction
    # construction
  lm_programs.call_module
    @ (module: module_state, inputs: map[string, string]) -> result[map[string, string], string]
    + binds inputs into the prompt and parses outputs from the completion
    - returns error when a required input is missing
    - returns error when the completion does not contain every declared output
    # execution
    -> std.llm.complete
  lm_programs.compile_chain
    @ (modules: list[module_state]) -> module_state
    + returns one composite module that threads outputs of each step into the next
    # composition
  lm_programs.optimize_instruction
    @ (module: module_state, examples: list[tuple[map[string, string], map[string, string]]]) -> result[module_state, string]
    + returns a module with a refined instruction derived from example input-output pairs
    - returns error when the examples list is empty
    # optimization
    -> std.llm.complete
