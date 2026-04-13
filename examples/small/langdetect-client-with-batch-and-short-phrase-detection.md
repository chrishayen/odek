# Requirement: "a language detection client supporting batch requests and short-phrase or single-word detection"

The project layer is a thin client over an HTTP transport primitive. Batch and single-phrase detection share one endpoint.

std
  std.http
    std.http.post_json
      @ (url: string, headers: map[string,string], body: string) -> result[string, string]
      + performs a POST with the given JSON body and returns the response body
      - returns error on non-2xx status or transport failure
      # http
  std.json
    std.json.encode_list_of_strings
      @ (items: list[string]) -> string
      + encodes a list of strings as a JSON array
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string,string], string]
      + parses a flat JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

langdetect
  langdetect.new_client
    @ (api_key: string, endpoint: string) -> langdetect_client
    + stores credentials and endpoint for subsequent calls
    - accepts an empty endpoint and uses a sensible default
    # construction
  langdetect.detect
    @ (client: langdetect_client, text: string) -> result[string, string]
    + returns the detected language code for a single phrase
    - returns error when the text is empty
    - returns error when the response cannot be parsed
    # single_detection
    -> std.http.post_json
    -> std.json.parse_object
  langdetect.detect_batch
    @ (client: langdetect_client, texts: list[string]) -> result[list[string], string]
    + returns one language code per input text in order
    - returns error when the input list is empty
    # batch_detection
    -> std.http.post_json
    -> std.json.encode_list_of_strings
