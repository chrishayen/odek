# Requirement: "opens urls, files, and executables with the system default handler"

Single entry point that dispatches a target to the host's default opener.

std
  std.proc
    std.proc.run
      fn (cmd: string, args: list[string]) -> result[i32, string]
      + runs the command and returns the exit code
      - returns error when the binary cannot be found
      # process

opener
  opener.open
    fn (target: string) -> result[void, string]
    + launches the host default handler for the target (url, file path, or executable)
    - returns error when the target string is empty
    ? caller is responsible for whatever the handler does next
    # dispatch
    -> std.proc.run
