# Requirement: "an INI file parser"

Parses and serializes INI documents as section-keyed maps of key-value strings.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads a file's contents as a string
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes a string to the given path
      - returns error when the parent directory does not exist
      # filesystem

configparser
  configparser.parse
    @ (raw: string) -> result[map[string, map[string, string]], string]
    + returns sections mapped to key-value pairs
    + treats keys before any section header as belonging to a default section
    - returns error on malformed key line
    ? lines beginning with ';' or '#' are treated as comments
    # parsing
  configparser.serialize
    @ (sections: map[string, map[string, string]]) -> string
    + returns the INI representation with section headers and key=value lines
    + emits sections in insertion order
    # serialization
  configparser.load_file
    @ (path: string) -> result[map[string, map[string, string]], string]
    + reads an INI file and returns its sections
    - returns error when the file is missing or malformed
    # io
    -> std.fs.read_all
  configparser.save_file
    @ (path: string, sections: map[string, map[string, string]]) -> result[void, string]
    + writes sections to the given path in INI format
    - returns error when the file cannot be written
    # io
    -> std.fs.write_all
