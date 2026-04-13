# Requirement: "a publishing platform library"

A content-management core: posts, authors, tags, drafts, and rendered output. Persistence and rendering go through std primitives.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.uuid
    std.uuid.new_v4
      @ () -> string
      + returns a random uuid in canonical form
      # identifiers
  std.hash
    std.hash.sha1_hex
      @ (data: bytes) -> string
      + returns lowercase hex sha1 digest
      # hashing
  std.text
    std.text.slugify
      @ (title: string) -> string
      + lowercases, trims, replaces non-alphanumerics with single hyphens
      + collapses runs of hyphens
      # text
  std.markdown
    std.markdown.to_html
      @ (source: string) -> string
      + converts markdown source to html
      # rendering

publishing
  publishing.new_store
    @ () -> store_state
    + creates an empty content store
    # construction
  publishing.create_author
    @ (state: store_state, name: string, email: string) -> result[tuple[string, store_state], string]
    + returns a new author id and updated state
    - returns error when email is empty
    # authors
    -> std.uuid.new_v4
  publishing.create_post
    @ (state: store_state, author_id: string, title: string, body_md: string) -> result[tuple[string, store_state], string]
    + returns a new post id; post starts in draft status
    - returns error when author does not exist
    # posts
    -> std.uuid.new_v4
    -> std.time.now_seconds
    -> std.text.slugify
  publishing.update_post
    @ (state: store_state, post_id: string, title: string, body_md: string) -> result[store_state, string]
    + returns updated state with new title and body
    - returns error when post does not exist
    # posts
    -> std.time.now_seconds
  publishing.publish_post
    @ (state: store_state, post_id: string) -> result[store_state, string]
    + transitions post from draft to published with publish timestamp
    - returns error when post is already published
    # publishing
    -> std.time.now_seconds
  publishing.unpublish_post
    @ (state: store_state, post_id: string) -> result[store_state, string]
    + reverts a published post to draft
    - returns error when post was not published
    # publishing
  publishing.tag_post
    @ (state: store_state, post_id: string, tag: string) -> result[store_state, string]
    + adds a tag to the post; idempotent for the same tag
    - returns error when post does not exist
    # tagging
    -> std.text.slugify
  publishing.list_published
    @ (state: store_state) -> list[string]
    + returns post ids ordered by publish time descending
    # queries
  publishing.list_by_tag
    @ (state: store_state, tag: string) -> list[string]
    + returns published post ids that carry the tag
    # queries
    -> std.text.slugify
  publishing.render_post
    @ (state: store_state, post_id: string) -> result[string, string]
    + returns rendered html for the post body
    - returns error when post does not exist
    # rendering
    -> std.markdown.to_html
  publishing.post_etag
    @ (state: store_state, post_id: string) -> result[string, string]
    + returns a content hash suitable for an http etag
    - returns error when post does not exist
    # caching
    -> std.hash.sha1_hex
