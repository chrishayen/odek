package tests

import "../src/widgets"
import "../src/core"
import "core:testing"

// ============================================================================
// TextInput creation and destruction
// ============================================================================

@(test)
test_text_input_create :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()
    testing.expect(t, ti != nil, "text_input_create should return non-nil")
    testing.expect(t, ti.visible == true, "TextInput should be visible by default")
    testing.expect(t, ti.enabled == true, "TextInput should be enabled by default")
    testing.expect(t, ti.focusable == true, "TextInput should be focusable")
    testing.expect(t, ti.cursor == .Text, "TextInput should have Text cursor")

    widgets.widget_destroy(ti)
}

@(test)
test_text_input_destroy_nil :: proc(t: ^testing.T) {
    // Should not crash
    widgets.widget_destroy(nil)
}

// ============================================================================
// Text operations
// ============================================================================

@(test)
test_text_input_set_text :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hello")
    text := widgets.text_input_get_text(ti)
    testing.expect(t, text == "Hello", "Text should be set correctly")
    testing.expect(t, ti.cursor_pos == 5, "Cursor should be at end")

    widgets.widget_destroy(ti)
}

@(test)
test_text_input_clear :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hello")
    widgets.text_input_clear(ti)

    text := widgets.text_input_get_text(ti)
    testing.expect(t, text == "", "Text should be empty after clear")
    testing.expect(t, ti.cursor_pos == 0, "Cursor should be at start")

    widgets.widget_destroy(ti)
}

@(test)
test_text_input_insert :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_insert(ti, "Hello")
    widgets.text_input_insert(ti, " World")

    text := widgets.text_input_get_text(ti)
    testing.expect(t, text == "Hello World", "Text should be concatenated")

    widgets.widget_destroy(ti)
}

// ============================================================================
// Cursor movement
// ============================================================================

@(test)
test_text_input_cursor_left :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hello")
    widgets.text_input_cursor_left(ti)

    testing.expect(t, ti.cursor_pos == 4, "Cursor should move left")

    widgets.widget_destroy(ti)
}

@(test)
test_text_input_cursor_right :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hello")
    ti.cursor_pos = 0
    widgets.text_input_cursor_right(ti)

    testing.expect(t, ti.cursor_pos == 1, "Cursor should move right")

    widgets.widget_destroy(ti)
}

@(test)
test_text_input_cursor_home :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hello")
    widgets.text_input_cursor_home(ti)

    testing.expect(t, ti.cursor_pos == 0, "Cursor should be at start")

    widgets.widget_destroy(ti)
}

@(test)
test_text_input_cursor_end :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hello")
    ti.cursor_pos = 0
    widgets.text_input_cursor_end(ti)

    testing.expect(t, ti.cursor_pos == 5, "Cursor should be at end")

    widgets.widget_destroy(ti)
}

// ============================================================================
// Deletion
// ============================================================================

@(test)
test_text_input_backspace :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hello")
    widgets.text_input_backspace(ti)

    text := widgets.text_input_get_text(ti)
    testing.expect(t, text == "Hell", "Backspace should delete last character")

    widgets.widget_destroy(ti)
}

@(test)
test_text_input_delete :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hello")
    ti.cursor_pos = 0
    widgets.text_input_delete(ti)

    text := widgets.text_input_get_text(ti)
    testing.expect(t, text == "ello", "Delete should remove character at cursor")

    widgets.widget_destroy(ti)
}

@(test)
test_text_input_backspace_at_start :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hello")
    ti.cursor_pos = 0
    widgets.text_input_backspace(ti)

    text := widgets.text_input_get_text(ti)
    testing.expect(t, text == "Hello", "Backspace at start should do nothing")

    widgets.widget_destroy(ti)
}

@(test)
test_text_input_delete_at_end :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hello")
    widgets.text_input_delete(ti)

    text := widgets.text_input_get_text(ti)
    testing.expect(t, text == "Hello", "Delete at end should do nothing")

    widgets.widget_destroy(ti)
}

// ============================================================================
// Insert in middle
// ============================================================================

@(test)
test_text_input_insert_middle :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_text(ti, "Hllo")
    ti.cursor_pos = 1
    widgets.text_input_insert(ti, "e")

    text := widgets.text_input_get_text(ti)
    testing.expect(t, text == "Hello", "Insert should work in middle")
    testing.expect(t, ti.cursor_pos == 2, "Cursor should advance after insert")

    widgets.widget_destroy(ti)
}

// ============================================================================
// Placeholder
// ============================================================================

@(test)
test_text_input_placeholder :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    widgets.text_input_set_placeholder(ti, "Enter text...")
    testing.expect(t, ti.placeholder == "Enter text...", "Placeholder should be set")

    widgets.widget_destroy(ti)
}

// ============================================================================
// Measure
// ============================================================================

@(test)
test_text_input_measure :: proc(t: ^testing.T) {
    ti := widgets.text_input_create()

    size := widgets.widget_measure(ti)
    testing.expect(t, size.width > 0, "Width should be > 0")
    testing.expect(t, size.height > 0, "Height should be > 0")

    widgets.widget_destroy(ti)
}
