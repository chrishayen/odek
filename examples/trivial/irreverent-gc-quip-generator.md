# Requirement: "a library that returns irreverent quips about a javascript engine's garbage collector"

A novelty library with a single function: pick a quip deterministically from a seed.

std: (all units exist)

gc_quips
  gc_quips.pick
    fn (seed: u64) -> string
    + returns one of a fixed set of quips selected by seed modulo the list length
    + same seed always returns the same quip
    # quip
