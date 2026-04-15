# Requirement: "an SMTP server protocol state machine"

A pure state machine: caller feeds lines, receives responses and transitions. No sockets.

std: (all units exist)

smtp_sm
  smtp_sm.new
    fn (hostname: string) -> smtp_state
    + creates a fresh session in the "expect HELO/EHLO" phase
    # construction
  smtp_sm.feed_line
    fn (state: smtp_state, line: string) -> tuple[smtp_response, smtp_state]
    + handles one command line and returns the reply plus the updated state
    - returns 503 "bad sequence" when a command is sent out of order
    - returns 500 "unknown command" on unrecognized verbs
    - returns 501 on syntax errors in MAIL FROM / RCPT TO arguments
    # protocol_dispatch
  smtp_sm.append_data_line
    fn (state: smtp_state, line: string) -> tuple[optional[smtp_response], smtp_state]
    + appends a DATA-phase line to the message body; returns none while collecting
    + returns 250 "ok" when the terminating "." line is seen
    ? leading-dot unescaping is applied as required by the protocol
    # data_collection
  smtp_sm.current_phase
    fn (state: smtp_state) -> string
    + returns one of "greet", "mail", "rcpt", "data", "quit"
    # inspection
  smtp_sm.pending_message
    fn (state: smtp_state) -> optional[smtp_message]
    + returns the fully assembled sender, recipients, and body after a successful DATA
    + returns none before DATA completes
    # extraction
  smtp_sm.reset
    fn (state: smtp_state) -> smtp_state
    + clears MAIL/RCPT/DATA but keeps the session open (RSET)
    # protocol_reset
