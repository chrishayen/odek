package tests

import "../src/core"
import "../src/render"
import "../src/widgets"
import "core:testing"

// ============================================================================
// Label creation and destruction
// ============================================================================

@(test)
test_label_create :: proc(t: ^testing.T) {
    l := widgets.label_create()
    testing.expect(t, l != nil, "label_create should return non-nil")
    testing.expect(t, l.visible == true, "Label should be visible by default")
    testing.expect(t, l.enabled == true, "Label should be enabled by default")
    testing.expect(t, l.dirty == true, "Label should be dirty by default")
    testing.expect(t, l.text == "", "Label text should be empty by default")
    testing.expect(t, l.font == nil, "Label font should be nil by default")
    testing.expect(t, l.wrap == true, "Label wrap should be true by default")
    testing.expect(t, l.h_align == .Start, "Label alignment should be Start by default")
    testing.expect(t, l.vtable == &widgets.label_vtable, "Label should have label vtable")

    widgets.widget_destroy(l)
}

@(test)
test_label_create_with_text :: proc(t: ^testing.T) {
    l := widgets.label_create("Hello World")
    testing.expect(t, l.text == "Hello World", "Label should have provided text")

    widgets.widget_destroy(l)
}

@(test)
test_label_destroy :: proc(t: ^testing.T) {
    l := widgets.label_create("Test")
    // Should not crash
    widgets.widget_destroy(l)
}

// ============================================================================
// Label property setters
// ============================================================================

@(test)
test_label_set_text :: proc(t: ^testing.T) {
    l := widgets.label_create()
    l.dirty = false

    widgets.label_set_text(l, "New text")
    testing.expect(t, l.text == "New text", "Text should be updated")
    testing.expect(t, l.dirty == true, "Label should be marked dirty")

    widgets.widget_destroy(l)
}

@(test)
test_label_set_text_same :: proc(t: ^testing.T) {
    l := widgets.label_create("Same")
    l.dirty = false

    widgets.label_set_text(l, "Same")
    testing.expect(t, l.dirty == false, "Label should not be marked dirty if text unchanged")

    widgets.widget_destroy(l)
}

@(test)
test_label_set_color :: proc(t: ^testing.T) {
    l := widgets.label_create()
    l.dirty = false

    widgets.label_set_color(l, core.COLOR_RED)
    testing.expect(t, l.color.r == core.COLOR_RED.r, "Color should be updated")
    testing.expect(t, l.dirty == true, "Label should be marked dirty")

    widgets.widget_destroy(l)
}

@(test)
test_label_set_align :: proc(t: ^testing.T) {
    l := widgets.label_create()
    l.dirty = false

    widgets.label_set_align(l, .Center)
    testing.expect(t, l.h_align == .Center, "Alignment should be updated")
    testing.expect(t, l.dirty == true, "Label should be marked dirty")

    widgets.widget_destroy(l)
}

@(test)
test_label_set_wrap :: proc(t: ^testing.T) {
    l := widgets.label_create()
    l.dirty = false

    widgets.label_set_wrap(l, false)
    testing.expect(t, l.wrap == false, "Wrap should be updated")
    testing.expect(t, l.dirty == true, "Label should be marked dirty")

    widgets.widget_destroy(l)
}

@(test)
test_label_set_wrap_same :: proc(t: ^testing.T) {
    l := widgets.label_create()
    l.dirty = false

    widgets.label_set_wrap(l, true)  // Same as default
    testing.expect(t, l.dirty == false, "Label should not be marked dirty if wrap unchanged")

    widgets.widget_destroy(l)
}

// ============================================================================
// Label in widget tree
// ============================================================================

@(test)
test_label_as_child :: proc(t: ^testing.T) {
    container := widgets.container_create()
    label := widgets.label_create("Test")

    widgets.widget_add_child(container, label)

    testing.expect(t, len(container.children) == 1, "Container should have 1 child")
    testing.expect(t, label.parent == container, "Label's parent should be container")

    widgets.widget_destroy(container)
}

// ============================================================================
// Word wrapping tests (with font)
// ============================================================================

@(test)
test_label_wrap_single_line :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    l := widgets.label_create("Short text", &font)
    l.rect = core.Rect{0, 0, 500, 100}  // Wide enough for text

    // Trigger line calculation
    widgets.label_layout(l)

    testing.expect(t, len(l.lines) == 1, "Short text should be one line")
    testing.expect(t, l.lines[0] == "Short text", "Line should contain full text")

    widgets.widget_destroy(l)
}

@(test)
test_label_wrap_multi_line :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    l := widgets.label_create("This is a longer text that should wrap to multiple lines", &font)
    l.rect = core.Rect{0, 0, 150, 100}  // Narrow width to force wrapping
    l.padding = widgets.Edges{0, 0, 0, 0}

    widgets.label_layout(l)

    testing.expect(t, len(l.lines) > 1, "Long text in narrow width should wrap to multiple lines")

    widgets.widget_destroy(l)
}

@(test)
test_label_wrap_newlines :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    l := widgets.label_create("Line one\nLine two\nLine three", &font)
    l.rect = core.Rect{0, 0, 500, 100}
    l.padding = widgets.Edges{0, 0, 0, 0}

    widgets.label_layout(l)

    testing.expect(t, len(l.lines) == 3, "Text with newlines should have 3 lines")
    if len(l.lines) >= 3 {
        testing.expect(t, l.lines[0] == "Line one", "First line should be 'Line one'")
        testing.expect(t, l.lines[1] == "Line two", "Second line should be 'Line two'")
        testing.expect(t, l.lines[2] == "Line three", "Third line should be 'Line three'")
    }

    widgets.widget_destroy(l)
}

@(test)
test_label_no_wrap :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    l := widgets.label_create("This is a longer text that would normally wrap", &font)
    l.rect = core.Rect{0, 0, 100, 100}  // Narrow width
    l.wrap = false

    widgets.label_layout(l)

    testing.expect(t, len(l.lines) == 1, "Text with wrap=false should be one line")

    widgets.widget_destroy(l)
}

// ============================================================================
// Label measurement tests
// ============================================================================

@(test)
test_label_measure_empty :: proc(t: ^testing.T) {
    l := widgets.label_create()
    l.padding = widgets.Edges{5, 10, 5, 10}

    size := widgets.label_measure(l)

    testing.expect(t, size.width == 20, "Empty label width should be left+right padding")
    testing.expect(t, size.height == 10, "Empty label height should be top+bottom padding")

    widgets.widget_destroy(l)
}

@(test)
test_label_measure_with_text :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    l := widgets.label_create("Hello", &font)
    l.padding = widgets.Edges{5, 10, 5, 10}
    l.wrap = false  // Single line for predictable measurement

    size := widgets.label_measure(l)

    text_width := render.text_measure(&font, "Hello")
    expected_width := text_width + 20  // 10 left + 10 right padding
    expected_height := font.line_height + 10  // 5 top + 5 bottom padding

    testing.expect(t, size.width == expected_width, "Width should be text width + padding")
    testing.expect(t, size.height == expected_height, "Height should be line height + padding")

    widgets.widget_destroy(l)
}

@(test)
test_label_measure_multiline :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    l := widgets.label_create("Line one\nLine two", &font)
    l.rect = core.Rect{0, 0, 300, 0}  // Set width constraint
    l.padding = widgets.Edges{5, 10, 5, 10}

    size := widgets.label_measure(l)

    expected_height := font.line_height * 2 + 10  // 2 lines + padding

    testing.expect(t, size.height == expected_height, "Height should be 2 lines + padding")

    widgets.widget_destroy(l)
}

// ============================================================================
// Label handle_event test
// ============================================================================

@(test)
test_label_handle_event :: proc(t: ^testing.T) {
    l := widgets.label_create("Test")

    event := core.Event{}
    result := widgets.label_handle_event(l, &event)

    testing.expect(t, result == false, "Label should not consume events")

    widgets.widget_destroy(l)
}

// ============================================================================
// Label cache invalidation tests
// ============================================================================

@(test)
test_label_cache_invalidation_on_text_change :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    l := widgets.label_create("Original", &font)
    l.rect = core.Rect{0, 0, 300, 100}

    widgets.label_layout(l)
    original_lines := len(l.lines)

    widgets.label_set_text(l, "New text that is different\nwith newline")
    widgets.label_layout(l)

    testing.expect(t, len(l.lines) > original_lines, "Lines should be recalculated after text change")

    widgets.widget_destroy(l)
}
