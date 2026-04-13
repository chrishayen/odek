# Requirement: "a library for reading and writing JSON from multiple sources"

A thin convenience layer over a JSON codec that can target strings, byte buffers, and file-like sources.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a generic value
      - returns error on malformed input
      # serialization
    std.json.encode
      @ (value: json_value) -> string
      + encodes a generic value to a JSON string
      # serialization
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads an entire file as text
      - returns error when the path cannot be opened
      # filesystem
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + writes content to a file, replacing prior contents
      - returns error when the path cannot be written
      # filesystem

ej
  ej.read_string
    @ (raw: string) -> result[json_value, string]
    + parses a JSON document from a string
    - returns error on malformed input
    # input
    -> std.json.parse
  ej.read_file
    @ (path: string) -> result[json_value, string]
    + reads a file and parses its contents as JSON
    - returns error when the file is missing or malformed
    # input
    -> std.fs.read_all
    -> std.json.parse
  ej.write_string
    @ (value: json_value) -> string
    + serializes a value to a JSON string
    # output
    -> std.json.encode
  ej.write_file
    @ (path: string, value: json_value) -> result[void, string]
    + writes a value as JSON to a file
    - returns error when the file cannot be written
    # output
    -> std.json.encode
    -> std.fs.write_all
