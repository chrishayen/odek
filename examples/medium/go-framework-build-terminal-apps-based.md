# Requirement: "a terminal application framework based on a model-update-view architecture"

An event-loop framework where user code supplies an initial model, an update function, and a view function. The framework coordinates input, model transitions, and output without performing IO directly.

std: (all units exist)

tui
  tui.new_program
    @ (initial_model: bytes) -> program_state
    + creates a program wrapping the user-supplied initial model
    # construction
  tui.send_message
    @ (state: program_state, message: bytes) -> program_state
    + enqueues a message for the next update cycle
    # messaging
  tui.step
    @ (state: program_state, update_fn_id: string) -> tuple[program_state, list[bytes]]
    + drains one message from the queue and applies the registered update function
    + returns any commands emitted by the update
    - returns unchanged state when the queue is empty
    # update_cycle
  tui.render
    @ (state: program_state, view_fn_id: string) -> string
    + invokes the registered view function and returns the rendered frame
    # rendering
  tui.quit
    @ (state: program_state) -> program_state
    + marks the program as finished so the caller can stop the loop
    # lifecycle
  tui.is_done
    @ (state: program_state) -> bool
    + returns true once quit has been requested
    # lifecycle
  tui.handle_key
    @ (state: program_state, key: string) -> program_state
    + translates a key event into a message and enqueues it
    # input
