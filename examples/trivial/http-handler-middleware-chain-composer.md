# Requirement: "a library for composing http handler middleware into a chain"

Left-to-right composition of middleware wrappers around a base handler.

std: (all units exist)

chain
  chain.compose
    fn (middlewares: list[middleware], base: handler) -> handler
    + returns a handler that invokes middlewares in list order, outermost first
    + returns the base handler unchanged when the list is empty
    # composition
