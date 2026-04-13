# Requirement: "a library returning the host location for a yearly event"

Best-effort interpretation: the requirement is a travel-venue blurb, not a software idea. Treat as a lookup of venue-per-year.

std: (all units exist)

event_venue
  event_venue.lookup
    @ (year: i32) -> result[string, string]
    + returns the host city for a known year
    - returns error when the year has no registered venue
    # lookup
