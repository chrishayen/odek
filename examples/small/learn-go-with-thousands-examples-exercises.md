# Requirement: "a library of exercises and quizzes with scoring"

Store exercises keyed by id, accept answers, and track per-learner scores.

std: (all units exist)

exerkit
  exerkit.new_bank
    @ () -> exercise_bank
    + creates an empty exercise bank
    # construction
  exerkit.add_exercise
    @ (bank: exercise_bank, id: string, question: string, answer: string) -> exercise_bank
    + registers an exercise
    # catalog
  exerkit.submit
    @ (bank: exercise_bank, id: string, answer: string) -> result[bool, string]
    + returns true when the answer is correct
    - returns error when the exercise id is unknown
    # grading
  exerkit.record_score
    @ (scores: map[string,i32], learner: string, delta: i32) -> map[string,i32]
    + adds delta to the learner's score, creating the entry if needed
    # scoring
