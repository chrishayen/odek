# Requirement: "strip leading whitespace from every line in a string"

Computes the common leading whitespace across non-empty lines and removes it from every line.

std: (all units exist)

stripindent
  stripindent.strip
    fn (text: string) -> string
    + removes the longest common leading whitespace prefix shared by non-empty lines
    + ignores blank lines when computing the shared prefix
    + returns the input unchanged when no common indent exists
    # dedent
