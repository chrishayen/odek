# Requirement: "a library that compresses and bundles linked or inline script and style assets into a single cached file"

Collects asset references from a markup document, concatenates their content, hashes the result, and returns a cache-addressed bundle.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the complete file contents
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path
      # filesystem
  std.crypto
    std.crypto.sha1_hex
      @ (data: bytes) -> string
      + returns the lowercase hex sha-1 digest
      # cryptography

asset_bundler
  asset_bundler.extract_script_srcs
    @ (markup: string) -> list[string]
    + returns the src attribute of every external script tag in order
    # parsing
  asset_bundler.extract_style_hrefs
    @ (markup: string) -> list[string]
    + returns the href of every linked stylesheet in order
    # parsing
  asset_bundler.extract_inline_scripts
    @ (markup: string) -> list[string]
    + returns the text content of every inline script block
    # parsing
  asset_bundler.extract_inline_styles
    @ (markup: string) -> list[string]
    + returns the text content of every inline style block
    # parsing
  asset_bundler.minify_js
    @ (source: string) -> string
    + strips comments and redundant whitespace from a script source
    # minification
  asset_bundler.minify_css
    @ (source: string) -> string
    + strips comments and redundant whitespace from a stylesheet source
    # minification
  asset_bundler.bundle
    @ (sources: list[string], base_dir: string) -> result[tuple[string, bytes], string]
    + returns (cache_name, bundle_bytes) where cache_name is a content-addressed filename
    - returns error when any referenced source cannot be read
    # bundling
    -> std.fs.read_all
    -> std.crypto.sha1_hex
  asset_bundler.rewrite_markup
    @ (markup: string, bundle_url: string, kind: string) -> string
    + replaces all referenced scripts or styles with a single reference to the bundle
    ? kind selects whether script or style references are rewritten
    # rewrite
  asset_bundler.persist
    @ (cache_dir: string, name: string, data: bytes) -> result[void, string]
    + writes the bundle to the cache directory under name
    # caching
    -> std.fs.write_all
