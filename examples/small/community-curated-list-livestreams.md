# Requirement: "a curated livestream directory"

Maintain a catalog of livestream entries with categories and a simple search. No networking; the caller supplies data.

std: (all units exist)

livestream_directory
  livestream_directory.new
    @ () -> directory_state
    + creates an empty directory
    # construction
  livestream_directory.add
    @ (state: directory_state, title: string, url: string, category: string) -> result[directory_state, string]
    + adds an entry with the given title, url, and category
    - returns error when an entry with the same url already exists
    # curation
  livestream_directory.remove
    @ (state: directory_state, url: string) -> directory_state
    + removes the entry with the given url if present
    # curation
  livestream_directory.by_category
    @ (state: directory_state, category: string) -> list[livestream_entry]
    + returns all entries whose category matches exactly
    # query
  livestream_directory.search
    @ (state: directory_state, term: string) -> list[livestream_entry]
    + returns entries whose title contains term as a case-insensitive substring
    # search
