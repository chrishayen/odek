# Requirement: "an internationalization library with plural form selection and interpolation"

Messages are looked up by key and locale, with plural selection using CLDR-style categories and ${var} interpolation.

std: (all units exist)

i18n
  i18n.new_catalog
    @ () -> catalog
    + creates an empty message catalog
    # construction
  i18n.add_message
    @ (cat: catalog, locale: string, key: string, forms: map[string, string]) -> catalog
    + registers a message with plural-form variants like "one" and "other"
    ? a non-plural message uses only the "other" form
    # registration
  i18n.plural_category
    @ (locale: string, n: i64) -> string
    + returns the CLDR plural category for n under the given locale
    ? categories include "zero" "one" "two" "few" "many" "other"
    # plurals
  i18n.translate
    @ (cat: catalog, locale: string, key: string, n: i64, vars: map[string, string]) -> result[string, string]
    + selects the plural form for n, interpolates ${var} placeholders, and returns the string
    - returns error when the key is missing in the given locale
    - returns error when a required form is missing
    # lookup
