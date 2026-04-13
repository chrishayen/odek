# Requirement: "a 2D/3D game engine with a GPU rendering backend"

Scene graph, transform math, mesh and shader resources, and a frame step that updates then draws.

std
  std.math
    std.math.mat4_identity
      @ () -> mat4
      + returns the 4x4 identity matrix
      # math
    std.math.mat4_multiply
      @ (a: mat4, b: mat4) -> mat4
      + returns a * b
      # math
    std.math.mat4_translate
      @ (m: mat4, x: f32, y: f32, z: f32) -> mat4
      + returns m post-multiplied by a translation
      # math
    std.math.mat4_perspective
      @ (fov_rad: f32, aspect: f32, near: f32, far: f32) -> mat4
      + returns a perspective projection matrix
      # math
  std.gpu
    std.gpu.create_shader
      @ (vertex_source: string, fragment_source: string) -> result[shader_handle, string]
      + compiles and links a shader program
      - returns error on compile or link failure with the driver log
      # gpu
    std.gpu.create_mesh
      @ (vertices: list[f32], indices: list[i32]) -> result[mesh_handle, string]
      + uploads a vertex/index buffer pair to the GPU
      - returns error when vertex stride is invalid
      # gpu
    std.gpu.draw_mesh
      @ (mesh: mesh_handle, shader: shader_handle, mvp: mat4) -> void
      + draws a mesh with the given shader and model-view-projection matrix
      # gpu
    std.gpu.clear_frame
      @ (r: f32, g: f32, b: f32) -> void
      + clears the framebuffer to the given color
      # gpu

engine
  engine.new_scene
    @ () -> scene
    + creates an empty scene with no nodes
    # construction
  engine.add_node
    @ (scene: scene, mesh: mesh_handle, shader: shader_handle, position: vec3) -> tuple[scene, node_id]
    + inserts a node and returns the updated scene and its id
    # scene_graph
  engine.set_node_position
    @ (scene: scene, id: node_id, position: vec3) -> result[scene, string]
    + moves a node to a new position
    - returns error when the id does not exist
    # scene_graph
  engine.new_camera
    @ (fov_rad: f32, aspect: f32, near: f32, far: f32) -> camera
    + creates a perspective camera
    # camera
    -> std.math.mat4_perspective
  engine.compute_mvp
    @ (cam: camera, node_pos: vec3) -> mat4
    + returns the model-view-projection matrix for a node
    # rendering
    -> std.math.mat4_identity
    -> std.math.mat4_translate
    -> std.math.mat4_multiply
  engine.render_frame
    @ (scene: scene, cam: camera, clear_color: vec3) -> void
    + clears the frame and draws every node
    # rendering
    -> std.gpu.clear_frame
    -> std.gpu.draw_mesh
