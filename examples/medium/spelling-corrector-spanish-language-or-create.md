# Requirement: "a spelling corrector"

Frequency-weighted edit-distance corrector: learns from a corpus, then ranks candidates by frequency.

std: (all units exist)

speller
  speller.new
    @ () -> speller_state
    + creates an empty corrector with no vocabulary
    # construction
  speller.train
    @ (state: speller_state, text: string) -> void
    + tokenizes the text into lowercase words and increments their frequencies
    ? non-letter characters act as word separators
    # training
  speller.load_vocabulary
    @ (state: speller_state, words: list[string], frequencies: list[i64]) -> result[void, string]
    + replaces the vocabulary with the given parallel arrays of words and counts
    - returns error when the arrays have different lengths
    # training
  speller.known
    @ (state: speller_state, word: string) -> bool
    + returns true when the word is in the vocabulary
    # lookup
  speller.edits1
    @ (word: string) -> list[string]
    + returns all strings at edit distance 1 (deletes, transposes, replaces, inserts)
    # candidate_generation
  speller.candidates
    @ (state: speller_state, word: string) -> list[string]
    + returns known candidates, preferring the word itself, then edits1, then edits2
    + returns the word alone when no known candidate exists
    # candidate_ranking
  speller.correct
    @ (state: speller_state, word: string) -> string
    + returns the highest-frequency known candidate within edit distance 2
    ? deterministic tie-break by lexicographic order
    # correction
