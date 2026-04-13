# Requirement: "a personal activity tracking library for media consumption and fitness"

Records users, their activity entries across several categories, and produces summaries. All storage is in-memory state; persistence is the caller's responsibility.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns the hex-encoded SHA-256 digest
      # hashing

activity_tracker
  activity_tracker.new
    @ () -> tracker_state
    + creates an empty tracker with no users or entries
    # construction
  activity_tracker.register_user
    @ (state: tracker_state, username: string, email: string) -> result[tuple[tracker_state, user_id], string]
    + creates a new user and returns their id
    - returns error when the username already exists
    # users
    -> std.hash.sha256_hex
  activity_tracker.lookup_user
    @ (state: tracker_state, username: string) -> optional[user_id]
    + returns the user id for the given username
    - returns none when no user matches
    # users
  activity_tracker.log_media
    @ (state: tracker_state, user: user_id, kind: string, title: string, progress: i32, total: i32) -> result[tuple[tracker_state, entry_id], string]
    + records a media entry (book, show, film, game) with progress
    - returns error when kind is not one of the supported categories
    - returns error when progress is greater than total
    # media
    -> std.time.now_seconds
  activity_tracker.log_workout
    @ (state: tracker_state, user: user_id, name: string, duration_sec: i32, sets: list[workout_set]) -> result[tuple[tracker_state, entry_id], string]
    + records a workout session with an ordered list of sets
    - returns error when duration is negative
    # fitness
    -> std.time.now_seconds
  activity_tracker.log_meal
    @ (state: tracker_state, user: user_id, description: string, calories: i32) -> result[tuple[tracker_state, entry_id], string]
    + records a meal with a calorie estimate
    - returns error when calories is negative
    # nutrition
    -> std.time.now_seconds
  activity_tracker.rate_entry
    @ (state: tracker_state, user: user_id, entry: entry_id, stars: i32) -> result[tracker_state, string]
    + attaches a 1..5 star rating to an entry
    - returns error when stars is outside 1..5
    - returns error when the entry does not belong to the user
    # rating
  activity_tracker.entries_for
    @ (state: tracker_state, user: user_id, kind: string) -> list[entry_summary]
    + returns the user's entries of a given kind, newest first
    - returns empty list when the user has none
    # query
  activity_tracker.daily_summary
    @ (state: tracker_state, user: user_id, day_unix: i64) -> daily_summary
    + aggregates workouts, meals, and media for a given calendar day
    # reports
  activity_tracker.streak
    @ (state: tracker_state, user: user_id, kind: string) -> i32
    + returns the number of consecutive days with at least one entry of the given kind, ending today
    - returns 0 when the most recent entry is not today
    # reports
    -> std.time.now_seconds
  activity_tracker.export_entries
    @ (state: tracker_state, user: user_id) -> list[entry_summary]
    + returns every entry for the user across all kinds
    # export
