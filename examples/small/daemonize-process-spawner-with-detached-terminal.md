# Requirement: "a library for spawning a new process detached from the controlling terminal"

Forks a child, detaches it from the terminal session, and optionally redirects its standard streams.

std
  std.os
    std.os.fork
      @ () -> result[i32, string]
      + returns 0 in the child and the child pid in the parent
      - returns error on failure
      # process
    std.os.setsid
      @ () -> result[void, string]
      + starts a new session detaching from the controlling terminal
      - returns error when the caller is already a session leader
      # process
    std.os.dup2
      @ (src_fd: i32, dst_fd: i32) -> result[void, string]
      + duplicates src_fd onto dst_fd
      - returns error on invalid descriptors
      # io
    std.os.open
      @ (path: string, flags: i32) -> result[i32, string]
      + opens a file and returns a file descriptor
      - returns error when the file cannot be opened
      # io

daemonize
  daemonize.spawn
    @ (stdout_path: string, stderr_path: string) -> result[daemon_handle, string]
    + forks, detaches from the controlling terminal, and redirects stdout and stderr
    + returns a handle containing the child pid in the parent
    - returns error when any step of detachment fails
    ? parent never sees the child's post-detach lifecycle; caller signals via pid
    # daemonization
    -> std.os.fork
    -> std.os.setsid
    -> std.os.open
    -> std.os.dup2
  daemonize.double_fork
    @ () -> result[i32, string]
    + performs the fork-setsid-fork sequence so the child cannot reacquire a terminal
    - returns error on any fork or setsid failure
    # daemonization
    -> std.os.fork
    -> std.os.setsid
