# Requirement: "minifiers for html, css, js, xml, json, and svg"

One minifier entry point per format. Format detection is explicit; the caller picks the right one.

std
  std.strings
    std.strings.trim
      fn (s: string) -> string
      + returns s without leading or trailing ascii whitespace
      # strings

minify
  minify.html
    fn (source: string) -> string
    + collapses insignificant whitespace, drops optional tags, and strips comments
    + preserves contents of pre, script, and style blocks
    # html
  minify.css
    fn (source: string) -> string
    + removes comments, collapses whitespace, and shortens zero units
    + preserves contents inside string literals
    # css
  minify.js
    fn (source: string) -> string
    + removes comments and unnecessary whitespace while preserving semantics
    + keeps template literals and regex literals intact
    # js
  minify.xml
    fn (source: string) -> string
    + collapses whitespace between elements but preserves significant whitespace
    # xml
  minify.json
    fn (source: string) -> result[string, string]
    + removes all whitespace outside string literals
    - returns error on invalid json
    # json
  minify.svg
    fn (source: string) -> string
    + applies xml minification plus path data compaction
    # svg
