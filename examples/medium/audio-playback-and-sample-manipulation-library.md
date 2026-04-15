# Requirement: "a library for audio playback and sample manipulation"

A composable stream graph over pcm samples. Sources feed through effects into a device sink.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the file does not exist
      # io
  std.audio_codec
    std.audio_codec.decode
      fn (data: bytes) -> result[audio_samples, string]
      + decodes audio bytes into pcm samples
      - returns error on unknown format signatures
      # codec
  std.audio_device
    std.audio_device.open_default
      fn (sample_rate: i32, channels: i32) -> result[audio_device, string]
      + opens the system default output device
      - returns error when no device is available
      # device
    std.audio_device.write
      fn (device: audio_device, samples: audio_samples) -> result[void, string]
      + queues samples for playback
      # device

beep
  beep.load
    fn (path: string) -> result[audio_stream, string]
    + loads an audio file into a pullable sample stream
    - returns error when decoding fails
    # loading
    -> std.fs.read_all
    -> std.audio_codec.decode
  beep.take
    fn (stream: audio_stream, count: i32) -> tuple[audio_samples, audio_stream]
    + pulls up to count samples from the stream
    # streaming
  beep.gain
    fn (stream: audio_stream, factor: f32) -> audio_stream
    + returns a stream whose samples are multiplied by factor
    # effects
  beep.mix
    fn (a: audio_stream, b: audio_stream) -> audio_stream
    + returns a stream summing a and b sample-for-sample
    ? ends when both inputs are exhausted; the shorter input contributes silence past its end
    # effects
  beep.resample
    fn (stream: audio_stream, from_rate: i32, to_rate: i32) -> audio_stream
    + returns a stream resampled to to_rate via linear interpolation
    # effects
  beep.play
    fn (device: audio_device, stream: audio_stream, chunk: i32) -> result[void, string]
    + pulls chunk samples at a time from stream and writes them until exhausted
    # playback
    -> std.audio_device.write
