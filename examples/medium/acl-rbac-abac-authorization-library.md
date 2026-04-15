# Requirement: "an authorization library supporting ACL, RBAC, and ABAC access control models"

A policy engine that evaluates a model definition against (subject, object, action) requests.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads the full contents of a file as text
      - returns error when the file does not exist
      # io
  std.text
    std.text.split_lines
      fn (input: string) -> list[string]
      + splits on newline, discarding trailing empty lines
      # text
    std.text.split_by
      fn (input: string, sep: string) -> list[string]
      + splits on a separator
      # text

casbin
  casbin.model_parse
    fn (text: string) -> result[model_def, string]
    + parses a model definition describing request, policy, and matcher shapes
    - returns error on missing sections
    # parsing
    -> std.text.split_lines
    -> std.text.split_by
  casbin.model_load
    fn (path: string) -> result[model_def, string]
    + reads and parses a model file
    # loading
    -> std.fs.read_all
  casbin.policy_new
    fn () -> policy_set
    + creates an empty policy set
    # construction
  casbin.policy_add
    fn (set: policy_set, rule: list[string]) -> policy_set
    + adds a rule (e.g. subject, object, action)
    # mutation
  casbin.policy_remove
    fn (set: policy_set, rule: list[string]) -> policy_set
    + removes a matching rule if present
    # mutation
  casbin.enforcer_new
    fn (model: model_def, policy: policy_set) -> enforcer_state
    + binds a model and policy for evaluation
    # construction
  casbin.enforce
    fn (enforcer: enforcer_state, subject: string, object: string, action: string) -> bool
    + evaluates the matcher and returns the permit/deny decision
    + matches ACL, RBAC, or ABAC shape depending on the model
    # evaluation
  casbin.add_role_for_user
    fn (enforcer: enforcer_state, user: string, role: string) -> enforcer_state
    + records a role assignment for RBAC evaluation
    # rbac
