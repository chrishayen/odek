# Requirement: "a library for fetching paginated data from any source"

A generic paginator that walks pages using caller-supplied fetch logic and collects items across pages.

std: (all units exist)

paginator
  paginator.new
    @ (page_size: i32) -> paginator_state
    + creates a paginator configured with a page size
    ? page_size tells the fetcher how many items to request per page
    # construction
  paginator.fetch_all
    @ (state: paginator_state, fetcher: page_fetcher) -> result[list[bytes], string]
    + invokes the fetcher repeatedly until it reports no next cursor and returns every item collected
    - returns error on the first fetcher failure
    # collection
  paginator.fetch_page
    @ (state: paginator_state, fetcher: page_fetcher, cursor: optional[string]) -> result[page_result, string]
    + fetches a single page at the given cursor and returns items plus the next cursor
    - returns error when the fetcher fails
    # fetching
  paginator.iter_pages
    @ (state: paginator_state, fetcher: page_fetcher) -> page_iterator
    + returns an iterator that yields one page at a time without buffering all items
    # streaming
