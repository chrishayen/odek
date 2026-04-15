# Requirement: "a collection of code snippets and tutorials organized for lookup"

A library that holds snippets and multi-step tutorials, looks them up by topic, and walks through tutorial steps one at a time.

std: (all units exist)

code_corpus
  code_corpus.new
    fn () -> corpus_state
    + creates an empty corpus
    # construction
  code_corpus.add_snippet
    fn (state: corpus_state, topic: string, title: string, body: string) -> corpus_state
    + adds a snippet under a topic
    # registration
  code_corpus.add_tutorial
    fn (state: corpus_state, topic: string, title: string, steps: list[string]) -> corpus_state
    + adds a tutorial under a topic
    # registration
  code_corpus.find_by_topic
    fn (state: corpus_state, topic: string) -> list[entry]
    + returns snippets and tutorials listed under a topic
    # lookup
  code_corpus.tutorial_step
    fn (state: corpus_state, title: string, index: i32) -> result[string, string]
    + returns the step at the given index of a tutorial
    - returns error when the tutorial is unknown
    - returns error when the index is out of range
    # tutorial_walk
