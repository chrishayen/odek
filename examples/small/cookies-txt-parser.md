# Requirement: "a parser for the cookies.txt file format"

Reads the tab-separated Netscape cookies.txt format and returns structured cookie records.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the file contents as a string
      - returns error when the file cannot be read
      # filesystem

cookies_txt
  cookies_txt.parse
    @ (text: string) -> result[list[cookie], string]
    + parses each non-comment line into domain, include_subdomains, path, secure, expires, name, value
    + skips blank lines and lines beginning with "#" (except the "#HttpOnly_" prefix which strips and sets http_only)
    - returns error when a data line does not have exactly seven tab-separated fields
    # parsing
  cookies_txt.load_file
    @ (path: string) -> result[list[cookie], string]
    + reads the file at path and parses it
    - returns error on read or parse failure
    # loading
    -> std.fs.read_all
