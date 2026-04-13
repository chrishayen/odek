# Requirement: "an HLS (M3U8) playlist parser and generator"

Two entry points for reading and writing; supporting runes classify the playlist variant and walk segments/streams.

std
  std.strings
    std.strings.split_lines
      @ (s: string) -> list[string]
      + splits on LF, tolerating CRLF
      # strings
    std.strings.join
      @ (xs: list[string], sep: string) -> string
      + joins elements with sep
      # strings
    std.strings.starts_with
      @ (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings

hls
  hls.parse
    @ (source: string) -> result[playlist, string]
    + parses #EXTM3U, #EXT-X-VERSION, and tag/value pairs
    - returns error when the first line is not #EXTM3U
    # parsing
    -> std.strings.split_lines
    -> std.strings.starts_with
  hls.generate
    @ (p: playlist) -> string
    + emits an M3U8 text with tags in canonical order
    + always starts with #EXTM3U
    # generation
    -> std.strings.join
  hls.is_master
    @ (p: playlist) -> bool
    + returns true when the playlist contains #EXT-X-STREAM-INF entries
    # classification
  hls.is_media
    @ (p: playlist) -> bool
    + returns true when the playlist contains #EXTINF segments
    # classification
  hls.segments
    @ (p: playlist) -> list[media_segment]
    + returns every media segment with duration and URI
    # query
  hls.variants
    @ (p: playlist) -> list[variant_stream]
    + returns every variant stream with bandwidth and URI
    # query
  hls.total_duration
    @ (p: playlist) -> f64
    + returns the sum of all segment durations in seconds
    # query
