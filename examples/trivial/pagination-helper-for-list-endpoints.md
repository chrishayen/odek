# Requirement: "a pagination helper for list endpoints"

Computes offset, limit, and total-pages metadata given a requested page and page size.

std: (all units exist)

paginate
  paginate.compute
    fn (total_items: i64, page: i32, per_page: i32) -> map[string, i64]
    + returns {offset, limit, page, per_page, total_pages, total_items}
    + clamps page to the range [1, total_pages] and per_page to at least 1
    - returns a zeroed page when total_items is 0 (total_pages = 0, offset = 0)
    # pagination
