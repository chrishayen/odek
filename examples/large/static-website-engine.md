# Requirement: "a static website engine"

Loads markdown content with front matter, runs each page through a template, and writes a full site tree.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file fully
      - returns error when unreadable
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data, creating parents as needed
      # filesystem
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns every file path under root
      # filesystem
  std.markdown
    std.markdown.to_html
      fn (source: string) -> string
      + renders commonmark to html
      # markdown
  std.yaml
    std.yaml.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a yaml mapping into a string-to-string map
      - returns error on invalid yaml
      # serialization
  std.template
    std.template.render
      fn (tpl: string, vars: map[string, string]) -> result[string, string]
      + expands {{name}} placeholders
      - returns error on unclosed expressions
      # templating

site
  site.load_content
    fn (dir: string) -> result[list[page], string]
    + walks dir, parses front matter plus body for each .md file
    - returns error on malformed front matter
    # content_loading
    -> std.fs.walk
    -> std.fs.read_all
    -> std.yaml.parse_object
  site.load_templates
    fn (dir: string) -> result[map[string, string], string]
    + reads every template file into a name-keyed map
    # template_loading
    -> std.fs.walk
    -> std.fs.read_all
  site.render_page
    fn (page: page, templates: map[string, string], globals: map[string, string]) -> result[string, string]
    + merges page vars with globals and renders through the page's template
    - returns error when the template is missing
    # rendering
    -> std.markdown.to_html
    -> std.template.render
  site.build_index
    fn (pages: list[page], templates: map[string, string]) -> result[string, string]
    + renders an index listing all pages, sorted by date descending
    # rendering
    -> std.template.render
  site.output_path
    fn (page: page, out_dir: string) -> string
    + maps an input slug to its output path under out_dir
    # layout
  site.write_site
    fn (pages: list[page], templates: map[string, string], out_dir: string) -> result[void, string]
    + renders every page and writes the full site tree
    # emission
    -> std.fs.write_all
