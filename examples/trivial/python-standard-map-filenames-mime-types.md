# Requirement: "a library that maps filenames to MIME types"

Single lookup function over a filename extension.

std: (all units exist)

mime
  mime.type_for_filename
    @ (filename: string) -> optional[string]
    + returns "text/html" for "index.html"
    + returns "application/json" for "data.json"
    + matching is case-insensitive on the extension
    - returns none for filenames with no known extension
    # mime_lookup
