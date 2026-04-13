# Requirement: "a library that generates coordinated color themes for multiple targets (editors, terminals, wallpapers, and others) from a single palette"

A palette drives a set of pluggable renderers, one per target.

std
  std.fs
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + creates or overwrites the file
      # filesystem
  std.text
    std.text.render_template
      @ (tmpl: string, vars: map[string,string]) -> string
      + substitutes {{name}} placeholders with the map values
      # templating

themer
  themer.normalize_palette
    @ (colors: map[string,string]) -> result[palette, string]
    + validates each value is "#rrggbb", returns a palette with canonical keys (bg, fg, accent, error, ...)
    - returns error when required keys are missing
    # palette
  themer.palette_to_vars
    @ (pal: palette) -> map[string,string]
    + returns a flat map suitable for template substitution
    # adapter
  themer.list_targets
    @ () -> list[string]
    + returns the names of supported targets
    # catalog
  themer.render_target
    @ (target: string, pal: palette) -> result[list[file_entry], string]
    + returns one or more (relative_path, content) pairs for the target
    - returns error when the target is not supported
    # rendering
    -> std.text.render_template
  themer.write_target
    @ (out_dir: string, files: list[file_entry]) -> result[void, string]
    + writes each rendered file under out_dir
    # io
    -> std.fs.write_all
  themer.generate_all
    @ (pal: palette, targets: list[string], out_dir: string) -> result[void, string]
    + renders and writes every requested target
    - returns error at the first failing target
    # orchestration
