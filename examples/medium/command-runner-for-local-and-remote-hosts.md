# Requirement: "a command runner that executes shell tasks on local and remote hosts"

A task is a name plus a command. Hosts are either local or remote over SSH. The runner fans out tasks across hosts and collects per-host results.

std
  std.process
    std.process.run
      fn (cmd: string, args: list[string], workdir: string) -> result[string, string]
      + runs a command locally and returns its stdout
      - returns error on non-zero exit
      # process
  std.ssh
    std.ssh.exec
      fn (host: string, user: string, port: i32, command: string) -> result[string, string]
      + runs a command on a remote host over SSH and returns its stdout
      - returns error on connection failure or non-zero exit
      # remote

command_runner
  command_runner.host_local
    fn (name: string) -> host
    + creates a local host entry
    # construction
  command_runner.host_remote
    fn (name: string, address: string, user: string, port: i32) -> host
    + creates a remote host entry
    # construction
  command_runner.task_new
    fn (name: string, command: string) -> task
    + creates a task record
    # construction
  command_runner.run_task
    fn (h: host, t: task) -> result[string, string]
    + runs the task against the host and returns its stdout
    - returns error on command failure
    # execution
    -> std.process.run
    -> std.ssh.exec
  command_runner.run_on_hosts
    fn (hosts: list[host], t: task) -> list[tuple[string, result[string, string]]]
    + runs the task against every host and returns labelled results
    # fan_out
