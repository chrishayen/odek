# Requirement: "a client for a cryptocurrency market data service"

Typed wrappers around a small set of HTTP endpoints that return ticker, coin, and exchange information.

std
  std.http
    std.http.get_json
      @ (url: string) -> result[string, string]
      + fetches a URL and returns the response body
      - returns error on non-2xx status codes
      # http
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.parse_array_of_objects
      @ (raw: string) -> result[list[map[string, string]], string]
      + parses a JSON array of objects into a list of maps
      - returns error when the root is not an array
      # serialization

market_client
  market_client.new
    @ (base_url: string) -> market_client_state
    + creates a client pointed at the given base URL
    # construction
  market_client.get_ticker
    @ (state: market_client_state, coin_id: string) -> result[map[string, string], string]
    + returns ticker fields (price, volume, change) for the given coin id
    - returns error when the coin id is unknown
    # ticker
    -> std.http.get_json
    -> std.json.parse_object
  market_client.list_coins
    @ (state: market_client_state) -> result[list[map[string, string]], string]
    + returns a list of coin summary objects
    - returns error when the response cannot be parsed
    # listing
    -> std.http.get_json
    -> std.json.parse_array_of_objects
  market_client.list_exchanges
    @ (state: market_client_state) -> result[list[map[string, string]], string]
    + returns a list of exchange summary objects
    # listing
    -> std.http.get_json
    -> std.json.parse_array_of_objects
