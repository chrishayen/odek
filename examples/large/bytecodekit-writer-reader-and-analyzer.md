# Requirement: "a toolkit for writing, reading, and analyzing virtual machine bytecode"

An assembler, disassembler, and a simple analyzer over instruction streams. Opcode metadata is authored once and consumed by each subsystem.

std
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + encodes bytes as lowercase hexadecimal
      # encoding
    std.encoding.hex_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes a hex string, accepting upper or lowercase
      - returns error on odd length or non-hex character
      # encoding

bytecodekit
  bytecodekit.opcode_info
    fn (mnemonic: string) -> optional[opcode_info]
    + returns metadata for the named opcode including byte value, operand size, and stack effect
    - returns none for unknown mnemonics
    # metadata
  bytecodekit.opcode_info_by_byte
    fn (byte_value: u8) -> optional[opcode_info]
    + returns metadata keyed by encoded byte value
    - returns none for undefined opcodes
    # metadata
  bytecodekit.assemble
    fn (source: string) -> result[bytes, string]
    + assembles textual mnemonics and hex literals into a bytecode blob
    - returns error with line number on unknown mnemonic or malformed operand
    -> std.encoding.hex_decode
    # assembly
  bytecodekit.disassemble
    fn (code: bytes) -> list[instruction]
    + returns the instruction list paired with program counter and operand bytes
    -> std.encoding.hex_encode
    # disassembly
  bytecodekit.format_listing
    fn (instructions: list[instruction]) -> string
    + renders a disassembled program as aligned text with pc, mnemonic, and operand
    # rendering
  bytecodekit.analyze_stack
    fn (instructions: list[instruction]) -> result[list[i32], string]
    + returns the stack depth after each instruction assuming linear flow
    - returns error when a pop would underflow
    # analysis
  bytecodekit.find_basic_blocks
    fn (instructions: list[instruction]) -> list[basic_block]
    + partitions the program into basic blocks using jump targets and fallthrough
    # analysis
  bytecodekit.reachable_instructions
    fn (instructions: list[instruction], entry: i64) -> list[i64]
    + returns the set of instruction pcs reachable from the entry point
    # analysis
  bytecodekit.validate
    fn (code: bytes) -> result[void, string]
    + verifies the bytecode decodes cleanly and every jump lands on a valid pc
    - returns error describing the first offending instruction
    # validation
  bytecodekit.compare
    fn (left: list[instruction], right: list[instruction]) -> list[diff_entry]
    + returns a list of differences between two disassembled programs
    # diff
