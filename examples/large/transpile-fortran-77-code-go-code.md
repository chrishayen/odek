# Requirement: "a library for transpiling fortran 77 source code to a modern structured language"

A source-to-source compiler: lex, parse, build symbol tables, lower to a target AST, and emit code.

std: (all units exist)

fortran_transpile
  fortran_transpile.tokenize
    @ (source: string) -> result[list[token], string]
    + produces tokens respecting fixed-form column rules (columns 1-5 labels, 6 continuation, 7-72 body)
    + recognizes keywords, identifiers, integer/real/double literals, operators, and strings
    - returns error with line and column on unterminated string literals
    # lexing
  fortran_transpile.parse
    @ (tokens: list[token]) -> result[program_ast, string]
    + parses program, subroutine, and function units with declarations and executable statements
    + handles DO, IF, GOTO, CALL, RETURN, CONTINUE, FORMAT, COMMON, DIMENSION
    - returns error with location on unexpected tokens
    # parsing
  fortran_transpile.resolve_symbols
    @ (ast: program_ast) -> result[program_ast, string]
    + annotates every identifier with its kind (local, parameter, common block, subroutine)
    + infers implicit types following I-N integer / else real rules
    - returns error when a label referenced by GOTO is undefined
    # semantics
  fortran_transpile.lower_arrays
    @ (ast: program_ast) -> program_ast
    + rewrites column-major 1-indexed array accesses into row-major 0-indexed form
    + collapses multi-dimensional index arithmetic
    # lowering
  fortran_transpile.lower_control_flow
    @ (ast: program_ast) -> result[program_ast, string]
    + converts numbered DO loops to structured loops
    + converts computed and assigned GOTO to switch-like dispatch
    - returns error on irreducible control flow
    # lowering
  fortran_transpile.lower_io
    @ (ast: program_ast) -> program_ast
    + rewrites READ, WRITE, and FORMAT statements as calls to a runtime io package
    # lowering
  fortran_transpile.build_target_ast
    @ (ast: program_ast) -> target_ast
    + returns a target-language AST where each subroutine becomes a function and common blocks become shared struct fields
    # codegen
  fortran_transpile.emit_source
    @ (target: target_ast) -> string
    + pretty-prints the target AST as source text with consistent indentation
    # codegen
  fortran_transpile.compile
    @ (source: string) -> result[string, string]
    + runs the full pipeline from source text to emitted target source
    - returns error with the first failing stage's message
    # pipeline
  fortran_transpile.collect_diagnostics
    @ (ast: program_ast) -> list[diagnostic]
    + returns warnings for constructs with no clean equivalent (EQUIVALENCE, Hollerith, EBCDIC literals)
    # diagnostics
