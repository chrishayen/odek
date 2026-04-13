# Requirement: "an http long-poll server for pub-sub messaging"

Clients poll a topic for events since a cursor; publishers append events that wake pending polls.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

longpoll
  longpoll.new_hub
    @ () -> hub_state
    + returns an empty hub with no topics
    # construction
  longpoll.publish
    @ (state: hub_state, topic: string, payload: bytes) -> hub_state
    + appends an event with a monotonic id and current timestamp
    + creates the topic on first publish
    # publishing
    -> std.time.now_millis
  longpoll.poll
    @ (state: hub_state, topic: string, since_id: i64, max_wait_ms: i32) -> poll_result
    + returns all events with id greater than since_id when any exist
    + returns an empty list and the unchanged cursor when timeout elapses with no events
    - returns an empty list when the topic does not exist
    # polling
    -> std.time.now_millis
  longpoll.expire
    @ (state: hub_state, max_age_ms: i64) -> hub_state
    + drops events older than max_age_ms from every topic
    # retention
    -> std.time.now_millis
