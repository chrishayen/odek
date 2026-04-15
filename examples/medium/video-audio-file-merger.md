# Requirement: "a library for merging video and audio files"

Concatenates multiple video files and optionally muxes in a separate audio track. Actual codec work is delegated to a pluggable transcoder; this library orchestrates the pipeline.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the full file
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating or replacing
      # filesystem
    std.fs.temp_path
      fn (prefix: string) -> string
      + returns a unique path in the temp directory
      # filesystem

vidmerger
  vidmerger.probe_container
    fn (data: bytes) -> result[media_info, string]
    + returns duration, video stream count, and audio stream count
    - returns error when the container header is unrecognized
    # introspection
  vidmerger.validate_compatible
    fn (infos: list[media_info]) -> result[void, string]
    + returns ok when all inputs share codec, resolution, and frame rate
    - returns error describing the first mismatch
    # validation
  vidmerger.build_concat_plan
    fn (paths: list[string]) -> result[concat_plan, string]
    + returns a plan listing inputs in order with cumulative offsets
    - returns error when the input list is empty
    # planning
    -> std.fs.read_all
    -> vidmerger.probe_container
    -> vidmerger.validate_compatible
  vidmerger.concat
    fn (plan: concat_plan, output_path: string) -> result[void, string]
    + writes a single video file concatenating the planned inputs
    - returns error when writing the output fails
    # muxing
    -> std.fs.write_all
    -> std.fs.temp_path
  vidmerger.mux_audio
    fn (video_path: string, audio_path: string, output_path: string) -> result[void, string]
    + combines a video stream with an external audio track
    - returns error when audio duration does not match video duration
    # muxing
    -> std.fs.read_all
    -> std.fs.write_all
  vidmerger.merge
    fn (video_paths: list[string], audio_path: optional[string], output_path: string) -> result[void, string]
    + concatenates videos and optionally layers an audio track
    # entrypoint
    -> vidmerger.build_concat_plan
    -> vidmerger.concat
    -> vidmerger.mux_audio
