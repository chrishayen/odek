# Requirement: "an environment variable reader with dotenv file support"

Reads key-value pairs from a dotenv-style file and exposes a lookup that falls back to the process environment.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns file contents as a string
      - returns error when the file does not exist
      # filesystem
  std.env
    std.env.get
      fn (key: string) -> optional[string]
      + returns the value of a process environment variable
      - returns none when unset
      # environment

dotenv
  dotenv.parse
    fn (raw: string) -> map[string, string]
    + parses KEY=value lines into a map
    + ignores blank lines and lines starting with '#'
    + strips matching single or double quotes around values
    - returns an empty map for empty input
    # parsing
  dotenv.load_file
    fn (path: string) -> result[map[string, string], string]
    + reads the file and parses it
    - returns error when the file cannot be read
    # loading
    -> std.fs.read_all
  dotenv.lookup
    fn (loaded: map[string, string], key: string) -> optional[string]
    + returns the dotenv value when present
    + falls back to the process environment when not in the map
    - returns none when neither source has the key
    # lookup
    -> std.env.get
