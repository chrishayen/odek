# Requirement: "an ACME client for fully-managed TLS certificate issuance and renewal"

Handles the ACME protocol flow: account registration, order creation, domain challenges, CSR submission, certificate retrieval, and renewal scheduling. Cryptographic primitives and HTTP transport live in std.

std
  std.http
    std.http.get
      fn (url: string) -> result[bytes, string]
      + performs an HTTP GET and returns the response body
      - returns error on network failure or non-2xx status
      # http
    std.http.post_json
      fn (url: string, body: bytes, headers: map[string,string]) -> result[bytes, string]
      + performs an HTTP POST with a JSON body
      - returns error on network failure
      # http
  std.crypto
    std.crypto.generate_rsa_key
      fn (bits: i32) -> result[bytes, string]
      + returns a new RSA private key in DER form
      - returns error for unsupported key sizes
      # cryptography
    std.crypto.sign_rsa
      fn (key: bytes, data: bytes) -> result[bytes, string]
      + returns an RSA-SHA256 signature over data
      # cryptography
    std.crypto.generate_csr
      fn (key: bytes, domains: list[string]) -> result[bytes, string]
      + returns a DER-encoded certificate signing request for the given domains
      # cryptography
  std.encoding
    std.encoding.base64url_encode
      fn (data: bytes) -> string
      + encodes bytes to base64url without padding
      # encoding
  std.json
    std.json.parse_object
      fn (raw: bytes) -> result[map[string,string], string]
      + parses a JSON object into a string map
      # serialization
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns wall-clock time in seconds
      # time

acme
  acme.new_client
    fn (directory_url: string) -> result[acme_client, string]
    + fetches the ACME directory and returns a configured client
    - returns error when the directory cannot be retrieved
    # construction
    -> std.http.get
    -> std.json.parse_object
  acme.register_account
    fn (client: acme_client, contact_email: string) -> result[account_state, string]
    + creates a new ACME account with a fresh key
    + returns the account URL and key material
    # accounts
    -> std.crypto.generate_rsa_key
    -> std.http.post_json
  acme.create_order
    fn (client: acme_client, account: account_state, domains: list[string]) -> result[order_state, string]
    + opens a new order for the requested domains
    - returns error when the server rejects the domain list
    # ordering
    -> std.http.post_json
  acme.get_challenges
    fn (client: acme_client, order: order_state) -> result[list[challenge], string]
    + returns the list of authorization challenges pending for an order
    # challenges
    -> std.http.get
  acme.build_http01_response
    fn (challenge: challenge, account: account_state) -> result[string, string]
    + returns the HTTP-01 key authorization token to serve at the challenge path
    # challenges
    -> std.encoding.base64url_encode
  acme.notify_challenge_ready
    fn (client: acme_client, challenge: challenge) -> result[void, string]
    + tells the server the challenge is ready to be verified
    - returns error when the server reports validation failure
    # challenges
    -> std.http.post_json
  acme.finalize_order
    fn (client: acme_client, order: order_state, domains: list[string]) -> result[bytes, string]
    + submits a CSR and returns the issued certificate chain as PEM bytes
    - returns error when the order is not ready
    # finalization
    -> std.crypto.generate_csr
    -> std.http.post_json
  acme.parse_certificate_expiry
    fn (pem: bytes) -> result[i64, string]
    + returns the NotAfter timestamp of the leaf certificate in seconds
    - returns error on malformed PEM
    # certificates
  acme.needs_renewal
    fn (pem: bytes, renew_before_seconds: i64) -> bool
    + returns true when the certificate expires within the window
    # renewal
    -> std.time.now_seconds
    -> acme.parse_certificate_expiry
  acme.renew
    fn (client: acme_client, account: account_state, domains: list[string]) -> result[bytes, string]
    + runs the full issuance flow and returns the new certificate PEM
    # renewal
    -> acme.create_order
    -> acme.finalize_order
