# Requirement: "an email server library providing message transport, mailbox access, and DKIM signing and verification"

A library that implements the core state machines and cryptographic helpers for running an email server. Network I/O is the caller's responsibility.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns 32 bytes of SHA-256 digest
      # cryptography
    std.crypto.rsa_sign_sha256
      fn (key: bytes, data: bytes) -> result[bytes, string]
      + returns an RSA-SHA256 signature
      - returns error when the key is malformed
      # cryptography
    std.crypto.rsa_verify_sha256
      fn (pubkey: bytes, data: bytes, sig: bytes) -> bool
      + returns true when the signature is valid
      # cryptography
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + returns base64 encoding
      # encoding
    std.encoding.base64_decode
      fn (s: string) -> result[bytes, string]
      + returns decoded bytes
      - returns error on invalid base64
      # encoding
  std.net
    std.net.dns_lookup_txt
      fn (name: string) -> result[list[string], string]
      + returns TXT record values for a DNS name
      - returns error when the lookup fails
      # dns
    std.net.dns_lookup_mx
      fn (domain: string) -> result[list[mx_record], string]
      + returns MX records sorted by preference
      # dns
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

mail_server
  mail_server.parse_message
    fn (raw: bytes) -> result[mail_message, string]
    + returns a structured message with headers and body
    - returns error on malformed message syntax
    # parsing
  mail_server.serialize_message
    fn (msg: mail_message) -> bytes
    + returns canonical on-wire representation
    # parsing
  mail_server.dkim_canonicalize
    fn (headers: list[header], body: bytes) -> bytes
    + returns the canonicalized input for DKIM signing (relaxed/simple modes)
    # dkim
  mail_server.dkim_sign
    fn (msg: mail_message, selector: string, domain: string, privkey: bytes) -> result[mail_message, string]
    + returns the message with a DKIM-Signature header added
    - returns error when signing fails
    # dkim
    -> mail_server.dkim_canonicalize
    -> std.crypto.sha256
    -> std.crypto.rsa_sign_sha256
    -> std.encoding.base64_encode
  mail_server.dkim_verify
    fn (msg: mail_message) -> result[bool, string]
    + returns true when all DKIM signatures validate
    - returns error when the DNS TXT record cannot be fetched
    # dkim
    -> mail_server.dkim_canonicalize
    -> std.net.dns_lookup_txt
    -> std.crypto.rsa_verify_sha256
    -> std.encoding.base64_decode
  mail_server.smtp_session_new
    fn () -> smtp_session
    + returns a session in the initial state
    # smtp
  mail_server.smtp_feed
    fn (session: smtp_session, line: string) -> tuple[smtp_reply, smtp_session]
    + returns (reply, new_state) advancing the SMTP command state machine
    - returns 5xx reply for protocol violations
    # smtp
  mail_server.smtp_take_message
    fn (session: smtp_session) -> optional[mail_message]
    + returns the assembled message when the session has DATA ready
    # smtp
  mail_server.imap_session_new
    fn (user: string) -> imap_session
    + returns an authenticated IMAP session
    # imap
  mail_server.imap_feed
    fn (session: imap_session, line: string) -> tuple[imap_reply, imap_session]
    + returns (reply, new_state) advancing the IMAP command state machine
    # imap
  mail_server.mailbox_new
    fn (name: string) -> mailbox_state
    + returns an empty mailbox
    # mailbox
  mail_server.mailbox_append
    fn (mbox: mailbox_state, msg: mail_message) -> mailbox_state
    + stores the message with an assigned UID
    # mailbox
    -> std.time.now_seconds
  mail_server.mailbox_fetch
    fn (mbox: mailbox_state, uid: i64) -> optional[mail_message]
    + returns the message with the given UID
    # mailbox
  mail_server.mailbox_expunge
    fn (mbox: mailbox_state, uid: i64) -> mailbox_state
    + marks the uid for removal and compacts the mailbox
    # mailbox
  mail_server.resolve_mx
    fn (domain: string) -> result[list[string], string]
    + returns hostnames of mail exchangers sorted by preference
    - returns error when the lookup fails
    # routing
    -> std.net.dns_lookup_mx
