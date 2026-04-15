# Requirement: "an FTP server framework"

Parses FTP protocol commands, tracks per-connection session state, and dispatches file operations through a pluggable backend the caller supplies. The transport loop is not in scope — the library handles command strings and produces reply strings.

std
  std.net
    std.net.parse_host_port
      fn (raw: string) -> result[tuple[string, i32], string]
      + parses a "host:port" string
      - returns error when port is missing or not numeric
      # networking
  std.text
    std.text.split_once
      fn (s: string, sep: string) -> tuple[string, string]
      + splits s on the first occurrence of sep
      + returns (s, "") when sep is not present
      # text
    std.text.trim_crlf
      fn (s: string) -> string
      + removes trailing CR and LF characters
      # text

ftp
  ftp.new_session
    fn () -> session_state
    + creates an unauthenticated session rooted at "/"
    # construction
  ftp.parse_command
    fn (raw: string) -> result[ftp_command, string]
    + parses a command line into a verb and argument
    - returns error when the verb is empty
    # parsing
    -> std.text.split_once
    -> std.text.trim_crlf
  ftp.handle_user
    fn (session: session_state, username: string) -> tuple[string, session_state]
    + records the username and returns "331 need password"
    # auth
  ftp.handle_pass
    fn (session: session_state, password: string, expected: string) -> tuple[string, session_state]
    + returns "230 logged in" when password matches and marks session authenticated
    - returns "530 login incorrect" when password does not match
    # auth
  ftp.handle_cwd
    fn (session: session_state, path: string) -> tuple[string, session_state]
    + updates the working directory and returns "250 ok"
    - returns "530 not logged in" when session is unauthenticated
    # navigation
  ftp.handle_pwd
    fn (session: session_state) -> string
    + returns "257" reply containing the current working directory
    # navigation
  ftp.handle_list
    fn (session: session_state, entries: list[string]) -> string
    + formats directory entries into an FTP list reply body
    ? the caller provides entries via the backend
    # listing
  ftp.handle_retr
    fn (session: session_state, path: string, data: bytes) -> string
    + returns "150 opening data connection" followed by a completion reply
    - returns "530 not logged in" when session is unauthenticated
    # transfer
  ftp.handle_stor
    fn (session: session_state, path: string) -> tuple[string, session_state]
    + returns "150 ok to send data" and marks the session as expecting an upload on path
    # transfer
  ftp.handle_quit
    fn (session: session_state) -> string
    + returns "221 goodbye"
    # lifecycle
  ftp.dispatch
    fn (session: session_state, raw: string, backend: backend_state) -> tuple[string, session_state]
    + parses a command and routes to the appropriate handler
    - returns "500 unknown command" when the verb is not recognized
    # dispatch
    -> std.text.trim_crlf
  ftp.format_pasv_reply
    fn (host: string, port: i32) -> string
    + formats the "227 entering passive mode" reply with the host and port
    # networking
    -> std.net.parse_host_port
