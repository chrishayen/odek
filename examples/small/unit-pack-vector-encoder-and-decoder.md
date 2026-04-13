# Requirement: "a compact encoder and decoder for unit-length 3D vectors"

Packs a unit vector by dropping the largest component and quantizing the other two plus two bits identifying which axis was dropped.

std: (all units exist)

unit_pack
  unit_pack.pack
    @ (x: f32, y: f32, z: f32) -> bytes
    + packs a unit vector into a 4-byte representation (2 bits axis selector + two 15-bit quantized components)
    ? the caller guarantees the input has magnitude near 1
    # encoding
  unit_pack.unpack
    @ (packed: bytes) -> result[tuple[f32, f32, f32], string]
    + reconstructs an approximate unit vector from its packed form
    - returns error when packed is not exactly 4 bytes
    # decoding
