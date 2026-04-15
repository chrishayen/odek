# Requirement: "a library for a lightweight database engine based on stack data structures"

Stores multiple named stacks with push/pop/peek. A thin routing layer maps verbs and paths to engine operations so callers can expose it over any transport.

std: (all units exist)

stackdb
  stackdb.new
    fn () -> db_state
    + creates an empty database with no stacks
    # construction
  stackdb.create_stack
    fn (state: db_state, name: string) -> result[db_state, string]
    + adds a new empty stack under the given name
    - returns error when a stack with that name already exists
    # stack_management
  stackdb.delete_stack
    fn (state: db_state, name: string) -> result[db_state, string]
    + removes the named stack and discards its contents
    - returns error when no stack with that name exists
    # stack_management
  stackdb.push
    fn (state: db_state, name: string, value: string) -> result[db_state, string]
    + appends value to the top of the named stack
    - returns error when the stack does not exist
    # operations
  stackdb.pop
    fn (state: db_state, name: string) -> result[tuple[string, db_state], string]
    + removes and returns the top value of the named stack
    - returns error when the stack is empty or does not exist
    # operations
  stackdb.peek
    fn (state: db_state, name: string) -> result[string, string]
    + returns the top value without removing it
    - returns error when the stack is empty or does not exist
    # operations
  stackdb.size
    fn (state: db_state, name: string) -> result[i32, string]
    + returns the number of elements in the named stack
    - returns error when the stack does not exist
    # operations
  stackdb.route
    fn (state: db_state, verb: string, path: string, body: optional[string]) -> result[tuple[string, db_state], string]
    + dispatches (verb, path) to the matching engine operation and returns a textual response
    - returns error on unknown verb or unmatched path
    # routing
