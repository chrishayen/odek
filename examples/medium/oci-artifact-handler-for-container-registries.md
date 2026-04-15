# Requirement: "a library for working with OCI artifacts in container registries"

Push and pull arbitrary artifacts described by an OCI manifest. Registry I/O goes through a thin http primitive.

std
  std.http
    std.http.get
      fn (url: string, headers: map[string,string]) -> result[bytes, string]
      + returns body bytes on 2xx
      - returns error on non-2xx or transport failure
      # http
    std.http.put
      fn (url: string, headers: map[string,string], body: bytes) -> result[void, string]
      + sends body and returns void on 2xx
      - returns error on non-2xx
      # http
  std.crypto
    std.crypto.sha256_hex
      fn (data: bytes) -> string
      + returns the hex-encoded sha256 digest
      # cryptography
  std.json
    std.json.encode
      fn (obj: map[string,string]) -> string
      + encodes a flat string map as JSON
      # serialization
    std.json.parse
      fn (raw: string) -> result[map[string,string], string]
      + parses a flat JSON object
      - returns error on malformed input
      # serialization

oci
  oci.make_descriptor
    fn (media_type: string, content: bytes) -> map[string,string]
    + returns a descriptor with mediaType, digest "sha256:<hex>", and size
    # descriptor
    -> std.crypto.sha256_hex
  oci.build_manifest
    fn (config: map[string,string], layers: list[map[string,string]]) -> string
    + returns a JSON manifest referencing the given config and layers
    # manifest_build
    -> std.json.encode
  oci.push_blob
    fn (registry: string, repo: string, content: bytes) -> result[map[string,string], string]
    + uploads content and returns its descriptor
    - returns error when the registry rejects the upload
    # push
    -> std.http.put
    -> std.crypto.sha256_hex
  oci.push_manifest
    fn (registry: string, repo: string, tag: string, manifest: string) -> result[void, string]
    + uploads the manifest under the given tag
    - returns error on non-2xx response
    # push
    -> std.http.put
  oci.pull_manifest
    fn (registry: string, repo: string, reference: string) -> result[string, string]
    + fetches the raw manifest for a tag or digest
    - returns error when the reference is not found
    # pull
    -> std.http.get
  oci.pull_blob
    fn (registry: string, repo: string, digest: string) -> result[bytes, string]
    + fetches a blob and verifies its digest
    - returns error when the computed digest does not match
    # pull
    -> std.http.get
    -> std.crypto.sha256_hex
