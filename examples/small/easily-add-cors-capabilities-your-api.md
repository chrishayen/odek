# Requirement: "a CORS header builder for HTTP APIs"

Given a request's origin and a policy config, produce the response headers the server should send. No HTTP server coupling.

std: (all units exist)

cors
  cors.new_policy
    @ (allowed_origins: list[string], allowed_methods: list[string], allow_credentials: bool, max_age_seconds: i32) -> cors_policy
    + creates a policy; "*" in allowed_origins means any origin
    ? empty allowed_methods defaults to GET, HEAD, POST
    # configuration
  cors.preflight_headers
    @ (policy: cors_policy, origin: string, requested_method: string) -> result[map[string, string], string]
    + returns Access-Control-Allow-* headers when the origin and method are allowed
    - returns error when the origin is not in the allow list
    - returns error when the method is not allowed
    # preflight
  cors.response_headers
    @ (policy: cors_policy, origin: string) -> map[string, string]
    + returns Access-Control-Allow-Origin and related headers for a simple request
    + returns an empty map when the origin is not allowed
    # simple_response
