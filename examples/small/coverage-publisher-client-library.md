# Requirement: "a client library for publishing code coverage reports to a coverage tracking service"

Parses a coverage profile, converts it to the upload payload, and submits it via a pluggable HTTP transport.

std
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

coverage_publisher
  coverage_publisher.parse_profile
    fn (raw: string) -> result[list[file_coverage], string]
    + parses a line-oriented coverage profile into per-file hit counts
    - returns error on malformed lines
    # parsing
  coverage_publisher.build_payload
    fn (token: string, service_name: string, files: list[file_coverage]) -> string
    + assembles the JSON upload payload from the profile and credentials
    # payload
    -> std.json.encode_object
  coverage_publisher.submit
    fn (transport_id: string, endpoint_url: string, payload: string) -> result[string, string]
    + posts the payload to endpoint_url via the transport and returns the server response
    - returns error on non-success status
    # upload
