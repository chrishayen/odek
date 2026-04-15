# Requirement: "convert an observable stream into a single-value future"

The result resolves with the last emitted value (or rejects with the stream's error).

std: (all units exist)

obs_to_future
  obs_to_future.from_observable
    fn (stream: observable) -> future[value, string]
    + resolves with the last value emitted before completion
    + resolves with the single value when exactly one is emitted
    - rejects with the stream's error when the stream errors
    - rejects with "empty stream" when the stream completes without emitting
    ? subscription starts immediately and is cancelled when the future is cancelled
    # adaptation
