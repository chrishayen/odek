# Requirement: "a scene editing library for manipulating scenes, images, and video timelines with a tool-based editing model"

The library exposes a scene graph, tool-driven edits, undo history, and media references. No rendering, no UI.

std
  std.math
    std.math.clamp_f64
      @ (v: f64, lo: f64, hi: f64) -> f64
      + returns v clamped to [lo, hi]
      # math
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

scene_editor
  scene_editor.new_scene
    @ () -> scene_state
    + returns an empty scene with no nodes and an empty history
    # construction
  scene_editor.add_node
    @ (state: scene_state, kind: string, parent: optional[string]) -> tuple[string, scene_state]
    + returns (node_id, new_state) after inserting a node of the given kind
    # scene_graph
  scene_editor.remove_node
    @ (state: scene_state, node_id: string) -> result[scene_state, string]
    + removes a node and its descendants
    - returns error when node_id does not exist
    # scene_graph
  scene_editor.set_transform
    @ (state: scene_state, node_id: string, x: f64, y: f64, z: f64) -> result[scene_state, string]
    + sets translation on the node
    - returns error when node_id does not exist
    # scene_graph
    -> std.math.clamp_f64
  scene_editor.attach_media
    @ (state: scene_state, node_id: string, media_id: string) -> result[scene_state, string]
    + attaches a media asset reference to the node
    - returns error when node_id does not exist
    # media
  scene_editor.register_media
    @ (state: scene_state, media_id: string, kind: string, duration_ms: i64) -> scene_state
    + stores metadata for an image or video asset
    # media
  scene_editor.add_keyframe
    @ (state: scene_state, node_id: string, time_ms: i64, property: string, value: f64) -> result[scene_state, string]
    + adds a keyframe for a named property
    - returns error when node_id does not exist
    # animation
  scene_editor.sample_at
    @ (state: scene_state, node_id: string, property: string, time_ms: i64) -> optional[f64]
    + returns the interpolated property value at the given time
    # animation
  scene_editor.push_tool_edit
    @ (state: scene_state, tool: string, payload: string) -> scene_state
    + applies a tool-driven edit and records it in history
    # tools
    -> std.time.now_millis
  scene_editor.undo
    @ (state: scene_state) -> result[scene_state, string]
    + reverts the most recent tool edit
    - returns error when history is empty
    # history
  scene_editor.redo
    @ (state: scene_state) -> result[scene_state, string]
    + reapplies the most recently undone edit
    - returns error when the redo stack is empty
    # history
  scene_editor.snapshot
    @ (state: scene_state) -> scene_snapshot
    + returns a serializable snapshot of the current scene graph
    # persistence
