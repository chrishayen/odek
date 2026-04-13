# Requirement: "an image server that stores, resizes, converts, and caches images"

An HTTP-adjacent library exposing request handlers for upload, fetch with transformations, and cache management. Image operations and HTTP bodies go through std seams.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the file does not exist
      # io
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path
      # io
    std.fs.remove
      @ (path: string) -> result[void, string]
      + deletes the file at path
      - returns error when the file does not exist
      # io
  std.image_codec
    std.image_codec.decode
      @ (data: bytes) -> result[image_buffer, string]
      + decodes common raster formats into an rgba image buffer
      - returns error on unknown format signatures
      # codec
    std.image_codec.encode
      @ (img: image_buffer, format: string, quality: i32) -> result[bytes, string]
      + encodes an image buffer to the given format and quality
      - returns error when format is unsupported
      # codec
  std.image_ops
    std.image_ops.resize
      @ (img: image_buffer, width: i32, height: i32) -> image_buffer
      + returns a resampled buffer using bilinear filtering
      # image
    std.image_ops.crop
      @ (img: image_buffer, x: i32, y: i32, width: i32, height: i32) -> result[image_buffer, string]
      + returns a buffer covering the given rectangle
      - returns error when the rectangle is outside the source
      # image
  std.crypto
    std.crypto.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex sha256 of data
      # hash

imgsrv
  imgsrv.new_server
    @ (origin_dir: string, cache_dir: string, default_quality: i32) -> imgsrv_state
    + creates a server state rooted at origin and cache directories
    # construction
  imgsrv.store_original
    @ (srv: imgsrv_state, name: string, data: bytes) -> result[string, string]
    + stores an uploaded image under name and returns a content hash
    - returns error when decode fails
    # storage
    -> std.image_codec.decode
    -> std.crypto.sha256_hex
    -> std.fs.write_all
  imgsrv.parse_transform
    @ (query: string) -> result[imgsrv_transform, string]
    + parses query parameters into a transform description
    - returns error when width or height are negative
    - returns error when format is unknown
    # request
  imgsrv.cache_key
    @ (name: string, tx: imgsrv_transform) -> string
    + returns a stable hashed filename for the transformed variant
    # cache
    -> std.crypto.sha256_hex
  imgsrv.get_variant
    @ (srv: imgsrv_state, name: string, tx: imgsrv_transform) -> result[bytes, string]
    + returns the bytes of the transformed image, computing and caching on miss
    - returns error when the original is not found
    # serving
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.image_codec.decode
    -> std.image_codec.encode
    -> std.image_ops.resize
    -> std.image_ops.crop
  imgsrv.invalidate
    @ (srv: imgsrv_state, name: string) -> result[i32, string]
    + removes every cached variant derived from name; returns the count removed
    # cache
    -> std.fs.remove
  imgsrv.delete_original
    @ (srv: imgsrv_state, name: string) -> result[void, string]
    + removes the original and all cached variants
    - returns error when name is unknown
    # storage
    -> std.fs.remove
  imgsrv.list_originals
    @ (srv: imgsrv_state) -> result[list[string], string]
    + returns the names of every stored original
    # introspection
  imgsrv.stats
    @ (srv: imgsrv_state) -> imgsrv_stats
    + returns counters for cache hits, misses, and bytes served
    # introspection
