# Requirement: "a library for calling native functions in shared libraries without an FFI compiler"

Dynamically loads a shared library, resolves symbols, marshals arguments, and invokes native functions at runtime.

std: (all units exist)

dynlib
  dynlib.load
    @ (path: string) -> result[library_handle, string]
    + opens the shared library at path
    - returns error when the library cannot be found or loaded
    # loader
  dynlib.close
    @ (lib: library_handle) -> void
    + unloads the shared library
    # loader
  dynlib.lookup
    @ (lib: library_handle, symbol: string) -> result[symbol_handle, string]
    + resolves a function symbol within the library
    - returns error when the symbol is not exported
    # loader
  dynlib.call
    @ (sym: symbol_handle, signature: call_signature, args: list[native_value]) -> result[native_value, string]
    + invokes the native function with marshalled arguments per the signature
    - returns error when the number of args does not match the signature
    - returns error when an arg's type does not match the signature
    ? uses platform calling conventions to lay out registers and the stack
    # invocation
  dynlib.signature
    @ (return_type: native_type, param_types: list[native_type]) -> call_signature
    + describes a native function's ABI signature
    # signatures
  dynlib.box_int
    @ (value: i64) -> native_value
    + wraps an integer as a native value
    # marshalling
  dynlib.box_string
    @ (value: string) -> native_value
    + wraps a string as a native-compatible pointer value
    # marshalling
  dynlib.unbox_int
    @ (value: native_value) -> result[i64, string]
    + extracts an integer from a native value
    - returns error when value is not an integer
    # marshalling
