# Requirement: "a typo-correction suggester that is aware of physical keyboard layouts"

Given a dictionary and a possibly misspelled input, returns ranked suggestions where adjacency on a configurable keyboard grid reduces the substitution cost.

std: (all units exist)

keyboard_suggest
  keyboard_suggest.load_layout
    fn (rows: list[string]) -> layout
    + builds a layout from row strings where each character's position determines its physical coordinates
    ? row 0 is the top row
    # layout
  keyboard_suggest.key_distance
    fn (layout: layout, a: string, b: string) -> f32
    + returns the Euclidean distance between two keys on the layout
    + returns 0.0 when a equals b
    # distance
  keyboard_suggest.weighted_edit_distance
    fn (layout: layout, from: string, to: string) -> f32
    + returns a Damerau-Levenshtein-style distance where substitution cost scales with key_distance
    + returns 0.0 for identical strings
    # distance
  keyboard_suggest.new_dictionary
    fn (words: list[string]) -> dictionary
    + builds a dictionary from a word list
    # construction
  keyboard_suggest.suggest
    fn (layout: layout, dict: dictionary, input: string, max_results: i32) -> list[tuple[string, f32]]
    + returns up to max_results words with the lowest weighted edit distance to input, sorted ascending
    + returns an empty list when the dictionary is empty
    # suggestion
