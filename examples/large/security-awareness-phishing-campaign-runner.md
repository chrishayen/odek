# Requirement: "a library for running phishing-style security awareness campaigns"

Authorized security teams author templates, launch campaigns against a tracked recipient list, and record which recipients interacted with the decoy.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when the file cannot be opened
      # filesystem
  std.net
    std.net.smtp_send
      @ (host: string, port: u16, from: string, to: string, message: string) -> result[void, string]
      + sends a message via SMTP
      - returns error when delivery fails
      # networking
  std.http
    std.http.serve
      @ (port: u16, handler: fn(http_request) -> http_response) -> result[void, string]
      + starts an HTTP server that invokes the handler for each request
      - returns error when the port cannot be bound
      # http
  std.crypto
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
  std.encoding
    std.encoding.base64url_encode
      @ (data: bytes) -> string
      + encodes bytes to base64url without padding
      # encoding
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

awareness
  awareness.load_template
    @ (path: string) -> result[template, string]
    + returns a template with subject, body, and tracked-link placeholders
    - returns error when the file cannot be read or has missing placeholders
    # templating
    -> std.fs.read_all
  awareness.load_recipients
    @ (path: string) -> result[list[recipient], string]
    + returns one recipient per line, each with an email address and optional display name
    - returns error when the file cannot be read
    # loading
    -> std.fs.read_all
  awareness.new_campaign
    @ (name: string, tmpl: template, recipients: list[recipient]) -> campaign_state
    + creates a campaign in the prepared state with one pending delivery per recipient
    # construction
    -> std.time.now_seconds
  awareness.mint_token
    @ (state: campaign_state, recipient_id: string) -> tuple[string, campaign_state]
    + returns a one-time tracking token bound to the recipient and records it in the campaign
    # tracking
    -> std.crypto.random_bytes
    -> std.encoding.base64url_encode
  awareness.render_message
    @ (tmpl: template, recipient: recipient, tracking_url: string) -> string
    + returns the fully substituted message body addressed to the recipient
    # templating
  awareness.send_campaign
    @ (state: campaign_state, smtp_host: string, smtp_port: u16, from: string, base_url: string) -> result[campaign_state, string]
    + mints a token and dispatches a message per recipient, marking each as sent or failed in the returned state
    - returns error when every delivery fails
    # execution
    -> std.net.smtp_send
  awareness.record_interaction
    @ (state: campaign_state, token: string, kind: string) -> result[campaign_state, string]
    + marks the recipient bound to the token as having interacted (open, click, submit)
    - returns error when the token is unknown
    # tracking
    -> std.time.now_seconds
  awareness.serve_tracker
    @ (state_ref: campaign_ref, port: u16) -> result[void, string]
    + starts an HTTP server that decodes tracking tokens from incoming requests and records interactions
    - returns error when the port cannot be bound
    # server
    -> std.http.serve
  awareness.summarize
    @ (state: campaign_state) -> campaign_summary
    + returns totals for sent, opened, clicked, and submitted across the campaign
    # reporting
