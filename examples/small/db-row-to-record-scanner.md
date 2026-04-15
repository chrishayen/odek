# Requirement: "a library for scanning database rows into typed records"

Scans flat rows and nested relations into record maps with configurable column-to-field naming.

std: (all units exist)

db_scan
  db_scan.column_to_field
    fn (column: string) -> string
    + lowercases the column and converts snake_case tokens to field names
    # naming
  db_scan.scan_row
    fn (columns: list[string], values: list[optional[bytes]]) -> result[map[string, optional[bytes]], string]
    + returns a field map for one row, using column_to_field for each key
    - returns error when columns and values have different lengths
    # scanning
  db_scan.scan_rows
    fn (columns: list[string], rows: list[list[optional[bytes]]]) -> result[list[map[string, optional[bytes]]], string]
    + returns one field map per row
    # scanning
  db_scan.group_parent_child
    fn (parent_rows: list[map[string, optional[bytes]]], parent_key: string, child_rows: list[map[string, optional[bytes]]], child_fk: string) -> list[map[string, optional[bytes]]]
    + attaches matching child rows to each parent under a "children" entry
    ? rows are joined by comparing parent[parent_key] to child[child_fk]
    # relations
