# Requirement: "parse a curated-links markdown file, enrich each link with repository metadata, and emit an updated file"

Walks the headings and bullet lists of a markdown document, extracts link entries, looks each one up against a metadata source, and renders a new document. std supplies file IO and a pluggable HTTP fetch used by the lookup.

std
  std.fs
    std.fs.read_text
      @ (path: string) -> result[string, string]
      + returns the file contents as UTF-8 text
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_text
      @ (path: string, data: string) -> result[void, string]
      + writes data to path
      - returns error when the path cannot be written
      # filesystem
  std.http
    std.http.get_json
      @ (url: string, headers: map[string, string]) -> result[string, string]
      + returns the response body as a string on HTTP 200
      - returns error on non-200 or transport failure
      # http_client

awesome
  awesome.parse_links_doc
    @ (source: string) -> links_doc
    + returns a tree of sections and link entries preserving heading order
    + attaches each bullet's URL, title, and description text to its section
    # parsing
  awesome.extract_entries
    @ (doc: links_doc) -> list[link_entry]
    + returns every link entry flattened in document order with its section path attached
    # extraction
  awesome.fetch_repo_info
    @ (url: string, fetch_fn: fn(string) -> result[string, string]) -> result[repo_info, string]
    + returns stars, last-updated timestamp, and short description for the repository at url
    - returns error when fetch_fn fails or the response is not recognized metadata
    ? fetch_fn is injected so callers can substitute a fake in tests
    # enrichment
    -> std.http.get_json
  awesome.enrich_entries
    @ (entries: list[link_entry], fetch_fn: fn(string) -> result[string, string]) -> list[enriched_entry]
    + returns an enriched entry for each input, with failures recorded on the entry itself
    # enrichment
  awesome.render_doc
    @ (doc: links_doc, enriched: list[enriched_entry]) -> string
    + returns a markdown document with each entry's line annotated with stars and last-updated
    + preserves original section headings and non-link paragraphs
    # rendering
  awesome.process_file
    @ (input_path: string, output_path: string, fetch_fn: fn(string) -> result[string, string]) -> result[void, string]
    + reads input_path, enriches, and writes the rendered result to output_path
    - returns error when reading or writing fails
    # orchestration
    -> std.fs.read_text
    -> std.fs.write_text
