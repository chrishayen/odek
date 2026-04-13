# Requirement: "a key-value storage library with a pluggable backend"

A generic KV facade over an injected backend. The backend is an opaque handle the caller supplies.

std: (all units exist)

keyv
  keyv.new
    @ (backend: kv_backend, namespace: string) -> keyv_store
    + builds a store that prefixes every key with namespace
    # construction
  keyv.get
    @ (store: keyv_store, key: string) -> result[optional[bytes], string]
    + returns the value for key, or none if absent
    - returns error when the backend is unreachable
    # reads
  keyv.set
    @ (store: keyv_store, key: string, value: bytes, ttl_ms: i64) -> result[void, string]
    + stores value under key; ttl_ms <= 0 means no expiry
    # writes
  keyv.delete
    @ (store: keyv_store, key: string) -> result[bool, string]
    + returns true when a key was removed, false when already absent
    # writes
  keyv.clear
    @ (store: keyv_store) -> result[void, string]
    + removes every key in the store's namespace
    # writes
