# Requirement: "a project template renderer"

Given a template directory tree with placeholder expressions and a variable map, render the tree to an output directory. Filesystem and placeholder substitution are separated.

std
  std.fs
    std.fs.list_tree
      @ (root: string) -> result[list[string], string]
      + returns all file paths under root, recursively
      - returns error when root does not exist
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when path is unreadable
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating parent directories as needed
      - returns error when the target is not writable
      # filesystem

template_renderer
  template_renderer.substitute
    @ (source: string, vars: map[string, string]) -> result[string, string]
    + replaces every "{{ name }}" occurrence with vars[name]
    - returns error when a referenced variable is missing from vars
    # substitution
  template_renderer.render_tree
    @ (template_dir: string, output_dir: string, vars: map[string, string]) -> result[i32, string]
    + copies every file from template_dir to output_dir with file contents and path segments substituted
    + returns the number of files written
    - returns error when any file cannot be read or written
    # rendering
    -> std.fs.list_tree
    -> std.fs.read_all
    -> std.fs.write_all
