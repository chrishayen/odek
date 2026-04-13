# Requirement: "a statically typed compiled scripting language with hot reload"

A full language front-end plus a reload-capable runtime. Lexer, parser, type checker, and bytecode emitter are the compilation pipeline; the runtime owns a module registry that swaps code while preserving live state.

std
  std.strings
    std.strings.split_lines
      @ (s: string) -> list[string]
      + splits s on newline characters
      # strings
  std.io
    std.io.read_all
      @ (path: string) -> result[string, string]
      + reads a text file into a string
      - returns error when path is unreadable
      # io
  std.hash
    std.hash.fnv1a_64
      @ (data: bytes) -> u64
      + returns the FNV-1a 64-bit hash of data
      # hashing

tiny_lang
  tiny_lang.tokenize
    @ (source: string) -> result[list[token], string]
    + returns tokens for identifiers, numbers, strings, operators, and keywords
    - returns error with line and column on an unrecognized character
    # lexing
    -> std.strings.split_lines
  tiny_lang.parse
    @ (tokens: list[token]) -> result[ast_module, string]
    + returns an AST module containing functions, type declarations, and top-level statements
    - returns error on unexpected tokens, annotated with position
    # parsing
  tiny_lang.type_check
    @ (module: ast_module) -> result[typed_module, string]
    + annotates every expression with its inferred type
    - returns error when a call's argument types do not match the declared signature
    - returns error when a variable is used before declaration
    # type_checking
  tiny_lang.emit_bytecode
    @ (module: typed_module) -> result[bytecode_module, string]
    + lowers a typed module to a linear bytecode sequence with a constant pool
    # codegen
  tiny_lang.compile
    @ (source: string) -> result[bytecode_module, string]
    + full pipeline: tokenize, parse, type check, emit
    - propagates the first stage's error
    # compilation
  tiny_lang.new_runtime
    @ () -> runtime_state
    + creates an empty runtime with no loaded modules
    # runtime
  tiny_lang.load_module
    @ (state: runtime_state, name: string, bytecode: bytecode_module) -> result[runtime_state, string]
    + registers bytecode under name and initializes its globals
    - returns error when name is already loaded
    # runtime
    -> std.hash.fnv1a_64
  tiny_lang.reload_module
    @ (state: runtime_state, name: string, bytecode: bytecode_module) -> result[runtime_state, string]
    + replaces the bytecode for name while preserving the globals whose type signature is unchanged
    + resets globals whose types changed to their new initializer
    - returns error when name is not currently loaded
    # hot_reload
    -> std.hash.fnv1a_64
  tiny_lang.call
    @ (state: runtime_state, module: string, func: string, args: list[value]) -> result[value, string]
    + invokes the named function and returns its result
    - returns error when module or func is unknown
    - returns error when arg arity or types do not match
    # invocation
  tiny_lang.load_from_file
    @ (state: runtime_state, name: string, path: string) -> result[runtime_state, string]
    + reads source from path, compiles it, and loads it under name
    - returns error when the file is missing or fails to compile
    # convenience
    -> std.io.read_all
