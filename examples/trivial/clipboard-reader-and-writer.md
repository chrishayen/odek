# Requirement: "a library to read from and write to the system clipboard"

Two operations, nothing more.

std: (all units exist)

clipboard
  clipboard.read_text
    fn () -> result[string, string]
    + returns the current clipboard contents as text
    - returns error when the system clipboard is unavailable
    - returns error when the clipboard holds a non-text payload
    # read
  clipboard.write_text
    fn (text: string) -> result[void, string]
    + replaces the clipboard contents with the given text
    - returns error when the system clipboard is unavailable
    # write
