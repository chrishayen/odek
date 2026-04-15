# Requirement: "convert HTML to Markdown-formatted text"

A simpler, rule-less HTML-to-Markdown converter: tokenize, walk, emit.

std
  std.html
    std.html.tokenize
      fn (source: string) -> result[list[html_token], string]
      + tokenizes text, start tags, end tags, and comments
      - returns error on unclosed tags
      # parsing
    std.html.decode_entities
      fn (text: string) -> string
      + decodes named and numeric HTML entities
      # parsing

html2text
  html2text.convert
    fn (html: string) -> result[string, string]
    + renders headings h1..h6 with leading "#" characters
    + renders paragraphs separated by blank lines
    + renders <a href="u">t</a> as "[t](u)"
    + renders <strong> as "**...**" and <em> as "*...*"
    + renders <ul>/<ol> as bulleted or numbered lists
    - returns error when the HTML cannot be tokenized
    # conversion
    -> std.html.tokenize
    -> std.html.decode_entities
