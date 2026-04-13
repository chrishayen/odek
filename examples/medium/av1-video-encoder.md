# Requirement: "an AV1 video encoder"

Takes a sequence of raw frames and produces an AV1 bitstream. Splits each frame into superblocks, picks intra or inter prediction, applies DCT and quantization, and emits entropy-coded symbols into a packed bitstream.

std
  std.dsp
    std.dsp.dct_8x8
      @ (block: list[i16]) -> list[i16]
      + returns the 8x8 forward DCT of the input block
      ? input and output are row-major 64-element vectors
      # signal
    std.dsp.idct_8x8
      @ (block: list[i16]) -> list[i16]
      + returns the 8x8 inverse DCT
      # signal
  std.bitstream
    std.bitstream.new_writer
      @ () -> bit_writer
      + returns an empty bit writer
      # bitstream
    std.bitstream.write_bits
      @ (w: bit_writer, value: u64, bit_count: i32) -> bit_writer
      + writes bit_count bits of value, msb first
      # bitstream
    std.bitstream.finish
      @ (w: bit_writer) -> bytes
      + flushes the writer to a byte buffer
      # bitstream

av1_encoder
  av1_encoder.new
    @ (width: i32, height: i32, quality: i32) -> encoder_state
    + returns an encoder configured for the given frame size and quality
    # construction
  av1_encoder.split_superblocks
    @ (frame: raw_frame) -> list[superblock]
    + partitions a frame into 64x64 superblocks
    # partitioning
  av1_encoder.choose_prediction
    @ (state: encoder_state, sb: superblock) -> prediction_mode
    + returns the prediction mode with the lowest rate-distortion cost
    + falls back to intra prediction for the first frame
    # mode_decision
  av1_encoder.transform_quantize
    @ (state: encoder_state, residual: list[i16]) -> list[i16]
    + applies DCT and quantization to a residual block
    # transform
    -> std.dsp.dct_8x8
  av1_encoder.encode_frame
    @ (state: encoder_state, frame: raw_frame) -> tuple[encoder_state, bytes]
    + returns updated state and the bitstream bytes for one frame
    # frame_encoding
    -> std.bitstream.new_writer
    -> std.bitstream.write_bits
    -> std.bitstream.finish
  av1_encoder.finalize
    @ (state: encoder_state) -> bytes
    + returns trailing sequence headers and end-of-stream markers
    # finalization
