# Requirement: "a procedural XML generation API"

A streaming writer that emits well-formed XML through explicit start/end/text calls.

std
  std.io
    std.io.buffer_new
      fn () -> byte_buffer
      + creates an empty growable byte buffer
      # io
    std.io.buffer_write_string
      fn (buf: byte_buffer, s: string) -> byte_buffer
      + appends UTF-8 bytes to the buffer
      # io
    std.io.buffer_contents
      fn (buf: byte_buffer) -> bytes
      + returns the accumulated bytes
      # io

xml_writer
  xml_writer.new
    fn (indent: string) -> writer_state
    + creates a writer; empty indent disables pretty-printing
    # construction
    -> std.io.buffer_new
  xml_writer.start_document
    fn (state: writer_state, version: string, encoding: string) -> writer_state
    + writes the XML declaration
    # document
  xml_writer.start_element
    fn (state: writer_state, name: string) -> result[writer_state, string]
    + opens a new element and pushes it on the element stack
    - returns error when name is empty or not a valid XML name
    # element
  xml_writer.write_attribute
    fn (state: writer_state, name: string, value: string) -> result[writer_state, string]
    + adds an attribute to the currently open start tag
    - returns error when called outside an open start tag
    # attribute
  xml_writer.write_text
    fn (state: writer_state, text: string) -> writer_state
    + escapes special characters and writes character data
    # text
  xml_writer.write_cdata
    fn (state: writer_state, text: string) -> result[writer_state, string]
    + writes a CDATA section unchanged
    - returns error when text contains the CDATA terminator
    # cdata
  xml_writer.end_element
    fn (state: writer_state) -> result[writer_state, string]
    + closes the innermost open element
    - returns error when the stack is empty
    # element
  xml_writer.finish
    fn (state: writer_state) -> result[bytes, string]
    + returns the full document bytes
    - returns error when there are still open elements
    # document
    -> std.io.buffer_contents
