# Requirement: "an embedded file access library"

Provides a read-only in-memory file system assembled from a build-time map of path to bytes.

std: (all units exist)

embed
  embed.new_fs
    fn (entries: map[string, bytes]) -> embedded_fs
    + creates a read-only file system from a path-to-content map
    # construction
  embed.read_file
    fn (fs: embedded_fs, path: string) -> result[bytes, string]
    + returns the bytes stored at the given path
    - returns error when the path is not present
    # lookup
  embed.list_dir
    fn (fs: embedded_fs, dir: string) -> list[string]
    + returns the direct children (file or subdir names) of dir
    + returns an empty list for unknown directories
    # listing
  embed.exists
    fn (fs: embedded_fs, path: string) -> bool
    + reports whether a file at the exact path is present
    # lookup
