# Requirement: "migrate repositories, issues, milestones, and labels between two code hosting services"

Reads resources from a source host and creates matching resources on a target host. Each resource kind has its own fetch-and-create pair; a top-level migrate drives them in dependency order.

std
  std.http
    std.http.get_json
      fn (url: string, headers: map[string,string]) -> result[bytes, string]
      + performs an authenticated GET and returns the raw JSON body
      - returns error on non-2xx status
      # http
    std.http.post_json
      fn (url: string, headers: map[string,string], body: bytes) -> result[bytes, string]
      + performs an authenticated POST with a JSON body
      # http
  std.json
    std.json.parse_array
      fn (raw: bytes) -> result[list[map[string,string]], string]
      + parses a JSON array of objects
      - returns error on malformed JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string,string]) -> bytes
      + encodes an object to JSON bytes
      # serialization

repomigrate
  repomigrate.new_source
    fn (base_url: string, token: string) -> source_client
    + binds to the source host with an auth token
    # construction
  repomigrate.new_target
    fn (base_url: string, token: string) -> target_client
    + binds to the target host with an auth token
    # construction
  repomigrate.list_repositories
    fn (source: source_client, owner: string) -> result[list[map[string,string]], string]
    + returns all repositories owned by the given account
    # fetch
    -> std.http.get_json
    -> std.json.parse_array
  repomigrate.create_repository
    fn (target: target_client, repo: map[string,string]) -> result[void, string]
    + creates a matching repository on the target
    - returns error when a repository with that name already exists
    # create
    -> std.http.post_json
    -> std.json.encode_object
  repomigrate.copy_labels
    fn (source: source_client, target: target_client, repo_name: string) -> result[i32, string]
    + copies all labels from source to target and returns the count
    # labels
    -> std.http.get_json
    -> std.http.post_json
    -> std.json.parse_array
    -> std.json.encode_object
  repomigrate.copy_milestones
    fn (source: source_client, target: target_client, repo_name: string) -> result[i32, string]
    + copies all milestones from source to target and returns the count
    # milestones
    -> std.http.get_json
    -> std.http.post_json
    -> std.json.parse_array
    -> std.json.encode_object
  repomigrate.copy_issues
    fn (source: source_client, target: target_client, repo_name: string) -> result[i32, string]
    + copies all issues, preserving their label and milestone references
    - returns error when a referenced label is missing on the target
    # issues
    -> std.http.get_json
    -> std.http.post_json
    -> std.json.parse_array
    -> std.json.encode_object
  repomigrate.migrate_repository
    fn (source: source_client, target: target_client, repo: map[string,string]) -> result[void, string]
    + creates the repository, then copies labels, milestones, and issues in order
    # entrypoint
    -> repomigrate.create_repository
    -> repomigrate.copy_labels
    -> repomigrate.copy_milestones
    -> repomigrate.copy_issues
