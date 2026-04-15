# Requirement: "an asynchronous task queue based on distributed message passing"

Producers enqueue jobs; workers subscribe and run them. Broker I/O is abstracted behind std transport primitives.

std
  std.json
    std.json.encode_object
      fn (obj: map[string,string]) -> string
      + encodes a flat string map as JSON
      # serialization
    std.json.parse_object
      fn (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string map
      - returns error on invalid JSON
      # serialization
  std.uuid
    std.uuid.new_v4
      fn () -> string
      + returns a random UUIDv4 string
      # identifiers
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

jobqueue
  jobqueue.app_new
    fn (broker_url: string) -> result[app_state, string]
    + creates an application wired to a message broker
    - returns error when the broker URL is unparseable
    # construction
  jobqueue.register_task
    fn (app: app_state, name: string, handler: task_handler) -> app_state
    + registers a named handler for later dispatch
    # registration
  jobqueue.delay
    fn (app: app_state, name: string, args: map[string,string]) -> result[string, string]
    + schedules a task for asynchronous execution and returns its id
    - returns error when the name is unknown
    # scheduling
    -> std.uuid.new_v4
    -> std.json.encode_object
    -> std.time.now_seconds
  jobqueue.receive
    fn (app: app_state) -> result[optional[task_envelope], string]
    + blocks for up to one poll interval waiting for a task
    # receive
    -> std.json.parse_object
  jobqueue.dispatch
    fn (app: app_state, envelope: task_envelope) -> result[string, string]
    + invokes the registered handler with the decoded args
    - returns error when execution raises
    # dispatch
  jobqueue.ack
    fn (app: app_state, envelope: task_envelope) -> result[void, string]
    + acknowledges successful processing to the broker
    # ack
  jobqueue.nack
    fn (app: app_state, envelope: task_envelope, requeue: bool) -> result[void, string]
    + rejects a task, optionally requeueing for retry
    # ack
  jobqueue.worker_loop
    fn (app: app_state) -> result[void, string]
    + receives, dispatches, and acks until interrupted
    # worker
