# Requirement: "a library for building chat bots with intent classification, entity extraction, sentiment analysis, and language identification"

A full NLP toolkit for bot authoring.

std
  std.strings
    std.strings.lower
      @ (s: string) -> string
      + lowercases ASCII
      # strings
    std.strings.tokenize
      @ (s: string) -> list[string]
      + splits on whitespace and punctuation boundaries
      # strings
  std.math
    std.math.log
      @ (x: f64) -> f64
      + natural log
      # math

botkit
  botkit.new
    @ () -> bot_state
    + creates an empty bot with no intents or entities
    # construction
  botkit.add_intent
    @ (state: bot_state, name: string, examples: list[string]) -> bot_state
    + registers an intent with training examples
    # training
    -> std.strings.tokenize
    -> std.strings.lower
  botkit.add_entity
    @ (state: bot_state, name: string, values: list[string]) -> bot_state
    + registers a gazetteer-style entity with known values
    # training
  botkit.train
    @ (state: bot_state) -> bot_state
    + computes per-intent token weights from the registered examples
    # training
    -> std.math.log
  botkit.classify_intent
    @ (state: bot_state, utterance: string) -> intent_result
    + returns the best-scoring intent and its confidence
    # inference
    -> std.strings.tokenize
    -> std.strings.lower
  botkit.extract_entities
    @ (state: bot_state, utterance: string) -> list[entity_span]
    + returns spans for every entity value found in the utterance
    # inference
    -> std.strings.lower
  botkit.sentiment
    @ (utterance: string) -> f32
    + returns a score from -1 (negative) to 1 (positive) using a lexicon
    # sentiment
    -> std.strings.tokenize
    -> std.strings.lower
  botkit.identify_language
    @ (text: string) -> string
    + returns an ISO 639-1 code based on character n-gram profiles
    - returns "und" when the text is too short to classify
    # language_id
  botkit.add_response
    @ (state: bot_state, intent: string, reply: string) -> bot_state
    + registers a canned reply for an intent
    # dialog
  botkit.respond
    @ (state: bot_state, utterance: string) -> bot_reply
    + classifies, extracts entities, and returns a reply plus structured metadata
    # dialog
