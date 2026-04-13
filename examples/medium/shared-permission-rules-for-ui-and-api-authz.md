# Requirement: "an authorization library that lets ui and api layers share the same permission rules"

A tiny permission engine: define abilities as (action, subject_type, condition) triples bound to a user; ask whether the user can perform an action on a subject. Works the same whether called from client or server.

std: (all units exist)

authz
  authz.new_ability_set
    @ () -> ability_set
    + creates an empty ability set
    # construction
  authz.allow
    @ (abilities: ability_set, action: string, subject_type: string) -> ability_set
    + adds an unconditional allow rule
    # rule
  authz.allow_when
    @ (abilities: ability_set, action: string, subject_type: string, condition_id: string) -> ability_set
    + adds a conditional allow rule whose predicate is resolved by condition_id
    ? conditions are registered separately; this stores the binding
    # rule
  authz.forbid
    @ (abilities: ability_set, action: string, subject_type: string) -> ability_set
    + adds an explicit deny that beats any allow
    # rule
  authz.register_condition
    @ (abilities: ability_set, condition_id: string, predicate_name: string) -> ability_set
    + binds a named predicate the host can evaluate against a subject
    # extension
  authz.can
    @ (abilities: ability_set, action: string, subject_type: string, subject_fields: map[string,string]) -> bool
    + returns true when some allow rule matches and no forbid rule matches
    - returns false when no rule matches
    # decision
  authz.check
    @ (abilities: ability_set, action: string, subject_type: string, subject_fields: map[string,string]) -> result[void, string]
    + returns ok when allowed, a descriptive error when not
    - returns error when an explicit forbid matched
    # decision
  authz.permitted_actions
    @ (abilities: ability_set, subject_type: string) -> list[string]
    + returns every action the current ruleset grants on subject_type with no condition
    # query
  authz.serialize
    @ (abilities: ability_set) -> string
    + encodes the ruleset as a portable string so it can cross the ui/api boundary
    # serialization
  authz.deserialize
    @ (encoded: string) -> result[ability_set, string]
    + restores an ability set produced by authz.serialize
    - returns error on malformed input
    # serialization
