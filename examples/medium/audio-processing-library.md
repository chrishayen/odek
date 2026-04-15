# Requirement: "an audio processing library"

Load PCM buffers, apply common effects, and mix tracks. Sample format is normalized to f32 interleaved.

std
  std.math
    std.math.sin
      fn (x: f64) -> f64
      + returns the sine of x
      # math
    std.math.log10
      fn (x: f64) -> f64
      + returns the base-10 logarithm of x
      # math

audio
  audio.buffer_new
    fn (sample_rate: i32, channels: i32, samples: list[f32]) -> result[audio_buffer, string]
    + wraps interleaved samples with the given rate and channel count
    - returns error when len(samples) is not a multiple of channels
    # construction
  audio.gain
    fn (buf: audio_buffer, db: f64) -> audio_buffer
    + scales every sample by 10^(db/20)
    ? clipping to [-1, 1] is the caller's responsibility
    # effects
    -> std.math.log10
  audio.mix
    fn (tracks: list[audio_buffer]) -> result[audio_buffer, string]
    + sums aligned tracks sample-by-sample
    - returns error when sample rates or channel counts disagree
    - returns error when tracks is empty
    # mixing
  audio.resample_linear
    fn (buf: audio_buffer, target_rate: i32) -> audio_buffer
    + returns a buffer at target_rate using linear interpolation per channel
    # resampling
  audio.generate_sine
    fn (sample_rate: i32, freq_hz: f64, duration_seconds: f64) -> audio_buffer
    + returns a mono buffer containing a sine wave of the given frequency and length
    # synthesis
    -> std.math.sin
  audio.fade
    fn (buf: audio_buffer, fade_in_seconds: f64, fade_out_seconds: f64) -> audio_buffer
    + applies linear fade-in and fade-out envelopes
    ? overlapping envelopes are allowed and compose multiplicatively
    # effects
