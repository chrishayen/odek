# Requirement: "a configurable fake SMTP server library that records received messages for test assertions"

Implements the SMTP state machine for a session, configurable to reject specific commands, and records captured envelopes for later inspection.

std: (all units exist)

fake_smtp
  fake_smtp.new
    @ (greeting: string) -> server_state
    + creates a server with the given banner and no captured messages
    # construction
  fake_smtp.configure_reject
    @ (s: server_state, verb: string, code: i32, text: string) -> void
    + makes the server reply with the given code+text whenever verb is issued
    ? verb is one of "HELO", "MAIL", "RCPT", "DATA"
    # configuration
  fake_smtp.new_session
    @ (s: server_state) -> session_state
    + starts a session in the initial state expecting HELO/EHLO
    # session
  fake_smtp.feed
    @ (sess: session_state, line: string) -> feed_response
    + advances the session state machine for one command line
    + returns (reply_code, reply_text, session_done)
    - returns a 503 reply when a command is issued out of order
    - returns a configured rejection when one applies
    # session
  fake_smtp.captured
    @ (s: server_state) -> list[captured_message]
    + returns all envelopes captured so far (from, to, data)
    # inspection
  fake_smtp.reset
    @ (s: server_state) -> void
    + clears captured messages while keeping configuration
    # inspection
