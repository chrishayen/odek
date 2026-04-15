# Requirement: "a high-throughput multipart object storage transfer library"

Upload and download large objects against a pluggable object storage backend using parallel part transfers. Network work lives behind std primitives.

std
  std.http
    std.http.put_bytes
      fn (url: string, data: bytes) -> result[string, string]
      + uploads bytes and returns the response body
      - returns error on network failure or non-2xx status
      # http
    std.http.get_bytes
      fn (url: string) -> result[bytes, string]
      + fetches bytes from the URL
      - returns error on network failure or non-2xx status
      # http
  std.crypto
    std.crypto.md5
      fn (data: bytes) -> bytes
      + returns the MD5 digest of the input
      # cryptography

object_transfer
  object_transfer.new
    fn (endpoint: string, bucket: string, part_size: i32, concurrency: i32) -> transfer_state
    + returns a transfer configuration
    # construction
  object_transfer.split_parts
    fn (total_size: i64, part_size: i32) -> list[part_range]
    + returns the ordered byte ranges for each part
    ? the final part absorbs any remainder
    # planning
  object_transfer.upload_part
    fn (state: transfer_state, key: string, part: part_range, data: bytes) -> result[part_etag, string]
    + uploads a single part and returns its etag
    - returns error on transport failure
    # upload
    -> std.crypto.md5
    -> std.http.put_bytes
  object_transfer.complete_upload
    fn (state: transfer_state, key: string, parts: list[part_etag]) -> result[string, string]
    + finalizes a multipart upload and returns the object URL
    - returns error when the storage backend rejects the manifest
    # upload
  object_transfer.download
    fn (state: transfer_state, key: string, total_size: i64) -> result[bytes, string]
    + fetches all parts in parallel and concatenates them in order
    - returns error when any part fails
    # download
    -> std.http.get_bytes
  object_transfer.verify_checksum
    fn (data: bytes, expected: bytes) -> bool
    + returns true when the data's MD5 matches the expected digest
    # integrity
    -> std.crypto.md5
