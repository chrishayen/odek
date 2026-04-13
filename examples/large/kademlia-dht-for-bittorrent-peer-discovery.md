# Requirement: "a Kademlia DHT for BitTorrent peer discovery"

Node id, routing table with k-buckets, and the four Kademlia RPCs wired to a bencode transport.

std
  std.hash
    std.hash.sha1
      @ (data: bytes) -> bytes
      + returns the 20-byte SHA-1 digest of the input
      # hashing
  std.random
    std.random.bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness
  std.net
    std.net.udp_send
      @ (host: string, port: u16, data: bytes) -> result[void, string]
      + sends a UDP datagram to host:port
      - returns error on resolve or send failure
      # networking
    std.net.udp_recv
      @ (timeout_ms: i32) -> result[udp_packet, string]
      + receives the next UDP datagram with sender address
      - returns error on timeout
      # networking
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns the current unix time in seconds
      # time

dht
  dht.new_node_id
    @ () -> bytes
    + generates a 20-byte random node id
    # identity
    -> std.random.bytes
  dht.distance
    @ (a: bytes, b: bytes) -> bytes
    + returns the XOR distance between two 20-byte ids
    # metric
  dht.new_table
    @ (self_id: bytes, k: i32) -> table_state
    + creates an empty routing table with 160 k-buckets of capacity k
    # construction
  dht.table_insert
    @ (table: table_state, id: bytes, host: string, port: u16) -> table_state
    + inserts or refreshes a contact in the appropriate bucket
    - drops the new contact when the bucket is full and the oldest is still alive
    # routing
    -> dht.distance
    -> std.time.now_seconds
  dht.table_closest
    @ (table: table_state, target: bytes, count: i32) -> list[contact]
    + returns the count contacts closest to target by XOR distance
    # routing
    -> dht.distance
  dht.encode_ping
    @ (transaction_id: bytes, self_id: bytes) -> bytes
    + encodes a ping query as bencoded bytes
    # rpc_codec
  dht.encode_find_node
    @ (transaction_id: bytes, self_id: bytes, target: bytes) -> bytes
    + encodes a find_node query
    # rpc_codec
  dht.encode_get_peers
    @ (transaction_id: bytes, self_id: bytes, info_hash: bytes) -> bytes
    + encodes a get_peers query
    # rpc_codec
  dht.encode_announce_peer
    @ (transaction_id: bytes, self_id: bytes, info_hash: bytes, port: u16, token: bytes) -> bytes
    + encodes an announce_peer query
    # rpc_codec
  dht.decode_message
    @ (data: bytes) -> result[dht_message, string]
    + parses a bencoded query or response into a tagged message
    - returns error on malformed bencoding or missing required keys
    # rpc_codec
  dht.iterative_find_node
    @ (table: table_state, target: bytes, alpha: i32) -> tuple[table_state, list[contact]]
    + performs an iterative find_node lookup against the alpha closest known contacts
    + converges when a round produces no contact closer than the current best
    # lookup
    -> dht.table_closest
    -> dht.encode_find_node
    -> std.net.udp_send
    -> std.net.udp_recv
