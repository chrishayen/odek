# Requirement: "a source code refactoring library"

Represents a project as a set of files, parses them into a symbol table, and exposes refactoring operations that return text edits.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the complete contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns every file path under root, depth-first
      - returns error when root is not a directory
      # filesystem

refactor
  refactor.open_project
    @ (root: string) -> result[project_state, string]
    + loads every source file under root into an in-memory project
    - returns error when root is not a directory
    # project
    -> std.fs.walk
    -> std.fs.read_all
  refactor.parse_module
    @ (source: string) -> result[module_ast, string]
    + returns the parsed AST for a module
    - returns error on syntax errors with line and column
    # parsing
  refactor.build_symbol_table
    @ (project: project_state) -> project_state
    + walks every module and records definitions, references, and scopes
    # analysis
  refactor.find_definition
    @ (project: project_state, path: string, line: i32, column: i32) -> optional[source_location]
    + returns the definition location for the identifier at the cursor
    # navigation
  refactor.find_references
    @ (project: project_state, path: string, line: i32, column: i32) -> list[source_location]
    + returns every reference to the identifier at the cursor
    # navigation
  refactor.rename
    @ (project: project_state, path: string, line: i32, column: i32, new_name: string) -> result[list[text_edit], string]
    + returns the edits needed to rename the symbol across the project
    - returns error when new_name collides with an existing symbol in scope
    # refactoring
  refactor.extract_variable
    @ (project: project_state, path: string, start: i32, end: i32, name: string) -> result[list[text_edit], string]
    + returns edits that extract the selected expression into a new local variable
    - returns error when the selection does not enclose a complete expression
    # refactoring
  refactor.extract_function
    @ (project: project_state, path: string, start: i32, end: i32, name: string) -> result[list[text_edit], string]
    + returns edits that extract the selected statements into a new function
    - returns error when the selection crosses function boundaries
    # refactoring
  refactor.inline_variable
    @ (project: project_state, path: string, line: i32, column: i32) -> result[list[text_edit], string]
    + returns edits that replace every reference with the assigned value and remove the declaration
    - returns error when the variable is reassigned
    # refactoring
  refactor.apply_edits
    @ (project: project_state, edits: list[text_edit]) -> project_state
    + returns the project with edits applied to in-memory sources
    # refactoring
