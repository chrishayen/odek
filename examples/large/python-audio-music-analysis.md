# Requirement: "a library for audio and music analysis"

Reads audio, runs spectral and temporal analysis primitives, and surfaces common features (tempo, onsets, spectrogram). The DSP primitives belong in std; the project layer composes them.

std
  std.math
    std.math.fft
      @ (samples: list[f64]) -> list[f64]
      + returns the magnitude spectrum of the input window
      ? length must be a power of two; caller pads
      # dsp
    std.math.ifft
      @ (spectrum: list[f64]) -> list[f64]
      + returns the inverse transform of a magnitude spectrum
      # dsp
    std.math.hann_window
      @ (size: i32) -> list[f64]
      + returns a Hann window of the given size
      # dsp
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the path does not exist
      # filesystem

audio
  audio.load_wav
    @ (path: string) -> result[audio_clip, string]
    + reads a wav file and returns a clip with samples and sample_rate
    - returns error on unsupported format or truncated file
    # io
    -> std.fs.read_all
  audio.resample
    @ (clip: audio_clip, target_rate: i32) -> audio_clip
    + returns a new clip resampled to target_rate using linear interpolation
    # preprocessing
  audio.to_mono
    @ (clip: audio_clip) -> audio_clip
    + averages channels into a single channel
    + returns the clip unchanged when already mono
    # preprocessing
  audio.stft
    @ (clip: audio_clip, window_size: i32, hop: i32) -> list[list[f64]]
    + returns a short-time Fourier transform matrix; outer list is frames
    ? windows shorter than window_size are zero-padded
    # spectral
    -> std.math.hann_window
    -> std.math.fft
  audio.mel_spectrogram
    @ (clip: audio_clip, n_mels: i32, window_size: i32, hop: i32) -> list[list[f64]]
    + returns a mel-scaled power spectrogram with n_mels bands
    # spectral
  audio.detect_onsets
    @ (clip: audio_clip) -> list[f64]
    + returns onset times in seconds by taking spectral-flux peaks
    # temporal
  audio.estimate_tempo
    @ (clip: audio_clip) -> f64
    + returns tempo in beats per minute based on autocorrelation of onset strength
    # temporal
  audio.chroma
    @ (clip: audio_clip) -> list[list[f64]]
    + returns a 12-bin chromagram frame sequence
    # harmonic
  audio.zero_crossing_rate
    @ (clip: audio_clip, window_size: i32, hop: i32) -> list[f64]
    + returns the zero-crossing rate per frame
    # temporal
  audio.rms_energy
    @ (clip: audio_clip, window_size: i32, hop: i32) -> list[f64]
    + returns the root-mean-square energy per frame
    # temporal
