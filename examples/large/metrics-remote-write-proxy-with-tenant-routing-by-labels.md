# Requirement: "a metrics remote-write proxy that routes samples to tenants based on label values"

Accepts an incoming metrics write stream, inspects labels on each sample, and forwards the batch to an upstream sink with a tenant identifier header derived from the labels.

std
  std.io
    std.io.http_post_bytes
      fn (url: string, headers: map[string,string], body: bytes) -> result[i32, string]
      + returns the HTTP status code on completed request
      - returns error on network failure or unreachable host
      # http
  std.encoding
    std.encoding.snappy_decompress
      fn (compressed: bytes) -> result[bytes, string]
      + decodes snappy framed bytes
      - returns error on corrupt input
      # compression
    std.encoding.snappy_compress
      fn (raw: bytes) -> bytes
      + compresses bytes using snappy framing
      # compression
    std.encoding.protobuf_decode_write_request
      fn (raw: bytes) -> result[write_request, string]
      + parses a remote-write request into timeseries with labels and samples
      - returns error on malformed protobuf
      # serialization
    std.encoding.protobuf_encode_write_request
      fn (req: write_request) -> bytes
      + serializes a write request back to protobuf
      # serialization
  std.log
    std.log.warn
      fn (message: string) -> void
      + writes a warning line to the log sink
      # logging

tenant_proxy
  tenant_proxy.new
    fn (tenant_label: string, default_tenant: string, upstream_url: string, header_name: string) -> proxy_config
    + stores routing configuration for later requests
    ? one tenant per incoming batch; mixed batches are split
    # construction
  tenant_proxy.extract_tenant
    fn (series_labels: map[string,string], cfg: proxy_config) -> string
    + returns the value of the configured tenant label when present
    - returns the default tenant when the label is missing or empty
    # routing
  tenant_proxy.group_by_tenant
    fn (req: write_request, cfg: proxy_config) -> map[string, write_request]
    + partitions timeseries into one sub-request per tenant value
    + preserves sample ordering within each group
    # partitioning
  tenant_proxy.handle_request
    fn (body: bytes, cfg: proxy_config) -> result[void, string]
    + decompresses, parses, groups, and forwards each group with the tenant header
    - returns error when any upstream POST returns a non-2xx status
    - returns error when the body is not a valid compressed write request
    # request_handling
    -> std.encoding.snappy_decompress
    -> std.encoding.protobuf_decode_write_request
    -> std.encoding.protobuf_encode_write_request
    -> std.encoding.snappy_compress
    -> std.io.http_post_bytes
    -> std.log.warn
