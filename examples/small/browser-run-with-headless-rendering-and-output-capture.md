# Requirement: "a library for running source code inside a headless browser and capturing output"

Starts a headless browser session, loads a script, and collects the console output and exit status. The browser driver is an injected primitive.

std
  std.browser
    std.browser.launch
      fn () -> result[browser_session, string]
      + starts a headless browser session
      - returns error when no browser is available
      # browser
    std.browser.eval_script
      fn (session: browser_session, script: string) -> result[browser_run, string]
      + evaluates the script in a fresh page and returns the console output and exit code
      - returns error when the script throws an uncaught exception
      # browser
    std.browser.close
      fn (session: browser_session) -> result[void, string]
      + closes the browser session
      # browser

browser_run
  browser_run.run
    fn (script: string) -> result[run_result, string]
    + launches a browser, evaluates the script, and returns captured stdout and exit code
    - returns error when the browser cannot be launched
    - returns error when the script throws
    # execution
    -> std.browser.launch
    -> std.browser.eval_script
    -> std.browser.close
