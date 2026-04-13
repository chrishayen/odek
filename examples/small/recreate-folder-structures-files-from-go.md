# Requirement: "a library for recreating a folder tree and its files on disk from an in-memory embedded filesystem"

Walks an abstract embedded filesystem and materializes it to a destination directory. Project layer is the walker; std provides filesystem primitives.

std
  std.fs
    std.fs.mkdirs
      @ (path: string) -> result[void, string]
      + creates the directory and all missing parents
      # filesystem
    std.fs.write_bytes
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating or truncating the file
      - returns error when a parent directory does not exist
      # filesystem

fs_unembed
  fs_unembed.list_entries
    @ (efs: embedded_fs, dir: string) -> result[list[embedded_entry], string]
    + lists the direct children of dir inside the embedded filesystem
    - returns error when dir is not a directory
    # introspection
  fs_unembed.materialize_entry
    @ (efs: embedded_fs, src: string, dst_root: string) -> result[void, string]
    + creates the corresponding directory or writes the file under dst_root
    # extraction
    -> std.fs.mkdirs
    -> std.fs.write_bytes
  fs_unembed.materialize_tree
    @ (efs: embedded_fs, src_root: string, dst_root: string) -> result[void, string]
    + recreates the entire src_root subtree under dst_root
    - returns error on the first entry that fails to write
    # extraction
