# Requirement: "render test results as readable sentences"

Converts camelCase test identifiers into human-readable sentences and groups them by subject.

std: (all units exist)

test_prose
  test_prose.humanize_name
    fn (test_name: string) -> string
    + splits a camelCase or snake_case identifier into a space-separated sentence
    + strips any leading "Test" or "test_" prefix
    - returns "" for an empty input
    # text_transform
  test_prose.format_result
    fn (subject: string, sentence: string, passed: bool) -> string
    + returns a line like " x subject should do the thing" with a pass or fail marker
    # formatting
  test_prose.group_by_subject
    fn (entries: list[test_entry]) -> map[string, list[test_entry]]
    + groups test entries by their subject field preserving order within each group
    # grouping
