# Requirement: "an HTML and XML parsing library"

A tolerant markup parser that produces a navigable tree for both HTML and XML, with basic query helpers.

std
  std.strings
    std.strings.to_lower
      @ (s: string) -> string
      + returns the lowercased form of s
      # text
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits s on sep returning all segments
      + returns [s] when sep is not found
      # text
    std.strings.trim
      @ (s: string) -> string
      + strips leading and trailing whitespace
      # text
  std.collections
    std.collections.hashmap_new
      @ () -> map[string, string]
      + returns an empty attribute map
      # collections

markup
  markup.tokenize
    @ (source: string) -> result[list[markup_token], string]
    + emits tokens for start tags, end tags, text, comments, and cdata
    + preserves entity references as-is
    - returns error on unterminated comment
    # lexing
  markup.parse_html
    @ (source: string) -> result[element, string]
    + builds a tree tolerating unclosed void elements (br, img, meta, etc.)
    + lowercases tag names for case-insensitive matching
    - returns error on malformed attribute list
    # html_parsing
    -> std.strings.to_lower
  markup.parse_xml
    @ (source: string) -> result[element, string]
    + builds a strict tree requiring balanced tags
    + preserves case of tag names and attributes
    - returns error when a tag is not closed
    - returns error when attribute values are unquoted
    # xml_parsing
  markup.get_attribute
    @ (el: element, name: string) -> optional[string]
    + returns the value of an attribute if present
    - returns none when attribute is absent
    # querying
  markup.find_by_tag
    @ (root: element, tag: string) -> list[element]
    + returns all descendants whose tag matches
    + returns [] when none found
    # querying
  markup.find_by_id
    @ (root: element, id: string) -> optional[element]
    + returns the first descendant whose id attribute equals id
    - returns none when no element has that id
    # querying
  markup.text_content
    @ (el: element) -> string
    + concatenates all descendant text nodes in document order
    # querying
    -> std.strings.trim
  markup.serialize
    @ (el: element) -> string
    + serializes an element tree back to markup form
    # serialization
