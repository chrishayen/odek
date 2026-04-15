# Requirement: "a standard library for the client-side web"

Wrappers around DOM, timers, storage and fetch exposed as a uniform surface. The browser host is abstracted behind opaque handles.

std: (all units exist)

webstd
  webstd.document_query
    fn (selector: string) -> optional[element_handle]
    + returns the first matching element or none
    # dom
  webstd.element_set_text
    fn (el: element_handle, text: string) -> void
    + replaces the element's text content with text
    # dom
  webstd.element_on_event
    fn (el: element_handle, event: string, handler: fn(event_obj) -> void) -> cleanup_handle
    + attaches handler for the named event and returns a disposer
    # dom
  webstd.set_timeout
    fn (delay_ms: i32, f: fn() -> void) -> timer_handle
    + schedules f to run once after delay_ms
    # timers
  webstd.clear_timer
    fn (t: timer_handle) -> void
    + cancels a timer if it has not yet fired
    # timers
  webstd.storage_get
    fn (key: string) -> optional[string]
    + returns the stored string for key or none
    # storage
  webstd.storage_set
    fn (key: string, value: string) -> result[void, string]
    + stores value under key
    - returns error when the quota is exceeded
    # storage
  webstd.fetch
    fn (url: string, method: string, body: optional[bytes], headers: map[string, string]) -> result[http_response, string]
    + returns status, headers and body on a successful request
    - returns error on a network failure
    # network
