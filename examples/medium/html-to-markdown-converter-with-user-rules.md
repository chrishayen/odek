# Requirement: "convert HTML to Markdown, extensible through user-supplied rules"

A tokenizing HTML parser feeds a walker that emits Markdown. Callers can register rules that override the default handling for a tag name.

std
  std.html
    std.html.tokenize
      fn (source: string) -> result[list[html_token], string]
      + tokenizes text, start tags, end tags, and comments
      - returns error on unclosed tags
      # parsing
    std.html.decode_entities
      fn (text: string) -> string
      + decodes named entities like &amp; and numeric entities like &#39;
      # parsing
  std.strings
    std.strings.trim
      fn (s: string) -> string
      + removes leading and trailing ASCII whitespace
      # strings

html2md
  html2md.new
    fn () -> converter_state
    + creates a converter with the built-in rule set
    # construction
  html2md.register_rule
    fn (state: converter_state, tag: string, rule: rule_fn) -> converter_state
    + installs a custom rule for the given tag, overriding any default
    ? rules receive the tag's attributes and its already-converted inner markdown
    # extensibility
  html2md.convert
    fn (state: converter_state, html: string) -> result[string, string]
    + renders headings h1..h6 as "#"..."######"
    + renders <a href="...">text</a> as "[text](url)"
    + renders <code> and <pre> as inline code and fenced blocks
    + renders <ul>/<ol> as bulleted or numbered lists
    + collapses whitespace inside block elements
    - returns error when the HTML cannot be tokenized
    # conversion
    -> std.html.tokenize
    -> std.html.decode_entities
    -> std.strings.trim
