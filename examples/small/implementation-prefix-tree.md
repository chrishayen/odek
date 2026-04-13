# Requirement: "an implementation of a prefix tree"

Classic trie with insert, lookup, and prefix scan.

std: (all units exist)

ptrie
  ptrie.new
    @ () -> trie_state
    + returns an empty trie
    # construction
  ptrie.insert
    @ (t: trie_state, key: string, value: string) -> trie_state
    + returns a new trie with key associated to value
    ? inserting an existing key overwrites its value
    # insert
  ptrie.lookup
    @ (t: trie_state, key: string) -> result[string, string]
    + returns the value associated with an exact key
    - returns error when the key is not present
    # lookup
  ptrie.has_prefix
    @ (t: trie_state, prefix: string) -> bool
    + returns true when at least one key in the trie begins with prefix
    # lookup
  ptrie.keys_with_prefix
    @ (t: trie_state, prefix: string) -> list[string]
    + returns every key in the trie that begins with prefix, in lexicographic order
    - returns [] when no keys match
    # scan
  ptrie.delete
    @ (t: trie_state, key: string) -> trie_state
    + returns a new trie with key removed and any now-empty branches pruned
    ? deleting a missing key is a no-op
    # delete
