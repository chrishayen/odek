# Requirement: "a mail server library supporting SMTP, IMAP, and mailbox storage"

Protocol parsers and a mailbox store. No networking loop here; callers feed frames in and get frames out.

std
  std.io
    std.io.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when path is not readable
      # io
    std.io.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to path
      - returns error when path is not writable
      # io
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns entries in a directory, sorted
      - returns error when path is not a directory
      # filesystem
    std.fs.make_dir
      fn (path: string) -> result[void, string]
      + creates a directory including parents
      - returns error on permission failure
      # filesystem
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns SHA-256 digest
      # cryptography

mail
  mail.parse_message
    fn (raw: bytes) -> result[message, string]
    + parses a message into headers and a MIME body tree
    - returns error on malformed headers or boundary handling
    # parsing
  mail.encode_message
    fn (m: message) -> bytes
    + serializes a message back to RFC 5322 bytes
    # encoding
  mail.smtp_new_session
    fn (hostname: string) -> smtp_session
    + creates a server session ready to greet a client
    # smtp
  mail.smtp_handle
    fn (session: smtp_session, line: string) -> tuple[list[string], smtp_session]
    + consumes one SMTP command line and returns response lines and the advanced session
    + handles EHLO, MAIL FROM, RCPT TO, DATA, QUIT, NOOP, RSET
    - returns a 5xx response on out-of-sequence commands
    # smtp
  mail.imap_new_session
    fn (hostname: string) -> imap_session
    + creates an IMAP server session in the not-authenticated state
    # imap
  mail.imap_handle
    fn (session: imap_session, line: string) -> tuple[list[string], imap_session]
    + consumes one tagged IMAP command and returns tagged responses and the advanced session
    + handles LOGIN, LIST, SELECT, FETCH, STORE, LOGOUT
    - returns a BAD response on unknown commands
    # imap
  mail.open_store
    fn (root: string) -> result[mail_store, string]
    + opens a maildir-style store rooted at root, creating directories if missing
    - returns error on permission failure
    # storage
    -> std.fs.make_dir
  mail.append
    fn (store: mail_store, mailbox: string, m: message) -> result[u64, string]
    + writes a message to the mailbox and returns its uid
    - returns error when the mailbox does not exist
    # storage
    -> std.io.write_all
    -> std.crypto.sha256
  mail.list_mailboxes
    fn (store: mail_store) -> result[list[string], string]
    + returns mailbox names in the store
    # storage
    -> std.fs.list_dir
  mail.fetch
    fn (store: mail_store, mailbox: string, uid: u64) -> result[message, string]
    + returns the message at uid
    - returns error when uid is not present
    # storage
    -> std.io.read_all
  mail.delete
    fn (store: mail_store, mailbox: string, uid: u64) -> result[void, string]
    + removes the message at uid
    - returns error when uid is not present
    # storage
  mail.authenticate
    fn (username: string, password: string, hashed: bytes) -> bool
    + returns true when hashing password yields the stored hash
    # auth
    -> std.crypto.sha256
