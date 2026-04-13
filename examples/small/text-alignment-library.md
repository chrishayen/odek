# Requirement: "a text alignment library"

Takes rows of cells and pads them so columns line up. The caller supplies the alignment per column.

std: (all units exist)

align
  align.column_widths
    @ (rows: list[list[string]]) -> list[i32]
    + returns the max character width of each column across all rows
    ? short rows are treated as if padded with empty strings
    # measurement
  align.pad_cell
    @ (cell: string, width: i32, how: string) -> string
    + pads cell with spaces to reach width
    + "left" pads on the right, "right" pads on the left, "center" splits the padding
    - returns the original cell unchanged when width is smaller than the cell length
    # padding
  align.format_rows
    @ (rows: list[list[string]], alignments: list[string], separator: string) -> list[string]
    + returns each row as a single string with cells padded and joined by separator
    ? alignments shorter than the column count default the remaining columns to "left"
    # formatting
