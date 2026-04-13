# Requirement: "a library that gathers and formats build and version information with an upgrade-available notice"

Formats version metadata and compares against a latest-release check.

std: (all units exist)

versioninfo
  versioninfo.build
    @ (name: string, semver: string, commit: string, built_at: string) -> versioninfo_state
    + constructs a version record from build metadata
    # construction
  versioninfo.format_short
    @ (info: versioninfo_state) -> string
    + returns "name vX.Y.Z"
    # formatting
  versioninfo.format_long
    @ (info: versioninfo_state) -> string
    + returns multi-line output with name, version, commit, and build timestamp
    # formatting
  versioninfo.format_json
    @ (info: versioninfo_state) -> string
    + returns JSON with the same fields as format_long
    # formatting
  versioninfo.upgrade_notice
    @ (current: string, latest: string) -> optional[string]
    + returns a user-facing message when latest is newer than current
    - returns absent when current is equal to or ahead of latest
    ? semver comparison on dot-separated integer parts
    # upgrade_check
