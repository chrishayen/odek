# Requirement: "a task runner that reads task definitions from a markdown document"

Parses a markdown document into named tasks (each task is a fenced code block under a heading), resolves dependencies between them, and produces an execution order. Actually running commands is the caller's responsibility — the library returns the script text.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of the file
      - returns error when the file does not exist
      # filesystem

mdtasks
  mdtasks.parse
    @ (markdown: string) -> result[task_doc, string]
    + extracts each heading followed by a fenced code block as a named task
    - returns error when a heading has no code block
    # parsing
  mdtasks.load
    @ (path: string) -> result[task_doc, string]
    + reads and parses the file
    - returns error when the file cannot be read
    # parsing
    -> std.fs.read_all
  mdtasks.tasks
    @ (doc: task_doc) -> list[string]
    + returns the list of task names in document order
    # query
  mdtasks.script
    @ (doc: task_doc, name: string) -> result[string, string]
    + returns the script body for the named task
    - returns error when no task with that name exists
    # query
  mdtasks.depends_on
    @ (doc: task_doc, name: string) -> result[list[string], string]
    + returns dependencies declared in a "Requires:" inline directive
    - returns error when the task does not exist
    # query
  mdtasks.order
    @ (doc: task_doc, target: string) -> result[list[string], string]
    + returns tasks in topological order ending with target
    - returns error when there is a cyclic dependency
    - returns error when target does not exist
    # planning
