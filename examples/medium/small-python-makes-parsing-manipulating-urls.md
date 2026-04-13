# Requirement: "a url parsing and manipulation library"

Parse a URL into components, read and write query parameters, and reassemble to a string.

std: (all units exist)

url_lib
  url_lib.parse
    @ (raw: string) -> result[url_parts, string]
    + parses scheme, host, port, path, query, and fragment
    - returns error when the scheme is missing
    # parsing
  url_lib.format
    @ (parts: url_parts) -> string
    + reassembles the parts into a canonical URL string
    # formatting
  url_lib.set_path
    @ (parts: url_parts, path: string) -> url_parts
    + returns parts with the path replaced
    # mutation
  url_lib.get_query
    @ (parts: url_parts, key: string) -> optional[string]
    + returns the first query value for the key
    - returns none when the key is absent
    # query
  url_lib.set_query
    @ (parts: url_parts, key: string, value: string) -> url_parts
    + returns parts with the query value added or replaced
    # query
  url_lib.remove_query
    @ (parts: url_parts, key: string) -> url_parts
    + returns parts with all values for the key removed
    # query
  url_lib.join
    @ (base: url_parts, reference: string) -> result[url_parts, string]
    + resolves a relative reference against the base
    - returns error when the reference is not a valid URL fragment
    # resolution
