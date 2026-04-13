# Requirement: "a graph-oriented live-coding language for real-time audio synthesis"

Parse a small textual DSL where each line defines a named node and its inputs, compile the result into an audio processing graph, then render sample buffers on demand with hot-swap support.

std
  std.audio
    std.audio.sine
      @ (freq_hz: f64, phase: f64, sample_rate: i32) -> f64
      + returns one sample of a sine wave at the given frequency and phase
      # dsp_primitive
    std.audio.saw
      @ (freq_hz: f64, phase: f64, sample_rate: i32) -> f64
      + returns one sample of a sawtooth wave
      # dsp_primitive
    std.audio.lowpass_step
      @ (input: f64, cutoff_hz: f64, state: f64, sample_rate: i32) -> f64
      + advances a one-pole lowpass filter by one sample
      # dsp_primitive
  std.parse
    std.parse.tokenize_line
      @ (line: string) -> list[string]
      + splits a line on whitespace, preserving quoted substrings
      # lexing

live_graph
  live_graph.new
    @ (sample_rate: i32) -> graph_state
    + creates an empty graph targeting the given sample rate
    # construction
  live_graph.parse_source
    @ (source: string) -> result[list[node_decl], string]
    + parses the DSL source into a list of named node declarations
    - returns error with line number on syntax failure
    # parsing
    -> std.parse.tokenize_line
  live_graph.compile
    @ (state: graph_state, decls: list[node_decl]) -> result[graph_state, string]
    + resolves node references and builds an execution plan
    - returns error on unknown node kind
    - returns error on cyclic references
    # compilation
  live_graph.hot_swap
    @ (state: graph_state, source: string) -> result[graph_state, string]
    + reparses and recompiles while preserving matching node identities and their runtime state
    - returns error when the new source fails to parse or compile
    # live_coding
  live_graph.render_block
    @ (state: graph_state, frames: i32) -> tuple[list[f32], graph_state]
    + renders the next block of interleaved stereo samples
    # rendering
    -> std.audio.sine
    -> std.audio.saw
    -> std.audio.lowpass_step
  live_graph.set_param
    @ (state: graph_state, node_name: string, param: string, value: f64) -> result[graph_state, string]
    + updates a node parameter without recompiling the graph
    - returns error on unknown node or parameter
    # live_coding
  live_graph.snapshot
    @ (state: graph_state) -> bytes
    + serializes the current graph and its runtime state for sharing
    # collaboration
  live_graph.restore
    @ (snapshot: bytes) -> result[graph_state, string]
    + reconstructs a graph state from a snapshot
    - returns error on malformed snapshot
    # collaboration
  live_graph.diff_sources
    @ (old_source: string, new_source: string) -> list[string]
    + returns the node names that differ between two source versions
    # collaboration
