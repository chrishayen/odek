# Requirement: "a lightweight text templating library with variable substitution and loops"

Templates use ${var} for substitution and %for/%endfor blocks for iteration. A small tokenizer feeds a recursive renderer.

std: (all units exist)

templating
  templating.tokenize
    @ (source: string) -> result[list[template_token], string]
    + produces text, substitution, loop-start, and loop-end tokens
    - returns error on an unclosed ${ substitution
    - returns error on a %for without a matching %endfor
    # tokenization
  templating.parse
    @ (tokens: list[template_token]) -> result[template_node, string]
    + builds a tree with nested loop nodes
    - returns error when %endfor does not match the innermost %for
    # parsing
  templating.render
    @ (node: template_node, context: map[string, template_value]) -> result[string, string]
    + substitutes ${var} from context and expands loops over list values
    - returns error when a referenced variable is missing
    - returns error when a loop target is not a list
    # rendering
  templating.compile
    @ (source: string) -> result[template_node, string]
    + tokenizes and parses in one call
    - returns error on any tokenize or parse failure
    # composition
