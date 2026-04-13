# Requirement: "a code generator that produces HTTP server and client stubs from a protobuf service definition"

Parses a protobuf descriptor, then emits matching server handler and client caller code keyed on HTTP routes.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns full file contents as text
      - returns error when path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + creates or overwrites the file
      # filesystem

httpgen
  httpgen.parse_proto
    @ (source: string) -> result[proto_file, string]
    + parses service, method, and message definitions into a proto_file
    - returns error on unexpected tokens or unbalanced braces
    # protobuf_parsing
  httpgen.extract_http_routes
    @ (file: proto_file) -> list[route_spec]
    + reads http option annotations on each method and returns verb + path + request/response types
    ? methods without an http option are skipped
    # route_extraction
  httpgen.render_server_handler
    @ (service: string, routes: list[route_spec]) -> string
    + emits server-side handler code: one decode/dispatch per route
    # codegen
  httpgen.render_client_caller
    @ (service: string, routes: list[route_spec]) -> string
    + emits client-side caller code: one typed function per route
    # codegen
  httpgen.generate
    @ (proto_path: string, out_server: string, out_client: string) -> result[void, string]
    + reads the proto file, parses it, writes the two outputs
    - returns error on parse failure
    # orchestration
    -> std.fs.read_all
    -> std.fs.write_all
