# Requirement: "a library for scanning sql rows into typed records using generics"

Maps sql column values onto the fields of a user-supplied record shape.

std: (all units exist)

row_scan
  row_scan.field_map
    fn (target_fields: list[string], columns: list[string]) -> result[list[i32], string]
    + returns the column index for each target field, preserving target order
    - returns error when a target field has no matching column
    # mapping
  row_scan.scan_one
    fn (columns: list[string], row: list[optional[bytes]], field_map: list[i32], target_fields: list[string]) -> result[map[string, optional[bytes]], string]
    + produces a field->value map for a single row
    - returns error when the row has fewer cells than expected columns
    # scanning
  row_scan.scan_all
    fn (columns: list[string], rows: list[list[optional[bytes]]], target_fields: list[string]) -> result[list[map[string, optional[bytes]]], string]
    + scans every row into a field map, stopping at the first error
    # scanning
