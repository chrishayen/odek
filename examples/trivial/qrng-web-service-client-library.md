# Requirement: "a client library for a quantum random number web service"

Fetches random integers from a remote quantum RNG endpoint.

std
  std.http
    std.http.get
      fn (url: string, headers: map[string,string]) -> result[bytes, string]
      + returns body bytes on 2xx
      - returns error on transport failure or non-2xx
      # http
  std.json
    std.json.parse
      fn (raw: string) -> result[map[string,string], string]
      + parses a flat JSON object
      - returns error on invalid JSON
      # serialization

qrng
  qrng.fetch_uints
    fn (endpoint: string, count: i32, byte_size: i32) -> result[list[u32], string]
    + returns count quantum-derived unsigned integers
    - returns error when count is not in [1, 1024]
    - returns error when the service response does not contain a data array
    # fetch
    -> std.http.get
    -> std.json.parse
