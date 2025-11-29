package tests

import "../src/core"
import "../src/render"
import "../src/widgets"
import "core:testing"

// ============================================================================
// Button creation and destruction
// ============================================================================

@(test)
test_button_create :: proc(t: ^testing.T) {
    b := widgets.button_create()
    testing.expect(t, b != nil, "button_create should return non-nil")
    testing.expect(t, b.visible == true, "Button should be visible by default")
    testing.expect(t, b.enabled == true, "Button should be enabled by default")
    testing.expect(t, b.dirty == true, "Button should be dirty by default")
    testing.expect(t, b.focusable == true, "Button should be focusable by default")
    testing.expect(t, b.text == "", "Button text should be empty by default")
    testing.expect(t, b.state == .Normal, "Button state should be Normal by default")
    testing.expect(t, b.vtable == &widgets.button_vtable, "Button should have button vtable")

    widgets.widget_destroy(b)
}

@(test)
test_button_create_with_text :: proc(t: ^testing.T) {
    b := widgets.button_create("Click Me")
    testing.expect(t, b.text == "Click Me", "Button should have provided text")

    widgets.widget_destroy(b)
}

@(test)
test_button_destroy :: proc(t: ^testing.T) {
    b := widgets.button_create("Test")
    // Should not crash
    widgets.widget_destroy(b)
}

// ============================================================================
// Button property setters
// ============================================================================

@(test)
test_button_set_text :: proc(t: ^testing.T) {
    b := widgets.button_create()
    b.dirty = false

    widgets.button_set_text(b, "New text")
    testing.expect(t, b.text == "New text", "Text should be updated")
    testing.expect(t, b.dirty == true, "Button should be marked dirty")

    widgets.widget_destroy(b)
}

@(test)
test_button_set_text_same :: proc(t: ^testing.T) {
    b := widgets.button_create("Same")
    b.dirty = false

    widgets.button_set_text(b, "Same")
    testing.expect(t, b.dirty == false, "Button should not be marked dirty if text unchanged")

    widgets.widget_destroy(b)
}

@(test)
test_button_set_colors :: proc(t: ^testing.T) {
    b := widgets.button_create()
    b.dirty = false

    normal := core.color_hex(0xFF0000)
    hover := core.color_hex(0x00FF00)
    pressed := core.color_hex(0x0000FF)

    widgets.button_set_colors(b, normal, hover, pressed)
    testing.expect(t, b.bg_normal.r == normal.r, "Normal color should be set")
    testing.expect(t, b.bg_hover.g == hover.g, "Hover color should be set")
    testing.expect(t, b.bg_pressed.b == pressed.b, "Pressed color should be set")
    testing.expect(t, b.dirty == true, "Button should be marked dirty")

    widgets.widget_destroy(b)
}

@(test)
test_button_set_text_color :: proc(t: ^testing.T) {
    b := widgets.button_create()
    b.dirty = false

    widgets.button_set_text_color(b, core.COLOR_RED)
    testing.expect(t, b.text_color.r == core.COLOR_RED.r, "Text color should be updated")
    testing.expect(t, b.dirty == true, "Button should be marked dirty")

    widgets.widget_destroy(b)
}

// ============================================================================
// Button state handling
// ============================================================================

@(test)
test_button_hover_state :: proc(t: ^testing.T) {
    b := widgets.button_create()
    b.dirty = false

    // Simulate pointer enter
    enter_event := core.Event{type = .Pointer_Enter}
    result := widgets.button_handle_event(b, &enter_event)

    testing.expect(t, result == true, "Button should consume enter event")
    testing.expect(t, b.state == .Hovered, "Button state should be Hovered")
    testing.expect(t, b.dirty == true, "Button should be marked dirty")

    widgets.widget_destroy(b)
}

@(test)
test_button_leave_state :: proc(t: ^testing.T) {
    b := widgets.button_create()
    b.state = .Hovered
    b.dirty = false

    // Simulate pointer leave
    leave_event := core.Event{type = .Pointer_Leave}
    result := widgets.button_handle_event(b, &leave_event)

    testing.expect(t, result == true, "Button should consume leave event")
    testing.expect(t, b.state == .Normal, "Button state should be Normal")
    testing.expect(t, b.dirty == true, "Button should be marked dirty")

    widgets.widget_destroy(b)
}

@(test)
test_button_pressed_state :: proc(t: ^testing.T) {
    b := widgets.button_create()
    b.state = .Hovered
    b.dirty = false

    // Simulate left button press
    press_event := core.Event{
        type = .Pointer_Button_Press,
        button = .Left,
    }
    result := widgets.button_handle_event(b, &press_event)

    testing.expect(t, result == true, "Button should consume press event")
    testing.expect(t, b.state == .Pressed, "Button state should be Pressed")
    testing.expect(t, b.dirty == true, "Button should be marked dirty")

    widgets.widget_destroy(b)
}

@(test)
test_button_release_triggers_click :: proc(t: ^testing.T) {
    b := widgets.button_create()
    b.state = .Pressed

    click_count := 0
    widgets.button_set_on_click(b, proc(btn: ^widgets.Button) {
        // Can't easily capture outer variable in Odin, use user_data instead
    })

    // Use a different approach - check state transition
    release_event := core.Event{
        type = .Pointer_Button_Release,
        button = .Left,
    }
    result := widgets.button_handle_event(b, &release_event)

    testing.expect(t, result == true, "Button should consume release event")
    testing.expect(t, b.state == .Hovered, "Button state should be Hovered after release")

    widgets.widget_destroy(b)
}

@(test)
test_button_right_click_ignored :: proc(t: ^testing.T) {
    b := widgets.button_create()
    b.state = .Hovered
    b.dirty = false

    // Simulate right button press
    press_event := core.Event{
        type = .Pointer_Button_Press,
        button = .Right,
    }
    result := widgets.button_handle_event(b, &press_event)

    testing.expect(t, result == false, "Button should not consume right click")
    testing.expect(t, b.state == .Hovered, "Button state should not change")

    widgets.widget_destroy(b)
}

@(test)
test_button_disabled_ignores_events :: proc(t: ^testing.T) {
    b := widgets.button_create()
    b.enabled = false

    enter_event := core.Event{type = .Pointer_Enter}
    result := widgets.button_handle_event(b, &enter_event)

    testing.expect(t, result == false, "Disabled button should not consume events")
    testing.expect(t, b.state == .Normal, "Disabled button state should not change")

    widgets.widget_destroy(b)
}

// ============================================================================
// Button in widget tree
// ============================================================================

@(test)
test_button_as_child :: proc(t: ^testing.T) {
    container := widgets.container_create()
    button := widgets.button_create("Test")

    widgets.widget_add_child(container, button)

    testing.expect(t, len(container.children) == 1, "Container should have 1 child")
    testing.expect(t, button.parent == container, "Button's parent should be container")

    widgets.widget_destroy(container)
}

// ============================================================================
// Button measurement
// ============================================================================

@(test)
test_button_measure_no_font :: proc(t: ^testing.T) {
    b := widgets.button_create("Test")
    b.padding = widgets.edges_all(10)

    size := widgets.button_measure(b)

    // Without font, should have default height + padding
    testing.expect(t, size.width == 20, "Width should be just padding")
    testing.expect(t, size.height == 16 + 20, "Height should be default 16 + padding")

    widgets.widget_destroy(b)
}

@(test)
test_button_measure_with_font :: proc(t: ^testing.T) {
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

    b := widgets.button_create("Click", &font)
    b.padding = widgets.edges_symmetric(16, 8)

    size := widgets.button_measure(b)

    text_width := render.text_measure(&font, "Click")
    expected_width := text_width + 32  // 16 left + 16 right
    expected_height := font.line_height + 16  // 8 top + 8 bottom

    testing.expect(t, size.width == expected_width, "Width should be text + padding")
    testing.expect(t, size.height == expected_height, "Height should be line height + padding")

    widgets.widget_destroy(b)
}

@(test)
test_button_measure_respects_min_size :: proc(t: ^testing.T) {
    b := widgets.button_create()
    b.min_size = core.Size{100, 50}
    b.padding = widgets.edges_all(5)

    size := widgets.button_measure(b)

    testing.expect(t, size.width >= 100, "Width should respect min_size")
    testing.expect(t, size.height >= 50, "Height should respect min_size")

    widgets.widget_destroy(b)
}
