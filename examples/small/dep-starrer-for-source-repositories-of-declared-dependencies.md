# Requirement: "a library that stars the source repositories of a project's declared dependencies"

Given a manifest listing dependencies, extract the repository URLs and return the ones that still need to be starred via a pluggable star-function.

std
  std.http
    std.http.get_json
      @ (url: string, auth_token: string) -> result[string, string]
      + performs an authenticated GET and returns the response body
      - returns error on non-2xx status
      # http

dep_starrer
  dep_starrer.extract_repos
    @ (manifest: string) -> list[string]
    + returns the repository URLs of all declared dependencies in the manifest
    + deduplicates repeated entries
    # parsing
  dep_starrer.is_already_starred
    @ (repo: string, auth_token: string) -> result[bool, string]
    + queries the host for the authenticated user's star status on the repo
    - returns error when the host rejects the token
    # query
    -> std.http.get_json
  dep_starrer.star_repo
    @ (repo: string, auth_token: string) -> result[void, string]
    + sends a star request for the repo
    - returns error when the host rejects the request
    # mutation
  dep_starrer.star_all_new
    @ (manifest: string, auth_token: string) -> result[list[string], string]
    + stars every repo in the manifest that is not already starred and returns the list of newly-starred repos
    - returns error on the first host failure
    # orchestration
