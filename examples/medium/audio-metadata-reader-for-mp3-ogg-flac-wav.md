# Requirement: "a library for reading music metadata from audio files (MP3, OGG, FLAC, WAV)"

A single entry point dispatches on file format and returns a common metadata record.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the file contents as bytes
      - returns error when the file cannot be read
      # filesystem

audio_meta
  audio_meta.detect_format
    @ (data: bytes) -> result[audio_format, string]
    + returns the format based on magic bytes at the start of the file
    - returns error when no known signature matches
    # detection
  audio_meta.read_id3v2
    @ (data: bytes) -> result[track_metadata, string]
    + reads ID3v2 frames for title, artist, album, year, track number, and duration
    - returns error when the ID3 header is missing or malformed
    # parsing
  audio_meta.read_vorbis_comments
    @ (data: bytes) -> result[track_metadata, string]
    + reads Vorbis comment blocks from an OGG or FLAC stream
    - returns error when the comment block cannot be located
    # parsing
  audio_meta.read_wav_info
    @ (data: bytes) -> result[track_metadata, string]
    + reads the RIFF INFO chunk from a WAV file
    - returns error when the RIFF header is invalid
    # parsing
  audio_meta.read
    @ (path: string) -> result[track_metadata, string]
    + returns metadata for the file at path, dispatching on detected format
    - returns error when the file cannot be read
    - returns error when the format is unsupported
    # reading
    -> std.fs.read_all
