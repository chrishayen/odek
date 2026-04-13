# Requirement: "a do-it-yourself video security library"

Captures frames from a camera source, runs motion detection, and triggers alerts with clip recording.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.fs
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, replacing any existing contents
      - returns error on io failure
      # storage
  std.image
    std.image.decode_frame
      @ (raw: bytes) -> result[frame, string]
      + decodes an encoded image frame into a pixel buffer
      - returns error on unsupported format
      # image
    std.image.to_grayscale
      @ (f: frame) -> frame
      + converts a color frame to a single-channel grayscale frame
      # image

vsec
  vsec.open_camera
    @ (source: string) -> result[camera_handle, string]
    + opens a camera or video source for frame reads
    - returns error when the source cannot be opened
    # capture
  vsec.read_frame
    @ (cam: camera_handle) -> result[frame, string]
    + reads and decodes the next frame from the camera
    - returns error at end of stream
    # capture
    -> std.image.decode_frame
  vsec.new_detector
    @ (threshold: f32, min_area: i32) -> detector_state
    + creates a motion detector with sensitivity and minimum region area
    # motion
  vsec.detect_motion
    @ (state: detector_state, current: frame) -> tuple[bool, detector_state]
    + compares the current grayscale frame against the rolling background model
    + returns true when a region exceeds the threshold and min_area
    # motion
    -> std.image.to_grayscale
  vsec.start_clip
    @ (alerts: alert_state, frame0: frame) -> alert_state
    + begins buffering frames for a new motion alert clip
    # recording
    -> std.time.now_millis
  vsec.append_frame
    @ (alerts: alert_state, f: frame) -> alert_state
    + appends a frame to the active clip buffer
    # recording
  vsec.finalize_clip
    @ (alerts: alert_state, output_dir: string) -> result[alert_state, string]
    + encodes the buffered frames and writes the clip to the output directory
    # recording
    -> std.fs.write_all
  vsec.notify
    @ (alerts: alert_state, sink: notification_sink) -> alert_state
    + delivers the most recent clip metadata to a pluggable notification sink
    # alerting
