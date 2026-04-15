# Requirement: "a static documentation generator that turns markdown files into a site"

Reads a source directory, renders markdown files to HTML, and writes the result to an output directory.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns every file path under root
      - returns error when root does not exist
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file as UTF-8
      - returns error when the file is missing
      # filesystem
    std.fs.write_all
      fn (path: string, data: string) -> result[void, string]
      + writes data to path, creating parent directories as needed
      - returns error when the path cannot be written
      # filesystem
  std.markdown
    std.markdown.render_html
      fn (source: string) -> string
      + renders markdown source to HTML
      # rendering

docgen
  docgen.load_source
    fn (root: string) -> result[list[source_page], string]
    + collects every .md file under root into source_page records
    - returns error when root cannot be walked
    # loading
    -> std.fs.walk
    -> std.fs.read_all
  docgen.render_page
    fn (page: source_page, template: string, nav_html: string) -> string
    + renders a single source page by substituting its HTML body and nav into template
    # rendering
    -> std.markdown.render_html
  docgen.build_nav
    fn (pages: list[source_page]) -> string
    + returns an HTML list of links grouped by directory in deterministic order
    # navigation
  docgen.build
    fn (source_root: string, output_root: string, template: string) -> result[i32, string]
    + renders every page, writes each to output_root mirroring source paths with .html, and returns the count
    - returns error on any read or write failure
    # build
    -> std.fs.write_all
