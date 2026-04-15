# Requirement: "a scaffolding tool for generating application project skeletons"

Generates a directory tree of starter files from a named template. File writes go through thin std primitives.

std
  std.fs
    std.fs.write_file
      fn (path: string, data: bytes) -> result[void, string]
      + creates parent directories as needed and writes the file
      - returns error when the path is unwritable
      # filesystem
    std.fs.make_dir
      fn (path: string) -> result[void, string]
      + creates the directory and all missing parents
      # filesystem

scaffolder
  scaffolder.list_templates
    fn () -> list[string]
    + returns the names of all registered project templates
    # discovery
  scaffolder.render
    fn (template: string, project_name: string) -> result[map[string, bytes], string]
    + returns a map of relative file paths to rendered file contents
    - returns error when the template name is unknown
    - returns error when project_name is empty
    # rendering
  scaffolder.write_out
    fn (files: map[string, bytes], target_dir: string) -> result[i32, string]
    + writes all files under target_dir and returns the count written
    - returns error when target_dir already contains files
    # materialization
    -> std.fs.make_dir
    -> std.fs.write_file
