# Requirement: "a database backup tool from different source drivers to different destinations"

A small orchestrator: pluggable source drivers produce backup byte streams, pluggable destinations consume them, and a runner glues a config to both.

std
  std.fs
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to a file, replacing any existing contents
      # filesystem
  std.net
    std.net.http_put
      fn (url: string, headers: map[string, string], body: bytes) -> result[u16, string]
      + uploads bytes via HTTP PUT and returns the status code
      - returns error on network failure
      # network
  std.time
    std.time.now_unix
      fn () -> i64
      + returns current unix time in seconds
      # time

onedump
  onedump.register_source
    fn (driver: string, handler: fn(map[string, string]) -> result[bytes, string]) -> void
    + registers a source driver that maps a config map to a byte stream
    # driver_registry
  onedump.register_destination
    fn (name: string, handler: fn(bytes, map[string, string]) -> result[void, string]) -> void
    + registers a destination handler keyed by name
    # driver_registry
  onedump.dump_local_file
    fn (data: bytes, config: map[string, string]) -> result[void, string]
    + writes the bytes to the path given by config["path"]
    - returns error when "path" is missing
    # destination
    -> std.fs.write_all
  onedump.dump_http
    fn (data: bytes, config: map[string, string]) -> result[void, string]
    + PUTs the bytes to the URL given by config["url"]
    - returns error when the server responds with a non-2xx status
    # destination
    -> std.net.http_put
  onedump.run_job
    fn (source_driver: string, source_config: map[string, string], destination: string, destination_config: map[string, string]) -> result[backup_receipt, string]
    + produces a backup from the source and sends it to the destination, returning a receipt with timestamp and byte count
    - returns error when the source driver is not registered
    - returns error when the destination handler is not registered
    # orchestration
    -> std.time.now_unix
  onedump.run_config
    fn (config: backup_config) -> result[list[backup_receipt], list[string]]
    + runs every job in a config, returning per-job receipts and a list of errors for any that failed
    # orchestration
