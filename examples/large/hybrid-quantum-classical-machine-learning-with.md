# Requirement: "a hybrid quantum-classical machine learning library with automatic differentiation"

Quantum circuits are parameterized. The library simulates circuits on classical hardware, computes expectation values as the loss, and differentiates parameters via the parameter-shift rule so classical optimizers can train them.

std
  std.math
    std.math.sin
      @ (x: f64) -> f64
      + returns the sine of x
      # math
    std.math.cos
      @ (x: f64) -> f64
      + returns the cosine of x
      # math
    std.math.pi
      @ () -> f64
      + returns pi
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

qml
  qml.new_circuit
    @ (num_qubits: i32) -> circuit
    + creates an empty circuit over the given number of qubits
    - errors when num_qubits is less than one
    # construction
  qml.add_rx
    @ (c: circuit, qubit: i32, param_index: i32) -> circuit
    + appends a parameterized RX rotation on the given qubit
    # gate
  qml.add_ry
    @ (c: circuit, qubit: i32, param_index: i32) -> circuit
    + appends a parameterized RY rotation
    # gate
  qml.add_cnot
    @ (c: circuit, control: i32, target: i32) -> circuit
    + appends a CNOT entangling gate
    # gate
  qml.simulate
    @ (c: circuit, params: list[f64]) -> list[complex]
    + returns the final state vector after applying all gates from the zero state
    # simulation
    -> std.complex.mul
    -> std.complex.add
    -> std.math.sin
    -> std.math.cos
  qml.expectation_z
    @ (state: list[complex], qubit: i32) -> f64
    + returns <Z> on the given qubit from the state vector
    # measurement
  qml.loss
    @ (c: circuit, params: list[f64], observable: observable) -> f64
    + returns the scalar expectation value of the observable for the circuit
    # loss
  qml.gradient
    @ (c: circuit, params: list[f64], observable: observable) -> list[f64]
    + returns the gradient of the loss with respect to each parameter using the parameter-shift rule
    ? parameter-shift evaluates loss at params + pi/2 and params - pi/2 per parameter
    # differentiation
    -> std.math.pi
  qml.sgd_step
    @ (params: list[f64], grad: list[f64], learning_rate: f64) -> list[f64]
    + returns params - learning_rate * grad
    # optimization
  qml.train
    @ (c: circuit, initial_params: list[f64], observable: observable, learning_rate: f64, steps: i32) -> list[f64]
    + runs gradient descent for the given number of steps and returns final parameters
    # training
