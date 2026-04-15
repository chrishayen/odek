# Requirement: "a translation module backed by dynamic json dictionaries"

Loads per-locale message dictionaries from json, looks up keys, and fills placeholders. Missing keys are added to the active locale so dictionaries grow over time.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file as text
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes text, creating parent directories
      # filesystem
  std.json
    std.json.parse_string_map
      fn (raw: string) -> result[map[string, string], string]
      + parses a flat string-to-string json object
      - returns error on malformed input
      # serialization
    std.json.encode_string_map
      fn (m: map[string, string]) -> string
      + encodes a flat string map as pretty-printed json
      # serialization

i18n
  i18n.new_catalog
    fn () -> catalog_state
    + creates an empty catalog with no locales loaded
    # construction
  i18n.load_locale
    fn (state: catalog_state, locale: string, path: string) -> result[catalog_state, string]
    + reads the json file and installs it as the dictionary for the locale
    - returns error when the file is missing or malformed
    # loading
    -> std.fs.read_all
    -> std.json.parse_string_map
  i18n.save_locale
    fn (state: catalog_state, locale: string, path: string) -> result[void, string]
    + serializes the dictionary for the locale and writes it to disk
    - returns error when the locale is not loaded
    # persistence
    -> std.json.encode_string_map
    -> std.fs.write_all
  i18n.translate
    fn (state: catalog_state, locale: string, key: string, vars: map[string, string]) -> tuple[string, catalog_state]
    + returns the template rendered with `{{ name }}` placeholders from vars
    + when the key is missing, falls back to the key itself and adds it to the dictionary
    ? dictionaries mutate on miss so tests can detect new strings
    # lookup
