# Requirement: "a graphics pack that returns illustrated characters by name and variant"

A keyed asset lookup where each character comes in multiple visual variants.

std: (all units exist)

graphics_pack
  graphics_pack.get
    fn (character: string, variant: string) -> optional[bytes]
    + returns the artwork bytes for a (character, variant) pair
    - returns none when either the character or the variant is unknown
    # catalog
