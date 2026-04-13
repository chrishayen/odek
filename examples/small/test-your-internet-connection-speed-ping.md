# Requirement: "a library that measures internet connection speed and round-trip latency"

Download/upload throughput and ping are three measurements. Time and HTTP primitives live in std.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.net
    std.net.http_get_bytes
      @ (url: string) -> result[bytes, string]
      + performs an HTTP GET and returns the response body
      - returns error on any network or status failure
      # network
    std.net.http_post_bytes
      @ (url: string, body: bytes) -> result[i64, string]
      + performs an HTTP POST and returns the number of bytes accepted
      - returns error on any network or status failure
      # network

speed_test
  speed_test.measure_ping
    @ (url: string, samples: i32) -> result[f64, string]
    + returns the mean round-trip time in milliseconds over the given number of HEAD-like probes
    - returns error when samples is less than one
    # latency
    -> std.time.now_millis
    -> std.net.http_get_bytes
  speed_test.measure_download
    @ (url: string) -> result[f64, string]
    + returns download throughput in megabits per second for a single body fetch
    - returns error when the response body is empty
    # download
    -> std.time.now_millis
    -> std.net.http_get_bytes
  speed_test.measure_upload
    @ (url: string, payload_bytes: i64) -> result[f64, string]
    + returns upload throughput in megabits per second for a single POST of the given size
    - returns error when payload_bytes is not positive
    # upload
    -> std.time.now_millis
    -> std.net.http_post_bytes
