# Requirement: "a library for reading and writing nested map values via a dot path"

Treats values as strings for simplicity; the caller serializes structured values.

std: (all units exist)

dot_prop
  dot_prop.split_path
    @ (path: string) -> list[string]
    + splits a dotted path into segments, treating "\." as a literal dot inside a segment
    ? empty path returns an empty list
    # path_parsing
  dot_prop.get
    @ (obj: map[string, string], path: string) -> result[string, string]
    + returns the value at the given dotted path
    - returns error when any intermediate segment is missing
    ? the outer map is flat; keys are the full encoded paths
    # read
    -> dot_prop.split_path
  dot_prop.set
    @ (obj: map[string, string], path: string, value: string) -> map[string, string]
    + stores value at the given dotted path and returns the updated map
    # write
    -> dot_prop.split_path
  dot_prop.has
    @ (obj: map[string, string], path: string) -> bool
    + returns true when the path exists in the object
    # existence
    -> dot_prop.split_path
  dot_prop.delete
    @ (obj: map[string, string], path: string) -> map[string, string]
    + removes the entry at the given path, returning the updated map
    ? no-op when the path is absent
    # delete
    -> dot_prop.split_path
