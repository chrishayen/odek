# Requirement: "a tv series file-renaming library"

Parses episode info out of messy filenames, looks up canonical titles via an injected metadata source, and produces new filenames.

std: (all units exist)

tv_renamer
  tv_renamer.parse_filename
    @ (filename: string) -> result[episode_ref, string]
    + extracts series name, season, and episode number from common patterns
    - returns error when no season/episode marker is found
    # parsing
  tv_renamer.lookup_title
    @ (source: metadata_source, ref: episode_ref) -> result[string, string]
    + returns the canonical episode title from the injected source
    - returns error when the series or episode is unknown
    # metadata_lookup
  tv_renamer.build_name
    @ (ref: episode_ref, title: string, template: string) -> string
    + substitutes series, season, episode, and title into the template
    ? pads season and episode to two digits
    # formatting
  tv_renamer.rename_batch
    @ (filenames: list[string], source: metadata_source, template: string) -> list[rename_plan_row]
    + returns one plan row per input with old, new, and optional error
    # batch_plan
