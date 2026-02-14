# Prompts Domain BDD

## Scope
Skills/rules, profile assignment, system prompt versioning, and prompt composition.

## Scenarios
- Given user-global skills/rules in an organization, when referenced from multiple projects in that org, then they are reusable.
- Given another organization, when requesting foreign org skills/rules, then access is denied.
- Given skill and rule profiles, when assigning them to a project, then each profile resolves only its own resource type.
- Given system prompt versions, when activating a version, then prompt resolution uses active version content.
- Given prompt preview requests, when composing prompt, then layer order is fixed:
  1) system prompt
  2) user global skills
  3) project skill profile skills
  4) user global rules
  5) project rule profile rules
- Given repeated instruction text across layers, when composing preview, then occurrences are preserved without deduplication.
- Given unchanged inputs, when previewing twice, then output is byte-identical.

## Invariants
- Skill and rule profiles are distinct resource types.
- Prompt keys include `chat_planning` and `worker`.
