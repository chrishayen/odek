# Requirement: "a library for managing a hosts file with named profiles"

Parse, edit, and write a hosts-style file where entries are grouped into named profiles delimited by marker comments.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns file contents
      - returns error when unreadable
      # filesystem
    std.fs.write_all
      fn (path: string, content: string) -> result[void, string]
      + writes content atomically
      - returns error on write failure
      # filesystem

hosts
  hosts.parse
    fn (raw: string) -> hosts_file
    + returns a structured view containing the default section and named profiles
    + preserves unknown lines verbatim
    # parsing
  hosts.render
    fn (file: hosts_file) -> string
    + serializes the structure back to canonical hosts format with profile markers
    # rendering
  hosts.add_entry
    fn (file: hosts_file, profile: string, ip: string, hostnames: list[string]) -> result[hosts_file, string]
    + returns a new file with the entry added to the named profile
    - returns error when the ip is not a valid address
    # edit
  hosts.enable_profile
    fn (file: hosts_file, profile: string) -> result[hosts_file, string]
    + uncomments every entry in the profile
    - returns error when the profile does not exist
    # toggle
  hosts.disable_profile
    fn (file: hosts_file, profile: string) -> result[hosts_file, string]
    + comments out every entry in the profile
    - returns error when the profile does not exist
    # toggle
  hosts.save
    fn (path: string, file: hosts_file) -> result[void, string]
    + renders and writes the file
    - returns error on write failure
    # persistence
    -> std.fs.write_all
