# Requirement: "a logs, metrics, and events router"

Routes observability records from source stages through transform stages to sink stages. Sources, transforms, and sinks are named components wired into a topology.

std
  std.serialization
    std.serialization.encode_json
      @ (value: map[string, string]) -> string
      + encodes a string-keyed map as JSON
      # serialization
    std.serialization.decode_json
      @ (raw: string) -> result[map[string, string], string]
      + decodes a JSON object into a string-keyed map
      - returns error on invalid JSON
      # serialization

obsrouter
  obsrouter.topology_new
    @ () -> topology_state
    + returns an empty topology
    # construction
  obsrouter.add_source
    @ (topo: topology_state, name: string, source_id: string) -> result[topology_state, string]
    + registers a source stage
    - returns error on duplicate name
    # topology
  obsrouter.add_transform
    @ (topo: topology_state, name: string, transform_id: string, inputs: list[string]) -> result[topology_state, string]
    + registers a transform reading from the named upstream stages
    - returns error on duplicate name or missing upstream
    # topology
  obsrouter.add_sink
    @ (topo: topology_state, name: string, sink_id: string, inputs: list[string]) -> result[topology_state, string]
    + registers a sink reading from the named upstream stages
    - returns error on duplicate name or missing upstream
    # topology
  obsrouter.validate
    @ (topo: topology_state) -> result[void, string]
    + checks that the graph is acyclic and every stage is reachable
    - returns error listing cycles or orphan stages
    # validation
  obsrouter.ingest
    @ (topo: topology_state, source_name: string, record: map[string, string]) -> result[topology_state, string]
    + pushes a record into the named source
    - returns error when the source is unknown
    # ingestion
    -> std.serialization.decode_json
  obsrouter.drain
    @ (topo: topology_state) -> map[string, list[map[string, string]]]
    + runs the topology until all records exit, returning per-sink outputs
    # execution
    -> std.serialization.encode_json
  obsrouter.classify
    @ (record: map[string, string]) -> string
    + returns "log", "metric", or "event" based on the record's shape
    # classification
