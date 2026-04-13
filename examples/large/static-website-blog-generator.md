# Requirement: "a static website and blog generator"

Builds pages and blog posts, paginates the post index, and emits an RSS feed.

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
      # filesystem
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns every regular file beneath root
      # filesystem
  std.time
    std.time.parse_iso8601
      @ (raw: string) -> result[i64, string]
      + returns a unix timestamp for an ISO-8601 date
      - returns error on malformed input
      # time
    std.time.format_rfc1123
      @ (epoch_seconds: i64) -> string
      + returns an RFC-1123 date string
      # time

blog_generator
  blog_generator.parse_post
    @ (raw: string) -> result[post, string]
    + extracts title, date, tags, and body from a post source
    - returns error when required metadata is missing
    # post_parse
    -> std.time.parse_iso8601
  blog_generator.render_post
    @ (p: post, template: string) -> string
    + fills title, date, and body into the post template
    # post_render
  blog_generator.render_page
    @ (source: string, template: string) -> string
    + fills body into a plain page template
    # page_render
  blog_generator.sort_posts_by_date
    @ (posts: list[post]) -> list[post]
    + returns posts newest first
    ? posts with equal dates keep input order
    # ordering
  blog_generator.paginate_index
    @ (posts: list[post], per_page: i32) -> list[list[post]]
    + splits sorted posts into fixed-size pages
    - returns one empty page when posts is empty and per_page > 0
    # pagination
  blog_generator.render_index_page
    @ (page: list[post], page_number: i32, total_pages: i32, template: string) -> string
    + fills a listing template with the page's posts and navigation
    # index_render
  blog_generator.render_rss_feed
    @ (site_title: string, site_url: string, posts: list[post]) -> string
    + returns a valid RSS 2.0 document with one item per post
    # rss
    -> std.time.format_rfc1123
  blog_generator.load_tree
    @ (content_root: string) -> result[tuple[list[post], list[tuple[string,string]]], string]
    + returns (posts, pages) read from the content tree
    - returns error when the content root is missing
    # load
    -> std.fs.walk
    -> std.fs.read_all
  blog_generator.build
    @ (content_root: string, out_root: string, templates: map[string,string], per_page: i32) -> result[i32, string]
    + renders posts, pages, index pagination, and feed under out_root
    + returns the total number of files written
    - returns error when a template is missing
    # build
    -> std.fs.write_all
