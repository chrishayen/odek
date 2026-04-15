# Requirement: "a database migration tool for managing schema evolution through incremental scripts"

Supports both pure-SQL scripts and scripts that delegate to a user-supplied callback.

std: (all units exist)

goose
  goose.register_script
    fn (version: i64, up: migration_fn, down: migration_fn) -> void
    + registers up/down callbacks for a version
    - returns silently when the version already exists, overwriting
    # registry
  goose.load_sql_scripts
    fn (sources: map[string, string]) -> result[void, string]
    + parses files with +goose Up / +goose Down markers and registers them
    - returns error on a script missing an Up section
    # loading
  goose.current_version
    fn (executor: sql_executor) -> result[i64, string]
    + returns the highest applied version, or 0 if none
    - returns error when the version table cannot be read
    # state
  goose.up_to
    fn (executor: sql_executor, target: i64) -> result[i32, string]
    + applies scripts in ascending order up to and including target
    - returns error and halts on the first script that fails
    # migration
  goose.down_to
    fn (executor: sql_executor, target: i64) -> result[i32, string]
    + rolls back scripts in descending order down to target
    - returns error when target exceeds current version
    # rollback
  goose.status
    fn (executor: sql_executor) -> result[list[script_status], string]
    + returns each registered script with applied/pending state
    # inspection
