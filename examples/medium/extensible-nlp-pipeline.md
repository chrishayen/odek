# Requirement: "an extensible natural language processing pipeline"

A tokenize-then-transform pipeline where plugins register by name and run in order over a shared document state.

std
  std.text
    std.text.split_sentences
      @ (text: string) -> list[string]
      + splits text on sentence-ending punctuation into non-empty sentences
      # text
    std.text.split_words
      @ (sentence: string) -> list[string]
      + splits a sentence into lowercased word tokens, dropping punctuation
      # text

nlp
  nlp.new_document
    @ (text: string) -> document_state
    + creates a document containing the raw text and empty annotations
    # construction
    -> std.text.split_sentences
    -> std.text.split_words
  nlp.register_plugin
    @ (pipeline: pipeline_state, name: string, fn: fn(document_state) -> document_state) -> pipeline_state
    + adds a plugin to the pipeline under the given name, preserving registration order
    - returns pipeline unchanged when a plugin with the same name already exists
    # registration
  nlp.new_pipeline
    @ () -> pipeline_state
    + creates an empty pipeline with no plugins registered
    # construction
  nlp.run
    @ (pipeline: pipeline_state, doc: document_state) -> document_state
    + applies each registered plugin in order and returns the final document state
    # execution
  nlp.get_annotation
    @ (doc: document_state, key: string) -> optional[string]
    + returns the annotation stored under key, if any
    # annotations
  nlp.set_annotation
    @ (doc: document_state, key: string, value: string) -> document_state
    + returns a new document with the annotation attached
    # annotations
