# Requirement: "a function that escapes regular expression special characters in a string"

Used to take user input and splice it safely into a regex pattern. One rune.

std: (all units exist)

regex_escape
  regex_escape.escape
    fn (input: string) -> string
    + prefixes each regex metacharacter in the input with a backslash
    + returns the input unchanged when it contains no metacharacters
    + returns "" for the empty string
    ? the escaped set includes . * + ? ^ $ { } ( ) | [ ] \ and forward slash
    # escaping
