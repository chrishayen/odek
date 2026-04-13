# Requirement: "a library that detects package dependency references in source files and maps them to documentation URLs"

Given a file's contents and a language hint, find every import-like reference and return ranges annotated with a resolved link. The project layer is a dispatcher and two concrete scanners.

std
  std.strings
    std.strings.contains
      @ (s: string, needle: string) -> bool
      + returns true when needle occurs in s
      # strings
    std.strings.index_of
      @ (s: string, needle: string, from: i32) -> i32
      + returns the byte offset of needle starting at or after from, or -1
      # strings

linkify
  linkify.scan
    @ (source: string, language: string) -> list[dependency_ref]
    + returns ranges for every recognized dependency in the source
    + dispatches to the scanner matching the language hint
    - returns an empty list when the language is unknown
    # dispatch
  linkify.scan_manifest
    @ (source: string) -> list[dependency_ref]
    + extracts dependency entries from a package manifest's dependency sections
    # manifest_scanning
    -> std.strings.index_of
  linkify.scan_imports
    @ (source: string) -> list[dependency_ref]
    + extracts module names from import and require statements in source code
    + handles both quoted and unquoted forms
    # import_scanning
    -> std.strings.contains
  linkify.resolve_url
    @ (ref: dependency_ref, registry: string) -> string
    + returns the documentation URL for a dependency name in the given registry
    # url_resolution
  linkify.annotate
    @ (refs: list[dependency_ref], registry: string) -> list[linked_ref]
    + maps each ref to a linked_ref containing its range and resolved URL
    # annotation
    -> linkify.resolve_url
