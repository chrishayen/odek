# Requirement: "a datastore testing helper for seeding, snapshotting, and verifying store contents"

Supports loading fixtures, comparing expected vs actual rows, and resetting state between tests.

std: (all units exist)

dsunit
  dsunit.load_fixture
    @ (source: string) -> result[fixture, string]
    + parses a fixture document of table-to-rows into an in-memory value
    - returns error on malformed documents
    # fixtures
  dsunit.seed
    @ (conn: connection, fx: fixture) -> result[i32, string]
    + inserts fixture rows and returns the count
    - returns error when a target table does not exist
    # seeding
  dsunit.snapshot
    @ (conn: connection, tables: list[string]) -> result[fixture, string]
    + returns a fixture capturing the current rows of the named tables
    # snapshot
  dsunit.diff
    @ (expected: fixture, actual: fixture) -> fixture_diff
    + returns added, missing, and changed rows keyed by table
    + empty diff when fixtures match
    # comparison
  dsunit.truncate_all
    @ (conn: connection, tables: list[string]) -> result[void, string]
    + empties each named table
    - returns error when a table does not exist
    # reset
  dsunit.assert_equal
    @ (expected: fixture, actual: fixture) -> result[void, string]
    - returns error with a rendered diff when fixtures differ
    # assertion
    -> dsunit.diff
