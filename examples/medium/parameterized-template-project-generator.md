# Requirement: "a library that creates projects from parameterized templates"

Walks a template directory, substitutes variables in file contents and names, and writes the result.

std
  std.fs
    std.fs.walk_files
      @ (root: string) -> result[list[string], string]
      + yields every regular file path under root recursively
      - returns error when root does not exist
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the full contents of a file as text
      # filesystem
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + writes content to path, creating parent directories as needed
      # filesystem

templater
  templater.render_string
    @ (text: string, vars: map[string,string]) -> result[string, string]
    + replaces every "{{ name }}" token with vars[name]
    - returns error when a token references an unknown variable
    # substitution
  templater.render_path
    @ (path: string, vars: map[string,string]) -> result[string, string]
    + applies substitution to a file path's segments
    # path_rewriting
    -> templater.render_string
  templater.generate
    @ (template_root: string, output_root: string, vars: map[string,string]) -> result[void, string]
    + walks the template, rewrites paths and contents, and writes them under output_root
    - returns error on the first substitution or I/O failure, reporting the source path
    # generation
    -> std.fs.walk_files
    -> std.fs.read_all
    -> templater.render_string
    -> templater.render_path
    -> std.fs.write_all
