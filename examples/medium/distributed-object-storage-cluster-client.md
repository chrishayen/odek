# Requirement: "a client library for a distributed object storage cluster"

Speaks to a master node for volume assignment, then uploads and retrieves blobs from the returned volume node over HTTP.

std
  std.http
    std.http.get
      @ (url: string) -> result[http_response, string]
      + performs an HTTP GET and returns status and body
      - returns error on network failure
      # network
    std.http.post_bytes
      @ (url: string, body: bytes, content_type: string) -> result[http_response, string]
      + performs an HTTP POST with a raw byte body
      - returns error on non-2xx response
      # network
    std.http.delete
      @ (url: string) -> result[http_response, string]
      + performs an HTTP DELETE
      # network
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object as a flat string map
      - returns error on invalid JSON
      # serialization

object_store_client
  object_store_client.new
    @ (master_url: string) -> client_state
    + stores the master endpoint used for volume assignments
    # construction
  object_store_client.assign
    @ (state: client_state) -> result[assignment, string]
    + asks the master for a new file id and volume node url
    - returns error when the master responds with a non-ok status
    # assignment
    -> std.http.get
    -> std.json.parse_object
  object_store_client.upload
    @ (state: client_state, data: bytes, content_type: string) -> result[string, string]
    + assigns an id, uploads to the volume node, returns the file id
    - returns error on assignment failure or upload failure
    # upload
    -> std.http.post_bytes
  object_store_client.download
    @ (state: client_state, file_id: string) -> result[bytes, string]
    + resolves the volume for the id and fetches the blob bytes
    - returns error when the id cannot be located
    # download
    -> std.http.get
  object_store_client.delete
    @ (state: client_state, file_id: string) -> result[void, string]
    + resolves the volume and deletes the blob
    - returns error when deletion fails
    # delete
    -> std.http.delete
