# Requirement: "a pull-style XML reader and writer"

The reader yields events one at a time; the writer builds a document incrementally.

std: (all units exist)

xml
  xml.new_reader
    @ (source: string) -> reader_state
    + creates a reader positioned at the start of source
    # construction
  xml.next_event
    @ (r: reader_state) -> result[optional[xml_event], string]
    + returns the next StartElement, EndElement, Text, Comment, or CData event
    + returns none when the end of input is reached
    - returns error on malformed markup or mismatched tags
    # parsing
  xml.skip_subtree
    @ (r: reader_state) -> result[void, string]
    + advances past the current element's children until its matching EndElement
    # navigation
    -> xml.next_event
  xml.unescape_entities
    @ (text: string) -> result[string, string]
    + resolves &amp;, &lt;, &gt;, &quot;, &apos;, and numeric character references
    - returns error on unknown named entities
    # escaping
  xml.new_writer
    @ () -> writer_state
    + creates a writer whose buffer begins with an XML declaration
    # construction
  xml.write_start
    @ (w: writer_state, name: string, attrs: map[string, string]) -> writer_state
    + appends "<name key=\"value\">" with attribute values escaped
    # writing
  xml.write_text
    @ (w: writer_state, text: string) -> writer_state
    + appends text with XML special characters escaped
    # writing
  xml.write_end
    @ (w: writer_state, name: string) -> writer_state
    + appends "</name>"
    # writing
  xml.finish
    @ (w: writer_state) -> string
    + returns the accumulated document
    # writing
