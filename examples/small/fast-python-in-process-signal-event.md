# Requirement: "an in-process signal and event dispatching system"

A registry maps named signals to ordered receiver identifiers; dispatching returns the list of receivers so the caller can invoke them.

std: (all units exist)

signals
  signals.new
    @ () -> signal_registry
    + creates an empty registry
    # construction
  signals.connect
    @ (registry: signal_registry, signal: string, receiver_id: string) -> signal_registry
    + appends a receiver to a signal's subscriber list
    ? re-connecting the same receiver to the same signal is a no-op
    # subscription
  signals.disconnect
    @ (registry: signal_registry, signal: string, receiver_id: string) -> signal_registry
    + removes a receiver from a signal
    + does nothing when the receiver is not subscribed
    # subscription
  signals.receivers_for
    @ (registry: signal_registry, signal: string) -> list[string]
    + returns receiver ids in the order they were connected
    + returns an empty list when the signal has no subscribers
    # dispatch
  signals.has_receivers
    @ (registry: signal_registry, signal: string) -> bool
    + returns true when at least one receiver is connected
    # dispatch
