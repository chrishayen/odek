# Requirement: "a library to manage the developer-tools debugging workflow"

Tracks active debugging sessions (target URL, breakpoint list, attach state) and exposes simple transitions for opening, attaching, and closing them.

std: (all units exist)

devtools_workflow
  devtools_workflow.new
    @ () -> workflow_state
    + creates an empty session registry
    # construction
  devtools_workflow.open_session
    @ (state: workflow_state, target_url: string) -> tuple[string, workflow_state]
    + allocates a new session id for the target and returns (id, new_state)
    # sessions
  devtools_workflow.attach
    @ (state: workflow_state, session_id: string) -> result[workflow_state, string]
    + marks the session as attached to the debugger
    - returns error when the session id is unknown
    - returns error when the session is already attached
    # sessions
  devtools_workflow.detach
    @ (state: workflow_state, session_id: string) -> result[workflow_state, string]
    + marks the session as detached but keeps its breakpoints
    - returns error when the session is not attached
    # sessions
  devtools_workflow.add_breakpoint
    @ (state: workflow_state, session_id: string, location: string) -> result[workflow_state, string]
    + records a breakpoint location on the session
    - returns error when the session id is unknown
    # breakpoints
  devtools_workflow.close_session
    @ (state: workflow_state, session_id: string) -> result[workflow_state, string]
    + removes the session and its breakpoints
    - returns error when the session id is unknown
    # sessions
