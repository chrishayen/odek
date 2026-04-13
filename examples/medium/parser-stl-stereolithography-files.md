# Requirement: "an STL (stereolithography) file parser"

Supports both the ASCII and binary STL variants; callers receive a uniform list of triangles.

std
  std.strings
    std.strings.starts_with
      @ (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings
  std.binary
    std.binary.read_u32_le
      @ (data: bytes, offset: i32) -> u32
      + reads a little-endian u32 at offset
      # binary
    std.binary.read_f32_le
      @ (data: bytes, offset: i32) -> f32
      + reads a little-endian IEEE-754 f32 at offset
      # binary

stl
  stl.parse
    @ (data: bytes) -> result[stl_mesh, string]
    + auto-detects binary vs ASCII STL and dispatches
    - returns error when neither format is recognized
    # parsing
  stl.parse_ascii
    @ (source: string) -> result[stl_mesh, string]
    + reads "solid ... endsolid" blocks with facet/normal/vertex keywords
    - returns error on malformed tokens
    # parsing
    -> std.strings.starts_with
  stl.parse_binary
    @ (data: bytes) -> result[stl_mesh, string]
    + reads the 80-byte header, u32 triangle count, and 50 bytes per triangle
    - returns error when the file is shorter than the declared triangle count requires
    # parsing
    -> std.binary.read_u32_le
    -> std.binary.read_f32_le
  stl.triangles
    @ (mesh: stl_mesh) -> list[triangle]
    + returns every triangle with its three vertices and normal
    # query
  stl.triangle_count
    @ (mesh: stl_mesh) -> i32
    + returns the number of triangles in the mesh
    # query
  stl.bounding_box
    @ (mesh: stl_mesh) -> tuple[vec3, vec3]
    + returns (min, max) corners of the axis-aligned bounding box
    # geometry
