# Requirement: "a library that maps programming concepts to short illustrative descriptions"

A lookup from a concept name to its description. One function to register, one to fetch.

std: (all units exist)

concepts
  concepts.new
    fn () -> concept_book
    + creates an empty concept book
    # construction
  concepts.add
    fn (book: concept_book, name: string, illustration: string) -> concept_book
    + stores an illustration keyed by concept name
    # store
  concepts.lookup
    fn (book: concept_book, name: string) -> optional[string]
    + returns the illustration for a concept if known
    - returns none for unknown concepts
    # lookup
