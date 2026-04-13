# Requirement: "a local game registry and launcher for DRM-free titles"

Keeps a registry of installed games, each pointing at an executable, and launches them through a pluggable process runner.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path
      # filesystem
    std.fs.exists
      @ (path: string) -> bool
      + returns true when the path exists
      # filesystem
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string map as JSON
      # serialization

game_registry
  game_registry.new
    @ () -> registry_state
    + constructs an empty registry
    # construction
  game_registry.add
    @ (state: registry_state, name: string, executable: string) -> result[registry_state, string]
    + returns a new state containing the entry
    - returns error when the executable path does not exist
    - returns error when a game with the same name is already registered
    # registration
    -> std.fs.exists
  game_registry.remove
    @ (state: registry_state, name: string) -> result[registry_state, string]
    + returns a new state without the entry
    - returns error when no such game is registered
    # registration
  game_registry.list
    @ (state: registry_state) -> list[tuple[string, string]]
    + returns every (name, executable) pair sorted by name
    # listing
  game_registry.save
    @ (state: registry_state, path: string) -> result[void, string]
    + writes the registry to disk as JSON
    - returns error when the path is not writable
    # persistence
    -> std.json.encode_object
    -> std.fs.write_all
  game_registry.load
    @ (path: string) -> result[registry_state, string]
    + returns the registry stored at path
    - returns error when the file is missing or malformed
    # persistence
    -> std.fs.read_all
    -> std.json.parse_object
  game_registry.launch
    @ (state: registry_state, name: string, runner: process_runner) -> result[void, string]
    + invokes runner on the registered executable
    - returns error when no such game is registered
    - returns error when the runner reports failure
    # launch
