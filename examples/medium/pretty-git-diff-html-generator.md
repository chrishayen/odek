# Requirement: "a library that renders unified diff text as styled HTML"

Parses unified-diff text into structured hunks and emits side-by-side or inline HTML.

std: (all units exist)

diffhtml
  diffhtml.parse_unified
    @ (raw: string) -> result[list[diff_file], string]
    + returns one entry per "diff --git" or "--- / +++" header with its hunks
    - returns error on malformed hunk headers
    # parsing
  diffhtml.classify_line
    @ (line: string) -> line_kind
    + returns context, added, removed, or header based on the leading character
    # parsing
  diffhtml.escape_html
    @ (s: string) -> string
    + replaces &, <, >, ", and ' with their HTML entity equivalents
    # rendering
  diffhtml.render_inline
    @ (files: list[diff_file]) -> string
    + returns a single HTML block showing diffs inline with class names by line kind
    # rendering
  diffhtml.render_side_by_side
    @ (files: list[diff_file]) -> string
    + returns a two-column HTML block with aligned old/new lines per hunk
    ? unchanged lines appear in both columns at the same row
    # rendering
