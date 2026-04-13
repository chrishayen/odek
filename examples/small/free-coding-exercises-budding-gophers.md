# Requirement: "a library providing coding-exercise prompts with expected solutions for learners"

Holds a catalog of exercises. Each exercise has a prompt and an expected output for an input; the library scores a learner-supplied solution.

std: (all units exist)

exercises
  exercises.new_catalog
    @ () -> catalog_state
    + returns an empty catalog
    # construction
  exercises.add_exercise
    @ (catalog: catalog_state, id: string, prompt: string, input: string, expected: string) -> catalog_state
    + adds an exercise with the given id, prompt, sample input, and expected output
    # catalog
  exercises.get_prompt
    @ (catalog: catalog_state, id: string) -> optional[string]
    + returns the prompt text for the exercise
    - returns none when the id is unknown
    # lookup
  exercises.score
    @ (catalog: catalog_state, id: string, learner_output: string) -> result[bool, string]
    + returns true when the learner output matches the expected output exactly
    - returns error when the id is unknown
    # grading
