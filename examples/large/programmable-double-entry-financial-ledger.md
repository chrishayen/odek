# Requirement: "a programmable double-entry financial ledger"

An append-only ledger of balanced entries with account balances queryable at any point in time.

std
  std.collections
    std.collections.map_new
      @ () -> map[string, i64]
      + creates an empty string-to-i64 map
      # collections
  std.uuid
    std.uuid.v4
      @ () -> string
      + returns a random UUIDv4 string
      # uuid
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

ledger
  ledger.new
    @ () -> ledger_state
    + creates an empty ledger with no accounts and no entries
    # construction
  ledger.open_account
    @ (state: ledger_state, name: string, currency: string, kind: string) -> result[tuple[string, ledger_state], string]
    + registers a new account and returns its generated id
    - returns error when kind is not one of asset, liability, equity, revenue, expense
    # accounts
    -> std.uuid.v4
  ledger.post_entry
    @ (state: ledger_state, lines: list[entry_line], memo: string) -> result[tuple[string, ledger_state], string]
    + appends an entry when debits equal credits and accounts exist
    - returns error when debits do not equal credits
    - returns error when any referenced account is unknown
    # posting
    -> std.uuid.v4
    -> std.time.now_millis
  ledger.balance
    @ (state: ledger_state, account_id: string) -> result[i64, string]
    + returns the current balance for an account
    - returns error when the account is unknown
    # balances
  ledger.balance_at
    @ (state: ledger_state, account_id: string, at_millis: i64) -> result[i64, string]
    + returns the balance summed over entries with timestamp <= at_millis
    # balances
  ledger.trial_balance
    @ (state: ledger_state) -> map[string, i64]
    + returns account_id -> balance for every account
    # reporting
    -> std.collections.map_new
  ledger.get_entry
    @ (state: ledger_state, entry_id: string) -> optional[entry]
    + returns the entry with the given id
    # introspection
  ledger.entries_for_account
    @ (state: ledger_state, account_id: string) -> list[entry]
    + returns every entry referencing the account in post order
    # introspection
  ledger.reverse_entry
    @ (state: ledger_state, entry_id: string, memo: string) -> result[tuple[string, ledger_state], string]
    + posts a new entry that flips every line of the original
    - returns error when the original entry is unknown
    # posting
    -> std.uuid.v4
    -> std.time.now_millis
  ledger.transfer
    @ (state: ledger_state, from_id: string, to_id: string, amount: i64, memo: string) -> result[tuple[string, ledger_state], string]
    + convenience wrapper that posts a two-line balanced entry
    - returns error when accounts have different currencies
    # posting
  ledger.validate
    @ (state: ledger_state) -> result[void, string]
    + returns ok when every entry balances and every referenced account exists
    - returns error describing the first inconsistency
    # invariants
