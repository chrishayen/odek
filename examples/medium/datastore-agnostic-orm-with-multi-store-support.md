# Requirement: "a datastore-agnostic tool that simplifies interaction with one or more backing stores"

Models are defined once; a pluggable adapter translates operations to a specific store.

std: (all units exist)

orm
  orm.define_model
    fn (name: string, attributes: list[attribute]) -> model
    + returns a model with the given name and attribute schema
    - returns a model whose create fails when no primary key is declared
    # schema
  orm.register_adapter
    fn (name: string, adapter: adapter) -> void
    + associates an adapter with a name
    # adapters
  orm.bind
    fn (model: model, adapter_name: string) -> result[bound_model, string]
    + returns a model bound to a registered adapter
    - returns error when the adapter name is unknown
    # binding
  orm.create
    fn (bound: bound_model, record: map[string, value]) -> result[map[string, value], string]
    + inserts a record and returns it with generated fields populated
    - returns error when a required attribute is missing
    # crud
  orm.find
    fn (bound: bound_model, criteria: map[string, value]) -> result[list[map[string, value]], string]
    + returns records matching all key/value pairs in criteria
    # crud
  orm.update
    fn (bound: bound_model, criteria: map[string, value], changes: map[string, value]) -> result[i32, string]
    + returns the number of records updated
    - returns error when a change targets an unknown attribute
    # crud
  orm.destroy
    fn (bound: bound_model, criteria: map[string, value]) -> result[i32, string]
    + returns the number of records removed
    # crud
