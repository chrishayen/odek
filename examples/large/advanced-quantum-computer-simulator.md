# Requirement: "a quantum computer simulator library"

Simulates a multi-qubit state vector under common gates and measurement.

std
  std.math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the non-negative square root
      - returns NaN for negative input
      # math
    std.math.sin
      @ (x: f64) -> f64
      + returns the sine of x in radians
      # math
    std.math.cos
      @ (x: f64) -> f64
      + returns the cosine of x in radians
      # math
  std.random
    std.random.uniform
      @ () -> f64
      + returns a uniform random f64 in [0.0, 1.0)
      # random

quantum_sim
  quantum_sim.new_register
    @ (n_qubits: i32) -> qstate
    + returns a register of n_qubits initialized to |0...0>
    - returns error-marker state when n_qubits is less than 1 or greater than 24
    # construction
  quantum_sim.apply_h
    @ (state: qstate, target: i32) -> qstate
    + applies a Hadamard gate to the target qubit
    # gates
    -> std.math.sqrt
  quantum_sim.apply_x
    @ (state: qstate, target: i32) -> qstate
    + applies a Pauli-X (bit flip) to the target qubit
    # gates
  quantum_sim.apply_z
    @ (state: qstate, target: i32) -> qstate
    + applies a Pauli-Z gate to the target qubit
    # gates
  quantum_sim.apply_rx
    @ (state: qstate, target: i32, theta: f64) -> qstate
    + applies a rotation around the X axis by angle theta
    # gates
    -> std.math.sin
    -> std.math.cos
  quantum_sim.apply_cnot
    @ (state: qstate, control: i32, target: i32) -> qstate
    + flips target when control is |1>
    - returns error-marker state when control == target
    # gates
  quantum_sim.apply_swap
    @ (state: qstate, a: i32, b: i32) -> qstate
    + swaps the amplitudes associated with qubits a and b
    # gates
  quantum_sim.probability
    @ (state: qstate, basis: i64) -> f64
    + returns the probability of observing basis state
    + probabilities across all basis states sum to 1
    # measurement
  quantum_sim.measure
    @ (state: qstate, target: i32) -> tuple[i32, qstate]
    + returns (bit, collapsed_state) sampling proportionally to amplitudes
    # measurement
    -> std.random.uniform
  quantum_sim.measure_all
    @ (state: qstate) -> tuple[i64, qstate]
    + returns a sampled bit pattern and the collapsed basis state
    # measurement
    -> std.random.uniform
  quantum_sim.tensor_product
    @ (a: qstate, b: qstate) -> qstate
    + returns the tensor product of two registers
    # composition
  quantum_sim.amplitudes
    @ (state: qstate) -> list[complex64]
    + returns the raw amplitude vector
    # introspection
