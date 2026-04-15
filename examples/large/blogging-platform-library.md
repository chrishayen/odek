# Requirement: "a blogging platform library"

Posts, tags, comments, and draft/publish workflow with slug generation and persistence via a pluggable store.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.text
    std.text.slugify
      fn (input: string) -> string
      + lowercases, strips diacritics, replaces non-alphanumerics with hyphens, and collapses repeats
      + returns "" for empty input
      # text
  std.crypto
    std.crypto.random_id
      fn (byte_len: i32) -> string
      + returns a hex-encoded cryptographically random identifier
      # randomness

blog
  blog.new
    fn () -> blog_state
    + creates an empty blog with no posts, tags, or comments
    # construction
  blog.create_post
    fn (state: blog_state, author: string, title: string, body: string) -> string
    + creates a draft post and returns its id
    + the slug is derived from the title
    # authoring
    -> std.text.slugify
    -> std.crypto.random_id
    -> std.time.now_seconds
  blog.update_post
    fn (state: blog_state, id: string, title: string, body: string) -> result[void, string]
    + updates title and body; regenerates the slug if the title changed
    - returns error when the post id is unknown
    # authoring
    -> std.text.slugify
  blog.publish
    fn (state: blog_state, id: string) -> result[void, string]
    + marks a draft as published and stamps the publish timestamp
    - returns error when the post is already published
    - returns error when the post id is unknown
    # workflow
    -> std.time.now_seconds
  blog.unpublish
    fn (state: blog_state, id: string) -> result[void, string]
    + moves a published post back to draft
    - returns error when the post id is unknown
    # workflow
  blog.get_by_slug
    fn (state: blog_state, slug: string) -> optional[post]
    + returns the published post with the given slug
    # retrieval
  blog.list_published
    fn (state: blog_state, limit: i32, offset: i32) -> list[post]
    + returns published posts in reverse chronological order
    # retrieval
  blog.add_tag
    fn (state: blog_state, post_id: string, tag: string) -> result[void, string]
    + attaches a tag to a post
    - returns error when the post id is unknown
    # tagging
  blog.list_by_tag
    fn (state: blog_state, tag: string) -> list[post]
    + returns published posts that carry the tag, newest first
    # tagging
  blog.add_comment
    fn (state: blog_state, post_id: string, author: string, body: string) -> result[string, string]
    + appends a comment and returns its id
    - returns error when the post id is unknown or the post is not published
    # comments
    -> std.crypto.random_id
    -> std.time.now_seconds
  blog.delete_comment
    fn (state: blog_state, comment_id: string) -> bool
    + removes a comment; returns true if one was removed
    # comments
  blog.delete_post
    fn (state: blog_state, id: string) -> bool
    + removes a post and all its comments; returns true if one was removed
    # authoring
