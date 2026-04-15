# Requirement: "stream a torrent video to a media-casting device"

The project bridges a torrent-backed byte source to an HTTP range server that a cast-capable device can pull from, then controls playback on the device through a pluggable control interface.

std
  std.http
    std.http.serve
      fn (addr: string, handler: request_handler) -> result[void, string]
      + binds the address and dispatches each request to the handler
      - returns error when the address cannot be bound
      # network
    std.http.parse_range_header
      fn (raw: string) -> result[tuple[i64, i64], string]
      + returns (start, end) for a "bytes=start-end" header
      - returns error on malformed input
      # network

torrent_caster
  torrent_caster.open_source
    fn (torrent_path: string, file_index: i32) -> result[torrent_source, string]
    + returns a seekable source for the selected file inside the torrent
    - returns error when the torrent cannot be loaded
    - returns error when the file index is out of range
    # source
  torrent_caster.read_range
    fn (source: torrent_source, start: i64, end: i64) -> result[bytes, string]
    + returns the bytes for the inclusive range
    - returns error when start > end or end exceeds the file size
    # source_read
  torrent_caster.handle_stream_request
    fn (source: torrent_source, range_header: optional[string]) -> result[tuple[i32, map[string,string], bytes], string]
    + returns a 206 partial-content response with Content-Range set
    + returns a 200 full response when no range header is present
    - returns a 416 when the range is unsatisfiable
    # http_handler
    -> std.http.parse_range_header
  torrent_caster.serve
    fn (source: torrent_source, addr: string) -> result[void, string]
    + serves the torrent file under a single endpoint until interrupted
    - returns error when the address cannot be bound
    # server
    -> std.http.serve
  torrent_caster.cast
    fn (controller: cast_controller, stream_url: string, mime: string) -> result[void, string]
    + tells the controller to start playback of the stream URL
    - returns error when the controller reports failure
    # playback
