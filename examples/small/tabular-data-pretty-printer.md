# Requirement: "a tabular data pretty-printer"

Formats a header and rows into an aligned text table with a selectable border style.

std: (all units exist)

tabulate
  tabulate.render
    fn (headers: list[string], rows: list[list[string]], style: i32) -> result[string, string]
    + returns a string containing an aligned table with the given style
    + column widths are the max of the header and the widest cell
    - returns error when a row has a different number of cells than headers
    ? style 0 is plain, style 1 is ascii-bordered, style 2 is grid
    # rendering
  tabulate.column_widths
    fn (headers: list[string], rows: list[list[string]]) -> list[i32]
    + returns the per-column maximum width across headers and rows
    # layout
