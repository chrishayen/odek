# Requirement: "a library for creating HTTP middlewares"

A middleware is a function that wraps a handler. The project exposes composition and a few common wrappers.

std: (all units exist)

middleware
  middleware.chain
    @ (middlewares: list[middleware_handle], final: handler_handle) -> handler_handle
    + composes middlewares in order so the first wraps the second and so on
    + an empty list returns the final handler unchanged
    # composition
  middleware.with_logging
    @ (inner: handler_handle) -> handler_handle
    + wraps a handler to record request method, path, and status
    # logging
  middleware.with_recover
    @ (inner: handler_handle) -> handler_handle
    + wraps a handler to catch panics and return a 500 response
    # resilience
  middleware.with_timeout
    @ (inner: handler_handle, millis: i32) -> handler_handle
    + wraps a handler to abort when it exceeds the given duration
    - returns a 504 response when the deadline elapses
    # timeouts
