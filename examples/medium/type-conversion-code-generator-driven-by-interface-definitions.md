# Requirement: "a type conversion code generator driven by interface definitions"

Given source and target type descriptions, matches fields by name and emits a function that copies the source into the target.

std: (all units exist)

converter
  converter.describe_type
    @ (name: string, fields: list[field_spec]) -> type_spec
    + builds a type spec with a name and a list of (field_name, field_type) pairs
    # specification
  converter.plan
    @ (src: type_spec, dst: type_spec) -> conversion_plan
    + returns a plan with one assignment per matched field name
    ? unmatched destination fields are left out of the plan; the caller decides what to do
    # planning
  converter.unmatched_fields
    @ (plan: conversion_plan) -> list[string]
    + returns destination field names that have no source counterpart
    # diagnostics
  converter.needs_cast
    @ (plan: conversion_plan) -> list[string]
    + returns field names where the source and destination types differ
    # diagnostics
  converter.render
    @ (plan: conversion_plan, function_name: string) -> string
    + returns the source text of a conversion function that applies the plan
    + inserts explicit casts for fields flagged by needs_cast
    # rendering
