# Requirement: "a quantum computing SDK for building, simulating, and sampling from quantum circuits"

Circuit construction, state-vector simulation, measurement sampling, and a shot-based execution API.

std
  std.math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the square root of x
      # math
    std.math.sin
      @ (x: f64) -> f64
      + returns sine of x
      # math
    std.math.cos
      @ (x: f64) -> f64
      + returns cosine of x
      # math
  std.complex
    std.complex.make
      @ (re: f64, im: f64) -> complex
      + constructs a complex number
      # complex
    std.complex.mul
      @ (a: complex, b: complex) -> complex
      + returns a * b
      # complex
    std.complex.add
      @ (a: complex, b: complex) -> complex
      + returns a + b
      # complex
    std.complex.abs_sq
      @ (a: complex) -> f64
      + returns |a|^2
      # complex
  std.random
    std.random.new
      @ (seed: u64) -> rng_state
      + creates a seeded PRNG
      # random
    std.random.next_f64
      @ (state: rng_state) -> tuple[f64, rng_state]
      + returns a value in [0, 1) and the advanced state
      # random

qsdk
  qsdk.new_circuit
    @ (num_qubits: i32) -> circuit
    + creates an empty circuit with the given qubit count
    - errors when num_qubits is less than one
    # construction
  qsdk.h
    @ (c: circuit, qubit: i32) -> circuit
    + appends a Hadamard gate
    # gate
  qsdk.x
    @ (c: circuit, qubit: i32) -> circuit
    + appends a Pauli-X gate
    # gate
  qsdk.rz
    @ (c: circuit, qubit: i32, angle: f64) -> circuit
    + appends an RZ rotation
    # gate
    -> std.math.sin
    -> std.math.cos
  qsdk.cnot
    @ (c: circuit, control: i32, target: i32) -> circuit
    + appends a controlled-NOT gate
    # gate
  qsdk.measure
    @ (c: circuit, qubit: i32, classical_bit: i32) -> circuit
    + appends a measurement instruction targeting a classical register bit
    # gate
  qsdk.simulate
    @ (c: circuit) -> list[complex]
    + returns the final state vector before any measurements
    # simulation
    -> std.complex.mul
    -> std.complex.add
    -> std.math.sqrt
  qsdk.probabilities
    @ (state: list[complex]) -> list[f64]
    + returns the probability of each computational basis state
    # measurement
    -> std.complex.abs_sq
  qsdk.sample
    @ (probs: list[f64], rng: rng_state) -> tuple[i32, rng_state]
    + draws one basis state index according to the distribution
    # sampling
    -> std.random.next_f64
  qsdk.run
    @ (c: circuit, shots: i32, seed: u64) -> map[string, i32]
    + simulates the circuit and returns a histogram over shot basis-state strings
    # execution
    -> std.random.new
  qsdk.expectation
    @ (state: list[complex], observable: pauli_string) -> f64
    + returns the expectation value of a Pauli observable
    # measurement
  qsdk.tensor_product
    @ (a: list[complex], b: list[complex]) -> list[complex]
    + returns the tensor product of two state vectors
    # linear_algebra
    -> std.complex.mul
