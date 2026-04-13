# Requirement: "a library for reading messages from a disposable inbox provider"

Lists and fetches messages for an inbox alias using HTTP calls and HTML parsing. The specific provider is abstracted behind a configurable base URL.

std
  std.http
    std.http.get
      @ (url: string, headers: map[string, string]) -> result[string, string]
      + returns the response body as text
      - returns error on non-2xx response or network failure
      # http_client
  std.html
    std.html.select_text
      @ (document: string, selector: string) -> list[string]
      + returns the inner text of every element matching the CSS selector
      # parsing

disposable_inbox
  disposable_inbox.new_client
    @ (base_url: string) -> inbox_client
    + returns a client bound to the given provider base URL
    # construction
  disposable_inbox.list_messages
    @ (client: inbox_client, alias: string) -> result[list[message_summary], string]
    + returns summary records for every message in the inbox
    - returns error when the alias is empty
    - returns error when the provider responds with a non-2xx status
    # listing
    -> std.http.get
    -> std.html.select_text
  disposable_inbox.fetch_message
    @ (client: inbox_client, alias: string, message_id: string) -> result[message_body, string]
    + returns the decoded subject and body for a single message
    - returns error when the message cannot be found
    # fetching
    -> std.http.get
    -> std.html.select_text
