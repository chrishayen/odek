# Requirement: "a library for transpiling a general-purpose language into embedded microcontroller code"

Parse the source, lower high-level features that the target cannot express, and emit microcontroller-flavored code.

std: (all units exist)

embed_transpile
  embed_transpile.parse
    @ (source: string) -> result[source_ast, string]
    + parses the supported subset (functions, integer and float types, arrays, structs, if, for, while)
    - returns error with line and column on syntax errors
    # parsing
  embed_transpile.check_supported
    @ (ast: source_ast) -> result[void, string]
    + verifies only the supported subset is used; rejects goroutines, channels, interfaces, closures, maps, and garbage-collected heap allocation
    - returns error naming the first unsupported construct
    # validation
  embed_transpile.lower_strings
    @ (ast: source_ast) -> source_ast
    + rewrites string operations as fixed-size character array operations
    # lowering
  embed_transpile.lower_slices
    @ (ast: source_ast) -> source_ast
    + rewrites dynamic slices to fixed-size arrays with explicit length parameters
    # lowering
  embed_transpile.map_stdlib_to_hal
    @ (ast: source_ast) -> source_ast
    + replaces stdlib IO calls with hardware-abstraction-layer equivalents (pinMode, digitalWrite, Serial.print)
    # mapping
  embed_transpile.emit_target
    @ (ast: source_ast) -> string
    + emits target source including a setup() and loop() entry point synthesized from the input's main function
    # codegen
  embed_transpile.compile
    @ (source: string) -> result[string, string]
    + runs parse, validation, lowering, mapping, and emission in sequence
    - returns error with the first failing stage's message
    # pipeline
