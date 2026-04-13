# Requirement: "an embeddable scripting language with concurrent event processing"

A tiny embeddable language: lex, parse, evaluate. Event handlers are registered by name and invoked concurrently by dispatching events to matching handlers.

std
  std.collections
    std.collections.list_append
      @ (items: list[string], item: string) -> list[string]
      + returns a new list with the item appended
      # collections
  std.strings
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits s on sep
      # strings

scriptlang
  scriptlang.tokenize
    @ (source: string) -> result[list[token], string]
    + recognizes identifiers, numbers, strings, operators, and keywords
    - returns error on unterminated string literal
    - returns error on unexpected character
    # lexing
  scriptlang.parse
    @ (tokens: list[token]) -> result[ast_node, string]
    + parses a program consisting of statements and expressions
    + handles function declarations and event-handler declarations
    - returns error on unexpected token
    - returns error on missing closing bracket
    # parsing
  scriptlang.new_interpreter
    @ () -> interpreter_state
    + returns a fresh interpreter with empty bindings and no handlers
    # construction
  scriptlang.define
    @ (state: interpreter_state, name: string, value: value) -> interpreter_state
    + binds the name to a value in the global scope
    + overwrites any existing binding
    # environment
  scriptlang.eval
    @ (state: interpreter_state, node: ast_node) -> result[tuple[value, interpreter_state], string]
    + evaluates an expression or statement and returns the resulting value
    - returns error on unbound identifier
    - returns error on type mismatch
    # evaluation
  scriptlang.register_handler
    @ (state: interpreter_state, event_kind: string, handler: ast_node) -> interpreter_state
    + associates the handler with the event kind
    + allows multiple handlers per kind; order of registration is preserved
    # events
  scriptlang.dispatch
    @ (state: interpreter_state, event_kind: string, payload: value) -> result[tuple[list[value], interpreter_state], string]
    + invokes every handler registered for the kind and returns their results
    + concurrent handlers observe a consistent snapshot of state at dispatch time
    - returns error when any handler raises
    # events
  scriptlang.run_program
    @ (state: interpreter_state, source: string) -> result[interpreter_state, string]
    + tokenizes, parses, and evaluates the source to update the interpreter
    - returns error on any lexing, parsing, or evaluation failure
    # facade
    -> scriptlang.tokenize
    -> scriptlang.parse
    -> scriptlang.eval
