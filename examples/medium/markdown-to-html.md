# Requirement: "a markdown to HTML converter"

Parsing vs. rendering is a legitimate split. Block-level and inline-level handling are distinct stages.

std: (all units exist)

markdown
  markdown.to_html
    fn (md: string) -> string
    + converts "# Title" to "<h1>Title</h1>"
    + converts "**bold**" to "<p><strong>bold</strong></p>"
    + handles nested lists and fenced code blocks
    + produces well-formed HTML
    ? subset: headings, paragraphs, bold, italic, links, lists, code spans, code fences
    ? raw HTML passthrough is NOT supported; angle brackets in text are escaped
    # rendering
  markdown.parse_block
    fn (line: string) -> markdown_block
    + classifies a line as heading / list-item / code-fence / paragraph
    # parsing
  markdown.render_inline
    fn (text: string) -> string
    + renders bold, italic, links, and code spans into HTML
    + escapes <, >, and & in non-code text
    # rendering
