# Requirement: "a script-based video editing library supporting many formats including animated images"

A timeline model that assembles clips, applies per-clip transforms, and renders to a destination container via pluggable codec primitives.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the file contents
      - returns error when the file does not exist
      # io
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to the path, replacing any existing content
      # io
  std.codec
    std.codec.decode_video
      @ (container: bytes) -> result[video_stream, string]
      + demuxes and decodes a video container into raw frames plus metadata
      - returns error when no supported video track is found
      # codec
    std.codec.encode_video
      @ (stream: video_stream, format: string) -> result[bytes, string]
      + encodes a frame stream into the named container format
      - returns error when format is not supported
      # codec
    std.codec.frame_count
      @ (stream: video_stream) -> i32
      + returns the number of frames in the stream
      # codec

video_edit
  video_edit.load_clip
    @ (path: string) -> result[clip, string]
    + returns a clip wrapping the decoded video stream
    - returns error when the file cannot be read or decoded
    # loading
    -> std.fs.read_all
    -> std.codec.decode_video
  video_edit.subclip
    @ (c: clip, start_sec: f64, end_sec: f64) -> result[clip, string]
    + returns a clip restricted to the given time range
    - returns error when end_sec <= start_sec
    # editing
  video_edit.resize
    @ (c: clip, width: i32, height: i32) -> clip
    + returns a clip whose frames are resampled to the new dimensions
    # transform
  video_edit.concat
    @ (clips: list[clip]) -> result[clip, string]
    + returns a single clip that plays the inputs back-to-back
    - returns error when the clip list is empty
    # editing
  video_edit.overlay_text
    @ (c: clip, text: string, start_sec: f64, end_sec: f64) -> clip
    + returns a clip with the text burned in over the given interval
    # transform
  video_edit.render
    @ (c: clip, path: string, format: string) -> result[void, string]
    + encodes the clip to the destination file in the given format
    - returns error on encoder failure
    # rendering
    -> std.codec.encode_video
    -> std.fs.write_all
