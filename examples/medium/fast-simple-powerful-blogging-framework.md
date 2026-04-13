# Requirement: "a static blog generator"

Reads markdown posts with front matter, renders them to HTML, and writes an index page.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path, creating directories as needed
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns file names (not full paths) in the directory
      - returns error when the directory does not exist
      # filesystem
  std.markdown
    std.markdown.render_html
      @ (source: string) -> string
      + renders CommonMark markdown to an HTML fragment
      # markdown

blog
  blog.parse_post
    @ (raw: string) -> result[post, string]
    + extracts front matter (title, date, slug) and body from a markdown file
    - returns error when front matter is missing or malformed
    # parsing
  blog.render_post
    @ (p: post) -> string
    + returns a full HTML page for a single post
    # rendering
    -> std.markdown.render_html
  blog.render_index
    @ (posts: list[post]) -> string
    + returns an HTML index listing posts in reverse chronological order
    # rendering
  blog.build_site
    @ (input_dir: string, output_dir: string) -> result[i32, string]
    + reads every post in input_dir, renders it, writes output, and returns the post count
    - returns error when any post fails to parse
    # pipeline
    -> std.fs.list_dir
    -> std.fs.read_all
    -> std.fs.write_all
