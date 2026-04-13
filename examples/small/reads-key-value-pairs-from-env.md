# Requirement: "a library that reads key-value pairs from a dotenv file and loads them into the environment"

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the file contents as a string
      - returns error when the file cannot be read
      # filesystem
  std.env
    std.env.set
      @ (key: string, value: string) -> void
      + sets a process environment variable
      # environment

dotenv_loader
  dotenv_loader.parse
    @ (raw: string) -> map[string, string]
    + parses KEY=value lines into a map
    + ignores blank lines and lines starting with '#'
    + strips matching surrounding single or double quotes from values
    - returns an empty map for empty input
    # parsing
  dotenv_loader.load
    @ (path: string, override: bool) -> result[i32, string]
    + reads the file, parses it, and sets each pair in the environment
    + returns the count of keys applied
    + when override is false, leaves already-set variables untouched
    - returns error when the file cannot be read
    # loading
    -> std.fs.read_all
    -> std.env.set
