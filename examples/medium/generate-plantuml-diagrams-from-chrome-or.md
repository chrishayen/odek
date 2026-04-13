# Requirement: "a library that generates sequence-diagram markup from a browser network inspection export"

Parses a HAR-style network trace and renders a sequence diagram with one arrow per request/response pair.

std
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses arbitrary JSON
      - returns error on malformed input
      # serialization
  std.url
    std.url.parse
      @ (raw: string) -> result[url_parts, string]
      + returns scheme, host, path, and query
      - returns error on malformed url
      # url

hoofli
  hoofli.load_har
    @ (raw: string) -> result[list[http_exchange], string]
    + returns one exchange per entry with method, url, status, and timings
    - returns error when the root does not contain log.entries
    # ingest
    -> std.json.parse_value
  hoofli.group_by_origin
    @ (exchanges: list[http_exchange]) -> list[origin_group]
    + returns exchanges grouped by scheme+host
    # grouping
    -> std.url.parse
  hoofli.render_sequence
    @ (groups: list[origin_group]) -> string
    + emits plantuml sequence-diagram markup with a participant per origin and one arrow per exchange
    + annotates each arrow with method and status
    # rendering
  hoofli.from_har
    @ (raw: string) -> result[string, string]
    + loads, groups, and renders in one call
    # orchestration
