# Requirement: "a library that inflects Russian names into grammatical cases"

Given a first, middle, or last name and a target case, return the inflected form using suffix rules.

std: (all units exist)

petrovich
  petrovich.load_rules
    fn (rules_yaml: string) -> result[rules_state, string]
    + parses a rules document describing suffix substitutions per name part and case
    - returns error when document has unknown sections
    # rules
  petrovich.inflect_first_name
    fn (state: rules_state, name: string, gender: gender_t, case: case_t) -> string
    + returns the inflected form according to the matching rule
    + returns the input unchanged when no rule matches
    # inflection
  petrovich.inflect_middle_name
    fn (state: rules_state, name: string, gender: gender_t, case: case_t) -> string
    + returns the patronymic inflected for the target case
    # inflection
  petrovich.inflect_last_name
    fn (state: rules_state, name: string, gender: gender_t, case: case_t) -> string
    + returns the last name inflected for the target case
    + handles indeclinable surnames by returning them unchanged
    # inflection
  petrovich.detect_gender
    fn (state: rules_state, middle_name: string) -> gender_t
    + returns masculine when the patronymic ends with a masculine suffix
    + returns feminine when the patronymic ends with a feminine suffix
    - returns androgynous when no suffix matches
    # detection
