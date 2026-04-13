# Requirement: "an asset catalog that returns artwork files by name"

A tiny catalog mapping asset names to their raw bytes.

std: (all units exist)

assets
  assets.lookup
    @ (name: string) -> optional[bytes]
    + returns the bytes of the asset with the given name
    - returns none when no asset by that name is registered
    # catalog
