# Requirement: "a SQL query builder that generates queries from record descriptions"

The library reflects record descriptions (field name + value) into INSERT, UPDATE, SELECT, and DELETE statements using positional parameters.

std: (all units exist)

qbuild
  qbuild.insert
    @ (table: string, record: map[string, string]) -> tuple[string, list[string]]
    + returns an INSERT statement with positional placeholders and the ordered values
    ? column order is stable so generated SQL is deterministic
    # insert
  qbuild.update
    @ (table: string, record: map[string, string], where: map[string, string]) -> tuple[string, list[string]]
    + returns an UPDATE statement with SET columns from record and WHERE from where
    - returns an empty SQL string when record is empty
    # update
  qbuild.select_by
    @ (table: string, columns: list[string], where: map[string, string]) -> tuple[string, list[string]]
    + returns a SELECT statement for the given columns filtered by equality predicates
    # select
  qbuild.delete_by
    @ (table: string, where: map[string, string]) -> tuple[string, list[string]]
    + returns a DELETE statement filtered by equality predicates
    - returns an empty SQL string when where is empty to avoid mass deletes
    # delete
