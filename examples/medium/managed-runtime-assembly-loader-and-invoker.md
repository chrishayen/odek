# Requirement: "a host-language bridge for loading and calling into managed runtime assemblies"

Load an assembly, resolve types and methods by name, invoke with boxed arguments. Runtime loading is a std seam.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the path does not exist
      # filesystem
  std.dynlib
    std.dynlib.load
      @ (path: string) -> result[dynlib_handle, string]
      + loads a native runtime host library
      - returns error when the file is missing or not a valid library
      # ffi
    std.dynlib.call
      @ (handle: dynlib_handle, symbol: string, args: list[bytes]) -> result[bytes, string]
      + invokes a symbol with boxed arguments and returns boxed result
      - returns error when the symbol is missing
      # ffi

managed_bridge
  managed_bridge.init_runtime
    @ () -> result[runtime_state, string]
    + initializes the managed runtime host
    - returns error when the runtime host library cannot be loaded
    # runtime
    -> std.dynlib.load
  managed_bridge.load_assembly
    @ (state: runtime_state, path: string) -> result[assembly_handle, string]
    + loads an assembly file into the runtime
    - returns error when the file is missing or not a valid assembly
    # loading
    -> std.fs.read_all
    -> std.dynlib.call
  managed_bridge.resolve_type
    @ (state: runtime_state, assembly: assembly_handle, full_name: string) -> result[type_handle, string]
    + looks up a type by its fully qualified name
    - returns error when the type does not exist in the assembly
    # reflection
    -> std.dynlib.call
  managed_bridge.resolve_method
    @ (state: runtime_state, type: type_handle, method_name: string, arg_type_names: list[string]) -> result[method_handle, string]
    + looks up a method by name and argument type signature
    - returns error when no matching overload exists
    # reflection
    -> std.dynlib.call
  managed_bridge.invoke
    @ (state: runtime_state, method: method_handle, target: optional[object_handle], args: list[bytes]) -> result[bytes, string]
    + invokes a method with boxed arguments and returns the boxed result
    - returns error when argument arity or types are incompatible
    # invocation
    -> std.dynlib.call
  managed_bridge.new_instance
    @ (state: runtime_state, type: type_handle, ctor_args: list[bytes]) -> result[object_handle, string]
    + constructs a new instance of a type
    - returns error when no matching constructor exists
    # invocation
    -> std.dynlib.call
  managed_bridge.shutdown
    @ (state: runtime_state) -> void
    + tears down the managed runtime host
    # runtime
