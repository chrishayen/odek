# Requirement: "middleware for monetizing API requests with per-request payments"

Before serving a request the middleware issues an invoice, waits for proof of payment, and only then calls the wrapped handler. Payment-network details sit behind a pluggable interface.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.crypto
    std.crypto.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex digest of SHA-256
      # hashing

paywall
  paywall.new
    @ (price: i64, client: payment_client) -> paywall_state
    + creates a paywall with the given per-request price in the smallest currency unit
    # construction
  paywall.issue_invoice
    @ (state: paywall_state, request_id: string) -> result[invoice, string]
    + returns an invoice with an amount, a payment hash, and an expiry timestamp
    - returns error when the backing client rejects the request
    # billing
    -> std.time.now_seconds
    -> std.crypto.sha256_hex
  paywall.verify_payment
    @ (state: paywall_state, invoice: invoice, preimage: string) -> result[bool, string]
    + returns true when the preimage hashes to the invoice payment hash and is not expired
    - returns error when the preimage does not match
    - returns error when the invoice has expired
    # verification
    -> std.time.now_seconds
    -> std.crypto.sha256_hex
  paywall.protect
    @ (state: paywall_state, handler: fn(request) -> response) -> fn(request) -> response
    + returns a handler that responds 402 until a valid preimage header is present
    + delegates to the inner handler once payment is verified
    # middleware
