# Requirement: "a general-purpose business framework"

A minimal ERP-style core: parties, products, orders, inventory, invoices, and payments. Persistence is a key-value store abstraction; the framework owns the business rules.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.uuid
    std.uuid.new_v4
      @ () -> string
      + returns a fresh random uuid as a canonical hyphenated string
      # identifiers
  std.kv
    std.kv.put
      @ (store: kv_state, key: string, value: bytes) -> kv_state
      + stores value under key, overwriting any previous value
      # storage
    std.kv.get
      @ (store: kv_state, key: string) -> optional[bytes]
      + returns the stored bytes when the key exists
      # storage
    std.kv.list_prefix
      @ (store: kv_state, prefix: string) -> list[string]
      + returns all keys starting with prefix in lexicographic order
      # storage

biz
  biz.new_framework
    @ (store: kv_state) -> framework_state
    + wraps a kv store with framework bookkeeping (sequences, indexes)
    # construction
  biz.create_party
    @ (fw: framework_state, name: string, kind: string) -> tuple[framework_state, string]
    + creates a party (customer, supplier, or employee) and returns its id
    - returns an error id when kind is outside the supported set
    # parties
    -> std.uuid.new_v4
    -> std.kv.put
  biz.create_product
    @ (fw: framework_state, sku: string, name: string, unit_price_cents: i64) -> tuple[framework_state, string]
    + creates a product catalog entry and returns its id
    - returns an error id when sku is already in use
    # catalog
    -> std.uuid.new_v4
    -> std.kv.put
  biz.adjust_stock
    @ (fw: framework_state, product_id: string, delta: i64) -> result[framework_state, string]
    + adjusts on-hand inventory for a product by delta (positive or negative)
    - returns error when the resulting stock would be negative
    # inventory
    -> std.kv.get
    -> std.kv.put
  biz.create_order
    @ (fw: framework_state, party_id: string, lines: list[order_line]) -> result[tuple[framework_state, string], string]
    + creates an order, reserves stock, and returns its id
    - returns error when any line references an unknown product
    - returns error when any line exceeds available stock
    # orders
    -> std.uuid.new_v4
    -> std.time.now_seconds
  biz.confirm_order
    @ (fw: framework_state, order_id: string) -> result[framework_state, string]
    + moves the order from draft to confirmed and deducts reserved stock
    - returns error when the order is not in draft state
    # orders
  biz.issue_invoice
    @ (fw: framework_state, order_id: string) -> result[tuple[framework_state, string], string]
    + creates an invoice tied to the order and returns its id
    - returns error when the order is not confirmed
    # invoices
    -> std.uuid.new_v4
    -> std.time.now_seconds
  biz.record_payment
    @ (fw: framework_state, invoice_id: string, amount_cents: i64) -> result[framework_state, string]
    + registers a payment against an invoice
    - returns error when amount would overpay the invoice balance
    # payments
    -> std.time.now_seconds
  biz.invoice_balance
    @ (fw: framework_state, invoice_id: string) -> result[i64, string]
    + returns the outstanding amount in cents for an invoice
    - returns error when the invoice id is unknown
    # invoices
  biz.list_parties_by_kind
    @ (fw: framework_state, kind: string) -> list[string]
    + returns the ids of every party matching the kind
    # parties
    -> std.kv.list_prefix
