# Requirement: "a handwritten-notes and sketch library"

A vector-stroke model for freeform sketches. The caller feeds pointer samples; the library assembles strokes and exposes them as geometry.

std: (all units exist)

sketchpad
  sketchpad.new
    @ () -> sketch_state
    + creates an empty sketch with no strokes
    # construction
  sketchpad.begin_stroke
    @ (state: sketch_state, color: u32, width: f32) -> sketch_state
    + starts a new stroke with the given color and pen width
    # input
  sketchpad.add_point
    @ (state: sketch_state, x: f32, y: f32, pressure: f32) -> result[sketch_state, string]
    + appends a sample to the currently open stroke
    - returns error when there is no open stroke
    # input
  sketchpad.end_stroke
    @ (state: sketch_state) -> sketch_state
    + closes the current stroke and commits it to the sketch
    + no-op when no stroke is open
    # input
  sketchpad.undo_last_stroke
    @ (state: sketch_state) -> sketch_state
    + removes the most recently committed stroke
    + returns unchanged state when there are no strokes
    # history
  sketchpad.clear
    @ (state: sketch_state) -> sketch_state
    + returns a fresh empty sketch
    # history
  sketchpad.bounding_box
    @ (state: sketch_state) -> optional[bbox]
    + returns the smallest axis-aligned box containing all stroke samples
    - returns none when the sketch is empty
    # geometry
  sketchpad.stroke_count
    @ (state: sketch_state) -> i32
    + returns the number of committed strokes
    # inspection
