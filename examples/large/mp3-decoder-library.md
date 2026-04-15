# Requirement: "an MP3 decoder library"

Decodes MPEG-1/2 Layer III audio frames into PCM samples. Decomposes the pipeline into its real subsystems: bitstream, Huffman, requantization, inverse transforms, synthesis.

std
  std.bits
    std.bits.reader_new
      fn (data: bytes) -> bit_reader
      + returns a bit reader positioned at the start of the byte buffer
      # bitstream
    std.bits.read_bits
      fn (reader: bit_reader, count: i32) -> tuple[u32, bit_reader]
      + reads up to 32 bits big-endian and returns (value, advanced_reader)
      # bitstream
    std.bits.byte_align
      fn (reader: bit_reader) -> bit_reader
      + advances the reader to the next byte boundary
      # bitstream

mp3_decoder
  mp3_decoder.find_sync
    fn (data: bytes, start: i64) -> optional[i64]
    + returns the byte offset of the next frame sync word
    - returns none when no sync word is found
    # framing
  mp3_decoder.parse_header
    fn (data: bytes, offset: i64) -> result[frame_header, string]
    + parses the 4-byte frame header at the offset
    - returns error on an invalid bitrate or sampling-rate field
    # framing
    -> std.bits.reader_new
    -> std.bits.read_bits
  mp3_decoder.frame_length
    fn (header: frame_header) -> i32
    + returns the total byte length of a frame given its header
    # framing
  mp3_decoder.read_side_info
    fn (header: frame_header, reader: bit_reader) -> tuple[side_info, bit_reader]
    + reads per-granule side information from the bitstream
    # side_info
    -> std.bits.read_bits
  mp3_decoder.huffman_decode
    fn (side: side_info, reader: bit_reader) -> tuple[list[i32], bit_reader]
    + decodes Huffman-coded frequency samples for one granule/channel
    # huffman
    -> std.bits.read_bits
  mp3_decoder.requantize
    fn (samples: list[i32], side: side_info) -> list[f32]
    + applies the MPEG-standard requantization formula to decoded samples
    # requantization
  mp3_decoder.reorder
    fn (samples: list[f32], side: side_info) -> list[f32]
    + reorders short-block samples into frequency order
    # reorder
  mp3_decoder.imdct
    fn (samples: list[f32], block_type: i32) -> list[f32]
    + applies the inverse modified discrete cosine transform for the block type
    # imdct
  mp3_decoder.overlap_add
    fn (current: list[f32], previous: list[f32]) -> tuple[list[f32], list[f32]]
    + returns (pcm_subband_input, new_overlap_buffer)
    # overlap
  mp3_decoder.synth_filterbank
    fn (subband: list[f32], state: synth_state) -> tuple[list[i16], synth_state]
    + produces 576 PCM samples per granule through the polyphase filterbank
    # synthesis
  mp3_decoder.decode_frame
    fn (data: bytes, offset: i64, state: synth_state) -> result[tuple[list[i16], i64, synth_state], string]
    + returns (pcm_samples, next_offset, new_state) for one full frame
    - returns error on a malformed frame
    # decoding
    -> std.bits.byte_align
