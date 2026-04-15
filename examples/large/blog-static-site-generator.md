# Requirement: "a static site and blog generator"

Reads content files and templates, renders them, and produces output files.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file as text
      - returns error when missing
      # fs
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes text to path
      - returns error on I/O failure
      # fs
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns filenames in the directory
      - returns error when directory does not exist
      # fs
  std.time
    std.time.parse_iso8601
      fn (value: string) -> result[i64, string]
      + parses an ISO-8601 timestamp into unix seconds
      - returns error on malformed input
      # time

blog
  blog.parse_front_matter
    fn (source: string) -> result[tuple[map[string, string], string], string]
    + splits leading key-value metadata from body text
    - returns error on malformed metadata block
    # parsing
  blog.load_post
    fn (path: string) -> result[post, string]
    + reads a content file and builds a post with metadata and body
    - returns error when required metadata fields are missing
    # content_ingest
    -> std.fs.read_all
    -> std.time.parse_iso8601
  blog.load_all_posts
    fn (content_dir: string) -> result[list[post], string]
    + loads every content file under the directory
    - returns error when directory cannot be listed
    # content_ingest
    -> std.fs.list_dir
  blog.render_markup
    fn (body: string) -> string
    + converts lightweight markup to HTML
    # rendering
  blog.load_template
    fn (path: string) -> result[template, string]
    + reads and compiles a template for later use
    - returns error on syntax errors in the template
    # templating
    -> std.fs.read_all
  blog.apply_template
    fn (tmpl: template, values: map[string, string]) -> result[string, string]
    + substitutes placeholders with values and returns the rendered string
    - returns error when a referenced placeholder is missing
    # templating
  blog.render_post
    fn (tmpl: template, post: post) -> result[string, string]
    + renders a single post using the given template
    - returns error when the template references unknown fields
    # rendering
  blog.render_index
    fn (tmpl: template, posts: list[post]) -> result[string, string]
    + renders a post index sorted by publish date descending
    # rendering
  blog.build_site
    fn (content_dir: string, template_dir: string, out_dir: string) -> result[i32, string]
    + writes every rendered page and returns the total files produced
    - returns error when output directory cannot be written
    # orchestration
    -> std.fs.write_all
