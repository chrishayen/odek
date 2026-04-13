# Requirement: "a library that exposes firewall rule management as an HTTP-style API"

Request handlers list, add, and remove firewall rules through a pluggable firewall backend. Authentication is delegated to a caller-supplied checker.

std
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a json object into a string map
      - returns error on invalid json
      # serialization
    std.json.encode_value
      @ (value: json_value) -> string
      + encodes a json value to a string
      # serialization

firewall_api
  firewall_api.parse_rule
    @ (raw: string) -> result[firewall_rule, string]
    + parses a json body into a firewall rule
    - returns error when required fields are missing
    - returns error when the port is out of range
    # parsing
    -> std.json.parse_object
  firewall_api.handle_list
    @ (backend: firewall_backend) -> http_response
    + returns a 200 response with every current rule as json
    # handler
    -> std.json.encode_value
  firewall_api.handle_add
    @ (backend: firewall_backend, body: string) -> http_response
    + parses the body, adds the rule via the backend, and returns 201
    - returns 400 on a malformed body
    - returns 409 when the rule already exists
    # handler
  firewall_api.handle_delete
    @ (backend: firewall_backend, rule_id: string) -> http_response
    + removes the rule and returns 204
    - returns 404 when the rule id does not exist
    # handler
  firewall_api.require_auth
    @ (req: http_request, checker: auth_checker, inner: request_handler) -> http_response
    + invokes inner only when checker accepts the request's credentials
    - returns 401 when credentials are missing or rejected
    # middleware
