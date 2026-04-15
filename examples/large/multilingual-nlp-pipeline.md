# Requirement: "a multilingual natural language processing pipeline"

A pipeline library that runs tokenization, sentence splitting, part-of-speech tagging, lemmatization, dependency parsing, and named-entity recognition. Models are loaded per language.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the path does not exist
      # filesystem
  std.text
    std.text.normalize_nfc
      fn (s: string) -> string
      + returns the NFC-normalized form of s
      # unicode
    std.text.is_letter
      fn (codepoint: i32) -> bool
      + returns true when the codepoint is a letter in any script
      # unicode

nlp_pipeline
  nlp_pipeline.load_model
    fn (language: string, path: string) -> result[model_state, string]
    + loads tokenizer, tagger, parser and ner weights for the given language
    - returns error when the language code is not supported
    - returns error when path is missing required files
    # model_loading
    -> std.fs.read_all
  nlp_pipeline.tokenize
    fn (m: model_state, text: string) -> list[token_span]
    + returns tokens with byte offsets covering text
    + handles language-specific clitics and contractions when the model supports them
    # tokenization
    -> std.text.normalize_nfc
    -> std.text.is_letter
  nlp_pipeline.split_sentences
    fn (m: model_state, tokens: list[token_span]) -> list[sentence_span]
    + returns sentence ranges over the token list
    # segmentation
  nlp_pipeline.tag_pos
    fn (m: model_state, sentence: sentence_span, tokens: list[token_span]) -> list[pos_tag]
    + returns a universal POS tag per token in the sentence
    # tagging
  nlp_pipeline.lemmatize
    fn (m: model_state, tokens: list[token_span], tags: list[pos_tag]) -> list[string]
    + returns the lemma per token using the tag as a hint
    # lemmatization
  nlp_pipeline.parse_dependencies
    fn (m: model_state, sentence: sentence_span, tokens: list[token_span], tags: list[pos_tag]) -> list[dep_edge]
    + returns head index and label per token; the root token has head index -1
    # parsing
  nlp_pipeline.recognize_entities
    fn (m: model_state, sentence: sentence_span, tokens: list[token_span]) -> list[entity_span]
    + returns entity spans with labels over the token range
    # ner
  nlp_pipeline.annotate
    fn (m: model_state, text: string) -> document
    + runs the full pipeline and returns a document with sentences, tokens, tags, lemmas, dependencies and entities
    # orchestration
    -> nlp_pipeline.tokenize
    -> nlp_pipeline.split_sentences
    -> nlp_pipeline.tag_pos
    -> nlp_pipeline.lemmatize
    -> nlp_pipeline.parse_dependencies
    -> nlp_pipeline.recognize_entities
  nlp_pipeline.document_to_json
    fn (doc: document) -> string
    + renders the document as a stable JSON string for inspection
    # serialization
