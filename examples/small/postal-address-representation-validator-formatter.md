# Requirement: "a library for postal address representation, validation, and formatting"

std: (all units exist)

postal_address
  postal_address.new
    fn (lines: list[string], city: string, region: string, postal_code: string, country_code: string) -> address_state
    + builds an address value with the given components
    ? country_code is expected to be an iso 3166-1 alpha-2 code
    # construction
  postal_address.validate
    fn (address: address_state) -> result[void, string]
    + confirms that required components for the country are present and well-formed
    - returns error when the country code is not recognized
    - returns error when the postal code does not match the country's pattern
    # validation
  postal_address.format_single_line
    fn (address: address_state) -> string
    + returns the address as a single comma-separated line
    # formatting
  postal_address.format_multi_line
    fn (address: address_state) -> list[string]
    + returns the address as a list of display lines in the country's conventional order
    # formatting
