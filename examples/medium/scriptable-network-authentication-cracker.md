# Requirement: "a scriptable network authentication cracker"

Runs user-supplied authentication scripts against a target with dictionary-driven credential generation.

std
  std.net
    std.net.http_request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an http request and returns the response
      - returns error on connection failure
      # networking
  std.fs
    std.fs.read_lines
      @ (path: string) -> result[list[string], string]
      + reads a file as a list of utf-8 lines
      - returns error when the file cannot be opened
      # io
  std.time
    std.time.sleep_millis
      @ (ms: i64) -> void
      + suspends the current thread for ms milliseconds
      # time

cracker
  cracker.load_wordlist
    @ (path: string) -> result[list[string], string]
    + loads a candidate password list from a file
    # wordlists
    -> std.fs.read_lines
  cracker.load_userlist
    @ (path: string) -> result[list[string], string]
    + loads a username list from a file
    # wordlists
    -> std.fs.read_lines
  cracker.new_attempt_iterator
    @ (users: list[string], passwords: list[string]) -> attempt_iterator
    + produces (user, password) pairs in user-major order
    # candidates
  cracker.compile_script
    @ (source: string) -> result[auth_script, string]
    + parses a scriptable authentication recipe into an executable plan
    - returns error on invalid script syntax
    # scripting
  cracker.run_attempt
    @ (script: auth_script, target: string, user: string, password: string) -> result[attempt_result, string]
    + executes the script against the target and returns success or failure
    # execution
    -> std.net.http_request
  cracker.run_campaign
    @ (script: auth_script, target: string, iter: attempt_iterator, delay_ms: i64) -> list[attempt_result]
    + iterates all candidate pairs with a delay between attempts
    # orchestration
    -> std.time.sleep_millis
