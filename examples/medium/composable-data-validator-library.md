# Requirement: "a composable data validation library"

Validators are values that can be combined; applying one to data returns either the coerced value or a list of errors.

std: (all units exist)

validator
  validator.string
    fn () -> validator_handle
    + returns a validator that accepts strings and rejects other types
    # primitives
  validator.number
    fn () -> validator_handle
    + returns a validator that accepts numeric values
    # primitives
  validator.boolean
    fn () -> validator_handle
    + returns a validator that accepts boolean values
    # primitives
  validator.optional
    fn (inner: validator_handle) -> validator_handle
    + returns a validator that also accepts null
    # combinators
  validator.list_of
    fn (element: validator_handle) -> validator_handle
    + returns a validator that checks every item against element
    # combinators
  validator.object
    fn (fields: map[string, validator_handle]) -> validator_handle
    + returns a validator that checks each named field with its validator
    + collects errors from all fields rather than stopping at the first
    # combinators
  validator.refine
    fn (inner: validator_handle, predicate: refine_fn, message: string) -> validator_handle
    + returns a validator that runs inner, then asserts predicate on the coerced value
    # combinators
  validator.validate
    fn (v: validator_handle, data: json_value) -> result[json_value, list[validation_error]]
    + returns the coerced value when data satisfies v, or the accumulated errors otherwise
    # execution
  validator.error_path
    fn (err: validation_error) -> string
    + returns a dotted path describing where the error occurred
    # reporting
