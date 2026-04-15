# Requirement: "a library for generating random temporary file and directory paths"

Returns fresh paths under the host's temporary directory. Does not create the file or directory.

std
  std.fs
    std.fs.temp_dir
      fn () -> string
      + returns the host's temporary directory path
      # filesystem
  std.random
    std.random.hex_token
      fn (n_bytes: i32) -> string
      + returns a random hex-encoded token of n_bytes of entropy
      # randomness

tempy
  tempy.temp_file
    fn (suffix: string) -> string
    + returns a fresh unique path inside the temp dir with the given suffix
    ? the file is not created; the caller does the write
    # temp_file
    -> std.fs.temp_dir
    -> std.random.hex_token
  tempy.temp_directory
    fn () -> string
    + returns a fresh unique directory path inside the temp dir
    ? the directory is not created
    # temp_dir
    -> std.fs.temp_dir
    -> std.random.hex_token
