# Requirement: "a programmatic source code generation library"

Builds an abstract representation of source elements (packages, functions, types) and renders them to text. No templates.

std: (all units exist)

codegen
  codegen.new_file
    @ (package_name: string) -> file_node
    + returns an empty file node bound to a package name
    # construction
  codegen.add_import
    @ (file: file_node, path: string) -> file_node
    + records an import; duplicates are ignored
    # imports
  codegen.add_function
    @ (file: file_node, name: string, params: list[param], return_type: string, body: list[statement]) -> file_node
    + appends a function declaration to the file
    # declarations
  codegen.add_struct
    @ (file: file_node, name: string, fields: list[field]) -> file_node
    + appends a struct type declaration with the given fields
    # declarations
  codegen.call
    @ (target: string, args: list[string]) -> statement
    + builds a call expression statement
    # statements
  codegen.assign
    @ (lhs: string, rhs: statement) -> statement
    + builds an assignment statement
    # statements
  codegen.render
    @ (file: file_node) -> string
    + returns the file serialized as source text with imports grouped and declarations ordered
    # rendering
