# Requirement: "a library for country, currency, phone, and top-level-domain code lookups"

Pure data-lookup surface: given a code in one scheme, return its metadata or convert between schemes.

std: (all units exist)

codes
  codes.country_by_alpha2
    fn (code: string) -> result[country_record, string]
    + returns the country record for a valid two-letter code
    - returns error when the code is unknown
    ? input is case-insensitive
    # country_lookup
  codes.country_by_alpha3
    fn (code: string) -> result[country_record, string]
    + returns the country record for a valid three-letter code
    - returns error when the code is unknown
    # country_lookup
  codes.country_by_numeric
    fn (code: i32) -> result[country_record, string]
    + returns the country record for a valid numeric code
    - returns error when the code is unknown
    # country_lookup
  codes.currency_by_code
    fn (code: string) -> result[currency_record, string]
    + returns the currency record for a valid three-letter currency code
    - returns error when the code is unknown
    # currency_lookup
  codes.calling_code_for_country
    fn (alpha2: string) -> result[string, string]
    + returns the international dialing prefix for a country
    - returns error when the country has no assigned calling code
    # phone_lookup
  codes.cctld_for_country
    fn (alpha2: string) -> result[string, string]
    + returns the country-code top-level domain for a country
    - returns error when the country has no ccTLD
    # tld_lookup
  codes.list_countries
    fn () -> list[country_record]
    + returns every known country record in alpha2 order
    # enumeration
