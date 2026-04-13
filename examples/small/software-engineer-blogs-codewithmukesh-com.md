# Requirement: "a minimal blog post model and store"

Input is a person-and-site description, not a library idea. Best-effort interpretation: an in-memory store for blog posts.

std: (all units exist)

blog_store
  blog_store.new
    @ () -> blog_store_state
    + creates an empty store with no posts
    # construction
  blog_store.add_post
    @ (store: blog_store_state, title: string, body: string, author: string) -> string
    + inserts a new post and returns its assigned id
    # write
  blog_store.get_post
    @ (store: blog_store_state, id: string) -> optional[blog_post]
    + returns the post with the given id, or none
    # read
  blog_store.list_by_author
    @ (store: blog_store_state, author: string) -> list[blog_post]
    + returns all posts by the given author, newest first
    # query
  blog_store.delete_post
    @ (store: blog_store_state, id: string) -> bool
    + removes a post; returns true if one was removed
    # write
