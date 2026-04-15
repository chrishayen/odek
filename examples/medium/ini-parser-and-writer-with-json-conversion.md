# Requirement: "an ini parser and writer with json conversion"

Parse ini text into sections, edit in memory, serialize back to ini or to json.

std
  std.json
    std.json.encode_value
      fn (value: json_value) -> string
      + encodes a json value as a string
      # serialization
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the file contents
      - returns error when missing
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the path
      # filesystem

ini
  ini.parse
    fn (raw: string) -> result[ini_doc, string]
    + parses ini text into sections with key-value pairs
    + treats lines starting with # or ; as comments
    - returns error on a value line outside any section
    # parsing
  ini.load_file
    fn (path: string) -> result[ini_doc, string]
    + reads and parses an ini file
    # loading
    -> std.fs.read_all
  ini.get
    fn (doc: ini_doc, section: string, key: string) -> optional[string]
    + returns the value for the section and key
    - returns none when absent
    # query
  ini.set
    fn (doc: ini_doc, section: string, key: string, value: string) -> ini_doc
    + sets or adds a key under the section
    # write
  ini.delete
    fn (doc: ini_doc, section: string, key: string) -> ini_doc
    + removes the key if present
    # write
  ini.render
    fn (doc: ini_doc) -> string
    + serializes the document back to ini text
    + preserves section and key insertion order
    # rendering
  ini.save_file
    fn (doc: ini_doc, path: string) -> result[void, string]
    + writes the document to the given path
    # persistence
    -> std.fs.write_all
  ini.to_json
    fn (doc: ini_doc) -> string
    + encodes the document as a nested json object of section-to-kv
    # conversion
    -> std.json.encode_value
