# Requirement: "a library for playing sounds and music"

One-shot sounds are loaded fully into memory; music streams from disk. Playback targets a std audio-device primitive.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the file does not exist
      # io
    std.fs.open_reader
      @ (path: string) -> result[byte_reader, string]
      + opens a streaming reader for the file at path
      - returns error when the file does not exist
      # io
  std.audio_codec
    std.audio_codec.decode_all
      @ (data: bytes) -> result[audio_samples, string]
      + decodes a complete audio file into pcm samples
      - returns error on unknown format signatures
      # codec
    std.audio_codec.decode_stream
      @ (reader: byte_reader) -> result[audio_stream, string]
      + wraps a reader in a pull-based decoded sample stream
      - returns error when the header is malformed
      # codec
  std.audio_device
    std.audio_device.open_default
      @ (sample_rate: i32, channels: i32) -> result[audio_device, string]
      + opens the system default output device
      - returns error when no device is available
      # device
    std.audio_device.write
      @ (device: audio_device, samples: audio_samples) -> result[void, string]
      + queues samples for playback
      # device

ears
  ears.load_sound
    @ (path: string) -> result[sound_handle, string]
    + loads a short clip into memory
    - returns error when decoding fails
    # sounds
    -> std.fs.read_all
    -> std.audio_codec.decode_all
  ears.play_sound
    @ (device: audio_device, sound: sound_handle) -> result[void, string]
    + writes the decoded samples to the device
    # sounds
    -> std.audio_device.write
  ears.open_music
    @ (path: string) -> result[music_handle, string]
    + opens a streaming music source
    - returns error when the file cannot be opened
    # music
    -> std.fs.open_reader
    -> std.audio_codec.decode_stream
  ears.pump_music
    @ (device: audio_device, music: music_handle, max_samples: i32) -> result[bool, string]
    + decodes and writes up to max_samples to the device; returns false when the stream is exhausted
    # music
    -> std.audio_device.write
  ears.close_music
    @ (music: music_handle) -> result[void, string]
    + releases the underlying stream
    # music
