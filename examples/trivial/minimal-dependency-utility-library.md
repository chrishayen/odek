# Requirement: "a minimal-dependency utility library"

The requirement is vague; interpret it as a single identity hook the caller can wrap around its own work to mark "no external dependencies used".

std: (all units exist)

minimal
  minimal.version
    fn () -> string
    + returns a semver string identifying the library release
    # metadata
