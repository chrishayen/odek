# Requirement: "a library for transcribing audio or video in any language on any platform"

A speech-to-text pipeline. The std layer provides generic audio decoding and model loading primitives; the project layer orchestrates demux, resample, inference, and output formatting.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads entire file contents
      - returns error when the path cannot be opened
      # io
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path
      - returns error when the path cannot be created
      # io
  std.audio
    std.audio.decode_container
      fn (data: bytes) -> result[audio_stream, string]
      + returns an audio stream from a container file (wav, mp3, mp4, mkv)
      - returns error when no audio track is present
      # audio
    std.audio.resample
      fn (stream: audio_stream, target_hz: i32, target_channels: i32) -> audio_stream
      + returns a new stream at the target sample rate and channel count
      # audio
    std.audio.to_mono_f32
      fn (stream: audio_stream) -> list[f32]
      + returns interleaved samples mixed to a single channel as float32 in [-1, 1]
      # audio
  std.text
    std.text.detect_language
      fn (sample: string) -> string
      + returns an ISO 639-1 code guessed from character frequency
      + returns "und" when the sample is too short
      # language

transcribe
  transcribe.load_model
    fn (path: string) -> result[model_handle, string]
    + loads a speech recognition model from disk
    - returns error when the file is not a recognized model format
    # model
    -> std.fs.read_all
  transcribe.prepare_audio
    fn (path: string) -> result[list[f32], string]
    + reads, demuxes, resamples to 16 kHz mono, and returns float samples
    - returns error when the file has no audio stream
    # preprocessing
    -> std.fs.read_all
    -> std.audio.decode_container
    -> std.audio.resample
    -> std.audio.to_mono_f32
  transcribe.run
    fn (model: model_handle, samples: list[f32], language_hint: optional[string]) -> result[transcript, string]
    + returns a transcript with segments, timestamps, and confidence scores
    - returns error when the model rejects the audio length
    ? language_hint=none triggers automatic language detection
    # inference
    -> std.text.detect_language
  transcribe.file
    fn (model: model_handle, path: string, language_hint: optional[string]) -> result[transcript, string]
    + convenience wrapper combining prepare_audio and run
    - returns error when either step fails
    # pipeline
  transcribe.to_srt
    fn (transcript: transcript) -> string
    + returns the transcript formatted as an SRT subtitle file
    # formatting
  transcribe.to_vtt
    fn (transcript: transcript) -> string
    + returns the transcript formatted as a WebVTT file
    # formatting
  transcribe.to_plain_text
    fn (transcript: transcript) -> string
    + returns the concatenated segment text separated by spaces
    # formatting
  transcribe.save
    fn (transcript: transcript, path: string, format: string) -> result[void, string]
    + writes the transcript in the requested format ("srt", "vtt", or "txt")
    - returns error when the format is unknown
    # io
    -> std.fs.write_all
  transcribe.translate_to_english
    fn (model: model_handle, samples: list[f32]) -> result[transcript, string]
    + returns an English transcript regardless of source language
    - returns error when the model does not support translation
    # translation
