# Requirement: "a library for detecting the file type of a byte buffer"

Matches a buffer's leading bytes against known magic-number signatures.

std: (all units exist)

filetype
  filetype.detect
    @ (buf: bytes) -> optional[file_type_info]
    + returns file type info when the buffer starts with a known magic signature
    - returns none when no signature matches
    - returns none for buffers shorter than the shortest known signature
    # detection
  filetype.mime_for
    @ (info: file_type_info) -> string
    + returns the MIME type string for a detected file type
    # mime
  filetype.extension_for
    @ (info: file_type_info) -> string
    + returns the canonical file extension for a detected file type
    # extension
