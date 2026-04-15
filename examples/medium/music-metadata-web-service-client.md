# Requirement: "a client library for a music metadata web service"

Typed lookups against a read-only music metadata API. HTTP and XML parsing live in std so the project layer is just request shaping and response decoding.

std
  std.http
    std.http.get
      fn (url: string) -> result[bytes, string]
      + performs a GET request and returns the body on 2xx
      - returns error on non-2xx or transport failure
      # networking
  std.xml
    std.xml.parse
      fn (raw: bytes) -> result[xml_node, string]
      + parses XML into a tree
      - returns error on malformed XML
      # serialization
    std.xml.find_text
      fn (node: xml_node, path: string) -> optional[string]
      + returns the text at a slash-separated element path
      # serialization
  std.url
    std.url.encode_query
      fn (params: map[string,string]) -> string
      + returns a URL-encoded query string
      # encoding

music_meta
  music_meta.lookup_artist
    fn (mbid: string) -> result[artist, string]
    + returns an artist record with name, sort_name, and country for the given id
    - returns error on unknown id
    # lookup
    -> std.http.get
    -> std.xml.parse
    -> std.xml.find_text
  music_meta.lookup_release
    fn (mbid: string) -> result[release, string]
    + returns a release record with title, date, and track count
    - returns error on unknown id
    # lookup
    -> std.http.get
    -> std.xml.parse
    -> std.xml.find_text
  music_meta.search_artist
    fn (query: string, limit: i32) -> result[list[artist], string]
    + returns up to limit artists matching the query
    - returns error when limit <= 0
    # search
    -> std.url.encode_query
    -> std.http.get
    -> std.xml.parse
