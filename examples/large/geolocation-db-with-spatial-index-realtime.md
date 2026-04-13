# Requirement: "a geolocation database with spatial index and realtime geofencing"

In-memory store of named objects at lat/lon coordinates, with radius queries and persistent geofences that fire enter/exit events on position updates.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.math
    std.math.sin
      @ (x: f64) -> f64
      + returns sine of x radians
      # math
    std.math.cos
      @ (x: f64) -> f64
      + returns cosine of x radians
      # math
    std.math.atan2
      @ (y: f64, x: f64) -> f64
      + returns the angle of the point (x, y) in radians
      # math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns non-negative square root
      # math

geodb
  geodb.new
    @ () -> geodb_state
    + creates an empty geolocation database
    # construction
  geodb.haversine_meters
    @ (lat_a: f64, lon_a: f64, lat_b: f64, lon_b: f64) -> f64
    + returns great-circle distance between two lat/lon points in meters
    ? uses the haversine formula with earth radius 6371008.8 m
    # distance
    -> std.math.sin
    -> std.math.cos
    -> std.math.atan2
    -> std.math.sqrt
  geodb.set_position
    @ (state: geodb_state, object_id: string, lat: f64, lon: f64) -> tuple[list[string], geodb_state]
    + stores or updates the object's position and returns the list of geofence events triggered
    + event strings are formatted as "enter:<fence_id>:<object_id>" or "exit:<fence_id>:<object_id>"
    # upsert
    -> std.time.now_millis
  geodb.delete
    @ (state: geodb_state, object_id: string) -> geodb_state
    + removes the object and any membership in geofences
    ? no-op when the object is unknown
    # delete
  geodb.get_position
    @ (state: geodb_state, object_id: string) -> result[tuple[f64, f64], string]
    + returns the object's (lat, lon)
    - returns error when the object is unknown
    # read
  geodb.nearby
    @ (state: geodb_state, lat: f64, lon: f64, radius_meters: f64) -> list[string]
    + returns ids of objects within the given radius, sorted ascending by distance
    # radius_query
  geodb.add_geofence
    @ (state: geodb_state, fence_id: string, lat: f64, lon: f64, radius_meters: f64) -> geodb_state
    + registers a circular geofence that will emit events on subsequent position updates
    # geofence_registration
  geodb.remove_geofence
    @ (state: geodb_state, fence_id: string) -> geodb_state
    + unregisters a geofence
    ? no-op when the fence is unknown
    # geofence_registration
  geodb.objects_inside
    @ (state: geodb_state, fence_id: string) -> result[list[string], string]
    + returns ids of objects currently inside the given geofence
    - returns error when the fence is unknown
    # query
  geodb.list_fences
    @ (state: geodb_state) -> list[string]
    + returns all registered geofence ids
    # introspection
