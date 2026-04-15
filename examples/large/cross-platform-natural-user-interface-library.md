# Requirement: "a cross-platform library for building natural-user-interface applications"

A widget tree with touch, gesture, and pointer events, a layout pass, and a backend-agnostic draw list consumable by any platform renderer.

std
  std.math
    std.math.clamp_f32
      fn (v: f32, lo: f32, hi: f32) -> f32
      + returns v constrained to [lo, hi]
      # math
    std.math.hypot_f32
      fn (dx: f32, dy: f32) -> f32
      + returns the Euclidean distance
      # math
  std.time
    std.time.now_millis
      fn () -> i64
      + returns the current monotonic time in milliseconds
      # time

nui
  nui.new_app
    fn () -> app_state
    + creates an empty application with an empty root
    # construction
  nui.set_root
    fn (state: app_state, root: widget) -> app_state
    + installs the root widget
    # construction
  nui.add_widget
    fn (parent: widget, child: widget) -> widget
    + appends a child to a parent widget
    # composition
  nui.layout
    fn (state: app_state, viewport_w: f32, viewport_h: f32) -> layout_tree
    + computes positions and sizes for every widget
    # layout
  nui.render
    fn (tree: layout_tree) -> list[draw_command]
    + emits a flat list of draw commands
    # rendering
  nui.dispatch_pointer
    fn (state: app_state, x: f32, y: f32, phase: pointer_phase) -> app_state
    + routes a pointer event through hit-testing to the target widget
    # events
    -> std.math.clamp_f32
  nui.recognize_gesture
    fn (state: gesture_state, touches: list[touch]) -> tuple[optional[gesture], gesture_state]
    + recognizes tap, swipe, and pinch from a sequence of touches
    + differentiates tap from swipe by distance threshold
    - returns none while a gesture is still being accumulated
    # gestures
    -> std.math.hypot_f32
    -> std.time.now_millis
  nui.animate
    fn (state: app_state, dt_millis: i64) -> app_state
    + advances running animations by dt and applies their outputs
    # animation
  nui.start_animation
    fn (state: app_state, target_id: string, prop: string, to: f32, duration_millis: i32) -> result[app_state, string]
    + schedules an animation toward a target value
    - returns error when target_id is unknown
    # animation
  nui.hit_test
    fn (tree: layout_tree, x: f32, y: f32) -> optional[string]
    + returns the id of the topmost widget containing the point
    # events
  nui.theme
    fn (state: app_state, theme: theme_spec) -> app_state
    + applies a theme to the widget tree
    # theming
