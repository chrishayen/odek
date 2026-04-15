# Requirement: "a library that repairs broken Unicode text"

Detects and fixes common mojibake (double-decoded text), stray replacement characters, and inconsistent normalization forms.

std
  std.unicode
    std.unicode.normalize_nfc
      fn (input: string) -> string
      + returns the NFC-normalized form of the input
      # unicode
    std.unicode.decode_utf8
      fn (data: bytes) -> result[string, string]
      + decodes UTF-8 bytes to a string
      - returns error on invalid byte sequences
      # unicode
    std.unicode.encode_latin1
      fn (input: string) -> result[bytes, string]
      + encodes a string as Latin-1
      - returns error when any code point exceeds 0xFF
      # unicode

unicode_fixer
  unicode_fixer.fix
    fn (input: string) -> string
    + applies the full repair pipeline: mojibake fix, replacement-char strip, NFC normalization
    + returns the input unchanged when nothing needs fixing
    # pipeline
  unicode_fixer.fix_mojibake
    fn (input: string) -> string
    + detects and undoes text that was UTF-8 bytes mis-decoded as Latin-1
    ? works by round-tripping through Latin-1 and re-decoding as UTF-8 when the result has fewer replacement characters
    # mojibake
    -> std.unicode.encode_latin1
    -> std.unicode.decode_utf8
  unicode_fixer.strip_replacements
    fn (input: string) -> string
    + removes U+FFFD replacement characters
    # sanitization
  unicode_fixer.normalize
    fn (input: string) -> string
    + returns the NFC form
    # normalization
    -> std.unicode.normalize_nfc
  unicode_fixer.diagnose
    fn (input: string) -> list[string]
    + returns a list of detected issues without modifying the input
    + returns an empty list when no issues are detected
    # diagnostics
