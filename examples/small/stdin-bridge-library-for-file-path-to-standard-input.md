# Requirement: "a library that adds standard input support to a program which otherwise accepts a file path argument"

Reads from standard input, writes to a temporary file, returns its path for callers that only know how to take a path.

std
  std.io
    std.io.read_stdin_all
      fn () -> result[bytes, string]
      + returns all bytes read from standard input until EOF
      - returns error when the stream cannot be read
      # io
  std.fs
    std.fs.write_temp_file
      fn (data: bytes, suffix: string) -> result[string, string]
      + writes data to a fresh temp file with the given suffix and returns its path
      - returns error when the temp directory is not writable
      # filesystem

stdin_bridge
  stdin_bridge.materialize
    fn (suffix: string) -> result[string, string]
    + reads all of stdin and returns a path to a temp file holding the bytes
    - returns error when stdin cannot be read or the temp file cannot be created
    # bridging
    -> std.io.read_stdin_all
    -> std.fs.write_temp_file
  stdin_bridge.cleanup
    fn (path: string) -> result[void, string]
    + removes the temp file created by materialize
    - returns error when the file has already been removed
    # cleanup
