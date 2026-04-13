# Requirement: "a library of string handling, formatting, and validation helpers"

A grab-bag utility library. Kept compact by picking genuinely reusable operations rather than every possible helper.

std: (all units exist)

strutil
  strutil.slugify
    @ (s: string) -> string
    + lowercases, replaces runs of non-alphanumerics with a single hyphen, and trims leading/trailing hyphens
    + unicode letters are folded to ASCII where possible
    # transformation
  strutil.truncate
    @ (s: string, max_chars: i32, ellipsis: string) -> string
    + returns s unchanged when its length is at most max_chars
    + otherwise returns the prefix plus ellipsis, never exceeding max_chars total characters
    # transformation
  strutil.pad_left
    @ (s: string, width: i32, fill: string) -> string
    + left-pads s to width using fill
    + returns s unchanged when it is already at least width characters
    # formatting
  strutil.is_email
    @ (s: string) -> bool
    + returns true when s matches a pragmatic local@domain.tld pattern
    - returns false for empty strings or strings missing an @ or domain
    # validation
  strutil.is_url
    @ (s: string) -> bool
    + returns true when s begins with http:// or https:// and has a host
    - returns false for schemes without a host
    # validation
  strutil.is_numeric
    @ (s: string) -> bool
    + returns true when every character is a decimal digit and s is non-empty
    - returns false for strings containing any non-digit
    # validation
  strutil.format_bytes
    @ (n: i64) -> string
    + renders a byte count using binary units (e.g. "1.5 KiB")
    # formatting
  strutil.levenshtein
    @ (a: string, b: string) -> i32
    + returns the edit distance between two strings
    + returns 0 for identical inputs
    # comparison
