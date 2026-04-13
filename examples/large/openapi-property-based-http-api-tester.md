# Requirement: "a property-based testing library for HTTP APIs described by an OpenAPI specification"

Reads an API schema, generates randomized requests that conform to the declared parameter and body schemas, executes them against a target host, and reports responses that violate the declared response contract.

std
  std.http
    std.http.send_request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs the HTTP request and returns status, headers, and body
      - returns error on transport failure
      # http_client
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses any JSON value (object, array, string, number, bool, null)
      - returns error on malformed JSON
      # serialization
    std.json.encode_value
      @ (value: json_value) -> string
      + renders a JSON value as a compact string
      # serialization
  std.random
    std.random.next_u64
      @ () -> u64
      + returns a uniformly random 64-bit integer
      # randomness

api_fuzzer
  api_fuzzer.load_spec
    @ (raw: string) -> result[api_spec, string]
    + parses an OpenAPI document (JSON) into an api_spec value
    - returns error when required top-level fields are missing
    # spec_loading
    -> std.json.parse_value
  api_fuzzer.list_operations
    @ (spec: api_spec) -> list[operation]
    + returns every path+method pair declared in the spec
    # introspection
  api_fuzzer.generate_value
    @ (schema: json_value, seed: u64) -> json_value
    + returns a random value satisfying the JSON Schema (type, enum, min, max, required)
    ? supports string, integer, number, boolean, array, object
    # generation
    -> std.random.next_u64
  api_fuzzer.build_request
    @ (op: operation, spec: api_spec, seed: u64) -> http_request
    + fills path parameters, query parameters, headers, and request body from the operation's schemas
    # generation
    -> api_fuzzer.generate_value
  api_fuzzer.check_response
    @ (op: operation, response: http_response) -> result[void, string]
    + returns ok when the response status is declared and the body matches the declared response schema
    - returns error naming the first contract violation
    # verification
    -> std.json.parse_value
  api_fuzzer.run_case
    @ (spec: api_spec, base_url: string, op: operation, seed: u64) -> case_result
    + builds a request, sends it, and records whether the response satisfied the contract
    # execution
    -> api_fuzzer.build_request
    -> std.http.send_request
    -> api_fuzzer.check_response
  api_fuzzer.run_campaign
    @ (spec: api_spec, base_url: string, cases_per_op: i32, seed: u64) -> list[case_result]
    + runs cases_per_op random cases against every operation in the spec
    # orchestration
    -> api_fuzzer.list_operations
    -> api_fuzzer.run_case
  api_fuzzer.format_report
    @ (results: list[case_result]) -> string
    + returns a human-readable summary with counts of passed, failed, and errored cases
    + lists each failure with method, path, and violation description
    # reporting
