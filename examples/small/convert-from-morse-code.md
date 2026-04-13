# Requirement: "convert to and from morse code"

Two entry points sharing an internal alphabet table. Words are separated by " / " in morse output.

std: (all units exist)

morse
  morse.encode
    @ (text: string) -> result[string, string]
    + encodes "SOS" as "... --- ..."
    + encodes "HI THERE" with " / " between the two words
    + is case-insensitive for input letters
    - returns error when a character has no morse representation
    # encoding
  morse.decode
    @ (code: string) -> result[string, string]
    + decodes "... --- ..." to "SOS"
    + decodes " / " as a word separator
    - returns error when a token is not a valid morse symbol
    # decoding
