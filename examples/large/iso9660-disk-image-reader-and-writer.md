# Requirement: "reading and creating ISO9660 disk images"

A full filesystem image reader and writer. The on-disk format has multiple descriptor types, path tables, and directory records, so std carries the byte-level primitives and the project layer exposes a file-tree API.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of the file at path
      - returns error when the file cannot be opened
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path, truncating any existing file
      - returns error when the path cannot be written
      # filesystem
  std.bytes
    std.bytes.read_u16_le
      @ (data: bytes, offset: i64) -> result[u16, string]
      + reads a little-endian u16 at the given offset
      - returns error when offset+2 exceeds data length
      # binary_io
    std.bytes.read_u32_le
      @ (data: bytes, offset: i64) -> result[u32, string]
      + reads a little-endian u32 at the given offset
      - returns error when offset+4 exceeds data length
      # binary_io
    std.bytes.write_u16_le
      @ (buf: bytes, offset: i64, value: u16) -> bytes
      + returns buf with a little-endian u16 written at offset
      # binary_io
    std.bytes.write_u32_le
      @ (buf: bytes, offset: i64, value: u32) -> bytes
      + returns buf with a little-endian u32 written at offset
      # binary_io
    std.bytes.slice
      @ (data: bytes, start: i64, end: i64) -> result[bytes, string]
      + returns the subslice data[start:end]
      - returns error when the range is out of bounds
      # binary_io
  std.encoding
    std.encoding.ascii_trim
      @ (data: bytes) -> string
      + returns the ASCII string with trailing spaces and nulls removed
      # encoding

iso9660
  iso9660.read_image
    @ (path: string) -> result[iso_image, string]
    + returns a parsed image when the file is a well-formed ISO9660 volume
    - returns error when the primary volume descriptor signature is missing
    # image_loading
    -> std.fs.read_all
  iso9660.parse_volume_descriptor
    @ (sector: bytes) -> result[iso_volume_descriptor, string]
    + parses a 2048-byte sector into a typed descriptor
    - returns error when the standard identifier is not "CD001"
    # parsing
    -> std.bytes.read_u16_le
    -> std.bytes.read_u32_le
    -> std.encoding.ascii_trim
  iso9660.parse_path_table
    @ (data: bytes) -> result[list[iso_path_entry], string]
    + returns path table entries in directory order
    - returns error when an entry's length field exceeds remaining bytes
    # parsing
    -> std.bytes.read_u16_le
    -> std.bytes.read_u32_le
  iso9660.parse_directory_record
    @ (data: bytes, offset: i64) -> result[iso_dir_record, string]
    + returns a parsed directory record and its length in bytes
    - returns error when the record length is zero or exceeds remaining bytes
    # parsing
    -> std.bytes.read_u32_le
  iso9660.list_directory
    @ (img: iso_image, dir_path: string) -> result[list[iso_dir_entry], string]
    + returns the immediate children of the given directory
    - returns error when the path does not resolve to a directory
    # navigation
  iso9660.read_file
    @ (img: iso_image, file_path: string) -> result[bytes, string]
    + returns the file contents at the given path
    - returns error when the path is missing or is a directory
    # file_access
    -> std.bytes.slice
  iso9660.new_builder
    @ (volume_id: string) -> iso_builder
    + returns an empty builder with the given volume identifier
    # construction
  iso9660.add_file
    @ (builder: iso_builder, path: string, content: bytes) -> result[iso_builder, string]
    + returns a builder with the file added under the given absolute path
    - returns error when the path is not absolute or contains invalid characters
    # authoring
  iso9660.add_directory
    @ (builder: iso_builder, path: string) -> result[iso_builder, string]
    + returns a builder with the directory added under the given absolute path
    - returns error when a parent directory is missing
    # authoring
  iso9660.build_image
    @ (builder: iso_builder) -> result[bytes, string]
    + returns the serialized ISO9660 image bytes
    - returns error when the builder contains no files or no root
    # serialization
    -> std.bytes.write_u16_le
    -> std.bytes.write_u32_le
  iso9660.write_image
    @ (builder: iso_builder, path: string) -> result[void, string]
    + writes the built image to the given path
    - returns error when the image cannot be built or written
    # serialization
    -> std.fs.write_all
