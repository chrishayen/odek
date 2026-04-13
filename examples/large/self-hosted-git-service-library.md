# Requirement: "a self-hosted git hosting service library"

Manages users, organizations, and bare repositories on disk. Provides a small API surface the caller can mount behind HTTP.

std
  std.fs
    std.fs.mkdir_all
      @ (path: string) -> result[void, string]
      + creates the directory and all parents
      - returns error when path exists and is not a directory
      # filesystem
    std.fs.remove_all
      @ (path: string) -> result[void, string]
      + recursively removes the path
      # filesystem
    std.fs.exists
      @ (path: string) -> bool
      + returns true when the path exists
      # filesystem
  std.crypto
    std.crypto.bcrypt_hash
      @ (password: string, cost: i32) -> result[string, string]
      + returns a bcrypt hash at the given cost
      - returns error when password exceeds the bcrypt length limit
      # password_hashing
    std.crypto.bcrypt_verify
      @ (password: string, hash: string) -> bool
      + returns true when the password matches the hash
      # password_hashing
  std.git
    std.git.init_bare
      @ (path: string) -> result[void, string]
      + initializes a bare repository at path
      - returns error when the directory is not empty
      # git
    std.git.list_refs
      @ (path: string) -> result[list[git_ref], string]
      + returns all refs in the bare repository
      # git
    std.git.read_blob
      @ (path: string, sha: string) -> result[bytes, string]
      + returns the blob bytes for the given object id
      - returns error when the object does not exist
      # git
    std.git.log
      @ (path: string, ref: string, limit: i32) -> result[list[git_commit], string]
      + returns up to limit commits reachable from ref, newest first
      # git
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

git_service
  git_service.new
    @ (data_root: string) -> service_state
    + creates a service rooted at data_root
    # construction
    -> std.fs.mkdir_all
  git_service.register_user
    @ (state: service_state, username: string, email: string, password: string) -> result[tuple[service_state, user_id], string]
    + creates a user with a bcrypt-hashed password
    - returns error when the username already exists
    - returns error when the password is empty
    # users
    -> std.crypto.bcrypt_hash
  git_service.authenticate
    @ (state: service_state, username: string, password: string) -> optional[user_id]
    + returns the user id when credentials verify
    - returns none on mismatch
    # auth
    -> std.crypto.bcrypt_verify
  git_service.create_organization
    @ (state: service_state, owner: user_id, name: string) -> result[tuple[service_state, org_id], string]
    + creates an organization owned by the user
    - returns error when the organization name is taken
    # organizations
  git_service.add_org_member
    @ (state: service_state, org: org_id, user: user_id, role: string) -> result[service_state, string]
    + adds a user to an organization with a role
    - returns error when role is not "member" or "admin"
    # organizations
  git_service.create_repository
    @ (state: service_state, owner_user: optional[user_id], owner_org: optional[org_id], name: string, is_private: bool) -> result[tuple[service_state, repo_id], string]
    + initializes a bare repository on disk and records it in the service
    - returns error when both owner_user and owner_org are none
    - returns error when a repository with that name already exists under the owner
    # repositories
    -> std.fs.mkdir_all
    -> std.git.init_bare
  git_service.delete_repository
    @ (state: service_state, actor: user_id, repo: repo_id) -> result[service_state, string]
    + deletes the repository on disk and removes it from the service
    - returns error when the actor is not an owner or admin
    # repositories
    -> std.fs.remove_all
  git_service.grant_access
    @ (state: service_state, repo: repo_id, user: user_id, level: string) -> result[service_state, string]
    + sets the access level ("read", "write", "admin") for the user on the repo
    - returns error when level is invalid
    # permissions
  git_service.check_access
    @ (state: service_state, repo: repo_id, user: user_id, required: string) -> bool
    + returns true when the user has at least the required access
    # permissions
  git_service.list_branches
    @ (state: service_state, repo: repo_id) -> result[list[string], string]
    + returns branch names for the repository
    - returns error when the repo id is unknown
    # browsing
    -> std.git.list_refs
  git_service.commit_history
    @ (state: service_state, repo: repo_id, ref: string, limit: i32) -> result[list[git_commit], string]
    + returns up to limit commits reachable from ref
    - returns error when the ref does not exist
    # browsing
    -> std.git.log
  git_service.open_issue
    @ (state: service_state, repo: repo_id, author: user_id, title: string, body: string) -> result[tuple[service_state, issue_id], string]
    + creates a new issue on the repository
    - returns error when title is empty
    # issues
    -> std.time.now_seconds
  git_service.close_issue
    @ (state: service_state, actor: user_id, issue: issue_id) -> result[service_state, string]
    + marks the issue as closed
    - returns error when the actor lacks write access
    # issues
