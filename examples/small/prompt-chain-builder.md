# Requirement: "a library for building a chain of interactive prompts"

Each step asks a question, reads an answer, and picks the next step based on the answer. The chain is a state machine over named steps.

std: (all units exist)

prompt_chain
  prompt_chain.new
    @ () -> chain_state
    + creates an empty chain with no steps
    # construction
  prompt_chain.add_step
    @ (state: chain_state, id: string, question: string, next_fn: transition_fn) -> result[chain_state, string]
    + registers a step with a prompt and a transition function from answer to next id
    - returns error when a step with that id already exists
    # configuration
  prompt_chain.run
    @ (state: chain_state, start_id: string, reader: input_reader, writer: output_writer) -> result[map[string, string], string]
    + walks the chain from start_id, collecting answers keyed by step id
    + terminates when a transition returns the empty id
    - returns error when start_id or any next id is unknown
    - returns error when the reader fails
    # execution
