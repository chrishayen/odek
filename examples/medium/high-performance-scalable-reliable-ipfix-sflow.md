# Requirement: "a collector for IPFIX, sFlow, and NetFlow network telemetry"

Decoders for three wire formats producing a unified flow-record type, and a reassembler that merges records into flows keyed by 5-tuple.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.net
    std.net.parse_ipv4
      @ (raw: bytes) -> result[string, string]
      + decodes 4 bytes as a dotted-quad address
      - returns error on length mismatch
      # networking

flow_collector
  flow_collector.decode_netflow_v5
    @ (wire: bytes) -> result[list[flow_record], string]
    + decodes a NetFlow v5 datagram into flow records
    - returns error on length mismatch or version mismatch
    # decoding
    -> std.net.parse_ipv4
  flow_collector.decode_ipfix
    @ (wire: bytes, templates: map[i32, ipfix_template]) -> result[tuple[list[flow_record], map[i32, ipfix_template]], string]
    + decodes an IPFIX message, updating the template map
    - returns error on unknown template or truncated record
    # decoding
    -> std.net.parse_ipv4
  flow_collector.decode_sflow_v5
    @ (wire: bytes) -> result[list[flow_record], string]
    + decodes an sFlow v5 datagram into flow records
    - returns error on malformed sample data
    # decoding
    -> std.net.parse_ipv4
  flow_collector.new_aggregator
    @ (flush_interval_sec: i64) -> aggregator_state
    + creates a flow aggregator with a flush cadence
    # construction
  flow_collector.ingest
    @ (state: aggregator_state, records: list[flow_record]) -> aggregator_state
    + merges records into existing flows keyed by 5-tuple
    # aggregation
    -> std.time.now_seconds
  flow_collector.flush_ready
    @ (state: aggregator_state) -> tuple[list[flow_record], aggregator_state]
    + returns flows whose flush interval has elapsed and removes them from state
    # flushing
    -> std.time.now_seconds
