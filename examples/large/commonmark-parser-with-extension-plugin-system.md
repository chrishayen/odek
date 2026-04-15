# Requirement: "a CommonMark-compliant markdown parser with an extension plugin system"

Two-phase parser: block structure first, then inline parsing inside each block. Extensions can register block and inline rules.

std: (all units exist)

markdown
  markdown.new
    fn () -> markdown_state
    + creates a parser state with the default CommonMark rules registered
    # construction
  markdown.use
    fn (state: markdown_state, plugin: fn(markdown_state) -> markdown_state) -> markdown_state
    + applies a plugin that registers additional rules
    # extension
  markdown.add_block_rule
    fn (state: markdown_state, name: string, rule: block_rule) -> markdown_state
    + registers a new block rule by name
    # extension
  markdown.add_inline_rule
    fn (state: markdown_state, name: string, rule: inline_rule) -> markdown_state
    + registers a new inline rule by name
    # extension
  markdown.tokenize_blocks
    fn (state: markdown_state, source: string) -> list[token]
    + returns a flat token stream describing block structure
    + handles headings, lists, block quotes, code blocks, fences, and paragraphs
    # parsing
  markdown.tokenize_inline
    fn (state: markdown_state, source: string) -> list[token]
    + returns a flat token stream for inline content
    + handles emphasis, strong, links, images, code spans, and autolinks
    # parsing
  markdown.parse
    fn (state: markdown_state, source: string) -> list[token]
    + returns the full token stream by running block then inline tokenization
    # parsing
  markdown.render_html
    fn (state: markdown_state, tokens: list[token]) -> string
    + renders tokens to HTML
    + escapes raw text according to CommonMark rules
    # rendering
  markdown.render
    fn (state: markdown_state, source: string) -> string
    + convenience that parses and renders source in one call
    # rendering
