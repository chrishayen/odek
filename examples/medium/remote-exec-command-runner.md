# Requirement: "a library for remotely executing commands on a remote host via a management protocol"

Open a session, run a command, and collect stdout, stderr, and exit code.

std
  std.http
    std.http.post
      fn (url: string, headers: map[string,string], body: bytes) -> result[bytes, string]
      + returns response body on 2xx
      - returns error on transport failure or non-2xx
      # http
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + encodes bytes as standard base64
      # encoding
    std.encoding.base64_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes standard base64
      - returns error on invalid input
      # encoding

remote_exec
  remote_exec.open_session
    fn (endpoint: string, username: string, password: string) -> result[session, string]
    + negotiates a session and returns an opaque handle
    - returns error on authentication failure
    # session
    -> std.http.post
  remote_exec.start_command
    fn (session: session, program: string, args: list[string]) -> result[string, string]
    + starts the command and returns a command id
    - returns error when the session is closed
    # execution
    -> std.http.post
  remote_exec.collect_output
    fn (session: session, command_id: string) -> result[command_output, string]
    + returns stdout, stderr, and exit code once the command finishes
    + decodes any base64 payload segments
    - returns error when the command id is unknown
    # execution
    -> std.http.post
    -> std.encoding.base64_decode
  remote_exec.run_once
    fn (endpoint: string, username: string, password: string, program: string, args: list[string]) -> result[command_output, string]
    + convenience that opens a session, runs one command, collects output, and closes
    - returns error from any step
    # convenience
  remote_exec.close_session
    fn (session: session) -> result[void, string]
    + releases the session on the remote host
    # session
    -> std.http.post
