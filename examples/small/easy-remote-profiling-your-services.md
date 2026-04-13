# Requirement: "a remote profiling endpoint library for services"

Exposes profiling endpoints that produce CPU, heap, and goroutine-style snapshots over HTTP, behind a shared secret.

std
  std.profile
    std.profile.cpu_snapshot
      @ (duration_seconds: i32) -> result[bytes, string]
      + captures a CPU profile for the given duration and returns the encoded snapshot
      # profiling
    std.profile.heap_snapshot
      @ () -> result[bytes, string]
      + captures a heap profile and returns the encoded snapshot
      # profiling
    std.profile.task_snapshot
      @ () -> result[bytes, string]
      + captures an in-flight task/goroutine snapshot
      # profiling

netbug
  netbug.handle
    @ (secret: string, path: string, query: map[string, string]) -> result[profile_response, string]
    + dispatches to cpu, heap, or task snapshot based on the path
    - returns error when the provided secret does not match
    - returns error for an unknown path
    ? cpu path reads duration_seconds from query, defaulting to 30
    # dispatch
    -> std.profile.cpu_snapshot
    -> std.profile.heap_snapshot
    -> std.profile.task_snapshot
