# Requirement: "a library for converting values between a host language and a dynamic guest language"

Bidirectional conversion between a statically-typed value model and a dynamic value model (numbers, strings, arrays, objects, null).

std: (all units exist)

interop
  interop.to_dynamic
    fn (value: host_value) -> dynamic_value
    + converts a host value into the dynamic representation
    ? integers become numbers; maps become objects; lists become arrays
    # conversion
  interop.from_dynamic
    fn (value: dynamic_value, expected: type_hint) -> result[host_value, string]
    + converts a dynamic value into the expected host type
    - returns error when the dynamic value cannot satisfy the expected type
    # conversion
  interop.encode_call
    fn (fn_name: string, args: list[host_value]) -> dynamic_value
    + builds a dynamic call payload {fn, args} for the guest
    # bridging
  interop.decode_return
    fn (payload: dynamic_value, expected: type_hint) -> result[host_value, string]
    + decodes a guest return payload into a host value
    - returns error on type mismatch
    # bridging
