package tests

import "../src/core"
import "../src/widgets"
import "core:testing"

// ============================================================================
// Widget Destruction and Global State Cleanup Tests
// ============================================================================
// These tests verify that destroying widgets properly clears global state
// pointers (hit_state.hovered, focus_manager.focused) to prevent segfaults
// when callbacks destroy widgets.

@(test)
test_destroy_hovered_widget :: proc(t: ^testing.T) {
    // Setup: Create a widget and set it as hovered
    hit_state := widgets.Hit_Test_State{}
    widgets.hit_state_set_global(&hit_state)

    widget := widgets.widget_create(&widgets.default_vtable)
    hit_state.hovered = widget

    // Act: Destroy the widget
    widgets.widget_destroy(widget)

    // Assert: Hit state should be cleared
    testing.expect(t, hit_state.hovered == nil,
        "Hit state should be cleared when hovered widget is destroyed")

    // Cleanup
    widgets.hit_state_set_global(nil)
}

@(test)
test_destroy_focused_widget :: proc(t: ^testing.T) {
    // Setup: Create a widget and set it as focused
    root := widgets.widget_create(&widgets.default_vtable)
    focus_manager := widgets.focus_manager_init(root)
    widgets.focus_manager_set_global(&focus_manager)

    focusable := widgets.widget_create(&widgets.default_vtable)
    focusable.focusable = true
    widgets.widget_add_child(root, focusable)
    widgets.focus_set(&focus_manager, focusable)

    testing.expect(t, focus_manager.focused == focusable,
        "Focus should be set before destroy")

    // Act: Remove and destroy the focused widget
    widgets.widget_remove_child(root, focusable)
    widgets.widget_destroy(focusable)

    // Assert: Focus should be cleared (THIS WILL FAIL INITIALLY)
    testing.expect(t, focus_manager.focused == nil,
        "Focus should be cleared when focused widget is destroyed")

    // Cleanup
    widgets.focus_manager_set_global(nil)
    widgets.widget_destroy(root)
}

@(test)
test_destroy_widget_with_focused_child :: proc(t: ^testing.T) {
    // Setup: Create a parent with a focused child
    root := widgets.widget_create(&widgets.default_vtable)
    focus_manager := widgets.focus_manager_init(root)
    widgets.focus_manager_set_global(&focus_manager)

    parent := widgets.container_create()
    widgets.widget_add_child(root, parent)

    child := widgets.widget_create(&widgets.default_vtable)
    child.focusable = true
    widgets.widget_add_child(parent, child)
    widgets.focus_set(&focus_manager, child)

    testing.expect(t, focus_manager.focused == child,
        "Child should be focused before destroy")

    // Act: Remove and destroy the parent (which destroys child too)
    widgets.widget_remove_child(root, parent)
    widgets.widget_destroy(parent)

    // Assert: Focus should be cleared (THIS WILL FAIL INITIALLY)
    testing.expect(t, focus_manager.focused == nil,
        "Focus should be cleared when parent of focused widget is destroyed")

    // Cleanup
    widgets.focus_manager_set_global(nil)
    widgets.widget_destroy(root)
}

@(test)
test_button_callback_destroys_self :: proc(t: ^testing.T) {
    // Setup: Create a structure similar to the todo app
    // where clicking a button destroys it
    hit_state := widgets.Hit_Test_State{}
    widgets.hit_state_set_global(&hit_state)

    root := widgets.container_create()
    focus_manager := widgets.focus_manager_init(root)
    widgets.focus_manager_set_global(&focus_manager)

    // Create a button that will be destroyed
    button := widgets.button_create("Delete")
    button.focusable = true
    widgets.widget_add_child(root, button)

    // Set button as hovered and focused (simulating user interaction)
    hit_state.hovered = button
    widgets.focus_set(&focus_manager, button)

    testing.expect(t, hit_state.hovered == button, "Button should be hovered")
    testing.expect(t, focus_manager.focused == button, "Button should be focused")

    // Act: Simulate what happens when button callback destroys the widget tree
    // (In real app, callback removes todo item and rebuilds UI)
    widgets.widget_remove_child(root, button)
    widgets.widget_destroy(button)

    // Assert: Both hit state and focus should be cleared
    testing.expect(t, hit_state.hovered == nil,
        "Hit state should be cleared after button destroyed")
    testing.expect(t, focus_manager.focused == nil,
        "Focus should be cleared after button destroyed (THIS WILL FAIL)")

    // Cleanup
    widgets.hit_state_set_global(nil)
    widgets.focus_manager_set_global(nil)
    widgets.widget_destroy(root)
}

@(test)
test_destroy_non_hovered_widget_keeps_hover :: proc(t: ^testing.T) {
    // Setup: Hover one widget, destroy a different one
    hit_state := widgets.Hit_Test_State{}
    widgets.hit_state_set_global(&hit_state)

    widget1 := widgets.widget_create(&widgets.default_vtable)
    widget2 := widgets.widget_create(&widgets.default_vtable)
    hit_state.hovered = widget1

    // Act: Destroy widget2 (not the hovered one)
    widgets.widget_destroy(widget2)

    // Assert: Hit state should still point to widget1
    testing.expect(t, hit_state.hovered == widget1,
        "Hit state should not change when non-hovered widget is destroyed")

    // Cleanup
    widgets.hit_state_set_global(nil)
    widgets.widget_destroy(widget1)
}

@(test)
test_destroy_non_focused_widget_keeps_focus :: proc(t: ^testing.T) {
    // Setup: Focus one widget, destroy a different one
    root := widgets.widget_create(&widgets.default_vtable)
    focus_manager := widgets.focus_manager_init(root)
    widgets.focus_manager_set_global(&focus_manager)

    widget1 := widgets.widget_create(&widgets.default_vtable)
    widget1.focusable = true
    widgets.widget_add_child(root, widget1)

    widget2 := widgets.widget_create(&widgets.default_vtable)
    widgets.widget_add_child(root, widget2)

    widgets.focus_set(&focus_manager, widget1)

    // Act: Destroy widget2 (not the focused one)
    widgets.widget_remove_child(root, widget2)
    widgets.widget_destroy(widget2)

    // Assert: Focus should still be on widget1
    testing.expect(t, focus_manager.focused == widget1,
        "Focus should not change when non-focused widget is destroyed")

    // Cleanup
    widgets.focus_manager_set_global(nil)
    widgets.widget_destroy(root)
}
