# Requirement: "a multi-agent framework for financial trading research"

Coordinates roles (analyst, researcher, trader, risk) that exchange messages, each backed by a pluggable language model. The framework owns the state machine; the LM and market data are injected.

std
  std.http
    std.http.post_json
      @ (url: string, headers: map[string, string], body: string) -> result[string, string]
      + posts a JSON body and returns the response body
      - returns error on non-2xx status
      # http
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string map as JSON
      # serialization
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

agents
  agents.model_client_new
    @ (endpoint: string, api_key: string) -> model_client
    + creates a language model client pointing at an endpoint
    # construction
  agents.model_complete
    @ (client: model_client, system: string, user: string) -> result[string, string]
    + returns the model's completion text
    - returns error when the endpoint is unreachable
    # model_inference
    -> std.http.post_json
    -> std.json.encode_object
    -> std.json.parse_object
  agents.market_snapshot
    @ (symbol: string, as_of: i64) -> market_snapshot
    + carries price, volume, and headline summaries
    # market_data
  agents.analyst_step
    @ (client: model_client, snapshot: market_snapshot) -> result[string, string]
    + produces a market commentary given the snapshot
    # analyst
    -> agents.model_complete
  agents.researcher_step
    @ (client: model_client, commentary: string, prior_research: list[string]) -> result[string, string]
    + produces research notes building on commentary and history
    # researcher
    -> agents.model_complete
  agents.trader_step
    @ (client: model_client, research: string, position: f64) -> result[trade_decision, string]
    + produces buy/sell/hold decision with size
    - returns error when the model output cannot be parsed
    # trader
    -> agents.model_complete
  agents.risk_check
    @ (decision: trade_decision, position: f64, limit: f64) -> result[trade_decision, string]
    + returns the decision clipped to respect the position limit
    - returns error when the decision would breach a hard stop
    # risk_management
  agents.round_new
    @ (symbol: string) -> round_state
    + creates an empty round for a symbol
    # orchestration
  agents.round_run
    @ (state: round_state, client: model_client, snapshot: market_snapshot, position: f64, limit: f64) -> result[tuple[round_state, trade_decision], string]
    + runs analyst -> researcher -> trader -> risk in sequence
    - returns error at the first failing step
    # orchestration
    -> agents.analyst_step
    -> agents.researcher_step
    -> agents.trader_step
    -> agents.risk_check
    -> std.time.now_seconds
