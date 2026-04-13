# Requirement: "a client for an image recognition service"

A thin typed layer over HTTP endpoints that accept image bytes or URLs and return tag predictions.

std
  std.http
    std.http.post_json
      @ (url: string, headers: map[string, string], body: string) -> result[string, string]
      + posts a JSON body to url and returns the response body
      - returns error on non-2xx status codes
      # http
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization

image_client
  image_client.new
    @ (base_url: string, api_key: string) -> image_client_state
    + creates a client that attaches the api key to each request
    # construction
  image_client.predict_url
    @ (state: image_client_state, image_url: string) -> result[list[tuple[string, f64]], string]
    + returns (tag, confidence) pairs ranked by confidence
    - returns error when image_url is unreachable
    # prediction
    -> std.http.post_json
    -> std.json.parse_object
  image_client.predict_bytes
    @ (state: image_client_state, image: bytes) -> result[list[tuple[string, f64]], string]
    + returns (tag, confidence) pairs for raw image bytes
    - returns error when the image is not decodable
    # prediction
    -> std.http.post_json
    -> std.json.parse_object
