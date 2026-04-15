# Requirement: "a framework for interactive command-line interfaces"

A library for prompting users at the terminal: ask questions, validate answers, and collect them into a result map. Transport of bytes to and from the TTY is the caller's concern; the framework produces prompt strings and parses response strings.

std
  std.strings
    std.strings.trim
      fn (s: string) -> string
      + removes leading and trailing ASCII whitespace
      # text
    std.strings.parse_i32
      fn (s: string) -> result[i32, string]
      + parses a signed decimal integer
      - returns error on non-numeric input
      # parsing

cliprompt
  cliprompt.new_session
    fn () -> session_state
    + creates an empty session with no questions and no answers
    # construction
  cliprompt.add_text
    fn (state: session_state, key: string, prompt: string, default_value: string) -> session_state
    + appends a free-text question with an optional default
    # question_definition
  cliprompt.add_choice
    fn (state: session_state, key: string, prompt: string, choices: list[string]) -> session_state
    + appends a multiple-choice question
    - returns unchanged state when choices is empty
    # question_definition
  cliprompt.add_confirm
    fn (state: session_state, key: string, prompt: string, default_yes: bool) -> session_state
    + appends a yes/no question
    # question_definition
  cliprompt.next_prompt
    fn (state: session_state) -> optional[string]
    + returns the rendered prompt string for the next unanswered question
    - returns none when every question has an answer
    # prompt_rendering
  cliprompt.submit_answer
    fn (state: session_state, raw: string) -> result[session_state, string]
    + validates raw against the current question's type and stores it
    - returns error "required" when raw is empty and no default is set
    - returns error when the input does not match an allowed choice
    # answer_validation
    -> std.strings.trim
    -> std.strings.parse_i32
  cliprompt.answers
    fn (state: session_state) -> map[string,string]
    + returns the collected answers keyed by question key
    # results
