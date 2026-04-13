# Requirement: "a client-side UI framework based on a model-update-view architecture"

Pure functional state transitions, a virtual node tree, and a diff-based patcher.

std: (all units exist)

elmish
  elmish.make_app
    @ (init: app_state, update: update_fn, view: view_fn) -> app
    + returns an app bundling the initial state and the two pure functions
    # construction
  elmish.dispatch
    @ (current: app, msg: message) -> app
    + applies update to produce the next state and recomputes the view tree
    + leaves current unchanged when update returns the same state
    # transition
  elmish.vnode_text
    @ (text: string) -> vnode
    + returns a text virtual node
    # virtual_dom
  elmish.vnode_element
    @ (tag: string, attrs: map[string,string], children: list[vnode]) -> vnode
    + returns an element virtual node
    # virtual_dom
  elmish.diff
    @ (old_tree: vnode, new_tree: vnode) -> list[patch]
    + returns the minimal patch list to transform old into new
    + empty list when the trees are identical
    # diffing
  elmish.apply_patches
    @ (target: dom_handle, patches: list[patch]) -> result[void, string]
    + applies each patch to the target dom
    - returns error when a path references a missing node
    # rendering
  elmish.on_event
    @ (node: vnode, event_name: string, message: message) -> vnode
    + returns a new vnode that dispatches message on the event
    # events
  elmish.subscribe
    @ (app: app, sink: message_sink) -> subscription
    + wires the sink to receive every message dispatched by the app
    # subscription
