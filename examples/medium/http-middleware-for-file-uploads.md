# Requirement: "an http middleware for handling file uploads"

Parses multipart requests, streams each part to a pluggable storage backend, enforces size and type limits.

std
  std.http
    std.http.read_multipart_part
      fn (req: http_request) -> result[optional[multipart_part], string]
      + returns the next part from the request body, or none at end
      - returns error on malformed multipart framing
      # http
  std.io
    std.io.copy
      fn (src: byte_reader, dst: byte_writer, limit_bytes: i64) -> result[i64, string]
      + copies up to limit_bytes, returning the number of bytes written
      - returns error when limit_bytes is exceeded
      # io
  std.mime
    std.mime.sniff
      fn (head: bytes) -> string
      + returns a best-guess content-type from the first bytes of a stream
      # mime

uploads
  uploads.new_config
    fn (max_file_bytes: i64, allowed_types: list[string]) -> upload_config
    + builds an upload policy
    # construction
  uploads.open_storage
    fn (root_dir: string) -> result[upload_storage, string]
    + opens a filesystem-backed storage rooted at a directory
    - returns error when the directory is not writable
    # storage
  uploads.handle_request
    fn (cfg: upload_config, storage: upload_storage, req: http_request) -> result[list[upload_record], string]
    + drains every file part and returns one record per stored file
    - returns error when any part exceeds max_file_bytes
    - returns error when a part's sniffed type is not in allowed_types
    # middleware
    -> std.http.read_multipart_part
    -> std.mime.sniff
    -> std.io.copy
  uploads.delete
    fn (storage: upload_storage, record: upload_record) -> result[void, string]
    + removes a previously stored upload
    - returns error when the record is not found
    # storage
