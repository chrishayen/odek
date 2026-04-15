# Requirement: "a 2D triangulation library that converts lines and polygons into triangle meshes suitable for a GPU pipeline"

Two core conversions: triangulating a polygon into filled triangles, and expanding a polyline into a stroked triangle strip. Output is a flat vertex buffer.

std: (all units exist)

triangulate
  triangulate.polygon_to_triangles
    fn (ring: list[f64]) -> result[list[i32], string]
    + returns an index list where each consecutive triple of indices forms a triangle
    + handles convex and simple concave polygons via ear clipping
    - returns error when the ring has fewer than three vertices
    - returns error when the ring self-intersects
    ? input is a flat list of alternating x,y coordinates
    # polygon_triangulation
  triangulate.polygon_with_holes
    fn (outer: list[f64], holes: list[list[f64]]) -> result[list[i32], string]
    + triangulates an outer ring with zero or more hole rings
    + hole rings are merged into the outer ring via bridge edges before ear clipping
    - returns error when any hole is not contained in the outer ring
    # polygon_triangulation
    -> triangulate.polygon_to_triangles
  triangulate.line_to_strip
    fn (points: list[f64], width: f64) -> list[f64]
    + returns a flat vertex buffer for a triangle strip representing a stroked polyline
    + each input point contributes two offset vertices perpendicular to the segment
    ? miter joins are used; degenerate points produce zero-length segments
    # stroke
  triangulate.line_to_indexed_mesh
    fn (points: list[f64], width: f64) -> tuple[list[f64], list[i32]]
    + returns (vertices, indices) describing the stroked polyline as indexed triangles
    + sharp corners are clamped to a bevel when the miter length exceeds twice the width
    # stroke
    -> triangulate.line_to_strip
