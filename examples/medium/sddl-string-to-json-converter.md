# Requirement: "convert a security descriptor definition string into a structured JSON document"

Parses the four sections (owner, primary group, DACL, SACL), expands ACE flags and rights, and emits a readable JSON object.

std
  std.json
    std.json.encode_value
      fn (value: json_value) -> string
      + encodes a generic value tree as JSON text
      # serialization

sddl
  sddl.split_sections
    fn (input: string) -> result[map[string, string], string]
    + extracts the O:, G:, D:, S: sections into a map keyed by letter
    - returns error on malformed section markers
    # parsing
  sddl.parse_sid
    fn (token: string) -> result[string, string]
    + expands well-known SID aliases like "BA" to their canonical names
    + passes through literal "S-1-..." SIDs unchanged
    - returns error on empty input
    # parsing
  sddl.parse_ace
    fn (token: string) -> result[ace_record, string]
    + parses the semicolon-delimited fields of an ACE
    + expands the access rights mask into a list of named rights
    + expands the flags field into a list of named flags
    - returns error when the field count is wrong
    # parsing
    -> sddl.parse_sid
  sddl.parse_acl
    fn (body: string) -> result[list[ace_record], string]
    + splits the ACL body into parenthesized ACE entries
    + parses each entry with parse_ace
    - returns error when parentheses are unbalanced
    # parsing
    -> sddl.parse_ace
  sddl.to_json
    fn (input: string) -> result[string, string]
    + returns a JSON object with "owner", "group", "dacl", "sacl" keys
    + omits sections that are absent from the input
    - returns error when any section fails to parse
    # conversion
    -> sddl.split_sections
    -> sddl.parse_sid
    -> sddl.parse_acl
    -> std.json.encode_value
