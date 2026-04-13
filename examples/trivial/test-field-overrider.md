# Requirement: "a helper for overriding private object fields in tests"

Abstracts over the host reflection facility; the caller supplies an object and a field name.

std: (all units exist)

test_fields
  test_fields.set_private
    @ (target: object_ref, field_name: string, value: object_ref) -> result[void, string]
    + overwrites the named field on target with value regardless of visibility
    - returns error when the field does not exist
    - returns error when value is not assignable to the field
    # test_helper
  test_fields.get_private
    @ (target: object_ref, field_name: string) -> result[object_ref, string]
    + returns the current value of the named field
    - returns error when the field does not exist
    # test_helper
