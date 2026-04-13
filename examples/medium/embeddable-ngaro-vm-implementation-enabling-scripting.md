# Requirement: "an embeddable stack-machine VM for a minimal concatenative language"

Simple linear memory, data and return stacks, and a fixed opcode set. Host callers can load an image and single-step or run to completion.

std: (all units exist)

stackvm
  stackvm.new
    @ (memory_size: i32) -> vm_state
    + creates a VM with a zeroed linear memory of the given size
    # construction
  stackvm.load_image
    @ (vm: vm_state, image: list[i32]) -> result[vm_state, string]
    + writes the image words into memory starting at address 0
    - returns error when the image is larger than memory
    # loading
  stackvm.step
    @ (vm: vm_state) -> result[vm_state, string]
    + fetches one instruction and executes it, advancing the program counter
    - returns error when the program counter points outside memory
    - returns error on stack underflow
    # execution
  stackvm.run
    @ (vm: vm_state, max_steps: i64) -> result[vm_state, string]
    + steps until the VM halts or max_steps is reached
    + returns the state with its halted flag set on normal termination
    # execution
  stackvm.push_data
    @ (vm: vm_state, value: i32) -> vm_state
    + pushes a value onto the data stack
    # stacks
  stackvm.pop_data
    @ (vm: vm_state) -> result[tuple[i32, vm_state], string]
    + pops the top of the data stack
    - returns error on empty stack
    # stacks
  stackvm.register_io
    @ (vm: vm_state, port: i32, handler: io_handler) -> vm_state
    + binds a host handler to an I/O port for the OUT instruction
    # host_binding
