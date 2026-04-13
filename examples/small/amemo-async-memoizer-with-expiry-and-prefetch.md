# Requirement: "memoize async-returning functions with expiry and prefetch"

Wraps a function so concurrent calls share the same in-flight result. Cached entries expire after a TTL, and prefetch refreshes an entry in the background before it expires.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current time in milliseconds
      # time

amemo
  amemo.new_cache
    @ (ttl_millis: i64, prefetch_millis: i64) -> amemo_cache
    + creates a cache with the given expiry and prefetch window
    ? prefetch_millis is how long before expiry a background refresh starts
    # construction
  amemo.get_or_compute
    @ (cache: amemo_cache, key: string, loader: async_fn) -> result[bytes, string]
    + returns a cached entry when present and fresh
    + invokes the loader and caches the result on miss
    + returns the pending result to concurrent callers sharing the same key
    - propagates the loader error when loading fails
    # read_through
    -> std.time.now_millis
  amemo.invalidate
    @ (cache: amemo_cache, key: string) -> amemo_cache
    + removes any entry stored under key
    # invalidation
  amemo.sweep_expired
    @ (cache: amemo_cache) -> amemo_cache
    + drops entries whose age exceeds the ttl
    # maintenance
    -> std.time.now_millis
