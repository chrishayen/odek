package tests

import "../src/core"
import "../src/widgets"
import "core:testing"

// ============================================================================
// Checkbox creation and destruction
// ============================================================================

@(test)
test_checkbox_create :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    testing.expect(t, cb != nil, "checkbox_create should return non-nil")
    testing.expect(t, cb.visible == true, "Checkbox should be visible by default")
    testing.expect(t, cb.enabled == true, "Checkbox should be enabled by default")
    testing.expect(t, cb.dirty == true, "Checkbox should be dirty by default")
    testing.expect(t, cb.focusable == true, "Checkbox should be focusable by default")
    testing.expect(t, cb.checked == false, "Checkbox should be unchecked by default")
    testing.expect(t, cb.vtable == &widgets.checkbox_vtable, "Checkbox should have checkbox vtable")

    widgets.widget_destroy(cb)
}

@(test)
test_checkbox_destroy :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    // Should not crash
    widgets.widget_destroy(cb)
}

// ============================================================================
// Checkbox state management
// ============================================================================

@(test)
test_checkbox_set_checked :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    cb.dirty = false

    widgets.checkbox_set_checked(cb, true)
    testing.expect(t, cb.checked == true, "Checkbox should be checked")
    testing.expect(t, cb.dirty == true, "Checkbox should be marked dirty")

    widgets.widget_destroy(cb)
}

@(test)
test_checkbox_set_checked_same :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    cb.checked = true
    cb.dirty = false

    widgets.checkbox_set_checked(cb, true)
    testing.expect(t, cb.dirty == false, "Checkbox should not be marked dirty if state unchanged")

    widgets.widget_destroy(cb)
}

@(test)
test_checkbox_toggle :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    testing.expect(t, cb.checked == false, "Initial state should be unchecked")

    widgets.checkbox_toggle(cb)
    testing.expect(t, cb.checked == true, "After toggle should be checked")

    widgets.checkbox_toggle(cb)
    testing.expect(t, cb.checked == false, "After second toggle should be unchecked")

    widgets.widget_destroy(cb)
}

@(test)
test_checkbox_is_checked :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    testing.expect(t, widgets.checkbox_is_checked(cb) == false, "Should return false when unchecked")

    cb.checked = true
    testing.expect(t, widgets.checkbox_is_checked(cb) == true, "Should return true when checked")

    widgets.widget_destroy(cb)
}

// ============================================================================
// Checkbox event handling
// ============================================================================

@(test)
test_checkbox_hover_state :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    cb.dirty = false

    enter_event := core.Event{type = .Pointer_Enter}
    result := widgets.checkbox_handle_event(cb, &enter_event)

    testing.expect(t, result == true, "Checkbox should consume enter event")
    testing.expect(t, cb.hovered_internal == true, "Checkbox should be hovered")
    testing.expect(t, cb.dirty == true, "Checkbox should be marked dirty")

    widgets.widget_destroy(cb)
}

@(test)
test_checkbox_leave_state :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    cb.hovered_internal = true
    cb.dirty = false

    leave_event := core.Event{type = .Pointer_Leave}
    result := widgets.checkbox_handle_event(cb, &leave_event)

    testing.expect(t, result == true, "Checkbox should consume leave event")
    testing.expect(t, cb.hovered_internal == false, "Checkbox should not be hovered")
    testing.expect(t, cb.dirty == true, "Checkbox should be marked dirty")

    widgets.widget_destroy(cb)
}

@(test)
test_checkbox_click_toggles :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    testing.expect(t, cb.checked == false, "Initial state should be unchecked")

    // Press
    press_event := core.Event{
        type = .Pointer_Button_Press,
        button = .Left,
    }
    widgets.checkbox_handle_event(cb, &press_event)
    testing.expect(t, cb.pressed == true, "Checkbox should be pressed")

    // Release - should toggle
    release_event := core.Event{
        type = .Pointer_Button_Release,
        button = .Left,
    }
    widgets.checkbox_handle_event(cb, &release_event)
    testing.expect(t, cb.checked == true, "Checkbox should be checked after click")

    widgets.widget_destroy(cb)
}

@(test)
test_checkbox_right_click_ignored :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    cb.dirty = false

    press_event := core.Event{
        type = .Pointer_Button_Press,
        button = .Right,
    }
    result := widgets.checkbox_handle_event(cb, &press_event)

    testing.expect(t, result == false, "Checkbox should not consume right click")
    testing.expect(t, cb.pressed == false, "Checkbox should not be pressed")

    widgets.widget_destroy(cb)
}

@(test)
test_checkbox_disabled_ignores_events :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    cb.enabled = false

    enter_event := core.Event{type = .Pointer_Enter}
    result := widgets.checkbox_handle_event(cb, &enter_event)

    testing.expect(t, result == false, "Disabled checkbox should not consume events")
    testing.expect(t, cb.hovered_internal == false, "Disabled checkbox state should not change")

    widgets.widget_destroy(cb)
}

@(test)
test_checkbox_space_toggles_when_focused :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    cb.focused = true
    testing.expect(t, cb.checked == false, "Initial state should be unchecked")

    key_event := core.Event{
        type = .Key_Press,
        keycode = u32(core.Keycode.Space),
    }
    result := widgets.checkbox_handle_event(cb, &key_event)

    testing.expect(t, result == true, "Checkbox should consume space key when focused")
    testing.expect(t, cb.checked == true, "Checkbox should be checked after space")

    widgets.widget_destroy(cb)
}

// ============================================================================
// Checkbox in widget tree
// ============================================================================

@(test)
test_checkbox_as_child :: proc(t: ^testing.T) {
    container := widgets.container_create()
    checkbox := widgets.checkbox_create()

    widgets.widget_add_child(container, checkbox)

    testing.expect(t, len(container.children) == 1, "Container should have 1 child")
    testing.expect(t, checkbox.parent == container, "Checkbox's parent should be container")

    widgets.widget_destroy(container)
}

// ============================================================================
// Checkbox measurement
// ============================================================================

@(test)
test_checkbox_measure :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    cb.box_size = 18
    cb.padding = widgets.edges_all(2)

    size := widgets.checkbox_measure(cb, -1)

    expected := 18 + 2 + 2  // box_size + padding
    testing.expect(t, size.width == i32(expected), "Width should be box_size + padding")
    testing.expect(t, size.height == i32(expected), "Height should be box_size + padding")

    widgets.widget_destroy(cb)
}

@(test)
test_checkbox_measure_respects_min_size :: proc(t: ^testing.T) {
    cb := widgets.checkbox_create()
    cb.min_size = core.Size{100, 50}

    size := widgets.checkbox_measure(cb, -1)

    testing.expect(t, size.width >= 100, "Width should respect min_size")
    testing.expect(t, size.height >= 50, "Height should respect min_size")

    widgets.widget_destroy(cb)
}
