# Requirement: "a static site generator that supports markdown and restructured text"

Walks a content tree, renders each source file through the matching parser, applies a template, and writes HTML to an output tree.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating parent directories
      - returns error when the parent directory cannot be created
      # filesystem
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns every regular file beneath root
      - returns error when root does not exist
      # filesystem
    std.fs.ext
      @ (path: string) -> string
      + returns the lowercased extension without the dot, or ""
      # filesystem
  std.text
    std.text.split_lines
      @ (input: string) -> list[string]
      + splits on "\n" and keeps no trailing empty element
      # text

site_generator
  site_generator.parse_front_matter
    @ (raw: string) -> tuple[map[string,string], string]
    + returns (metadata, body) when the document begins with "---"
    + returns (empty, raw) when no front matter is present
    # metadata
    -> std.text.split_lines
  site_generator.render_markdown
    @ (body: string) -> string
    + converts headings, paragraphs, links, and code spans to HTML
    ? supports a commonmark-compatible subset; not every edge case
    # markdown
  site_generator.render_restructured
    @ (body: string) -> string
    + converts sections, paragraphs, and inline markup to HTML
    ? supports the core directives; not the full spec
    # rest
  site_generator.select_renderer
    @ (path: string) -> optional[string]
    + returns "markdown" for .md and .markdown
    + returns "rest" for .rst
    - returns none for unsupported extensions
    # dispatch
    -> std.fs.ext
  site_generator.apply_template
    @ (template: string, meta: map[string,string], body_html: string) -> string
    + substitutes {{title}}, {{body}}, and arbitrary meta keys
    - leaves unreferenced placeholders untouched
    # templating
  site_generator.output_path_for
    @ (source_root: string, source_path: string, out_root: string) -> string
    + mirrors the relative path under out_root with a .html extension
    # paths
  site_generator.render_one
    @ (source_path: string, template: string) -> result[tuple[map[string,string], string], string]
    + returns (metadata, final_html) for a supported source file
    - returns error for an unsupported extension
    # pipeline
    -> std.fs.read_all
  site_generator.build
    @ (source_root: string, out_root: string, template: string) -> result[i32, string]
    + walks source_root, renders every supported file, and writes output
    + returns the number of pages written
    - returns error when the source root is missing
    # build
    -> std.fs.walk
    -> std.fs.write_all
