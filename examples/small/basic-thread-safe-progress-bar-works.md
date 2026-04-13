# Requirement: "a thread-safe progress bar library"

The bar holds mutable counter state behind a lock. Rendering returns a string; the caller decides where to print it.

std
  std.sync
    std.sync.new_mutex
      @ () -> mutex_handle
      + returns a fresh unlocked mutex
      # concurrency
    std.sync.with_lock
      @ (m: mutex_handle, f: fn() -> void) -> void
      + runs f while holding m, releasing on return
      # concurrency

progress_bar
  progress_bar.new
    @ (total: i64, width: i32) -> progress_bar_state
    + creates a bar with the given total count and display width
    ? total of zero is allowed; render will show a full bar immediately
    # construction
    -> std.sync.new_mutex
  progress_bar.add
    @ (bar: progress_bar_state, delta: i64) -> void
    + atomically increments the current count by delta, saturating at total
    - negative delta leaves the count unchanged
    # progress_update
    -> std.sync.with_lock
  progress_bar.render
    @ (bar: progress_bar_state) -> string
    + returns a string like "[####------] 40%" reflecting the current count
    # rendering
    -> std.sync.with_lock
