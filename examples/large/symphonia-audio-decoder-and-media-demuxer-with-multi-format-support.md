# Requirement: "an audio decoding and media demuxing library supporting multiple container and codec formats"

A demuxer dispatches by container; a decoder dispatches by codec. Both are pluggable registries on top of generic byte-reader primitives.

std
  std.io
    std.io.byte_reader_new
      fn (data: bytes) -> byte_reader
      + wraps a byte slice in a seekable reader
      # io
    std.io.read_u32_be
      fn (reader: byte_reader) -> result[u32, string]
      + reads four bytes as a big-endian unsigned integer
      - returns error at end of stream
      # io
    std.io.read_bytes
      fn (reader: byte_reader, n: i32) -> result[bytes, string]
      + reads exactly n bytes or fails
      # io
    std.io.seek
      fn (reader: byte_reader, offset: i64) -> result[void, string]
      + repositions the read head
      - returns error when offset is out of range
      # io
  std.bits
    std.bits.bitreader_new
      fn (data: bytes) -> bit_reader
      + creates a reader that yields individual bits
      # bits
    std.bits.read_bits
      fn (reader: bit_reader, n: i32) -> result[u32, string]
      + reads up to 32 bits as an unsigned integer
      - returns error when fewer than n bits remain
      # bits

symphonia
  symphonia.probe_format
    fn (data: bytes) -> result[container_kind, string]
    + sniffs the container from header magic bytes
    - returns error when no format matches
    # probing
    -> std.io.byte_reader_new
    -> std.io.read_bytes
  symphonia.demuxer_new
    fn (kind: container_kind, data: bytes) -> result[demuxer_state, string]
    + opens a demuxer for the identified container
    # construction
    -> std.io.byte_reader_new
  symphonia.read_packet
    fn (demuxer: demuxer_state) -> result[optional[packet], string]
    + returns the next elementary stream packet or none at EOF
    - returns error on truncated stream
    # demuxing
    -> std.io.read_u32_be
    -> std.io.read_bytes
    -> std.io.seek
  symphonia.streams
    fn (demuxer: demuxer_state) -> list[stream_info]
    + lists codec and timing metadata for each track
    # metadata
  symphonia.decoder_new
    fn (codec: codec_kind) -> result[decoder_state, string]
    + constructs a decoder for the given codec identifier
    - returns error when the codec is unsupported
    # construction
  symphonia.decode_packet
    fn (decoder: decoder_state, pkt: packet) -> result[audio_frame, string]
    + decodes a single compressed packet into PCM samples
    - returns error on malformed bitstream
    # decoding
    -> std.bits.bitreader_new
    -> std.bits.read_bits
  symphonia.decoder_reset
    fn (decoder: decoder_state) -> decoder_state
    + clears internal state after a seek
    # lifecycle
  symphonia.format_info
    fn (frame: audio_frame) -> audio_format
    + returns sample rate, channel count, and bit depth
    # metadata
