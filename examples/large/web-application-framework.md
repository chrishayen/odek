# Requirement: "a framework for creating web applications"

A component-based web app framework. Components declare local state and a view function. The framework maintains a virtual tree, computes diffs on updates, and emits DOM patch instructions.

std
  std.collections
    std.collections.list_zip_index
      fn (items: list[string]) -> list[tuple[i32,string]]
      + pairs each element with its zero-based index
      # collections
  std.ids
    std.ids.new_id
      fn () -> string
      + returns a unique opaque id
      # identity
  std.strings
    std.strings.escape_html
      fn (s: string) -> string
      + escapes &, <, >, ", ' for safe embedding in HTML text
      # text

webapp
  webapp.new
    fn (root_component: string) -> app_state
    + creates an app whose root is the named component and whose tree is empty
    # construction
  webapp.mount_component
    fn (state: app_state, parent_id: string, component_type: string, initial: map[string,string]) -> tuple[app_state, string]
    + instantiates a component under parent_id and returns its assigned id
    # component_lifecycle
    -> std.ids.new_id
  webapp.unmount_component
    fn (state: app_state, component_id: string) -> app_state
    + removes the component and all descendants
    # component_lifecycle
  webapp.set_state
    fn (state: app_state, component_id: string, key: string, value: string) -> app_state
    + updates a single state key and marks the component dirty
    - returns unchanged state when component_id is not mounted
    # state_update
  webapp.render_node
    fn (state: app_state, component_id: string) -> vnode
    + calls the component's view to produce a virtual DOM node
    # rendering
  webapp.diff
    fn (before: vnode, after: vnode) -> list[patch]
    + produces a minimal list of patches to transform before into after
    ? children are matched by key when present, else by index
    # diffing
    -> std.collections.list_zip_index
  webapp.apply_patches
    fn (state: app_state, patches: list[patch]) -> app_state
    + applies a patch list to the live tree
    # reconciliation
  webapp.render_to_html
    fn (node: vnode) -> string
    + serializes a vnode to an HTML string for server-side rendering
    # ssr
    -> std.strings.escape_html
  webapp.dispatch_event
    fn (state: app_state, component_id: string, event_name: string, payload: string) -> result[app_state, string]
    + routes a DOM-level event to the component's handler
    - returns error when component_id has no handler for event_name
    # event_handling
  webapp.tick
    fn (state: app_state) -> tuple[app_state, list[patch]]
    + re-renders every dirty component and returns the combined patch list
    # scheduler
