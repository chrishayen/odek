# Requirement: "a datastore-agnostic tool that simplifies interaction with one or more backing stores"

Models are defined once; a pluggable adapter translates operations to a specific store.

std: (all units exist)

orm
  orm.define_model
    @ (name: string, attributes: list[attribute]) -> model
    + returns a model with the given name and attribute schema
    - returns a model whose create fails when no primary key is declared
    # schema
  orm.register_adapter
    @ (name: string, adapter: adapter) -> void
    + associates an adapter with a name
    # adapters
  orm.bind
    @ (model: model, adapter_name: string) -> result[bound_model, string]
    + returns a model bound to a registered adapter
    - returns error when the adapter name is unknown
    # binding
  orm.create
    @ (bound: bound_model, record: map[string, value]) -> result[map[string, value], string]
    + inserts a record and returns it with generated fields populated
    - returns error when a required attribute is missing
    # crud
  orm.find
    @ (bound: bound_model, criteria: map[string, value]) -> result[list[map[string, value]], string]
    + returns records matching all key/value pairs in criteria
    # crud
  orm.update
    @ (bound: bound_model, criteria: map[string, value], changes: map[string, value]) -> result[i32, string]
    + returns the number of records updated
    - returns error when a change targets an unknown attribute
    # crud
  orm.destroy
    @ (bound: bound_model, criteria: map[string, value]) -> result[i32, string]
    + returns the number of records removed
    # crud
