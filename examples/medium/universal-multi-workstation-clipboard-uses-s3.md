# Requirement: "a shared clipboard that stores and retrieves clipboard snippets through a remote object-storage backend"

Two surfaces: copy (write a snippet under a key) and paste (read the current snippet). The object store is abstracted as a small put/get interface.

std
  std.object_store
    std.object_store.put_object
      @ (endpoint: string, bucket: string, key: string, data: bytes) -> result[void, string]
      + uploads an object to a remote store
      - returns error on auth or network failure
      # storage
    std.object_store.get_object
      @ (endpoint: string, bucket: string, key: string) -> result[bytes, string]
      + downloads an object from a remote store
      - returns error when the key does not exist
      # storage
    std.object_store.list_objects
      @ (endpoint: string, bucket: string, prefix: string) -> result[list[string], string]
      + returns keys under the given prefix
      # storage
  std.time
    std.time.now_millis
      @ () -> i64
      + returns the current unix time in milliseconds
      # time

clipboard
  clipboard.new
    @ (endpoint: string, bucket: string, user_id: string) -> clipboard_state
    + constructs a clipboard bound to a user namespace in a remote bucket
    # construction
  clipboard.copy
    @ (state: clipboard_state, content: bytes) -> result[string, string]
    + writes the snippet under "{user_id}/latest" and a timestamped history key, returns the history key
    - returns error on storage failure
    # copy
    -> std.time.now_millis
    -> std.object_store.put_object
  clipboard.paste
    @ (state: clipboard_state) -> result[bytes, string]
    + reads the latest snippet
    - returns error when nothing has been copied
    # paste
    -> std.object_store.get_object
  clipboard.history
    @ (state: clipboard_state, limit: i32) -> result[list[string], string]
    + returns the most recent history keys in descending order
    # history
    -> std.object_store.list_objects
  clipboard.paste_at
    @ (state: clipboard_state, history_key: string) -> result[bytes, string]
    + reads a specific historical snippet
    - returns error when the key does not exist
    # paste
    -> std.object_store.get_object
