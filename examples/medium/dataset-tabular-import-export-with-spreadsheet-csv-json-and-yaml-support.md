# Requirement: "tabular datasets with import and export in spreadsheet, CSV, JSON, and YAML formats"

A dataset is a headered table of string cells. Readers and writers are per-format; the dataset type is format-agnostic.

std
  std.csv
    std.csv.parse
      @ (raw: string) -> result[list[list[string]], string]
      + parses a CSV document into rows of cells
      - returns error on unterminated quoted field
      # serialization
    std.csv.encode
      @ (rows: list[list[string]]) -> string
      + encodes rows, quoting cells that contain commas or quotes
      # serialization
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document
      # serialization
    std.json.encode
      @ (v: json_value) -> string
      + serializes a json value
      # serialization
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[yaml_value, string]
      + parses a YAML document
      # serialization
    std.yaml.encode
      @ (v: yaml_value) -> string
      + serializes a yaml value
      # serialization

dataset
  dataset.new
    @ (headers: list[string]) -> dataset_state
    + creates an empty dataset with the given column headers
    # construction
  dataset.append_row
    @ (ds: dataset_state, row: list[string]) -> result[dataset_state, string]
    + returns a dataset with the row appended
    - returns error when row length does not match header count
    # mutation
  dataset.from_csv
    @ (raw: string) -> result[dataset_state, string]
    + parses CSV and treats the first row as headers
    - returns error when CSV is malformed or empty
    # import
    -> std.csv.parse
  dataset.to_csv
    @ (ds: dataset_state) -> string
    + returns CSV with headers on the first row
    # export
    -> std.csv.encode
  dataset.from_json
    @ (raw: string) -> result[dataset_state, string]
    + parses a JSON array of objects, inferring headers from keys
    - returns error when the root is not an array of objects
    # import
    -> std.json.parse
  dataset.to_json
    @ (ds: dataset_state) -> string
    + encodes rows as a JSON array of header-keyed objects
    # export
    -> std.json.encode
  dataset.from_yaml
    @ (raw: string) -> result[dataset_state, string]
    + parses a YAML sequence of mappings
    - returns error on type mismatch
    # import
    -> std.yaml.parse
  dataset.to_yaml
    @ (ds: dataset_state) -> string
    + encodes rows as a YAML sequence of mappings
    # export
    -> std.yaml.encode
