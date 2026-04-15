# Requirement: "a library that produces clickable hyperlinks for terminal output"

Wraps text and a URL in the OSC 8 escape sequence that supporting terminals render as a link.

std: (all units exist)

term_link
  term_link.format
    fn (text: string, url: string) -> string
    + returns the text wrapped in OSC 8 escape sequences referencing the URL
    + returns the plain text when url is empty
    # formatting
