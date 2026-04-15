# Requirement: "a WAV audio encoding and decoding library"

Reads and writes PCM WAV files. Handles the RIFF header, format chunk, and data chunk.

std
  std.bytes
    std.bytes.read_u16_le
      fn (data: bytes, offset: i32) -> u16
      + reads a little-endian unsigned 16-bit integer
      # byte_reading
    std.bytes.read_u32_le
      fn (data: bytes, offset: i32) -> u32
      + reads a little-endian unsigned 32-bit integer
      # byte_reading
    std.bytes.read_i16_le
      fn (data: bytes, offset: i32) -> i16
      + reads a little-endian signed 16-bit integer
      # byte_reading
    std.bytes.write_u16_le
      fn (value: u16) -> bytes
      + returns two bytes in little-endian order
      # byte_writing
    std.bytes.write_u32_le
      fn (value: u32) -> bytes
      + returns four bytes in little-endian order
      # byte_writing
    std.bytes.write_i16_le
      fn (value: i16) -> bytes
      + returns two bytes in little-endian order
      # byte_writing
    std.bytes.slice
      fn (data: bytes, start: i32, length: i32) -> bytes
      + returns a sub-range of the input
      - returns empty on out-of-range slices
      # byte_reading

wav
  wav.parse_header
    fn (data: bytes) -> result[wav_format, string]
    + returns channels, sample rate, bits per sample, and data offset
    - returns error when the RIFF magic is missing
    - returns error when the format chunk is malformed
    # parsing
    -> std.bytes.read_u16_le
    -> std.bytes.read_u32_le
  wav.decode
    fn (data: bytes) -> result[wav_audio, string]
    + returns the parsed format and the decoded 16-bit PCM samples
    - returns error on unsupported bit depths
    - returns error when the data chunk size exceeds the file
    # decoding
    -> std.bytes.read_i16_le
    -> std.bytes.slice
  wav.encode
    fn (audio: wav_audio) -> bytes
    + returns a complete WAV file with RIFF header and PCM data
    + produces correct chunk sizes for the sample count
    # encoding
    -> std.bytes.write_u16_le
    -> std.bytes.write_u32_le
    -> std.bytes.write_i16_le
  wav.samples_per_channel
    fn (audio: wav_audio, channel: i32) -> list[i16]
    + returns samples for a single channel from interleaved data
    - returns empty when the channel index is out of range
    # utility
  wav.duration_seconds
    fn (audio: wav_audio) -> f64
    + returns total duration computed from sample count and rate
    # utility
