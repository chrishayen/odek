# Requirement: "a library to notify a team across multiple channels when an application is deployed"

Formats a deploy notification and dispatches it through a caller-supplied sink function. Notification destinations are abstract; the library never embeds a specific service.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

deploy_notify
  deploy_notify.format_message
    @ (app: string, version: string, actor: string) -> string
    + returns a human-readable one-line deploy summary including a timestamp
    # formatting
    -> std.time.now_seconds
  deploy_notify.new_dispatcher
    @ () -> dispatcher_state
    + creates an empty dispatcher with no registered sinks
    # construction
  deploy_notify.register_sink
    @ (state: dispatcher_state, name: string) -> dispatcher_state
    + adds a named notification sink to the dispatcher
    # configuration
  deploy_notify.dispatch
    @ (state: dispatcher_state, message: string) -> list[tuple[string, bool]]
    + returns per-sink delivery results as (sink_name, success)
    ? the actual delivery is performed by the host; this returns intent
    # dispatch
