# Requirement: "a library to unwrap 3d models into flat papercraft layouts"

Takes a triangle mesh, groups connected coplanar-ish faces, unfolds each group into 2D, and packs them onto pages with cut and fold annotations.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when the file cannot be opened
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes the contents to the given path, replacing any existing file
      # filesystem
  std.math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the square root of a non-negative number
      # math
    std.math.acos
      @ (x: f64) -> f64
      + returns the arc-cosine in radians
      # math
  std.text
    std.text.split_lines
      @ (s: string) -> list[string]
      + splits on newline and drops a trailing empty segment
      # strings

papercraft
  papercraft.parse_obj
    @ (source: string) -> result[mesh, string]
    + returns a mesh of triangles with shared vertex indices parsed from a Wavefront-style text format
    - returns error on malformed vertex or face lines
    # loading
    -> std.text.split_lines
  papercraft.build_adjacency
    @ (m: mesh) -> adjacency
    + returns a map from each edge to the one or two triangles that share it
    # topology
  papercraft.dihedral_angle
    @ (m: mesh, tri_a: i32, tri_b: i32) -> f64
    + returns the angle in radians between the two triangle normals across their shared edge
    # geometry
    -> std.math.acos
  papercraft.find_patches
    @ (m: mesh, adj: adjacency, flat_threshold: f64) -> list[patch]
    + partitions the mesh into patches by growing from seed triangles across edges whose dihedral angle is below the threshold
    # partitioning
  papercraft.unfold_patch
    @ (m: mesh, p: patch) -> unfolded_patch
    + returns an unfolded 2D layout of the patch by rotating each triangle into the plane of its parent across the shared edge
    ? the first triangle of the patch is placed in a canonical orientation
    # unfolding
    -> std.math.sqrt
  papercraft.classify_edges
    @ (patches: list[patch], adj: adjacency) -> edge_classification
    + labels every mesh edge as cut, fold-mountain, or fold-valley based on whether it crosses a patch boundary and the dihedral sign
    # annotation
  papercraft.pack_pages
    @ (unfolded: list[unfolded_patch], page_width: f64, page_height: f64, margin: f64) -> list[page]
    + arranges the unfolded patches onto pages using a shelf-packing layout
    - scales down any patch that exceeds the page interior before packing
    # layout
  papercraft.render_svg
    @ (pages: list[page], classification: edge_classification) -> list[string]
    + returns one SVG document per page with cut and fold edges drawn in distinct styles
    # rendering
  papercraft.unwrap
    @ (source: string, page_width: f64, page_height: f64) -> result[list[string], string]
    + parses a mesh, partitions, unfolds, packs, and returns the rendered SVG pages
    - returns error when the mesh cannot be parsed
    # orchestration
    -> std.fs.read_all
