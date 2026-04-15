# Requirement: "a one-way string transliterator with language-specific rules"

Converts non-ASCII text to an ASCII approximation, with a default Unicode-to-ASCII table that can be overridden per language.

std: (all units exist)

translit
  translit.new
    fn () -> translit_state
    + creates a transliterator with the default Unicode fallback mapping
    # construction
  translit.register_language
    fn (state: translit_state, language: string, rules: map[string, string]) -> translit_state
    + adds a rule set that is applied when the given language is selected
    ? later registrations of the same language replace earlier ones
    # configuration
  translit.transliterate
    fn (state: translit_state, text: string) -> string
    + returns the ASCII approximation using the default mapping
    + leaves unmapped characters unchanged
    # transliteration
  translit.transliterate_lang
    fn (state: translit_state, text: string, language: string) -> string
    + applies the language-specific rules first, then the default mapping for remaining characters
    - falls back to the default mapping when the language is not registered
    # transliteration
