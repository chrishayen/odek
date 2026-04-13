# Requirement: "a mail hosting and management platform library"

Core domain logic for administering mail domains, mailboxes, aliases, and quotas. Transport and UI live above this layer.

std
  std.crypto
    std.crypto.bcrypt_hash
      @ (password: string, cost: i32) -> result[string, string]
      + returns a bcrypt hash of the password
      - returns error when cost is out of range
      # cryptography
    std.crypto.bcrypt_verify
      @ (password: string, hash: string) -> bool
      + returns true when the password matches the hash
      # cryptography
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix timestamp in seconds
      # time
  std.text
    std.text.email_valid
      @ (address: string) -> bool
      + returns true when the address has a local part, "@", and a domain with at least one dot
      # validation

mail_platform
  mail_platform.new
    @ () -> platform_state
    + creates an empty platform with no domains or mailboxes
    # construction
  mail_platform.add_domain
    @ (state: platform_state, domain: string) -> result[platform_state, string]
    + registers a new mail domain
    - returns error when the domain already exists
    # domains
  mail_platform.remove_domain
    @ (state: platform_state, domain: string) -> result[platform_state, string]
    + removes a domain and all of its mailboxes and aliases
    - returns error when the domain does not exist
    # domains
  mail_platform.create_mailbox
    @ (state: platform_state, address: string, password: string, quota_bytes: i64) -> result[platform_state, string]
    + creates a mailbox with a hashed password and a quota
    - returns error when the address is malformed
    - returns error when the domain is not registered
    - returns error when the mailbox already exists
    # mailboxes
    -> std.text.email_valid
    -> std.crypto.bcrypt_hash
  mail_platform.delete_mailbox
    @ (state: platform_state, address: string) -> result[platform_state, string]
    + removes a mailbox
    - returns error when the mailbox does not exist
    # mailboxes
  mail_platform.set_password
    @ (state: platform_state, address: string, password: string) -> result[platform_state, string]
    + replaces a mailbox password with a new bcrypt hash
    - returns error when the mailbox does not exist
    # mailboxes
    -> std.crypto.bcrypt_hash
  mail_platform.authenticate
    @ (state: platform_state, address: string, password: string) -> bool
    + returns true when the password verifies against the stored hash
    # mailboxes
    -> std.crypto.bcrypt_verify
  mail_platform.add_alias
    @ (state: platform_state, alias: string, target: string) -> result[platform_state, string]
    + redirects mail for alias to the target address
    - returns error when the target mailbox does not exist
    - returns error when the alias is already used
    # aliases
    -> std.text.email_valid
  mail_platform.remove_alias
    @ (state: platform_state, alias: string) -> result[platform_state, string]
    + removes an alias
    - returns error when the alias does not exist
    # aliases
  mail_platform.record_usage
    @ (state: platform_state, address: string, bytes: i64) -> result[platform_state, string]
    + updates the mailbox's recorded storage usage
    - returns error when usage would exceed the quota
    # quota
    -> std.time.now_seconds
  mail_platform.list_mailboxes
    @ (state: platform_state, domain: string) -> list[string]
    + returns the addresses in the given domain
    # queries
