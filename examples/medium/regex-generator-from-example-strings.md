# Requirement: "a library that generates a regular expression from example strings"

Builds a character-level trie from examples and emits an alternation-based regex that matches them.

std: (all units exist)

regex_gen
  regex_gen.new_trie
    fn () -> trie_node
    + returns an empty trie root
    # construction
  regex_gen.insert
    fn (root: trie_node, word: string) -> trie_node
    + inserts word's characters into the trie, marking the terminal node
    # insertion
  regex_gen.build
    fn (samples: list[string]) -> trie_node
    + inserts every sample into a fresh trie
    + returns an empty-accepting trie when samples is empty
    # build
    -> regex_gen.insert
  regex_gen.escape_char
    fn (ch: string) -> string
    + returns a regex-safe form of a single character, escaping metacharacters
    # escaping
  regex_gen.render
    fn (root: trie_node) -> string
    + walks the trie and emits a regex using groups and alternation
    + collapses single-character branches into character classes
    ? output is anchored with ^ and $
    # rendering
    -> regex_gen.escape_char
