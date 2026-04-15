# Requirement: "a library for building native desktop applications on a major OS"

Provides a thin object model over windows, menus, and event loops.

std: (all units exist)

nativeapp
  nativeapp.new_app
    fn (name: string) -> app_state
    + creates an application object with the given display name
    # construction
  nativeapp.new_window
    fn (state: app_state, title: string, width: i32, height: i32) -> tuple[app_state, window_id]
    + creates a window of the given size and returns its id
    - returns error via state when dimensions are non-positive
    # window_management
  nativeapp.set_menu
    fn (state: app_state, items: list[menu_item]) -> app_state
    + installs a menu bar built from the supplied items
    # menu
  nativeapp.on_event
    fn (state: app_state, kind: string, handler: callback) -> app_state
    + registers a handler for a named window or menu event
    # events
  nativeapp.run
    fn (state: app_state) -> void
    + enters the event loop and blocks until the app exits
    # lifecycle
