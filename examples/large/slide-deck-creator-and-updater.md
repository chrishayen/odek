# Requirement: "a library for creating and updating presentation slide-deck files"

Build a presentation in memory, add slides and shapes, save. Zip packaging and XML generation belong in std since they are generic primitives.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, overwriting
      - returns error on permission failure
      # filesystem
  std.zip
    std.zip.create
      @ () -> zip_state
      + creates an empty zip archive in memory
      # archive
    std.zip.add_entry
      @ (state: zip_state, name: string, data: bytes) -> zip_state
      + adds a named entry with uncompressed data
      # archive
    std.zip.read_entry
      @ (state: zip_state, name: string) -> result[bytes, string]
      + returns the bytes of a named entry
      - returns error when entry is missing
      # archive
    std.zip.finalize
      @ (state: zip_state) -> bytes
      + serializes the archive to a byte buffer
      # archive
    std.zip.open
      @ (raw: bytes) -> result[zip_state, string]
      + parses an existing zip archive
      - returns error on corrupt input
      # archive
  std.xml
    std.xml.build_element
      @ (tag: string, attrs: map[string, string], children: list[string]) -> string
      + returns a serialized xml element
      # serialization
    std.xml.parse
      @ (raw: string) -> result[xml_node, string]
      + parses xml into a node tree
      - returns error on malformed input
      # parsing

slides
  slides.new
    @ (title: string) -> deck_state
    + creates an empty deck with the given title
    # construction
  slides.add_slide
    @ (state: deck_state, layout: string) -> tuple[deck_state, i32]
    + appends a new slide with the given layout and returns (state, slide_index)
    ? layout names are "title", "title_content", "blank"
    # slides
  slides.set_title
    @ (state: deck_state, slide_index: i32, text: string) -> result[deck_state, string]
    + sets the title placeholder on a slide
    - returns error when slide_index is out of range
    # content
  slides.add_text_box
    @ (state: deck_state, slide_index: i32, text: string, x: i32, y: i32, width: i32, height: i32) -> result[deck_state, string]
    + adds a text box with absolute positioning in EMU units
    - returns error when slide_index is out of range
    # shapes
  slides.add_image
    @ (state: deck_state, slide_index: i32, image_path: string, x: i32, y: i32, width: i32, height: i32) -> result[deck_state, string]
    + embeds an image from disk on a slide
    - returns error when image_path cannot be read
    # shapes
    -> std.fs.read_all
  slides.add_bullet_list
    @ (state: deck_state, slide_index: i32, bullets: list[string]) -> result[deck_state, string]
    + adds a bulleted list to the content placeholder
    - returns error when slide_index is out of range
    # content
  slides.slide_count
    @ (state: deck_state) -> i32
    + returns the number of slides in the deck
    # introspection
  slides.save
    @ (state: deck_state, path: string) -> result[void, string]
    + serializes the deck to a package file at path
    - returns error on write failure
    # io
    -> std.xml.build_element
    -> std.zip.create
    -> std.zip.add_entry
    -> std.zip.finalize
    -> std.fs.write_all
  slides.load
    @ (path: string) -> result[deck_state, string]
    + loads an existing deck from disk
    - returns error when the file is missing or malformed
    # io
    -> std.fs.read_all
    -> std.zip.open
    -> std.zip.read_entry
    -> std.xml.parse
  slides.remove_slide
    @ (state: deck_state, slide_index: i32) -> result[deck_state, string]
    + removes the slide at the given index
    - returns error when slide_index is out of range
    # slides
