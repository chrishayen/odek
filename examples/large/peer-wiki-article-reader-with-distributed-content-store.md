# Requirement: "a peer-to-peer encyclopedia article reader backed by a distributed content store"

Reads encyclopedia articles by title through a content-addressed peer-to-peer store. The project layer wires lookup, chunk fetching, and assembly; std provides the protocol primitives.

std
  std.hash
    std.hash.sha1
      fn (data: bytes) -> bytes
      + returns the 20-byte SHA-1 digest of data
      # hashing
  std.net
    std.net.tcp_connect
      fn (host: string, port: u16) -> result[conn, string]
      + opens a TCP connection to host:port
      - returns error when the host is unreachable
      # networking
    std.net.read_bytes
      fn (c: conn, n: i32) -> result[bytes, string]
      + reads exactly n bytes from the connection
      - returns error on early close
      # networking
  std.encoding
    std.encoding.bencode_decode
      fn (raw: bytes) -> result[map[string, string], string]
      + decodes a bencoded dictionary into string entries
      - returns error on malformed input
      # serialization

peer_wiki
  peer_wiki.open_index
    fn (index_path: string) -> result[wiki_index, string]
    + loads a content index mapping article titles to chunk hashes
    - returns error when the index file is missing
    # indexing
    -> std.encoding.bencode_decode
  peer_wiki.lookup_article
    fn (idx: wiki_index, title: string) -> result[list[bytes], string]
    + returns the chunk hashes that compose the article for the given title
    - returns error when the title is not present in the index
    # lookup
  peer_wiki.fetch_chunk
    fn (hash: bytes, peers: list[string]) -> result[bytes, string]
    + requests a chunk by hash from the first responsive peer and verifies its digest
    - returns error when no peer serves the chunk
    # chunk_fetch
    -> std.net.tcp_connect
    -> std.net.read_bytes
    -> std.hash.sha1
  peer_wiki.read_article
    fn (idx: wiki_index, title: string, peers: list[string]) -> result[string, string]
    + returns the full article text for the given title by fetching and joining its chunks
    - returns error when any chunk cannot be retrieved
    # article_read
  peer_wiki.announce_piece
    fn (hash: bytes, tracker_url: string) -> result[void, string]
    + notifies the tracker that this peer holds the given chunk
    - returns error when the tracker rejects the announce
    # seeding
