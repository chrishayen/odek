# Requirement: "a library for managing SSH host alias configuration"

Parses SSH config files into host entries and supports add, update, and delete operations. The library returns the updated config text; writing to disk is the caller's responsibility.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads the full file contents as a UTF-8 string
      - returns error when the file does not exist or cannot be read
      # filesystem

ssh_alias
  ssh_alias.parse
    fn (text: string) -> result[list[host_entry], string]
    + parses an SSH config into host entries, each with an alias and option map
    - returns error on malformed Host or unterminated option lines
    # parsing
  ssh_alias.render
    fn (entries: list[host_entry]) -> string
    + serializes host entries back to SSH config text
    + produces stable ordering so unchanged entries round-trip exactly
    # serialization
  ssh_alias.add
    fn (entries: list[host_entry], alias: string, options: map[string, string]) -> result[list[host_entry], string]
    + appends a new host entry
    - returns error when alias already exists
    # mutation
  ssh_alias.update
    fn (entries: list[host_entry], alias: string, options: map[string, string]) -> result[list[host_entry], string]
    + replaces the options of an existing alias
    - returns error when alias is not found
    # mutation
  ssh_alias.remove
    fn (entries: list[host_entry], alias: string) -> result[list[host_entry], string]
    + removes the entry with the given alias
    - returns error when alias is not found
    # mutation
  ssh_alias.load_from_file
    fn (path: string) -> result[list[host_entry], string]
    + convenience that reads a file and parses it
    - returns error when the file cannot be read or parsed
    # loading
    -> std.fs.read_all
