# Requirement: "a greeting library for a programming language learning course"

The input was a course description, not a library idea. Interpreted as a minimal greeter for a learner's first program.

std: (all units exist)

learner
  learner.welcome
    @ (name: string) -> string
    + returns a welcome message addressed to the learner
    + returns a generic welcome when name is empty
    # greeting
