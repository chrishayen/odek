# Requirement: "a file move utility that works across devices"

Tries a cheap rename first; on cross-device errors falls back to copy-then-delete.

std
  std.fs
    std.fs.rename
      @ (src: string, dst: string) -> result[void, string]
      + renames within a single filesystem
      - returns cross-device error when src and dst live on different devices
      # filesystem
    std.fs.copy_file
      @ (src: string, dst: string) -> result[void, string]
      + copies file contents and permissions
      - returns error when src cannot be read or dst cannot be written
      # filesystem
    std.fs.remove
      @ (path: string) -> result[void, string]
      + deletes a file
      # filesystem

movefile
  movefile.move
    @ (src: string, dst: string) -> result[void, string]
    + moves a file, using rename when possible and copy-then-delete otherwise
    - returns error when src does not exist
    - leaves src in place when the fallback copy fails
    # file_move
    -> std.fs.rename
    -> std.fs.copy_file
    -> std.fs.remove
