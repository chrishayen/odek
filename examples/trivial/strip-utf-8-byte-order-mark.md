# Requirement: "strip a UTF-8 byte order mark from a string or byte buffer"

Removes a leading UTF-8 BOM (EF BB BF) if present.

std: (all units exist)

stripbom
  stripbom.from_string
    @ (text: string) -> string
    + removes the leading U+FEFF character when present
    + returns the input unchanged when no BOM is present
    # bom
  stripbom.from_bytes
    @ (data: bytes) -> bytes
    + removes the leading EF BB BF byte sequence when present
    + returns the input unchanged when no BOM is present
    # bom
