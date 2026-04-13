# Requirement: "a curated course list library"

Returns a fixed list of courses. No filtering, no ranking.

std: (all units exist)

courses
  courses.all
    @ () -> list[tuple[string, string]]
    + returns every course as (title, url) pairs
    ? list is hardcoded at compile time
    # catalogue
