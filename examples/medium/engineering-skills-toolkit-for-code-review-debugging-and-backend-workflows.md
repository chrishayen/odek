# Requirement: "an engineering skills toolkit for code review, debugging, and backend workflows"

Catalogs named "skills" (review checklists, debug playbooks, workflow snippets) and runs them against a code blob. Loading skill content is a std primitive.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the path does not exist
      # filesystem
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization

skills
  skills.new_catalog
    fn () -> catalog_state
    + creates an empty skills catalog
    # construction
  skills.load_from_dir
    fn (state: catalog_state, dir: string) -> result[catalog_state, string]
    + scans a directory for skill manifests and loads them
    - returns error when the directory is missing
    # loading
    -> std.fs.read_all
    -> std.json.parse_object
  skills.register
    fn (state: catalog_state, name: string, category: string, body: string) -> catalog_state
    + adds a skill in-memory under a category
    # registry
  skills.list_by_category
    fn (state: catalog_state, category: string) -> list[string]
    + returns skill names in a category
    # registry
  skills.find
    fn (state: catalog_state, name: string) -> optional[string]
    + returns the body of a skill by name
    # registry
  skills.run_review
    fn (state: catalog_state, code: string) -> list[string]
    + runs every skill tagged "review" against the code and collects their findings
    # review
  skills.run_debug
    fn (state: catalog_state, error_message: string, context: string) -> list[string]
    + runs every skill tagged "debug" and returns suggested next actions
    # debug
