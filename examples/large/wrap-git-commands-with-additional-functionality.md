# Requirement: "a library that wraps git commands with higher-level operations against a hosted repository provider"

Runs git against a local working copy and calls a remote provider API for hosted features.

std
  std.process
    std.process.run
      @ (program: string, args: list[string], cwd: string) -> result[process_output, string]
      + executes a process and captures stdout, stderr, and exit code
      - returns error when the program cannot be launched
      # process
  std.http
    std.http.request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an HTTP request and returns status, headers, and body
      - returns error on transport failure
      # http
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses any JSON value
      - returns error on malformed input
      # serialization
    std.json.encode_value
      @ (v: json_value) -> string
      + serializes a JSON value
      # serialization

githelper
  githelper.new
    @ (cwd: string, api_base: string, token: string) -> git_session
    + creates a session bound to a working copy and a provider API endpoint
    # construction
  githelper.current_branch
    @ (session: git_session) -> result[string, string]
    + returns the current branch name
    - returns error when not inside a working copy
    # local_git
    -> std.process.run
  githelper.status
    @ (session: git_session) -> result[list[file_status], string]
    + returns the porcelain status of the working copy
    - returns error when git fails
    # local_git
    -> std.process.run
  githelper.create_branch
    @ (session: git_session, name: string) -> result[void, string]
    + creates and checks out a new branch
    - returns error when the branch already exists
    # local_git
    -> std.process.run
  githelper.commit
    @ (session: git_session, message: string) -> result[string, string]
    + stages all tracked changes and commits, returning the new sha
    - returns error when there is nothing to commit
    # local_git
    -> std.process.run
  githelper.push
    @ (session: git_session, remote: string, branch: string) -> result[void, string]
    + pushes a branch to a remote
    - returns error when the push is rejected
    # local_git
    -> std.process.run
  githelper.open_pull_request
    @ (session: git_session, title: string, body: string, head: string, base: string) -> result[string, string]
    + creates a pull request via the provider API and returns its URL
    - returns error on non-success API response
    # remote_api
    -> std.http.request
    -> std.json.encode_value
    -> std.json.parse_value
  githelper.list_pull_requests
    @ (session: git_session, state: string) -> result[list[pull_request], string]
    + lists pull requests filtered by state
    - returns error on API failure
    # remote_api
    -> std.http.request
    -> std.json.parse_value
  githelper.fork
    @ (session: git_session, owner: string, repo: string) -> result[string, string]
    + forks a repository via the provider API and returns the new clone URL
    - returns error when the user lacks permission
    # remote_api
    -> std.http.request
    -> std.json.parse_value
  githelper.clone
    @ (session: git_session, url: string, dest: string) -> result[void, string]
    + clones a repository to a destination directory
    - returns error when the destination already exists
    # local_git
    -> std.process.run
