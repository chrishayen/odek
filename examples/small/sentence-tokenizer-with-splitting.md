# Requirement: "a sentence tokenizer that splits text into sentences"

Splits on terminal punctuation while avoiding common abbreviation false positives.

std: (all units exist)

sentences
  sentences.default_abbreviations
    fn () -> list[string]
    + returns a built-in list of common abbreviations that should not terminate a sentence
    # configuration
  sentences.tokenize
    fn (text: string) -> list[string]
    + splits text into sentences on '.', '!', '?' followed by whitespace and an uppercase letter
    + preserves terminal punctuation on each returned sentence
    + treats quoted punctuation as part of the preceding sentence
    - does not split after a token listed in default_abbreviations
    # tokenization
  sentences.tokenize_with
    fn (text: string, abbreviations: list[string]) -> list[string]
    + behaves like tokenize but uses the provided abbreviation list instead of the defaults
    # tokenization
