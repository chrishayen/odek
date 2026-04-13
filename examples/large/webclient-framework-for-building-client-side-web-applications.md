# Requirement: "a framework for building client-side web applications"

A component-and-hooks framework for the client: components own local state, side effects run after render, and a virtual tree is diffed into DOM patches.

std
  std.ids
    std.ids.new_id
      @ () -> string
      + returns a unique opaque id
      # identity
  std.strings
    std.strings.escape_html
      @ (s: string) -> string
      + escapes HTML-special characters for safe text embedding
      # text
  std.collections
    std.collections.list_equal
      @ (a: list[string], b: list[string]) -> bool
      + returns true when a and b have equal length and equal elements
      # collections

webclient
  webclient.new
    @ (root_type: string) -> app_state
    + creates an app with the given root component type and an empty tree
    # construction
  webclient.mount
    @ (state: app_state, parent_id: string, component_type: string, props: map[string,string]) -> tuple[app_state, string]
    + instantiates a component under parent_id and returns its id
    # component_lifecycle
    -> std.ids.new_id
  webclient.unmount
    @ (state: app_state, component_id: string) -> app_state
    + removes the component and all descendants, running cleanup for each effect
    # component_lifecycle
  webclient.use_state
    @ (state: app_state, component_id: string, slot: i32, initial: string) -> tuple[app_state, string]
    + registers the slot's state and returns its current value
    ? slot identifies the nth use_state call within the component
    # hooks
  webclient.set_state
    @ (state: app_state, component_id: string, slot: i32, value: string) -> app_state
    + updates a state slot and marks the component dirty
    # hooks
  webclient.use_effect
    @ (state: app_state, component_id: string, slot: i32, deps: list[string], effect_id: string) -> app_state
    + registers an effect that runs when deps change from the previous render
    # hooks
    -> std.collections.list_equal
  webclient.render_virtual
    @ (state: app_state, component_id: string) -> vnode
    + calls the component's view to produce a virtual node
    # rendering
  webclient.diff_tree
    @ (before: vnode, after: vnode) -> list[patch]
    + returns a minimal patch list to transform before into after
    # diffing
  webclient.apply_patches
    @ (state: app_state, patches: list[patch]) -> app_state
    + applies patches to the live tree
    # reconciliation
  webclient.flush_effects
    @ (state: app_state) -> tuple[app_state, list[string]]
    + returns the effect_ids whose deps changed since the last flush
    # scheduler
  webclient.dispatch_event
    @ (state: app_state, target_id: string, event_name: string, payload: string) -> result[app_state, string]
    + routes a DOM event to the target component's handler
    - returns error when target_id has no handler for event_name
    # event_handling
  webclient.render_to_string
    @ (root: vnode) -> string
    + serializes a virtual tree to an HTML string
    # ssr
    -> std.strings.escape_html
