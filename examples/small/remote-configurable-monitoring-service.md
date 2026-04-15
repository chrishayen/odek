# Requirement: "a lightweight, remotely configurable monitoring service"

Runs user-defined probes on a schedule. The probe set is loaded from a remote config URL so operators can change it without redeploying.

std
  std.http
    std.http.get_string
      fn (url: string) -> result[string, string]
      + fetches url and returns the body as a string on 2xx
      - returns error on non-2xx or network failure
      # http_client
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

monitor
  monitor.load_probes
    fn (config_url: string) -> result[list[probe], string]
    + fetches the config and decodes it into a list of probes
    - returns error when the response is not a valid probe list
    # config_loading
    -> std.http.get_string
  monitor.run_probe
    fn (p: probe) -> probe_result
    + executes the probe and returns status, latency_ms, and the checked-at timestamp
    - marks the result failed when the probe times out or returns a non-ok status
    # probe_execution
    -> std.time.now_millis
    -> std.http.get_string
  monitor.run_all
    fn (probes: list[probe]) -> list[probe_result]
    + runs every probe once and returns the full result set in input order
    # scheduling
