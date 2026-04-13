# Requirement: "a library for folder-based operations (copy, move, search) on a hierarchical secret store"

Treats a secret store as a tree of paths and operates on whole subtrees. The project layer depends on an abstract secret_client so it can be backed by any store.

std: (all units exist)

secret_tree
  secret_tree.list_tree
    @ (client: secret_client, root: string) -> result[list[string], string]
    + returns every leaf path under the given root, recursively
    - returns error when the root cannot be listed
    # traversal
  secret_tree.read_tree
    @ (client: secret_client, root: string) -> result[map[string, map[string, string]], string]
    + returns a map from leaf path to its key/value pairs
    - returns error when any leaf cannot be read
    # reading
    -> secret_tree.list_tree
  secret_tree.copy_tree
    @ (client: secret_client, src: string, dst: string) -> result[i32, string]
    + copies every secret under src to the corresponding path under dst
    + returns the number of secrets copied
    - returns error when src does not exist
    # copying
    -> secret_tree.read_tree
  secret_tree.move_tree
    @ (client: secret_client, src: string, dst: string) -> result[i32, string]
    + copies src to dst then deletes src, only deleting on successful copy
    - returns error when the copy phase fails, leaving src intact
    # moving
    -> secret_tree.copy_tree
  secret_tree.delete_tree
    @ (client: secret_client, root: string) -> result[i32, string]
    + deletes every leaf under root and returns how many were removed
    # deletion
    -> secret_tree.list_tree
  secret_tree.search_tree
    @ (client: secret_client, root: string, needle: string) -> result[list[match], string]
    + returns every leaf whose path, key, or value contains needle
    # search
    -> secret_tree.read_tree
