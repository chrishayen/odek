package tests

import "../src/core"
import "../src/widgets"
import "core:testing"

// ============================================================================
// Dropdown creation and destruction
// ============================================================================

@(test)
test_dropdown_create :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    testing.expect(t, dd != nil, "dropdown_create should return non-nil")
    testing.expect(t, dd.visible == true, "Dropdown should be visible by default")
    testing.expect(t, dd.enabled == true, "Dropdown should be enabled by default")
    testing.expect(t, dd.dirty == true, "Dropdown should be dirty by default")
    testing.expect(t, dd.focusable == true, "Dropdown should be focusable by default")
    testing.expect(t, dd.is_open == false, "Dropdown should be closed by default")
    testing.expect(t, dd.selected_index == 0, "Dropdown should have index 0 selected by default")
    testing.expect(t, dd.vtable == &widgets.dropdown_vtable, "Dropdown should have dropdown vtable")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_destroy :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "Option 1")
    widgets.dropdown_add_option(dd, "Option 2")
    // Should not crash
    widgets.widget_destroy(dd)
}

// ============================================================================
// Dropdown options management
// ============================================================================

@(test)
test_dropdown_add_option :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    dd.dirty = false

    widgets.dropdown_add_option(dd, "Option 1")
    testing.expect(t, len(dd.options) == 1, "Dropdown should have 1 option")
    testing.expect(t, dd.options[0] == "Option 1", "Option should match")
    testing.expect(t, dd.dirty == true, "Dropdown should be marked dirty")

    widgets.dropdown_add_option(dd, "Option 2")
    testing.expect(t, len(dd.options) == 2, "Dropdown should have 2 options")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_set_options :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "Old Option")
    dd.dirty = false

    widgets.dropdown_set_options(dd, {"A", "B", "C"})
    testing.expect(t, len(dd.options) == 3, "Dropdown should have 3 options")
    testing.expect(t, dd.options[0] == "A", "First option should be A")
    testing.expect(t, dd.options[1] == "B", "Second option should be B")
    testing.expect(t, dd.options[2] == "C", "Third option should be C")
    testing.expect(t, dd.dirty == true, "Dropdown should be marked dirty")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_set_options_resets_index :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "A")
    widgets.dropdown_add_option(dd, "B")
    widgets.dropdown_add_option(dd, "C")
    dd.selected_index = 2

    // Set fewer options - index should be clamped
    widgets.dropdown_set_options(dd, {"X", "Y"})
    testing.expect(t, dd.selected_index == 1, "Selected index should be clamped to valid range")

    widgets.widget_destroy(dd)
}

// ============================================================================
// Dropdown selection
// ============================================================================

@(test)
test_dropdown_set_selected :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "A")
    widgets.dropdown_add_option(dd, "B")
    widgets.dropdown_add_option(dd, "C")
    dd.dirty = false

    widgets.dropdown_set_selected(dd, 1)
    testing.expect(t, dd.selected_index == 1, "Selected index should be 1")
    testing.expect(t, dd.dirty == true, "Dropdown should be marked dirty")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_set_selected_invalid :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "A")
    widgets.dropdown_add_option(dd, "B")
    dd.dirty = false

    widgets.dropdown_set_selected(dd, 5)  // Invalid index
    testing.expect(t, dd.selected_index == 0, "Selected index should not change for invalid index")
    testing.expect(t, dd.dirty == false, "Dropdown should not be marked dirty")

    widgets.dropdown_set_selected(dd, -1)  // Negative index
    testing.expect(t, dd.selected_index == 0, "Selected index should not change for negative index")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_set_selected_same :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "A")
    dd.selected_index = 0
    dd.dirty = false

    widgets.dropdown_set_selected(dd, 0)
    testing.expect(t, dd.dirty == false, "Dropdown should not be marked dirty if index unchanged")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_get_selected :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "A")
    widgets.dropdown_add_option(dd, "B")
    dd.selected_index = 1

    testing.expect(t, widgets.dropdown_get_selected(dd) == 1, "Should return selected index")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_get_selected_text :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "Alpha")
    widgets.dropdown_add_option(dd, "Beta")
    dd.selected_index = 1

    testing.expect(t, widgets.dropdown_get_selected_text(dd) == "Beta", "Should return selected text")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_get_selected_text_empty :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    testing.expect(t, widgets.dropdown_get_selected_text(dd) == "", "Should return empty string for no options")

    widgets.widget_destroy(dd)
}

// ============================================================================
// Dropdown event handling
// ============================================================================

@(test)
test_dropdown_click_opens :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "A")
    testing.expect(t, dd.is_open == false, "Should be closed initially")

    press_event := core.Event{
        type = .Pointer_Button_Press,
        button = .Left,
    }
    result := widgets.dropdown_handle_event(dd, &press_event)

    testing.expect(t, result == true, "Should consume click event")
    testing.expect(t, dd.is_open == true, "Should be open after click")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_click_closes :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "A")
    dd.is_open = true

    press_event := core.Event{
        type = .Pointer_Button_Press,
        button = .Left,
        pointer_y = 0,  // Click outside options
    }
    result := widgets.dropdown_handle_event(dd, &press_event)

    testing.expect(t, result == true, "Should consume click event")
    testing.expect(t, dd.is_open == false, "Should be closed after click")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_disabled_ignores_events :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "A")
    dd.enabled = false

    press_event := core.Event{
        type = .Pointer_Button_Press,
        button = .Left,
    }
    result := widgets.dropdown_handle_event(dd, &press_event)

    testing.expect(t, result == false, "Disabled dropdown should not consume events")
    testing.expect(t, dd.is_open == false, "Disabled dropdown should not open")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_escape_closes :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "A")
    dd.is_open = true
    dd.focused = true

    key_event := core.Event{
        type = .Key_Press,
        keycode = u32(core.Keycode.Escape),
    }
    result := widgets.dropdown_handle_event(dd, &key_event)

    testing.expect(t, result == true, "Should consume escape key")
    testing.expect(t, dd.is_open == false, "Should be closed after escape")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_space_opens :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "A")
    dd.focused = true
    testing.expect(t, dd.is_open == false, "Should be closed initially")

    key_event := core.Event{
        type = .Key_Press,
        keycode = u32(core.Keycode.Space),
    }
    result := widgets.dropdown_handle_event(dd, &key_event)

    testing.expect(t, result == true, "Should consume space key")
    testing.expect(t, dd.is_open == true, "Should be open after space")

    widgets.widget_destroy(dd)
}

// ============================================================================
// Dropdown in widget tree
// ============================================================================

@(test)
test_dropdown_as_child :: proc(t: ^testing.T) {
    container := widgets.container_create()
    dropdown := widgets.dropdown_create()

    widgets.widget_add_child(container, dropdown)

    testing.expect(t, len(container.children) == 1, "Container should have 1 child")
    testing.expect(t, dropdown.parent == container, "Dropdown's parent should be container")

    widgets.widget_destroy(container)
}

// ============================================================================
// Dropdown measurement
// ============================================================================

@(test)
test_dropdown_measure :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "Short")
    widgets.dropdown_add_option(dd, "A Very Long Option Text")

    size := widgets.dropdown_measure(dd, -1)

    testing.expect(t, size.width > 0, "Width should be positive")
    testing.expect(t, size.height > 0, "Height should be positive")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_measure_respects_min_size :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    dd.min_size = core.Size{200, 40}

    size := widgets.dropdown_measure(dd, -1)

    testing.expect(t, size.width >= 200, "Width should respect min_size")
    testing.expect(t, size.height >= 40, "Height should respect min_size")

    widgets.widget_destroy(dd)
}

// ============================================================================
// Dropdown panel positioning
// ============================================================================

@(test)
test_dropdown_panel_opens_downward_when_space :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "Option 1")
    widgets.dropdown_add_option(dd, "Option 2")
    widgets.dropdown_add_option(dd, "Option 3")

    // Position dropdown near top of window (plenty of space below)
    trigger_rect := core.Rect{x = 10, y = 10, width = 100, height = 30}
    window_height: i32 = 600

    panel_rect := widgets.dropdown_get_panel_rect(dd, trigger_rect, window_height)

    // Panel should be below the trigger (y > trigger bottom)
    testing.expect(t, panel_rect.y > trigger_rect.y + trigger_rect.height, "Panel should open below trigger")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_panel_opens_upward_when_no_space :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "Option 1")
    widgets.dropdown_add_option(dd, "Option 2")
    widgets.dropdown_add_option(dd, "Option 3")

    // Position dropdown near bottom of window (not enough space below)
    trigger_rect := core.Rect{x = 10, y = 550, width = 100, height = 30}
    window_height: i32 = 600

    panel_rect := widgets.dropdown_get_panel_rect(dd, trigger_rect, window_height)

    // Panel should be above the trigger (y < trigger top)
    testing.expect(t, panel_rect.y < trigger_rect.y, "Panel should open above trigger when near bottom")

    widgets.widget_destroy(dd)
}

@(test)
test_dropdown_full_rect_includes_upward_panel :: proc(t: ^testing.T) {
    dd := widgets.dropdown_create()
    widgets.dropdown_add_option(dd, "Option 1")
    widgets.dropdown_add_option(dd, "Option 2")
    widgets.dropdown_add_option(dd, "Option 3")

    // Set up dropdown near bottom of window
    dd.rect = core.Rect{x = 10, y = 550, width = 100, height = 30}
    dd.is_open = true

    // Set window context so dropdown knows window height
    widgets.window_context_set(800, 600)

    full_rect := widgets.dropdown_get_full_rect(dd)

    // Full rect should start above the trigger (includes upward panel)
    testing.expect(t, full_rect.y < dd.rect.y, "Full rect should include upward panel")

    widgets.widget_destroy(dd)
}
