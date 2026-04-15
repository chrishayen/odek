# Requirement: "a library that mirrors a local process into a remote container so it sees the remote environment"

Intercepts system calls on a local process and forwards filesystem, network, and environment operations to an agent running inside a remote container, so the local binary behaves as if it were executing there.

std
  std.net
    std.net.dial_tcp
      fn (host: string, port: u16) -> result[conn_state, string]
      + returns a connected socket
      # networking
    std.net.send
      fn (conn: conn_state, payload: bytes) -> result[void, string]
      + writes the payload
      # networking
    std.net.recv
      fn (conn: conn_state, max: u32) -> result[bytes, string]
      + reads up to max bytes
      # networking
  std.process
    std.process.intercept_syscalls
      fn (pid: i32) -> result[syscall_hook, string]
      + installs a syscall hook on the target process
      - returns error when the process cannot be attached
      # process_control
    std.process.resume_syscall
      fn (hook: syscall_hook, result_value: i64) -> void
      + resumes the suspended syscall with the given return value
      # process_control

process_mirror
  process_mirror.new_session
    fn (remote_host: string, remote_port: u16, target_pid: i32) -> result[mirror_session, string]
    + attaches to the target process and connects to the remote agent
    - returns error when the agent is unreachable
    # construction
    -> std.net.dial_tcp
    -> std.process.intercept_syscalls
  process_mirror.classify_syscall
    fn (call_number: i32) -> syscall_category
    + returns whether the call should be mirrored, kept local, or blocked
    # routing
  process_mirror.forward_fs_call
    fn (session: mirror_session, call_number: i32, args: list[i64]) -> result[i64, string]
    + forwards a filesystem syscall to the remote agent and returns the result
    - returns error on protocol failure
    # remote_fs
    -> std.net.send
    -> std.net.recv
  process_mirror.forward_net_call
    fn (session: mirror_session, call_number: i32, args: list[i64]) -> result[i64, string]
    + forwards a socket syscall so network traffic uses the remote container's network namespace
    # remote_net
    -> std.net.send
    -> std.net.recv
  process_mirror.forward_env_read
    fn (session: mirror_session, name: string) -> optional[string]
    + returns an environment variable as seen by the remote agent
    # remote_env
    -> std.net.send
    -> std.net.recv
  process_mirror.handle_intercept
    fn (session: mirror_session, hook: syscall_hook, call_number: i32, args: list[i64]) -> mirror_session
    + dispatches an intercepted syscall to the correct forwarder or lets it run locally
    # dispatch
    -> std.process.resume_syscall
  process_mirror.mount_remote_dir
    fn (session: mirror_session, remote_path: string, local_mountpoint: string) -> result[mirror_session, string]
    + advertises remote files at local_mountpoint for the target process
    # filesystem_bridge
    -> std.net.send
  process_mirror.detach
    fn (session: mirror_session) -> result[void, string]
    + releases the syscall hook and closes the agent connection
    # teardown
  process_mirror.agent_serve_step
    fn (listener: conn_state) -> result[void, string]
    + handles one incoming forwarding request inside the remote agent
    # agent
    -> std.net.recv
    -> std.net.send
  process_mirror.agent_handle_fs_request
    fn (req: bytes) -> bytes
    + executes a filesystem request in the remote agent and returns the response
    # agent
  process_mirror.agent_handle_net_request
    fn (req: bytes) -> bytes
    + executes a socket request in the remote agent and returns the response
    # agent
