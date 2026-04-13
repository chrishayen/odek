# Requirement: "a utility that converts a callback-style function into a promise-returning function"

Wraps a function whose last argument is a node-style `(err, result)` callback so that calling it instead returns a promise.

std: (all units exist)

promisify
  promisify.wrap
    @ (fn: callable) -> callable
    + returns a new function that, when called, invokes fn and resolves the promise with the callback's result
    + rejects the promise when the callback receives a non-null error
    ? the wrapped function's last parameter must be an error-first callback
    # promisification
