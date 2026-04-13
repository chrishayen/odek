# Requirement: "a cross-platform library for building natural-user-interface applications"

A widget tree with touch, gesture, and pointer events, a layout pass, and a backend-agnostic draw list consumable by any platform renderer.

std
  std.math
    std.math.clamp_f32
      @ (v: f32, lo: f32, hi: f32) -> f32
      + returns v constrained to [lo, hi]
      # math
    std.math.hypot_f32
      @ (dx: f32, dy: f32) -> f32
      + returns the Euclidean distance
      # math
  std.time
    std.time.now_millis
      @ () -> i64
      + returns the current monotonic time in milliseconds
      # time

nui
  nui.new_app
    @ () -> app_state
    + creates an empty application with an empty root
    # construction
  nui.set_root
    @ (state: app_state, root: widget) -> app_state
    + installs the root widget
    # construction
  nui.add_widget
    @ (parent: widget, child: widget) -> widget
    + appends a child to a parent widget
    # composition
  nui.layout
    @ (state: app_state, viewport_w: f32, viewport_h: f32) -> layout_tree
    + computes positions and sizes for every widget
    # layout
  nui.render
    @ (tree: layout_tree) -> list[draw_command]
    + emits a flat list of draw commands
    # rendering
  nui.dispatch_pointer
    @ (state: app_state, x: f32, y: f32, phase: pointer_phase) -> app_state
    + routes a pointer event through hit-testing to the target widget
    # events
    -> std.math.clamp_f32
  nui.recognize_gesture
    @ (state: gesture_state, touches: list[touch]) -> tuple[optional[gesture], gesture_state]
    + recognizes tap, swipe, and pinch from a sequence of touches
    + differentiates tap from swipe by distance threshold
    - returns none while a gesture is still being accumulated
    # gestures
    -> std.math.hypot_f32
    -> std.time.now_millis
  nui.animate
    @ (state: app_state, dt_millis: i64) -> app_state
    + advances running animations by dt and applies their outputs
    # animation
  nui.start_animation
    @ (state: app_state, target_id: string, prop: string, to: f32, duration_millis: i32) -> result[app_state, string]
    + schedules an animation toward a target value
    - returns error when target_id is unknown
    # animation
  nui.hit_test
    @ (tree: layout_tree, x: f32, y: f32) -> optional[string]
    + returns the id of the topmost widget containing the point
    # events
  nui.theme
    @ (state: app_state, theme: theme_spec) -> app_state
    + applies a theme to the widget tree
    # theming
