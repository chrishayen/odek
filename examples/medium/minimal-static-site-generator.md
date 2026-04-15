# Requirement: "a minimal static site generator"

Reads markdown sources with optional front matter, renders each through a template, and writes HTML to an output directory.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns file contents
      - returns error when the file is missing
      # filesystem
    std.fs.write_all
      fn (path: string, data: string) -> result[void, string]
      + writes data to path, creating parent directories as needed
      # filesystem
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns entries in a directory
      # filesystem
  std.text
    std.text.markdown_to_html
      fn (source: string) -> string
      + converts a markdown string to HTML
      # markdown

site_gen
  site_gen.parse_front_matter
    fn (source: string) -> tuple[map[string, string], string]
    + splits a document into key-value front matter and body
    + returns empty map when no front matter delimiter is present
    # front_matter
  site_gen.render_template
    fn (template: string, vars: map[string, string]) -> string
    + replaces {{key}} placeholders with values from vars
    + leaves unknown placeholders empty
    # templating
  site_gen.build_page
    fn (source: string, template: string) -> string
    + parses front matter, renders markdown body, injects into template
    # page_build
    -> std.text.markdown_to_html
  site_gen.build_site
    fn (src_dir: string, out_dir: string, template_path: string) -> result[i32, string]
    + builds every .md file under src_dir into out_dir and returns the count written
    - returns error when the template cannot be read
    # site_build
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.fs.list_dir
