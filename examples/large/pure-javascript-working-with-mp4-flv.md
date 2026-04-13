# Requirement: "a library for reading MP4 and FLV video files and producing MPEG-TS segments for HLS streaming"

Parse container boxes and tags, extract samples, repackage into MPEG-TS segments, and emit an HLS playlist.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the file contents
      - returns error when the path cannot be opened
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes the bytes to the file
      # filesystem

video
  video.parse_mp4_boxes
    @ (raw: bytes) -> result[list[mp4_box], string]
    + returns the top-level boxes of the isobmff file
    - returns error on truncated box headers
    # parsing
  video.extract_mp4_tracks
    @ (boxes: list[mp4_box]) -> result[list[track], string]
    + returns the tracks with codec metadata and sample tables
    - returns error when moov or trak is missing
    # parsing
  video.read_mp4_samples
    @ (raw: bytes, track: track) -> result[list[sample], string]
    + returns the sample offsets, sizes, durations, and keyframe flags
    - returns error when the sample table is malformed
    # parsing
  video.parse_flv_tags
    @ (raw: bytes) -> result[list[flv_tag], string]
    + returns the flv tags in file order
    - returns error when the flv header is missing
    # parsing
  video.extract_flv_tracks
    @ (tags: list[flv_tag]) -> result[list[track], string]
    + returns the audio and video tracks with codec metadata
    # parsing
  video.split_segments
    @ (samples: list[sample], target_duration_ms: i64) -> list[segment]
    + returns segments each starting at a keyframe with approximately the target duration
    + returns a single segment for very short inputs
    # segmentation
  video.encode_ts_segment
    @ (segment: segment, video_track: track, audio_track: optional[track]) -> bytes
    + returns the mpeg-ts bytes for the segment using pes and adaptation fields
    # muxing
  video.build_hls_playlist
    @ (segment_names: list[string], segment_durations_ms: list[i64], target_duration_s: i32) -> string
    + returns the m3u8 playlist body referencing the segments
    + includes the EXT-X-ENDLIST tag
    # playlist
  video.repackage_to_hls
    @ (source_path: string, output_dir: string, target_duration_ms: i64) -> result[list[string], string]
    + returns the names of the written segment files
    + writes an index.m3u8 alongside the segments
    - returns error when the source container is unsupported
    # pipeline
    -> std.fs.read_all
    -> std.fs.write_all
  video.pcr_from_timestamp
    @ (pts_90khz: i64) -> i64
    + returns the program clock reference ticks for the given pts
    # muxing
  video.sample_to_pes
    @ (sample: sample, stream_id: u8, pts_90khz: i64, dts_90khz: i64) -> bytes
    + returns the pes packet wrapping the sample bytes
    # muxing
