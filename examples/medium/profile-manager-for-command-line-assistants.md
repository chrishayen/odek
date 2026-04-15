# Requirement: "a profile manager library for switching between named configuration profiles for command-line assistants"

Stores profiles on disk and applies them by rewriting a target config file. No GUI, no product names.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating parent directories as needed
      - returns error when the destination is not writable
      # filesystem
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

profile_manager
  profile_manager.load_store
    fn (path: string) -> result[profile_store, string]
    + returns the decoded store from disk
    - returns error when the file cannot be read or parsed
    # persistence
    -> std.fs.read_all
    -> std.json.parse_object
  profile_manager.save_store
    fn (store: profile_store, path: string) -> result[void, string]
    + writes the store to disk atomically
    - returns error when the destination cannot be written
    # persistence
    -> std.json.encode_object
    -> std.fs.write_all
  profile_manager.upsert_profile
    fn (store: profile_store, name: string, config: map[string, string]) -> profile_store
    + inserts or replaces the named profile
    # editing
  profile_manager.remove_profile
    fn (store: profile_store, name: string) -> result[profile_store, string]
    + removes the named profile
    - returns error when the profile does not exist
    # editing
  profile_manager.activate
    fn (store: profile_store, name: string, target_path: string) -> result[profile_store, string]
    + writes the profile's config to the target path and marks it active
    - returns error when the profile does not exist
    # switching
    -> std.json.encode_object
    -> std.fs.write_all
  profile_manager.active_profile
    fn (store: profile_store) -> optional[string]
    + returns the currently active profile name, if any
    # introspection
  profile_manager.list_profiles
    fn (store: profile_store) -> list[string]
    + returns all profile names in alphabetical order
    # introspection
