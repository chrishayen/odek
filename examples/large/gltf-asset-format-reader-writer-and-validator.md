# Requirement: "a 3D asset format reader, writer, and validator"

Parses, validates, and serializes a JSON-based 3D scene format with binary buffer references.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads file bytes
      - returns error when missing
      # fs
    std.fs.write_all
      fn (path: string, contents: bytes) -> result[void, string]
      + writes bytes to path
      - returns error on I/O failure
      # fs
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses any JSON value into a generic tree
      - returns error on malformed JSON
      # serialization
    std.json.serialize_value
      fn (value: json_value) -> string
      + serializes a JSON tree back to text
      # serialization
  std.encoding
    std.encoding.base64_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes standard base64
      - returns error on invalid input
      # encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + encodes bytes as standard base64
      # encoding

gltf
  gltf.parse_json
    fn (raw: string) -> result[scene, string]
    + parses a JSON document into a scene graph
    - returns error on schema violations
    # parsing
    -> std.json.parse_value
  gltf.parse_binary
    fn (data: bytes) -> result[scene, string]
    + parses a binary container with embedded JSON chunk and buffer chunk
    - returns error on invalid chunk layout
    - returns error on magic number mismatch
    # parsing
  gltf.load_file
    fn (path: string) -> result[scene, string]
    + loads a scene from disk, choosing JSON or binary by magic bytes
    - returns error on I/O or parse failure
    # ingest
    -> std.fs.read_all
  gltf.validate
    fn (scene: scene) -> result[void, list[string]]
    + returns list of schema violations when any are present
    + checks that indices point to valid nodes, meshes, and accessors
    # validation
  gltf.resolve_buffer
    fn (scene: scene, buffer_index: i32) -> result[bytes, string]
    + returns decoded bytes for the referenced buffer, handling data URIs
    - returns error on index out of range
    - returns error on unsupported URI scheme
    # buffer_resolution
    -> std.encoding.base64_decode
  gltf.serialize_json
    fn (scene: scene) -> result[string, string]
    + serializes a scene back to JSON text
    - returns error when scene references missing resources
    # serialization
    -> std.json.serialize_value
  gltf.serialize_binary
    fn (scene: scene) -> result[bytes, string]
    + serializes a scene to the binary container format
    - returns error when scene references missing resources
    # serialization
    -> std.encoding.base64_encode
  gltf.save_file
    fn (scene: scene, path: string, binary: bool) -> result[void, string]
    + writes the scene to disk in the requested variant
    - returns error on I/O failure
    # persistence
    -> std.fs.write_all
  gltf.list_meshes
    fn (scene: scene) -> list[string]
    + returns the name of every mesh in the scene
    # inspection
  gltf.get_node_transform
    fn (scene: scene, node_index: i32) -> result[matrix4, string]
    + returns the 4x4 world transform for the node
    - returns error when node index is out of range
    # transforms
