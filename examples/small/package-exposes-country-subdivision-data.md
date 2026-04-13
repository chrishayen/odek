# Requirement: "a country and subdivision reference data library"

Reads ISO country records and their subdivisions from a preloaded dataset and exposes lookups by code and name.

std: (all units exist)

countries
  countries.load
    @ (raw: string) -> result[countries_db, string]
    + parses a JSON-encoded country dataset into an in-memory database
    - returns error on malformed JSON
    # loading
  countries.find_by_alpha2
    @ (db: countries_db, code: string) -> optional[country]
    + looks up a country by its ISO 3166-1 alpha-2 code
    ? lookup is case-insensitive
    # lookup
  countries.find_by_name
    @ (db: countries_db, name: string) -> optional[country]
    + looks up a country by its common or official name
    # lookup
  countries.subdivisions_of
    @ (db: countries_db, code: string) -> list[subdivision]
    + returns all subdivisions for the given alpha-2 country code
    + returns an empty list when the country is unknown
    # lookup
  countries.list_all
    @ (db: countries_db) -> list[country]
    + returns every country in the database
    # listing
