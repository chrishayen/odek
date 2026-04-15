# Requirement: "a terminal hint overlay that highlights text matches and returns the one selected by a typed label"

Scans a screen buffer for matching tokens, assigns short labels, and resolves a typed prefix to the selected match. std provides only regex compilation and matching.

std
  std.regex
    std.regex.compile
      fn (pattern: string) -> result[regex_handle, string]
      + compiles a pattern once for reuse
      - returns error on invalid regex syntax
      # parsing
    std.regex.find_all
      fn (re: regex_handle, text: string) -> list[regex_match]
      + returns all non-overlapping matches with byte offsets and captured text
      # parsing

hinter
  hinter.scan
    fn (screen: string, patterns: list[string]) -> result[list[hint], string]
    + returns hints for every match of any pattern, in reading order
    - returns error when any pattern fails to compile
    # scanning
    -> std.regex.compile
    -> std.regex.find_all
  hinter.assign_labels
    fn (hints: list[hint], alphabet: string) -> list[labeled_hint]
    + assigns short unique labels drawn from alphabet using a minimal-length scheme
    ? earlier hints get shorter labels; labels are prefix-free
    # labeling
  hinter.resolve
    fn (labeled: list[labeled_hint], typed: string) -> hint_resolution
    + returns "unique" with the hint when typed matches exactly one label
    + returns "ambiguous" when typed is a prefix of multiple labels
    - returns "none" when typed matches no label
    # selection
