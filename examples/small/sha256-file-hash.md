# Requirement: "compute the sha256 of a file"

One project function that wires together two std primitives. Both primitives are generic and reused by any file-hashing or file-checksumming use case.

std
  std.hash
    std.hash.sha256
      fn (data: bytes) -> bytes
      + returns the 32-byte sha256 digest of the input
      + the known digest for empty input is returned for empty input
      # hashing
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file at path into bytes
      - returns error when file does not exist
      - returns error when file is not readable
      # filesystem

file_hash
  file_hash.sha256_of_file
    fn (path: string) -> result[bytes, string]
    + returns the sha256 digest of the file at path
    - returns error when the file cannot be read
    # hashing
    -> std.fs.read_all
    -> std.hash.sha256
