# Requirement: "an automated reference audio mastering library"

Given a target mix and a reference track, matches loudness, spectral balance, and dynamics of the target toward the reference. Heavy DSP primitives live in std; the project layer is the matching pipeline.

std
  std.audio
    std.audio.read_wav
      @ (path: string) -> result[audio_buffer, string]
      + decodes a PCM wav file into sample frames with metadata
      - returns error on unsupported bit depth or format
      # audio_io
    std.audio.write_wav
      @ (path: string, buf: audio_buffer) -> result[void, string]
      + encodes samples to a PCM wav file
      # audio_io
    std.audio.resample
      @ (buf: audio_buffer, target_rate: i32) -> audio_buffer
      + returns a buffer resampled to target_rate with a low-pass antialias filter
      # dsp
  std.dsp
    std.dsp.fft
      @ (samples: list[f32]) -> list[complex]
      + returns the forward FFT of the input window
      ? input length must be a power of two
      # dsp
    std.dsp.ifft
      @ (spectrum: list[complex]) -> list[f32]
      + returns the inverse FFT as real samples
      # dsp
    std.dsp.window_hann
      @ (n: i32) -> list[f32]
      + returns a Hann window of length n
      # dsp
    std.dsp.stft
      @ (samples: list[f32], window: list[f32], hop: i32) -> list[list[complex]]
      + returns overlapping FFT frames using the supplied window and hop
      # dsp
    std.dsp.istft
      @ (frames: list[list[complex]], window: list[f32], hop: i32) -> list[f32]
      + reconstructs a signal from STFT frames using overlap-add
      # dsp
  std.loudness
    std.loudness.lufs_integrated
      @ (buf: audio_buffer) -> f64
      + returns the integrated loudness in LUFS per ITU-R BS.1770
      # metering

mastering
  mastering.analyze_reference
    @ (reference: audio_buffer) -> reference_profile
    + returns a profile capturing loudness, average magnitude spectrum, and peak envelope
    # analysis
    -> std.dsp.stft
    -> std.loudness.lufs_integrated
  mastering.match_loudness
    @ (target: audio_buffer, profile: reference_profile) -> audio_buffer
    + returns the target scaled so its integrated loudness matches the reference
    # loudness_matching
    -> std.loudness.lufs_integrated
  mastering.match_spectrum
    @ (target: audio_buffer, profile: reference_profile) -> audio_buffer
    + applies a matching EQ derived from the ratio of reference and target magnitude spectra
    ? smoothing is applied across bands to avoid narrow spikes
    # spectral_matching
    -> std.dsp.stft
    -> std.dsp.istft
  mastering.match_dynamics
    @ (target: audio_buffer, profile: reference_profile) -> audio_buffer
    + applies a multiband compressor tuned so target crest factor approaches reference
    # dynamics_matching
  mastering.brickwall_limiter
    @ (buf: audio_buffer, ceiling_db: f64) -> audio_buffer
    + returns a buffer whose peaks do not exceed the ceiling
    - leaves samples unchanged when no peak exceeds the ceiling
    # limiting
  mastering.master
    @ (target: audio_buffer, reference: audio_buffer) -> result[audio_buffer, string]
    + runs the full analyze-match-limit pipeline and returns the mastered target
    - returns error when sample rates differ and resampling fails
    # pipeline
    -> std.audio.resample
