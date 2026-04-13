# Requirement: "a library to simulate bad network connections"

Given a packet and a network-condition profile, the library decides whether to drop, delay, duplicate, or pass the packet through. Randomness and time both come through thin std primitives.

std
  std.random
    std.random.next_unit_f64
      @ () -> f64
      + returns a pseudo-random number in [0.0, 1.0)
      # randomness
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

netsim
  netsim.new_profile
    @ (drop_rate: f64, duplicate_rate: f64, min_delay_millis: i32, max_delay_millis: i32) -> result[profile, string]
    + returns a profile when rates are in [0,1] and delays are non-negative with min <= max
    - returns error when any rate is out of range
    - returns error when max_delay_millis is less than min_delay_millis
    # construction
  netsim.apply
    @ (p: profile, packet: bytes) -> list[scheduled_packet]
    + returns an empty list when the packet is dropped
    + returns one scheduled_packet with a release time when passed through
    + returns two scheduled_packets when the packet is duplicated
    # simulation
    -> std.random.next_unit_f64
    -> std.time.now_millis
