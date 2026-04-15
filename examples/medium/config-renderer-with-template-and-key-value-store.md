# Requirement: "a library that renders application config files from templates and a key-value store"

Watches a key-value data source, renders configured templates with the latest values, writes output files atomically, and runs reload hooks when outputs change.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file's entire contents
      - returns error when the path does not exist
      # filesystem
    std.fs.write_atomic
      fn (path: string, data: bytes) -> result[void, string]
      + writes to a temp file and renames over the target
      - returns error on write or rename failure
      # filesystem
    std.fs.sha256_file
      fn (path: string) -> result[bytes, string]
      + returns the SHA-256 digest of a file
      - returns error when the file cannot be read
      # filesystem
  std.process
    std.process.run
      fn (command: string, args: list[string]) -> result[i32, string]
      + runs a command and returns its exit code
      - returns error when the command cannot be launched
      # process

config_renderer
  config_renderer.new
    fn (store_fetch_fn: string) -> renderer_state
    + creates a renderer backed by a pluggable function that fetches values by key
    # construction
  config_renderer.add_template
    fn (state: renderer_state, src: string, dst: string, reload_cmd: string) -> renderer_state
    + registers a template source, output path, and optional reload command
    # registration
  config_renderer.render
    fn (state: renderer_state, template_src: string, values: map[string, string]) -> result[string, string]
    + substitutes "{{key}}" occurrences with values
    - returns error on unknown keys referenced by the template
    # templating
    -> std.fs.read_all
  config_renderer.apply_one
    fn (state: renderer_state, src: string) -> result[bool, string]
    + renders the template, writes the output if it differs, and returns true when the file changed
    - returns error when fetching values or writing fails
    # application
    -> std.fs.write_atomic
    -> std.fs.sha256_file
  config_renderer.apply_all
    fn (state: renderer_state) -> result[list[string], string]
    + applies every registered template and returns the list of destinations that changed
    - returns error on the first template that fails
    # application
  config_renderer.reload_changed
    fn (state: renderer_state, changed: list[string]) -> result[void, string]
    + runs the reload command for each changed destination
    - returns error on the first reload command that fails
    # reload
    -> std.process.run
