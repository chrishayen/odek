# Requirement: "a terminal UI bluetooth manager"

Wraps a bluetooth control API and renders scan results, pair/connect state, and device actions to a text UI. The adapter and UI driver are injected so tests can run headless.

std
  std.tui
    std.tui.new_screen
      @ (width: i32, height: i32) -> screen
      + creates an off-screen character buffer
      # tui_primitive
    std.tui.draw_text
      @ (s: screen, row: i32, col: i32, text: string) -> screen
      + writes text starting at the given position, clipping at the edge
      # tui_primitive
    std.tui.draw_box
      @ (s: screen, row: i32, col: i32, w: i32, h: i32, title: string) -> screen
      + draws a single-line-bordered box with a title
      # tui_primitive
    std.tui.render
      @ (s: screen) -> string
      + returns the screen as a newline-delimited string
      # tui_primitive
  std.event
    std.event.decode_key
      @ (raw: bytes) -> optional[key_event]
      + decodes ASCII and common escape sequences to key events
      # input
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

bluetooth_manager
  bluetooth_manager.new
    @ (adapter: bluetooth_adapter) -> manager_state
    + creates a manager bound to the given adapter handle
    # construction
  bluetooth_manager.start_scan
    @ (state: manager_state) -> result[manager_state, string]
    + begins device discovery and records the start time
    - returns error when adapter is powered off
    # scan_control
    -> std.time.now_millis
  bluetooth_manager.stop_scan
    @ (state: manager_state) -> manager_state
    + ends device discovery
    # scan_control
  bluetooth_manager.ingest_device
    @ (state: manager_state, dev: device) -> manager_state
    + adds or updates a device record in the device list
    ? duplicates by address are merged
    # discovery
  bluetooth_manager.pair
    @ (state: manager_state, addr: string) -> result[manager_state, string]
    + marks a device as pairing and forwards the request to the adapter
    - returns error when the address is not in the device list
    # pairing
  bluetooth_manager.connect
    @ (state: manager_state, addr: string) -> result[manager_state, string]
    + marks a paired device as connecting
    - returns error when the device is not yet paired
    # connection
  bluetooth_manager.disconnect
    @ (state: manager_state, addr: string) -> manager_state
    + marks a connected device as disconnected
    # connection
  bluetooth_manager.trust
    @ (state: manager_state, addr: string, trusted: bool) -> manager_state
    + updates the trusted flag of a device
    # trust
  bluetooth_manager.selected_index
    @ (state: manager_state) -> i32
    + returns the currently highlighted row in the device list
    # ui_state
  bluetooth_manager.handle_key
    @ (state: manager_state, key: key_event) -> manager_state
    + moves the selection, toggles scan, or triggers the action under cursor
    # input_handling
    -> std.event.decode_key
  bluetooth_manager.render
    @ (state: manager_state, width: i32, height: i32) -> string
    + returns the rendered TUI frame as a string
    # rendering
    -> std.tui.new_screen
    -> std.tui.draw_box
    -> std.tui.draw_text
    -> std.tui.render
