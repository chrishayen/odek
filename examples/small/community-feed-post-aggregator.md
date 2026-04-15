# Requirement: "a developer community post aggregator"

Collect posts from community members and expose a combined feed ordered by recency. The source fetcher is a pluggable callable; this library owns merging and ordering.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

community_feed
  community_feed.new
    fn () -> feed_state
    + creates an empty feed with no registered members
    # construction
  community_feed.add_member
    fn (state: feed_state, handle: string) -> feed_state
    + registers a member handle whose posts should be included
    # membership
  community_feed.ingest
    fn (state: feed_state, handle: string, posts: list[post]) -> result[feed_state, string]
    + appends the given posts for handle, deduplicating by post id
    - returns error when handle is not a registered member
    # ingestion
    -> std.time.now_seconds
  community_feed.latest
    fn (state: feed_state, limit: i32) -> list[post]
    + returns up to limit posts across all members, newest first by published_at
    # query
