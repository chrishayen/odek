# Requirement: "a library for building server-rendered web applications with hypermedia-driven interactivity"

Wires routes to handlers that return HTML fragments and supports partial page swaps.

std
  std.strings
    std.strings.concat
      @ (parts: list[string]) -> string
      + joins a list of strings efficiently
      # strings

hypermedia
  hypermedia.new
    @ () -> server_state
    + creates an empty server with no routes
    # construction
  hypermedia.route
    @ (state: server_state, method: string, path: string, handler: callback) -> server_state
    + registers a handler that returns an html fragment for a method and path
    # routing
  hypermedia.render_page
    @ (title: string, head: string, body: string) -> string
    + wraps a body fragment in a full html document
    # rendering
    -> std.strings.concat
  hypermedia.render_fragment
    @ (tag: string, attrs: map[string, string], children: list[string]) -> string
    + emits a single html element with attributes and children
    # rendering
    -> std.strings.concat
  hypermedia.handle
    @ (state: server_state, method: string, path: string, headers: map[string, string]) -> result[hypermedia_response, string]
    + dispatches to the matching handler and chooses full-page or partial response based on headers
    - returns error when no route matches
    # dispatch
  hypermedia.trigger
    @ (name: string, detail: map[string, string]) -> string
    + produces a response header that instructs the client to fire a named event
    # interactivity
