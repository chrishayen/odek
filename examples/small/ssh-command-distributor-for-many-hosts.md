# Requirement: "a library for distributing commands over ssh to many hosts"

Runs the same shell command on a list of remote hosts and aggregates the results.

std
  std.ssh
    std.ssh.run
      fn (host: string, user: string, command: string) -> result[ssh_output, string]
      + executes command and returns stdout, stderr, and exit_code
      - returns error on connection or authentication failure
      # ssh

fanout
  fanout.new_task
    fn (hosts: list[string], user: string, command: string) -> task
    + creates a task that will run command on every host as user
    # construction
  fanout.run
    fn (t: task) -> list[host_result]
    + runs the command on every host and returns one result per host in input order
    + a failed connection is captured as a host_result with a non-zero exit and error message
    # execution
    -> std.ssh.run
  fanout.summarize
    fn (results: list[host_result]) -> summary
    + returns counts of succeeded, failed, and total hosts
    # reporting
