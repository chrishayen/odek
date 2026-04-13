# Requirement: "a block-based system modeling and simulation framework"

Define blocks with typed ports, wire them into a graph, and step the simulation with a fixed timestep.

std: (all units exist)

pathsim
  pathsim.new_system
    @ () -> system_state
    + creates an empty system with no blocks and no connections
    # construction
  pathsim.add_block
    @ (system: system_state, name: string, step: fn(list[f64]) -> list[f64], input_count: i32, output_count: i32) -> system_state
    + registers a block with a step function mapping input port values to output port values
    # modeling
  pathsim.connect
    @ (system: system_state, src_block: string, src_port: i32, dst_block: string, dst_port: i32) -> result[system_state, string]
    + wires an output port of one block to an input port of another
    - returns error when a referenced block does not exist
    - returns error when the port index is out of range
    - returns error when the destination input is already connected
    # wiring
  pathsim.topological_order
    @ (system: system_state) -> result[list[string], string]
    + returns an execution order such that producers run before consumers
    - returns error when the system contains a cycle with no unit-delay break
    # scheduling
  pathsim.step
    @ (system: system_state, dt: f64) -> system_state
    + advances the simulation by dt seconds, running each block once in topological order
    # execution
    -> pathsim.topological_order
  pathsim.run_until
    @ (system: system_state, dt: f64, t_end: f64) -> system_state
    + repeatedly steps the system until simulation time reaches t_end
    # execution
    -> pathsim.step
  pathsim.add_probe
    @ (system: system_state, block: string, port: i32) -> system_state
    + marks an output port so its value is recorded at every step
    # observability
  pathsim.probe_history
    @ (system: system_state, block: string, port: i32) -> list[f64]
    + returns the recorded values for a probed port in chronological order
    - returns an empty list when the port was never probed
    # observability
  pathsim.integrator_block
    @ (initial: f64) -> block_spec
    + returns a block spec for an integrator: out += in * dt
    # library
  pathsim.gain_block
    @ (k: f64) -> block_spec
    + returns a block spec that multiplies its input by k
    # library
  pathsim.sum_block
    @ (input_count: i32) -> block_spec
    + returns a block spec that sums its inputs
    # library
