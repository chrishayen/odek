# Requirement: "a simple config handling library"

A small library for loading, getting, and setting configuration values from a key-value store backed by a JSON file.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns file contents as a string
      - returns error when file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes string contents to path, creating parent directories
      - returns error when path is not writable
      # filesystem
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

config
  config.load
    @ (path: string) -> result[config_state, string]
    + returns an empty config when the file does not exist
    + returns a populated config when the file exists
    - returns error on malformed JSON
    # loading
    -> std.fs.read_all
    -> std.json.parse_object
  config.get
    @ (state: config_state, key: string) -> optional[string]
    + returns the value when the key exists
    - returns none when the key is missing
    # lookup
  config.set
    @ (state: config_state, key: string, value: string) -> config_state
    + returns a new state with the key set to value
    + overwrites any previous value for the key
    # mutation
  config.save
    @ (state: config_state, path: string) -> result[void, string]
    + writes the current config to disk as JSON
    - returns error when path is not writable
    # persistence
    -> std.json.encode_object
    -> std.fs.write_all
