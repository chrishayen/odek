# Requirement: "a manager for authorized_keys files across multiple remote hosts"

A library that models the desired set of public keys per host and computes a diff against the current state. Remote transport is injected by the caller.

std: (all units exist)

authkeys
  authkeys.parse_file
    fn (content: string) -> list[auth_entry]
    + returns one entry per non-empty, non-comment line with options, key type, key body and comment
    + tolerates trailing whitespace and mixed line endings
    # parsing
  authkeys.render_file
    fn (entries: list[auth_entry]) -> string
    + returns a newline-separated file preserving entry order
    + ensures a trailing newline
    # rendering
  authkeys.key_fingerprint
    fn (entry: auth_entry) -> string
    + returns the base64 sha256 fingerprint of the key body
    # identification
  authkeys.plan_changes
    fn (current: list[auth_entry], desired: list[auth_entry]) -> change_plan
    + returns added, removed and unchanged entries keyed by fingerprint
    + entries with the same fingerprint but different comments are classified as unchanged
    # diffing
  authkeys.apply_plan
    fn (host: string, plan: change_plan, fetch: fn(string) -> result[string, string], push: fn(string, string) -> result[void, string]) -> result[apply_result, string]
    + reads current file via fetch, applies plan, writes the result via push
    - returns error when fetch fails
    - returns error when push fails
    # orchestration
