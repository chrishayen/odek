# Requirement: "a QR code generator (ASCII and PNG) for SEPA credit transfer payloads"

Builds the SEPA payload string, encodes it as a QR code, and renders to either ASCII or PNG. Std holds the generic QR and PNG primitives.

std
  std.qr
    std.qr.encode
      fn (data: string, ecc_level: string) -> result[list[list[bool]], string]
      + returns a square matrix of modules for the given payload
      - returns error when payload exceeds capacity for the chosen ecc level
      # qr_encoding
  std.image
    std.image.write_png
      fn (pixels: list[list[bool]], scale: i32) -> bytes
      + returns PNG bytes with each true module drawn as a black square
      ? scale is pixels per module
      # image_encoding

sepa_qr
  sepa_qr.build_payload
    fn (beneficiary: string, iban: string, amount_cents: i64, reference: string) -> result[string, string]
    + returns the SEPA credit transfer payload following the EPC QR guidelines
    - returns error when iban is empty or amount is negative
    # payload_construction
  sepa_qr.render_ascii
    fn (beneficiary: string, iban: string, amount_cents: i64, reference: string) -> result[string, string]
    + returns a multi-line string where each module is two characters wide
    - returns error when the payload cannot be encoded
    # ascii_rendering
    -> std.qr.encode
  sepa_qr.render_png
    fn (beneficiary: string, iban: string, amount_cents: i64, reference: string, scale: i32) -> result[bytes, string]
    + returns PNG bytes encoding the SEPA payment QR
    - returns error when scale is less than 1
    # png_rendering
    -> std.qr.encode
    -> std.image.write_png
