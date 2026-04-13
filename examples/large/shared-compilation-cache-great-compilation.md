# Requirement: "a shared compilation cache for compiler invocations"

Normalizes a compiler command into a cache key, looks it up in a local or remote store, and returns cached artifacts on hit or records new ones on miss.

std
  std.crypto
    std.crypto.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex SHA-256 digest of data
      # cryptography
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the path is missing
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, creating parents as needed
      # filesystem
    std.fs.exists
      @ (path: string) -> bool
      + returns true when a path exists
      # filesystem
  std.proc
    std.proc.run_capture
      @ (cmd: string, args: list[string]) -> result[tuple[i32, bytes, bytes], string]
      + runs a subprocess and returns (exit_code, stdout, stderr)
      - returns error when the binary cannot be launched
      # process
  std.http
    std.http.get
      @ (url: string) -> result[bytes, string]
      + fetches a URL and returns the response body
      - returns error on non-2xx or network failure
      # http
    std.http.put
      @ (url: string, body: bytes) -> result[i32, string]
      + uploads bytes with PUT and returns the status code
      - returns error on network failure
      # http

ccache
  ccache.new_local
    @ (root: string, max_bytes: i64) -> cache_state
    + creates a local cache rooted at a directory with an eviction budget
    # construction
  ccache.new_remote
    @ (base_url: string) -> cache_state
    + creates a cache backed by a remote store reachable at base_url
    # construction
  ccache.parse_invocation
    @ (cmd: string, args: list[string]) -> result[invocation, string]
    + classifies a compiler invocation into inputs, outputs, and normalized flags
    - returns error when the compiler is not recognized
    - returns error when the invocation cannot be cached (e.g. links stdin)
    # parsing
  ccache.compute_key
    @ (inv: invocation) -> string
    + returns the hex cache key derived from compiler, flags, and source hashes
    # keying
    -> std.fs.read_all
    -> std.crypto.sha256_hex
  ccache.lookup
    @ (state: cache_state, key: string) -> result[list[artifact], string]
    + returns the cached artifacts for key
    - returns error "miss" when the key is absent
    # lookup
    -> std.fs.exists
    -> std.fs.read_all
    -> std.http.get
  ccache.store
    @ (state: cache_state, key: string, artifacts: list[artifact]) -> result[cache_state, string]
    + stores artifacts under key in the cache
    # store
    -> std.fs.write_all
    -> std.http.put
  ccache.evict
    @ (state: cache_state) -> cache_state
    + drops least-recently-used entries until the cache is within max_bytes
    # eviction
  ccache.run
    @ (state: cache_state, cmd: string, args: list[string]) -> result[tuple[i32, bytes, bytes], string]
    + parses the invocation, checks the cache, runs the compiler on miss, and records the result
    + returns the cached outputs on hit without invoking the compiler
    - returns error when the invocation cannot be cached
    # orchestration
    -> std.proc.run_capture
  ccache.stats
    @ (state: cache_state) -> cache_stats
    + returns hit count, miss count, and total bytes stored
    # stats
  ccache.purge
    @ (state: cache_state) -> cache_state
    + removes every cached entry
    # purge
