# Requirement: "a library that treats filesystem paths as their own type instead of using strings"

A small wrapper type with the operations users need most: join, parent, extension, stringify.

std: (all units exist)

pathtype
  pathtype.from_string
    @ (raw: string) -> path
    + returns a normalized path (collapses redundant separators, resolves "." components)
    # construction
  pathtype.to_string
    @ (p: path) -> string
    + returns the canonical string form
    # conversion
  pathtype.join
    @ (base: path, segment: string) -> path
    + returns base with segment appended using the platform separator
    + treats an absolute segment as replacing base
    # composition
  pathtype.parent
    @ (p: path) -> optional[path]
    + returns the parent directory
    + returns none for the root path
    # navigation
  pathtype.extension
    @ (p: path) -> optional[string]
    + returns the extension without the leading dot
    + returns none when the final component has no dot
    # inspection
  pathtype.is_absolute
    @ (p: path) -> bool
    + returns true when the path starts at the root
    # inspection
