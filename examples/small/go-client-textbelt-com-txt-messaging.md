# Requirement: "a client for a text-message sending service"

A minimal library that posts a message to a phone number via an HTTP endpoint.

std
  std.http
    std.http.post_form
      @ (url: string, fields: map[string, string]) -> result[string, string]
      + posts form-encoded fields to url and returns the response body
      - returns error on non-2xx status codes
      # http
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization

sms_client
  sms_client.new
    @ (endpoint: string, api_key: string) -> sms_client_state
    + creates a client pointed at endpoint and carrying api_key
    # construction
  sms_client.send
    @ (state: sms_client_state, phone: string, message: string) -> result[string, string]
    + returns a message id when the service accepts the send
    - returns error when the phone number is malformed
    - returns error when the quota is exhausted
    # sending
    -> std.http.post_form
    -> std.json.parse_object
