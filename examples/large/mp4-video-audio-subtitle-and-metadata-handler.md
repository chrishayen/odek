# Requirement: "a library for working with mp4 files containing video, audio, subtitles, or metadata"

An ISO base media file format reader/writer. The std layer holds generic binary and I/O primitives; the project layer parses and serializes the box hierarchy and exposes track-level accessors.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads entire file contents
      - returns error when the path cannot be opened
      # io
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path, truncating existing content
      - returns error when the path cannot be created
      # io
  std.binary
    std.binary.read_u32_be
      fn (data: bytes, offset: i64) -> result[u32, string]
      + reads a big-endian unsigned 32-bit integer at the given offset
      - returns error when offset+4 exceeds length
      # binary
    std.binary.read_u64_be
      fn (data: bytes, offset: i64) -> result[u64, string]
      + reads a big-endian unsigned 64-bit integer at the given offset
      - returns error when offset+8 exceeds length
      # binary
    std.binary.write_u32_be
      fn (value: u32) -> bytes
      + encodes a u32 as 4 big-endian bytes
      # binary
    std.binary.slice
      fn (data: bytes, start: i64, length: i64) -> result[bytes, string]
      + returns a sub-slice of the input
      - returns error when start+length exceeds data length
      # binary

mp4
  mp4.parse_box_header
    fn (data: bytes, offset: i64) -> result[box_header, string]
    + reads a 4-byte size and 4-byte type, recognizing the 64-bit largesize form
    - returns error when remaining bytes are insufficient
    # parsing
    -> std.binary.read_u32_be
    -> std.binary.read_u64_be
  mp4.parse_file
    fn (data: bytes) -> result[mp4_file, string]
    + walks the box tree and returns a structured representation
    - returns error when a child box size overflows its parent
    # parsing
    -> std.binary.slice
  mp4.load_file
    fn (path: string) -> result[mp4_file, string]
    + reads a file from disk and parses it
    - returns error when the file is not a valid mp4 container
    # io
    -> std.fs.read_all
  mp4.list_tracks
    fn (file: mp4_file) -> list[track_info]
    + returns one entry per track with id, type (video/audio/subtitle), duration, and language
    + returns an empty list when the file has no moov box
    # tracks
  mp4.read_metadata
    fn (file: mp4_file) -> map[string, string]
    + returns tag keys and values from the udta/meta box
    + returns an empty map when no metadata is present
    # metadata
  mp4.extract_samples
    fn (file: mp4_file, track_id: u32) -> result[list[sample], string]
    + returns the sample table (offset, size, duration, keyframe flag) for the track
    - returns error when the track id is not found
    # samples
  mp4.write_file
    fn (file: mp4_file, path: string) -> result[void, string]
    + serializes the box tree back to disk
    - returns error when a required box is missing
    # serialization
    -> std.binary.write_u32_be
    -> std.fs.write_all
  mp4.add_subtitle_track
    fn (file: mp4_file, language: string, cues: list[subtitle_cue]) -> result[mp4_file, string]
    + returns a new file with an additional subtitle track
    - returns error when the language code is not 3 letters
    # tracks
  mp4.set_metadata
    fn (file: mp4_file, key: string, value: string) -> mp4_file
    + returns a new file with the given metadata tag set
    # metadata
  mp4.duration_seconds
    fn (file: mp4_file) -> f64
    + returns the overall movie duration in seconds computed from the mvhd box
    + returns 0 when mvhd timescale is zero
    # tracks
