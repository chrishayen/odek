# Requirement: "a high-level audio manipulation library"

In-memory PCM audio segments with slicing, concatenation, gain, and fade operations. Decoding and encoding are delegated to std codec primitives.

std
  std.audio
    std.audio.decode_wav
      @ (data: bytes) -> result[pcm_buffer, string]
      + decodes a little-endian PCM WAV file into sample_rate, channels, and samples
      - returns error on unsupported format or truncated header
      # decoding
    std.audio.encode_wav
      @ (buf: pcm_buffer) -> bytes
      + encodes a PCM buffer as a little-endian WAV file
      # encoding
  std.math
    std.math.db_to_linear
      @ (db: f64) -> f64
      + converts a decibel value to a linear amplitude ratio
      # math

audio
  audio.from_wav
    @ (data: bytes) -> result[audio_segment, string]
    + creates a segment from WAV bytes
    - returns error when decoding fails
    # construction
    -> std.audio.decode_wav
  audio.to_wav
    @ (seg: audio_segment) -> bytes
    + renders the segment back to WAV bytes
    # serialization
    -> std.audio.encode_wav
  audio.duration_ms
    @ (seg: audio_segment) -> i64
    + returns the segment duration in milliseconds
    # inspection
  audio.slice
    @ (seg: audio_segment, start_ms: i64, end_ms: i64) -> result[audio_segment, string]
    + returns a copy of the samples between start and end
    - returns error when start or end falls outside the segment
    # editing
  audio.concat
    @ (a: audio_segment, b: audio_segment) -> result[audio_segment, string]
    + returns a new segment containing a followed by b
    - returns error when sample rate or channel count differ
    # editing
  audio.apply_gain
    @ (seg: audio_segment, gain_db: f64) -> audio_segment
    + scales every sample by the gain expressed in dB
    + samples are clipped to the native range to avoid wraparound
    # effects
    -> std.math.db_to_linear
  audio.fade
    @ (seg: audio_segment, fade_in_ms: i64, fade_out_ms: i64) -> audio_segment
    + applies linear fade-in and fade-out envelopes
    # effects
