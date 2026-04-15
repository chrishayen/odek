# Decomposition examples

A corpus of worked decompositions for the odek decomposer. Each file shows one requirement and its expected decomposition in the same text-DSL format used by `cmd/decompose.md`'s `token_validator` example.

## Directory layout

Files are sorted by complexity tier:

- **`trivial/`** — 1–2 runes, `std` is usually empty. Single operations or pairs of pure functions.
- **`small/`** — 2–5 runes. `std` may contain thin general-purpose wrappers (e.g. `std.io.print_string`).
- **`medium/`** — 5–10 runes. `std` contains real substantive utilities — encoding, hashing, parsing, timing primitives a few different projects would reuse.
- **`large/`** — 10+ runes. `std` is a rich set of subsystems (e.g. `std.http`, `std.jwt`, `std.bcrypt`, `std.websocket`, `std.sql`).

## File naming

Kebab-case after the requirement. Filenames are grepable — an agent scanning the directory should recognize `csv-reader.md`, `jwt-hs256.md`, `http-server-with-routing.md` without opening them.

## File format

```markdown
# Requirement: "<the user's requirement in quotes>"

<optional prose note about the decomposition's intent>

std
  std.package
    std.package.unit
      fn (args) -> return_type
      + positive test
      - negative test
      ? assumption
      # responsibility_tag

project_name
  project_name.unit
    fn (args) -> return_type
    + test
    - test
    # tag
    -> std.package.unit        (reference — no redefinition)
```

When `std` is empty, use `std: (all units exist)` on one line and skip the std tree.

## Notes for the agent reading these

- Every rune name appears exactly once. `-> std.path.unit` is a reference, not a redefinition.
- `std` runes never mention a specific feature by name — they describe generic capabilities.
- Thin general-purpose stdlib wrappers (`std.io.print_string`) are fine. Feature-specific filler in std (`std.io.print_hello_world_greeting`) is not.
- If the user asks for "hello world", produce ONE greeter — never `say_hello` + `say_hello_with_name` + `say_hello_multiple_times`. Restraint is a feature.
