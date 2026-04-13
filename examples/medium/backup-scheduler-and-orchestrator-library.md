# Requirement: "a scheduling and orchestration library for a backup tool"

Manages backup repositories, schedules snapshot jobs, and tracks their status. Execution of the underlying backup binary is abstracted behind a runner interface.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.id
    std.id.generate
      @ () -> string
      + returns a new unique identifier
      # id_generation

backup_orchestrator
  backup_orchestrator.new
    @ () -> orchestrator_state
    + returns an empty orchestrator
    # construction
  backup_orchestrator.add_repository
    @ (state: orchestrator_state, name: string, path: string) -> result[orchestrator_state, string]
    + registers a repository
    - returns error when the name is already used
    # configuration
  backup_orchestrator.schedule_job
    @ (state: orchestrator_state, repo_name: string, cron: string, paths: list[string]) -> result[tuple[string, orchestrator_state], string]
    + returns the new job id and updated state
    - returns error when the repository is unknown
    - returns error when the cron expression is invalid
    # scheduling
    -> std.id.generate
  backup_orchestrator.due_jobs
    @ (state: orchestrator_state, at_time: i64) -> list[string]
    + returns ids of jobs whose next run time is at or before at_time
    # scheduling
  backup_orchestrator.mark_started
    @ (state: orchestrator_state, job_id: string) -> result[orchestrator_state, string]
    + transitions the job to running and records the start time
    # execution
    -> std.time.now_seconds
  backup_orchestrator.mark_finished
    @ (state: orchestrator_state, job_id: string, success: bool, message: string) -> result[orchestrator_state, string]
    + records the outcome and computes the next scheduled run
    # execution
    -> std.time.now_seconds
  backup_orchestrator.job_history
    @ (state: orchestrator_state, job_id: string) -> list[job_run]
    + returns previous runs for the job in chronological order
    # inspection
