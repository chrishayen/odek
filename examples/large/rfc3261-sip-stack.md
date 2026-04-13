# Requirement: "an RFC 3261 compliant SIP stack"

A core SIP message library with transaction and dialog state machines. Transport I/O is delegated to the caller.

std
  std.parsing
    std.parsing.split_lines
      @ (raw: string) -> list[string]
      + splits on CRLF, skipping empty trailing lines
      # parsing
    std.parsing.trim
      @ (s: string) -> string
      + returns s with leading and trailing ASCII whitespace removed
      # parsing
  std.random
    std.random.hex_string
      @ (byte_count: i32) -> string
      + returns a hex-encoded random string of the requested byte length
      # randomness

sip
  sip.parse_message
    @ (raw: string) -> result[sip_message, string]
    + parses a request or response into a typed message with headers and body
    - returns error when the start line is malformed
    - returns error when a header is missing a colon
    # parsing
    -> std.parsing.split_lines
    -> std.parsing.trim
  sip.serialize_message
    @ (msg: sip_message) -> string
    + returns the wire representation with CRLF separators
    # serialization
  sip.new_request
    @ (method: string, uri: string, from: string, to: string, call_id: string) -> sip_message
    + builds a request with required headers and a fresh branch parameter
    # request_building
    -> std.random.hex_string
  sip.response_for
    @ (request: sip_message, status: u16, reason: string) -> sip_message
    + returns a response copying Via, From, To, Call-ID, and CSeq from the request
    # response_building
  sip.client_transaction_new
    @ (request: sip_message) -> client_transaction
    + creates a transaction in the Calling state with the given INVITE or non-INVITE request
    # transactions
  sip.client_transaction_receive
    @ (txn: client_transaction, response: sip_message) -> client_transaction
    + advances the state machine by one response (Calling->Proceeding->Completed->Terminated)
    - returns unchanged transaction when the response does not match the outstanding request
    # transactions
  sip.server_transaction_new
    @ (request: sip_message) -> server_transaction
    + creates a transaction in the Trying state
    # transactions
  sip.server_transaction_send
    @ (txn: server_transaction, response: sip_message) -> server_transaction
    + advances the server state machine for provisional and final responses
    # transactions
  sip.dialog_from_2xx
    @ (request: sip_message, response: sip_message) -> result[dialog_state, string]
    + creates a dialog from a 2xx response to an INVITE, storing local/remote tags and remote target
    - returns error when the response lacks a To tag
    # dialog
  sip.dialog_build_request
    @ (dialog: dialog_state, method: string) -> tuple[sip_message, dialog_state]
    + builds an in-dialog request and increments the local CSeq
    # dialog
