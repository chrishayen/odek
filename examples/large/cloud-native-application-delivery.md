# Requirement: "a library for declarative application delivery across multiple environments"

Describe an application as components and traits, resolve it against a target environment, and roll out revisions.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns file contents
      - returns error when unreadable
      # filesystem
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns unix time
      # time
  std.hash
    std.hash.sha1_hex
      @ (data: bytes) -> string
      + returns the hex sha1
      # hashing

delivery
  delivery.parse_application
    @ (raw: string) -> result[application, string]
    + parses an application with name, components, and traits
    - returns error when required fields are missing
    # parsing
  delivery.validate_application
    @ (app: application) -> result[void, string]
    + checks that every trait targets a declared component
    - returns error listing unknown component references
    # validation
  delivery.resolve_component
    @ (comp: component, env: environment) -> result[resolved_component, string]
    + applies environment variable overrides and defaults
    - returns error on missing required variable
    # resolution
  delivery.plan_revision
    @ (app: application, env: environment, previous: optional[revision]) -> result[revision_plan, string]
    + returns the set of components to create, update, or delete
    - returns error when the application fails validation
    # planning
  delivery.hash_revision
    @ (plan: revision_plan) -> string
    + returns a stable content hash of the plan
    # identity
    -> std.hash.sha1_hex
  delivery.apply_plan
    @ (plan: revision_plan, target: delivery_target) -> result[revision, string]
    + submits the plan to the target and returns the new revision record
    - returns error when the target rejects the plan
    # rollout
    -> std.time.now_seconds
  delivery.rollback
    @ (target: delivery_target, to_revision: string) -> result[void, string]
    + re-applies the plan captured in the named revision
    - returns error when the revision is unknown
    # rollout
  delivery.diff_revisions
    @ (before: revision, after: revision) -> list[revision_change]
    + returns added, removed, and modified components
    # diff
  delivery.load_from_file
    @ (path: string) -> result[application, string]
    + reads an application manifest from disk
    - returns error on parse or read failure
    # loading
    -> std.fs.read_all
  delivery.register_trait
    @ (registry: trait_registry, name: string, handler: trait_handler) -> trait_registry
    + returns a registry with the new trait handler installed
    # extensibility
  delivery.apply_traits
    @ (comp: resolved_component, registry: trait_registry) -> result[resolved_component, string]
    + applies every trait handler in declaration order
    - returns error when a handler fails
    # traits
