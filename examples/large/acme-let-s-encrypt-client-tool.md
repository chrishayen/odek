# Requirement: "an ACME certificate issuance client with automatic renewal"

Implements the ACME v2 flow: account registration, order creation, challenge completion, certificate fetching, and a renewal scheduler. HTTP and crypto live in std.

std
  std.http
    std.http.post_json
      @ (url: string, headers: map[string, string], body: string) -> result[http_response, string]
      + sends a POST with the given JSON body and returns status, headers, and body
      - returns error on network failure
      # http
    std.http.get
      @ (url: string, headers: map[string, string]) -> result[http_response, string]
      + sends a GET and returns status, headers, and body
      - returns error on network failure
      # http
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a flat JSON object
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a flat JSON object
      # serialization
  std.crypto
    std.crypto.rsa_generate
      @ (bits: i32) -> result[key_pair, string]
      + generates an RSA key pair of the given size
      - returns error for sizes below 2048
      # cryptography
    std.crypto.sign_rs256
      @ (key: key_pair, data: bytes) -> result[bytes, string]
      + returns an RSASSA-PKCS1-v1_5 SHA-256 signature over data
      # cryptography
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the SHA-256 digest of data
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

acme
  acme.account_new
    @ (contact_email: string) -> result[account, string]
    + creates a new account with a freshly generated key pair
    # account
    -> std.crypto.rsa_generate
  acme.account_register
    @ (directory_url: string, acct: account) -> result[account, string]
    + registers the account against the directory's new-account endpoint
    - returns error when the directory cannot be fetched
    - returns error when the server rejects the registration
    # account
    -> std.http.get
    -> std.http.post_json
    -> std.json.parse_object
  acme.order_new
    @ (acct: account, domains: list[string]) -> result[order, string]
    + submits a new certificate order for the given domain identifiers
    - returns error when the server rejects the order
    # order
    -> std.http.post_json
    -> std.json.encode_object
  acme.challenge_http01_key_auth
    @ (acct: account, token: string) -> result[string, string]
    + returns the key authorization string that must be served at /.well-known/acme-challenge/{token}
    # challenge
    -> std.crypto.sha256
    -> std.encoding.base64url_encode
  acme.challenge_dns01_record
    @ (acct: account, token: string) -> result[string, string]
    + returns the TXT record value that must be placed at _acme-challenge.{domain}
    # challenge
    -> std.crypto.sha256
    -> std.encoding.base64url_encode
  acme.challenge_complete
    @ (acct: account, challenge_url: string) -> result[void, string]
    + notifies the server that the challenge has been placed and waits for validation to complete
    - returns error when the server reports the challenge as invalid
    # challenge
    -> std.http.post_json
    -> std.json.parse_object
  acme.order_finalize
    @ (acct: account, ord: order, csr: bytes) -> result[order, string]
    + submits the CSR to finalize the order
    - returns error when finalization fails
    # order
    -> std.http.post_json
    -> std.crypto.sign_rs256
  acme.certificate_fetch
    @ (acct: account, ord: order) -> result[bytes, string]
    + downloads the issued certificate chain in PEM form
    - returns error when the certificate is not yet issued
    # issuance
    -> std.http.get
  acme.needs_renewal
    @ (cert_not_after: i64, renew_before_days: i32) -> bool
    + returns true when the certificate expires within renew_before_days
    # renewal
    -> std.time.now_seconds
  acme.renew
    @ (acct: account, domains: list[string]) -> result[bytes, string]
    + runs a full order and returns the new certificate chain
    # renewal
    -> acme.order_new
    -> acme.order_finalize
    -> acme.certificate_fetch
