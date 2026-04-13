# Requirement: "a tool that bakes an on-disk directory tree into a static in-memory virtual filesystem module"

Walks a directory, reads each file, and emits a source text that declares a static map from path to bytes along with accessor functions. The output is a language-neutral pseudo-source rendered as a string.

std
  std.fs
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns all file paths under root in depth-first order
      - returns error when root does not exist
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file's contents as bytes
      - returns error when the file cannot be opened
      # filesystem
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      + encodes bytes as a lowercase hex string
      # encoding

vfsbake
  vfsbake.collect
    @ (root: string) -> result[map[string, bytes], string]
    + walks root and returns a map from repo-relative path to file contents
    - returns error when any file cannot be read
    # collection
    -> std.fs.walk
    -> std.fs.read_all
  vfsbake.escape_path
    @ (path: string) -> string
    + normalizes path separators to forward slashes and trims a leading "./"
    # path_normalization
  vfsbake.emit_entry
    @ (path: string, data: bytes) -> string
    + renders one entry as a path literal bound to a hex-encoded byte literal
    # emission
    -> std.encoding.hex_encode
    -> vfsbake.escape_path
  vfsbake.emit_module
    @ (files: map[string, bytes], module_name: string) -> string
    + renders the full module body: a named static map populated with all entries followed by read and exists accessors
    # emission
    -> vfsbake.emit_entry
  vfsbake.bake
    @ (root: string, module_name: string) -> result[string, string]
    + collects the tree and returns the rendered module source
    - returns error when collection fails
    # pipeline
    -> vfsbake.collect
    -> vfsbake.emit_module
