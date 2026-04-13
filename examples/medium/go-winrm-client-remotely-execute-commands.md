# Requirement: "a remote command execution client that runs commands on Windows machines over SOAP"

The project layer is a small client; the bulk of the work is SOAP envelope construction and HTTP transport in std.

std
  std.http
    std.http.post
      @ (url: string, headers: map[string, string], body: string) -> result[string, string]
      + sends an HTTP POST and returns the response body
      - returns error on non-2xx status
      # http_client
  std.xml
    std.xml.build_element
      @ (tag: string, attrs: map[string, string], children: list[string]) -> string
      + produces a serialized XML element with escaped text
      # xml_serialization
    std.xml.extract_text
      @ (xml: string, path: string) -> optional[string]
      + returns the text content of the first element matching a slash-separated path
      # xml_parsing

winrm
  winrm.connect
    @ (endpoint: string, username: string, password: string) -> winrm_session
    + constructs a session value holding credentials and endpoint
    # construction
  winrm.open_shell
    @ (session: winrm_session) -> result[string, string]
    + returns a shell id after sending the "Create Shell" SOAP envelope
    - returns error when the server rejects authentication
    # shell_lifecycle
    -> std.xml.build_element
    -> std.http.post
    -> std.xml.extract_text
  winrm.run_command
    @ (session: winrm_session, shell_id: string, command: string, args: list[string]) -> result[command_handle, string]
    + issues a "Command" SOAP envelope and returns the command id
    # command_execution
    -> std.xml.build_element
    -> std.http.post
    -> std.xml.extract_text
  winrm.receive_output
    @ (session: winrm_session, handle: command_handle) -> result[tuple[string, string, i32], string]
    + polls the server and returns (stdout, stderr, exit_code) once the command completes
    # command_output
    -> std.xml.build_element
    -> std.http.post
    -> std.xml.extract_text
  winrm.close_shell
    @ (session: winrm_session, shell_id: string) -> result[void, string]
    + sends a "Delete Shell" SOAP envelope to release server resources
    # shell_lifecycle
    -> std.xml.build_element
    -> std.http.post
