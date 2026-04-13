# Requirement: "linear-time suffix array construction"

Builds a suffix array with the SA-IS algorithm. The project layer is classification, bucket placement, LMS sorting, and the recursive driver.

std: (all units exist)

suffix_array
  suffix_array.build
    @ (text: bytes) -> list[i32]
    + returns the suffix array of text in O(n) time
    + returns empty list for empty input
    ? input is treated as bytes with an implicit sentinel smaller than all characters
    # entry_point
  suffix_array.classify_types
    @ (text: bytes) -> list[bool]
    + returns a bool per position: true for S-type, false for L-type
    ? last position is S-type by definition
    # sais_classification
  suffix_array.find_lms_positions
    @ (types: list[bool]) -> list[i32]
    + returns positions that are S-type and whose predecessor is L-type
    # sais_lms
  suffix_array.bucket_sizes
    @ (text: bytes, alphabet_size: i32) -> list[i32]
    + returns the count of occurrences per symbol
    # sais_buckets
  suffix_array.induced_sort
    @ (text: bytes, types: list[bool], lms: list[i32], buckets: list[i32]) -> list[i32]
    + places LMS suffixes then induces L-type then S-type suffixes into the suffix array
    # sais_induction
  suffix_array.name_lms_substrings
    @ (text: bytes, sa: list[i32], lms: list[i32], types: list[bool]) -> tuple[list[i32], i32]
    + assigns a name to each LMS substring and returns (reduced_text, distinct_name_count)
    # sais_naming
