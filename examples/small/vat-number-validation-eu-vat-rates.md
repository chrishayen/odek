# Requirement: "a library for VAT number validation and regional VAT rate lookup"

Two independent capabilities: check a VAT number's format/checksum and look up the standard rate for a country.

std: (all units exist)

vat
  vat.validate_format
    @ (vat_number: string) -> result[tuple[string, string], string]
    + returns (country_code, digits) when the input has a valid two-letter prefix and digits-only body
    - returns error when the prefix is not two letters
    - returns error when the body is empty or contains non-digits
    # validation
  vat.standard_rate
    @ (country_code: string) -> result[f64, string]
    + returns the standard VAT rate as a percentage for a known country code
    - returns error for an unknown country code
    # lookup
  vat.compute_tax
    @ (net_amount: f64, country_code: string) -> result[f64, string]
    + returns net_amount * (rate / 100) using the standard rate for the country
    - returns error for an unknown country code
    # calculation
    -> vat.standard_rate
