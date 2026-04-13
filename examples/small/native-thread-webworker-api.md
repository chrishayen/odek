# Requirement: "a library for a lightweight web worker api using native threads"

Spawns worker threads that run user scripts and exchange messages with the parent.

std
  std.thread
    std.thread.spawn
      @ (fn: thread_entry) -> result[thread_handle, string]
      + starts a new os thread running fn
      - returns error when the thread cannot be created
      # threading
    std.thread.join
      @ (h: thread_handle) -> result[void, string]
      + waits for the thread to finish
      # threading

webworker
  webworker.new
    @ (script: string) -> result[worker, string]
    + spawns a worker that runs script and awaits messages
    - returns error on thread creation failure
    # lifecycle
    -> std.thread.spawn
  webworker.post_message
    @ (w: worker, message: string) -> result[void, string]
    + enqueues message for delivery to the worker
    - returns error when the worker has been terminated
    # messaging
  webworker.receive
    @ (w: worker) -> result[optional[string], string]
    + returns the next message posted by the worker, or none when the queue is empty
    - returns error when the worker has crashed
    # messaging
  webworker.terminate
    @ (w: worker) -> result[void, string]
    + signals the worker to stop and joins its thread
    - returns error on join failure
    # lifecycle
    -> std.thread.join
