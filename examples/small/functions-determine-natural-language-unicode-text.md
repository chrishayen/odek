# Requirement: "a natural-language detector for unicode text"

N-gram frequency compared against a bundled set of language profiles. Script detection narrows the candidates first.

std: (all units exist)

lang_detect
  lang_detect.script_of
    @ (text: string) -> string
    + returns the dominant Unicode script name (e.g. "Latin", "Cyrillic", "Han")
    + returns "Unknown" when text has no letters
    # preselection
  lang_detect.profile_of
    @ (text: string) -> map[string, f64]
    + returns a normalized trigram frequency profile of the input
    + ignores non-letter characters and lowercases letters
    # profiling
  lang_detect.distance
    @ (profile: map[string, f64], reference: map[string, f64]) -> f64
    + returns an out-of-place distance between two trigram profiles
    + identical profiles return 0.0
    # scoring
  lang_detect.detect
    @ (text: string) -> result[string, string]
    + returns the ISO 639-1 code of the most likely language
    - returns error when text is shorter than the minimum sample length
    # detection
    -> lang_detect.script_of
    -> lang_detect.profile_of
    -> lang_detect.distance
