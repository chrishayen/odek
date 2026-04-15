# Requirement: "a simple text template library with placeholders and source-control metadata"

A tiny template engine that substitutes named placeholders, with a helper that augments the value map with metadata read from a source-control repository.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the file contents as a string
      - returns error when the file cannot be read
      # filesystem
  std.process
    std.process.run
      fn (command: string, args: list[string]) -> result[string, string]
      + runs a command and returns its standard output
      - returns error when the command exits non-zero
      # process

tmpl
  tmpl.render
    fn (template: string, values: map[string, string]) -> result[string, string]
    + substitutes {{name}} placeholders with matching values
    + leaves whitespace around the name optional inside the braces
    - returns error when a placeholder has no matching key
    # rendering
  tmpl.repo_metadata
    fn (repo_path: string) -> result[map[string, string], string]
    + returns a map with commit hash, branch, and tag for the repository at repo_path
    - returns error when the path is not a repository
    # metadata
    -> std.process.run
  tmpl.render_file
    fn (template_path: string, values: map[string, string]) -> result[string, string]
    + reads the template file and renders it with values
    - returns error when the file cannot be read
    - returns error when rendering fails
    # rendering
    -> std.fs.read_all
