# Requirement: "a script-based alerting manager"

Alerts are defined as scripts evaluated on a schedule; when a script emits alerts, they fan out to pluggable notification channels.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.script
    std.script.eval
      @ (source: string, inputs: map[string,string]) -> result[map[string,string], string]
      + evaluates a small scripting expression with named inputs and returns emitted outputs
      - returns error on syntax or runtime failures
      # scripting
  std.http
    std.http.post_json
      @ (url: string, body: string) -> result[i32, string]
      + posts a JSON body to a URL and returns the status code
      - returns error on network failure
      # http

alerting
  alerting.new_manager
    @ () -> manager_state
    + creates an empty manager with no rules and no channels
    # construction
  alerting.register_channel
    @ (state: manager_state, name: string, kind: string, url: string) -> manager_state
    + registers a named notification channel with a sink kind and endpoint
    # channels
  alerting.add_rule
    @ (state: manager_state, name: string, script: string, interval_sec: i32, channels: list[string]) -> result[manager_state, string]
    + adds a scheduled rule that evaluates script every interval_sec and routes alerts to named channels
    - returns error when referenced channels are not registered
    # rule_registration
  alerting.tick
    @ (state: manager_state, now: i64) -> tuple[manager_state, list[string]]
    + runs any rule whose next_run <= now and returns the ids of fired alerts
    # scheduling
    -> std.script.eval
  alerting.evaluate_rule
    @ (state: manager_state, rule_name: string, inputs: map[string,string]) -> result[list[alert_event], string]
    + evaluates a single rule with provided inputs and returns the produced alerts without dispatching
    - returns error when the rule name is unknown
    # evaluation
    -> std.script.eval
  alerting.dispatch
    @ (state: manager_state, alerts: list[alert_event]) -> list[dispatch_result]
    + sends each alert to every channel named on its rule and returns per-channel results
    # dispatch
    -> std.http.post_json
  alerting.list_rules
    @ (state: manager_state) -> list[string]
    + returns the names of all registered rules
    # introspection
  alerting.remove_rule
    @ (state: manager_state, name: string) -> manager_state
    + removes a rule by name; no-op if it does not exist
    # rule_management
  alerting.history
    @ (state: manager_state, rule_name: string, limit: i32) -> list[alert_event]
    + returns the most recent alerts fired by a rule, newest first
    # history
  alerting.ack
    @ (state: manager_state, alert_id: string) -> manager_state
    + marks an alert as acknowledged so suppression rules can skip it
    # acknowledgement
