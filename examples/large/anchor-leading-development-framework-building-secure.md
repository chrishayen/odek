# Requirement: "a development framework for building on-chain programs on a smart-contract platform"

A framework that lets callers declare program accounts and instructions, generates serialization code, and dispatches incoming instruction data to typed handlers.

std
  std.encoding
    std.encoding.leb128_encode_u64
      @ (n: u64) -> bytes
      + encodes n as unsigned LEB128
      # encoding
    std.encoding.leb128_decode_u64
      @ (data: bytes) -> result[tuple[u64, i32], string]
      + decodes a LEB128 value and returns the value and bytes consumed
      - returns error on truncated input
      # encoding
  std.hash
    std.hash.sha256
      @ (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest of data
      # hashing

onchain_framework
  onchain_framework.declare_program
    @ (name: string, program_id: bytes) -> program_def
    + creates an empty program definition with no instructions or accounts
    # declaration
  onchain_framework.declare_account
    @ (p: program_def, name: string, fields: list[tuple[string, string]]) -> program_def
    + registers a named account layout with typed fields
    # declaration
  onchain_framework.declare_instruction
    @ (p: program_def, name: string, arg_types: list[string], account_refs: list[string]) -> program_def
    + registers an instruction with its argument types and referenced account names
    # declaration
  onchain_framework.instruction_discriminator
    @ (p: program_def, name: string) -> result[bytes, string]
    + returns the 8-byte discriminator derived from SHA-256 of the instruction name
    - returns error when the instruction is not declared
    # dispatch
    -> std.hash.sha256
  onchain_framework.serialize_args
    @ (p: program_def, instruction: string, values: list[bytes]) -> result[bytes, string]
    + serializes values according to the declared argument types of the instruction, prefixed by its discriminator
    - returns error when the arity does not match
    # serialization
    -> std.encoding.leb128_encode_u64
  onchain_framework.deserialize_args
    @ (p: program_def, data: bytes) -> result[tuple[string, list[bytes]], string]
    + reads the discriminator and argument bytes, returning the instruction name and per-arg slices
    - returns error when the discriminator matches no declared instruction
    # deserialization
    -> std.encoding.leb128_decode_u64
  onchain_framework.derive_pda
    @ (program_id: bytes, seeds: list[bytes]) -> bytes
    + returns a program-derived address for the given seeds
    # addressing
    -> std.hash.sha256
  onchain_framework.dispatch
    @ (p: program_def, handlers: map[string, fn(list[bytes]) -> result[void, string]], data: bytes) -> result[void, string]
    + deserializes data and invokes the handler registered for the matching instruction
    - returns error when no handler is registered for the instruction
    # dispatch
  onchain_framework.validate_account
    @ (p: program_def, account_name: string, data: bytes) -> result[void, string]
    + verifies that data conforms to the declared layout for the named account
    - returns error when the size or any field is malformed
    # validation
