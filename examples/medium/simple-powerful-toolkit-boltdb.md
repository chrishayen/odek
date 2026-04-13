# Requirement: "a toolkit for an embedded key-value store with typed buckets"

A thin object-mapper over a bucketed key-value store: save/load typed records by key, scan, and delete.

std
  std.kv
    std.kv.open
      @ (path: string) -> result[kv_handle, string]
      + opens or creates the store file at the given path
      - returns error when the path cannot be opened
      # storage
    std.kv.close
      @ (handle: kv_handle) -> void
      + flushes and closes the store
      # storage
    std.kv.put
      @ (handle: kv_handle, bucket: string, key: bytes, value: bytes) -> result[void, string]
      + writes the key-value pair inside the given bucket, creating the bucket if needed
      # storage
    std.kv.get
      @ (handle: kv_handle, bucket: string, key: bytes) -> result[optional[bytes], string]
      + returns the stored bytes or none when absent
      # storage
    std.kv.delete
      @ (handle: kv_handle, bucket: string, key: bytes) -> result[void, string]
      + removes the key from the bucket; no error when absent
      # storage
    std.kv.scan_prefix
      @ (handle: kv_handle, bucket: string, prefix: bytes) -> result[list[tuple[bytes, bytes]], string]
      + returns all entries whose key starts with prefix, in key order
      # storage
  std.json
    std.json.encode
      @ (value: map[string, string]) -> string
      + encodes a string map as a json object
      # serialization
    std.json.decode
      @ (raw: string) -> result[map[string, string], string]
      + parses a json object into a string map
      - returns error on malformed input
      # serialization

toolkit
  toolkit.save
    @ (handle: kv_handle, bucket: string, key: string, fields: map[string, string]) -> result[void, string]
    + serializes the record and stores it under the key
    - returns error when the underlying store rejects the write
    # persistence
    -> std.json.encode
    -> std.kv.put
  toolkit.load
    @ (handle: kv_handle, bucket: string, key: string) -> result[optional[map[string, string]], string]
    + returns the deserialized record or none when missing
    - returns error when stored bytes are not valid json
    # persistence
    -> std.kv.get
    -> std.json.decode
  toolkit.delete
    @ (handle: kv_handle, bucket: string, key: string) -> result[void, string]
    + removes the record
    # persistence
    -> std.kv.delete
  toolkit.find_prefix
    @ (handle: kv_handle, bucket: string, key_prefix: string) -> result[list[tuple[string, map[string, string]]], string]
    + returns all records whose keys share the prefix, in key order
    - returns error when any stored value fails to decode
    # queries
    -> std.kv.scan_prefix
    -> std.json.decode
