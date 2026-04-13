# Requirement: "a file picker with thumbnails and search"

Browses a directory tree, produces thumbnail images for known file types, and filters entries by a search query.

std
  std.fs
    std.fs.list_dir
      @ (dir: string) -> result[list[dir_entry], string]
      + returns entries with name, size, and kind
      - returns error when the directory cannot be read
      # io
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file's raw bytes
      # io
  std.image
    std.image.decode
      @ (data: bytes) -> result[image, string]
      + decodes a raster image from bytes
      - returns error on unknown format
      # imaging
    std.image.resize
      @ (img: image, width: i32, height: i32) -> image
      + resizes an image to the given dimensions
      # imaging
    std.image.encode_png
      @ (img: image) -> bytes
      + encodes an image as PNG
      # imaging

pikeru
  pikeru.open
    @ (root: string) -> result[browser_state, string]
    + creates a browser rooted at the given directory
    # navigation
    -> std.fs.list_dir
  pikeru.navigate
    @ (state: browser_state, child: string) -> result[browser_state, string]
    + descends into a subdirectory of the current view
    - returns error when child is not a directory
    # navigation
    -> std.fs.list_dir
  pikeru.parent
    @ (state: browser_state) -> result[browser_state, string]
    + ascends to the parent directory
    - returns error when already at the root
    # navigation
    -> std.fs.list_dir
  pikeru.search
    @ (state: browser_state, query: string) -> list[dir_entry]
    + returns entries in the current view whose name contains query, case-insensitive
    # filtering
  pikeru.generate_thumbnail
    @ (path: string, size: i32) -> result[bytes, string]
    + returns a PNG-encoded thumbnail for an image file
    - returns error when the file is not a decodable image
    # thumbnails
    -> std.fs.read_all
    -> std.image.decode
    -> std.image.resize
    -> std.image.encode_png
  pikeru.select
    @ (state: browser_state, name: string) -> result[string, string]
    + returns the absolute path of the selected entry
    - returns error when no entry with that name exists
    # selection
