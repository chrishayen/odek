# Requirement: "a recommendation engine based on item-to-item collaborative filtering"

Ingests user/item ratings and recommends items using cosine similarity between item vectors. Project layer owns the rating store and recommendation query; std provides math.

std
  std.math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the square root of x
      ? behavior on negative x is caller's responsibility
      # math

collab_filter
  collab_filter.new_store
    @ () -> store_state
    + creates an empty rating store
    # construction
  collab_filter.record_rating
    @ (store: store_state, user: string, item: string, rating: f64) -> store_state
    + stores (or overwrites) the user's rating for the item
    + rating is normalized into the range [-1,1] on insert
    # data_ingestion
  collab_filter.cosine_similarity
    @ (store: store_state, item_a: string, item_b: string) -> f64
    + returns cosine similarity between two item rating vectors
    + returns 0 when the items share no raters
    # similarity
    -> std.math.sqrt
  collab_filter.similar_items
    @ (store: store_state, item: string, top_n: i32) -> list[tuple[string, f64]]
    + returns the top_n items ranked by similarity (descending)
    + excludes the query item itself
    # recommendation
  collab_filter.recommend_for_user
    @ (store: store_state, user: string, top_n: i32) -> list[tuple[string, f64]]
    + returns items the user has not rated, scored by weighted similarity to their rated items
    - returns empty list when the user has no ratings
    # recommendation
  collab_filter.remove_rating
    @ (store: store_state, user: string, item: string) -> store_state
    + removes the user's rating for the item if present
    # data_ingestion
