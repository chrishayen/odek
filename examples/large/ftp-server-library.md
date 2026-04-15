# Requirement: "an FTP server library"

A protocol-level FTP server: command parsing, session state, and handlers for authentication, directory listing, and file transfer. The caller owns the network listener and the filesystem adapter.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[file_entry], string]
      + returns directory entries with name, size, and mtime
      - returns error when the path is not a directory
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path, replacing any existing file
      - returns error when the parent directory does not exist
      # filesystem
    std.fs.delete
      fn (path: string) -> result[void, string]
      + removes the file at the given path
      - returns error when the path does not exist
      # filesystem
  std.time
    std.time.format_rfc3339
      fn (unix_seconds: i64) -> string
      + formats a unix timestamp as an RFC 3339 date string
      # time

ftp
  ftp.parse_command
    fn (line: string) -> result[ftp_command, string]
    + parses an FTP command line into a verb and argument
    - returns error on an empty line
    - returns error when the verb contains non-letter characters
    # protocol
  ftp.format_response
    fn (code: i32, message: string) -> string
    + returns the wire representation "NNN message\r\n"
    + multi-line messages are formatted per RFC 959
    # protocol
  ftp.new_session
    fn () -> session_state
    + creates an unauthenticated session with working directory "/"
    # session
  ftp.handle_user
    fn (session: session_state, username: string) -> tuple[session_state, string]
    + records the username and returns a "331 Password required" response
    # auth
  ftp.handle_pass
    fn (session: session_state, password: string, auth: auth_adapter) -> tuple[session_state, string]
    + returns a "230 Login successful" response when the adapter accepts the password
    - returns a "530 Login incorrect" response when the adapter rejects it
    # auth
  ftp.handle_cwd
    fn (session: session_state, path: string) -> tuple[session_state, string]
    + updates the session working directory and returns a "250" response
    - returns a "550" response when the directory does not exist
    # navigation
  ftp.handle_list
    fn (session: session_state) -> tuple[session_state, string, bytes]
    + returns the data-channel payload for an LIST of the current directory
    - returns a "550" control-channel response when the directory is inaccessible
    # listing
    -> std.fs.list_dir
    -> std.time.format_rfc3339
  ftp.handle_retr
    fn (session: session_state, path: string) -> tuple[session_state, string, bytes]
    + returns file bytes for the data channel and a "150"/"226" control sequence
    - returns a "550" control response when the file is missing
    # transfer
    -> std.fs.read_all
  ftp.handle_stor
    fn (session: session_state, path: string, data: bytes) -> tuple[session_state, string]
    + writes the uploaded bytes and returns "226 Transfer complete"
    - returns "553" when the write fails
    # transfer
    -> std.fs.write_all
  ftp.handle_dele
    fn (session: session_state, path: string) -> tuple[session_state, string]
    + removes the file and returns "250 Requested file action okay"
    - returns "550" when the file cannot be deleted
    # transfer
    -> std.fs.delete
  ftp.dispatch
    fn (session: session_state, command: ftp_command, auth: auth_adapter) -> tuple[session_state, string]
    + routes a parsed command to its handler and returns (new_session, control_response)
    - returns a "500 Unknown command" response for unrecognized verbs
    - returns a "530 Not logged in" response for commands that require auth
    # dispatch
