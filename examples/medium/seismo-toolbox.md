# Requirement: "a seismology toolbox"

Represents a seismic trace as an opaque handle; operations are DSP transforms and simple event detection.

std: (all units exist)

seismo
  seismo.new_trace
    fn (samples: list[f64], sample_rate_hz: f64) -> trace_state
    + creates a trace from evenly sampled values
    ? sample_rate_hz must be positive; caller's responsibility
    # construction
  seismo.detrend
    fn (trace: trace_state) -> trace_state
    + subtracts a best-fit linear trend from the samples
    # filtering
  seismo.bandpass
    fn (trace: trace_state, low_hz: f64, high_hz: f64) -> result[trace_state, string]
    + applies a zero-phase bandpass filter between low_hz and high_hz
    - returns error when low_hz >= high_hz or high_hz > nyquist
    # filtering
  seismo.normalize
    fn (trace: trace_state) -> trace_state
    + scales samples so the maximum absolute value is 1.0
    # filtering
  seismo.sta_lta
    fn (trace: trace_state, short_window_sec: f64, long_window_sec: f64) -> list[f64]
    + returns the short-term over long-term average ratio per sample
    ? classic trigger for onset detection
    # detection
  seismo.detect_onsets
    fn (ratios: list[f64], threshold: f64) -> list[i32]
    + returns sample indices where the ratio first crosses threshold
    # detection
  seismo.peak_amplitude
    fn (trace: trace_state) -> f64
    + returns the maximum absolute sample value
    # analysis
  seismo.rms
    fn (trace: trace_state) -> f64
    + returns the root-mean-square of the samples
    # analysis
