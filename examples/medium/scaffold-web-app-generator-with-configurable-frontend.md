# Requirement: "a project scaffolder that generates a web application skeleton with a backend and a configurable frontend framework"

Takes a project name and frontend choice, writes a tree of template files to disk.

std
  std.fs
    std.fs.ensure_dir
      @ (path: string) -> result[void, string]
      + creates the directory and any missing parents
      # filesystem
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + creates or overwrites the file
      - returns error when the parent does not exist
      # filesystem
  std.text
    std.text.render_template
      @ (tmpl: string, vars: map[string,string]) -> string
      + substitutes {{name}} placeholders with the map values
      + leaves unknown placeholders unchanged
      # templating

scaffold
  scaffold.list_frontends
    @ () -> list[string]
    + returns the supported frontend framework names
    # catalog
  scaffold.build_plan
    @ (project: string, frontend: string) -> result[list[file_entry], string]
    + returns every file to create as (relative path, rendered content)
    - returns error when frontend is not in the supported list
    # planning
    -> std.text.render_template
  scaffold.write_plan
    @ (root: string, plan: list[file_entry]) -> result[void, string]
    + ensures directories and writes each file under root
    - returns error on the first file that cannot be written
    # materialize
    -> std.fs.ensure_dir
    -> std.fs.write_all
  scaffold.generate
    @ (project: string, frontend: string, root: string) -> result[void, string]
    + builds the plan and writes it
    # orchestration
