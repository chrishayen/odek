# Requirement: "a library that computes distance metrics between sequences"

A selection of string and sequence distance metrics with a uniform interface.

std: (all units exist)

sequence_distance
  sequence_distance.hamming
    fn (a: string, b: string) -> result[i32, string]
    + returns the number of positions at which the strings differ
    - returns error when lengths differ
    # edit_distance
  sequence_distance.levenshtein
    fn (a: string, b: string) -> i32
    + returns the minimum single-character edits needed to transform a into b
    + returns 0 for identical strings
    # edit_distance
  sequence_distance.damerau_levenshtein
    fn (a: string, b: string) -> i32
    + like levenshtein but also counts adjacent transpositions as one edit
    # edit_distance
  sequence_distance.jaro
    fn (a: string, b: string) -> f64
    + returns a similarity score between 0.0 (none) and 1.0 (identical)
    # similarity
  sequence_distance.jaro_winkler
    fn (a: string, b: string) -> f64
    + jaro with a common-prefix bonus up to four characters
    # similarity
    -> sequence_distance.jaro
  sequence_distance.jaccard
    fn (a: list[string], b: list[string]) -> f64
    + returns intersection size divided by union size over unique items
    + returns 1.0 when both inputs are empty
    # set_similarity
  sequence_distance.lcs_length
    fn (a: string, b: string) -> i32
    + returns the length of the longest common subsequence
    # subsequence
