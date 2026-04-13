# Requirement: "a bindless graphics api"

A low-level graphics abstraction where resources (buffers, textures, samplers) are addressed by stable handles rather than bind slots. Shaders access resources by handle index.

std
  std.collections
    std.collections.slab_new
      @ () -> slab_state
      + returns an empty slab that issues stable generational handles
      # collections
    std.collections.slab_insert
      @ (slab: slab_state, payload: bytes) -> tuple[u64, slab_state]
      + returns (handle, new_slab) where the handle encodes index and generation
      # collections
    std.collections.slab_get
      @ (slab: slab_state, handle: u64) -> result[bytes, string]
      + returns the payload when the generation matches
      - returns error when the handle is stale or out of range
      # collections
    std.collections.slab_remove
      @ (slab: slab_state, handle: u64) -> result[slab_state, string]
      + frees the slot and bumps its generation
      - returns error when the handle is stale
      # collections

gfx
  gfx.device_new
    @ () -> device_state
    + creates an empty device with no resources
    # device
  gfx.buffer_create
    @ (device: device_state, size_bytes: u64, usage: u32) -> tuple[u64, device_state]
    + returns a buffer handle and the updated device
    # resources
    -> std.collections.slab_insert
  gfx.buffer_destroy
    @ (device: device_state, handle: u64) -> result[device_state, string]
    + releases the buffer and invalidates the handle
    - returns error on a stale handle
    # resources
    -> std.collections.slab_remove
  gfx.texture_create
    @ (device: device_state, width: u32, height: u32, format: u32) -> tuple[u64, device_state]
    + returns a texture handle
    # resources
    -> std.collections.slab_insert
  gfx.texture_destroy
    @ (device: device_state, handle: u64) -> result[device_state, string]
    + releases the texture and invalidates the handle
    - returns error on a stale handle
    # resources
    -> std.collections.slab_remove
  gfx.sampler_create
    @ (device: device_state, filter: u32, wrap: u32) -> tuple[u64, device_state]
    + returns a sampler handle
    # resources
    -> std.collections.slab_insert
  gfx.pipeline_create
    @ (device: device_state, vertex_shader: bytes, fragment_shader: bytes) -> result[tuple[u64, device_state], string]
    + returns a pipeline handle when both shader blobs validate
    - returns error on unrecognized shader magic bytes
    # pipeline
  gfx.cmd_buffer_new
    @ (device: device_state) -> cmd_buffer_state
    + returns an empty command buffer bound to the device
    # commands
  gfx.cmd_set_pipeline
    @ (cmd: cmd_buffer_state, pipeline: u64) -> result[cmd_buffer_state, string]
    + records a pipeline-bind command
    - returns error on stale pipeline handle
    # commands
  gfx.cmd_draw_bindless
    @ (cmd: cmd_buffer_state, vertex_count: u32, resource_handles: list[u64]) -> result[cmd_buffer_state, string]
    + records a draw that references resources by handle index
    - returns error when any handle is stale
    # commands
    -> std.collections.slab_get
  gfx.submit
    @ (device: device_state, cmd: cmd_buffer_state) -> result[device_state, string]
    + enqueues the command buffer for execution
    - returns error when the command buffer is empty
    # execution
