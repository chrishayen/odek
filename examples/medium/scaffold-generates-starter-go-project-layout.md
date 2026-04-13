# Requirement: "a project scaffolding library that materializes a starter directory tree from a template"

Template files live in memory as a list of relative paths and string bodies. The library renders placeholders and writes the tree to disk.

std
  std.fs
    std.fs.make_dir_all
      @ (path: string) -> result[void, string]
      + creates a directory and any missing parents
      - returns error when the path cannot be created
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to a file
      - returns error on write failure
      # filesystem
    std.fs.exists
      @ (path: string) -> bool
      + reports whether a path exists
      # filesystem
  std.path
    std.path.join
      @ (parts: list[string]) -> string
      + joins segments using the platform separator
      # path
    std.path.dir
      @ (path: string) -> string
      + returns the directory portion of a path
      # path

scaffold
  scaffold.template_entry
    @ (relative_path: string, body: string) -> template_entry
    + creates a single template file record
    # construction
  scaffold.render_body
    @ (body: string, vars: map[string,string]) -> string
    + substitutes "{{name}}" placeholders with their values
    ? unknown placeholders are left as-is
    # templating
  scaffold.render_path
    @ (relative_path: string, vars: map[string,string]) -> string
    + substitutes placeholders in the destination path
    # templating
  scaffold.materialize_entry
    @ (out_root: string, entry: template_entry, vars: map[string,string]) -> result[void, string]
    + creates any missing parent directories and writes the rendered file
    - returns error when the destination path escapes out_root
    # write
    -> std.path.join
    -> std.path.dir
    -> std.fs.make_dir_all
    -> std.fs.write_all
  scaffold.materialize
    @ (out_root: string, entries: list[template_entry], vars: map[string,string], overwrite: bool) -> result[i32, string]
    + writes every entry and returns the number of files created
    - returns error when any destination already exists and overwrite is false
    - returns error on write failure
    # pipeline
    -> std.fs.exists
