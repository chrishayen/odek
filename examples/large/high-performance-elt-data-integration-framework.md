# Requirement: "an elt data integration framework with pluggable architecture"

Extract-Load-Transform pipelines where sources read from external systems and destinations write to warehouses. Plugins are registered by name and referenced in pipeline definitions.

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
  std.io
    std.io.read_stream_chunk
      @ (stream: stream_handle) -> result[optional[bytes], string]
      + returns the next chunk of bytes or none at end of stream
      - returns error on I/O failure
      # io

elt
  elt.framework_new
    @ () -> framework_state
    + returns an empty framework with no registered plugins
    # construction
  elt.register_source
    @ (framework: framework_state, name: string, source: source_plugin) -> result[framework_state, string]
    + registers a source plugin under the given name
    - returns error on duplicate name
    # plugins
  elt.register_destination
    @ (framework: framework_state, name: string, destination: destination_plugin) -> result[framework_state, string]
    + registers a destination plugin under the given name
    - returns error on duplicate name
    # plugins
  elt.register_transform
    @ (framework: framework_state, name: string, transform: transform_plugin) -> result[framework_state, string]
    + registers a transform plugin under the given name
    - returns error on duplicate name
    # plugins
  elt.pipeline_new
    @ (source_name: string, destination_name: string) -> pipeline_spec
    + creates a pipeline spec referencing registered plugins by name
    # pipelines
  elt.pipeline_add_transform
    @ (spec: pipeline_spec, transform_name: string) -> pipeline_spec
    + appends a transform stage to the pipeline spec
    # pipelines
  elt.pipeline_validate
    @ (framework: framework_state, spec: pipeline_spec) -> result[void, string]
    + checks that every referenced plugin is registered
    - returns error listing any missing plugin names
    # validation
  elt.extract_batch
    @ (framework: framework_state, source_name: string, cursor: string) -> result[tuple[list[map[string, string]], string], string]
    + pulls a batch of records from the source and returns (records, next_cursor)
    - returns error when the source is unknown
    # extraction
    -> std.serialization.decode_json
  elt.apply_transforms
    @ (framework: framework_state, spec: pipeline_spec, records: list[map[string, string]]) -> result[list[map[string, string]], string]
    + runs the pipeline's transforms in order
    - returns error when a transform fails
    # transformation
  elt.load_batch
    @ (framework: framework_state, destination_name: string, records: list[map[string, string]]) -> result[i64, string]
    + writes records to the destination and returns the number accepted
    - returns error when the destination is unknown
    # loading
    -> std.serialization.encode_json
  elt.run_once
    @ (framework: framework_state, spec: pipeline_spec, cursor: string) -> result[string, string]
    + runs extract, transform, and load once; returns the updated cursor
    - returns error when any stage fails
    # execution
