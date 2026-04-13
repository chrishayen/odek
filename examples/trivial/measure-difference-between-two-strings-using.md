# Requirement: "compute Levenshtein distance between two strings"

A single pure function. No std dependencies beyond what the host language provides.

std: (all units exist)

leven
  leven.distance
    @ (a: string, b: string) -> i32
    + returns the minimum number of single-character edits between the two strings
    + returns 0 when the strings are equal
    + returns len(b) when a is empty
    # edit_distance
