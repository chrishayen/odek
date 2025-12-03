package tests

import "../src/core"
import "../src/widgets"
import "core:testing"

// ============================================================================
// Toggle group creation and destruction
// ============================================================================

@(test)
test_toggle_group_create :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)

	testing.expect(t, g != nil, "toggle_group_create should return non-nil")
	testing.expect(t, g.visible == true, "Toggle group should be visible by default")
	testing.expect(t, g.enabled == true, "Toggle group should be enabled by default")
	testing.expect(t, g.dirty == true, "Toggle group should be dirty by default")
	testing.expect(t, g.focusable == true, "Toggle group should be focusable by default")
	testing.expect(t, g.selected_index == 0, "Toggle group should have first item selected by default")
	testing.expect(t, len(g.options) == 3, "Toggle group should have 3 options")
	testing.expect(t, g.vtable == &widgets.toggle_group_vtable, "Toggle group should have toggle_group vtable")

	widgets.widget_destroy(g)
}

@(test)
test_toggle_group_destroy :: proc(t: ^testing.T) {
	options := []string{"A", "B"}
	g := widgets.toggle_group_create(options)
	// Should not crash
	widgets.widget_destroy(g)
}

// ============================================================================
// Toggle group state management
// ============================================================================

@(test)
test_toggle_group_set_selected :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)
	g.dirty = false

	widgets.toggle_group_set_selected(g, 2)
	testing.expect(t, g.selected_index == 2, "Selected index should be 2")
	testing.expect(t, g.dirty == true, "Toggle group should be marked dirty")

	widgets.widget_destroy(g)
}

@(test)
test_toggle_group_set_selected_same :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)
	g.selected_index = 1
	g.dirty = false

	widgets.toggle_group_set_selected(g, 1)
	testing.expect(t, g.dirty == false, "Toggle group should not be marked dirty if selection unchanged")

	widgets.widget_destroy(g)
}

@(test)
test_toggle_group_set_selected_out_of_bounds :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)
	g.dirty = false

	widgets.toggle_group_set_selected(g, 10)
	testing.expect(t, g.selected_index == 0, "Selected index should remain unchanged for out of bounds")
	testing.expect(t, g.dirty == false, "Toggle group should not be marked dirty for invalid index")

	widgets.toggle_group_set_selected(g, -1)
	testing.expect(t, g.selected_index == 0, "Selected index should remain unchanged for negative index")

	widgets.widget_destroy(g)
}

@(test)
test_toggle_group_get_selected :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)

	testing.expect(t, widgets.toggle_group_get_selected(g) == 0, "Should return 0 initially")

	g.selected_index = 2
	testing.expect(t, widgets.toggle_group_get_selected(g) == 2, "Should return 2 after change")

	widgets.widget_destroy(g)
}

@(test)
test_toggle_group_get_selected_text :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)

	testing.expect(t, widgets.toggle_group_get_selected_text(g) == "One", "Should return 'One' initially")

	g.selected_index = 1
	testing.expect(t, widgets.toggle_group_get_selected_text(g) == "Two", "Should return 'Two' after change")

	widgets.widget_destroy(g)
}

// ============================================================================
// Toggle group event handling
// ============================================================================

@(test)
test_toggle_group_hover_state :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)
	g.rect = core.Rect{0, 0, 150, 30}
	g.dirty = false

	enter_event := core.Event{type = .Pointer_Enter, pointer_x = 10, pointer_y = 15}
	result := widgets.toggle_group_handle_event(g, &enter_event)

	testing.expect(t, result == true, "Toggle group should consume enter event")
	testing.expect(t, g.hovered_index >= 0, "Toggle group should have hovered item")
	testing.expect(t, g.dirty == true, "Toggle group should be marked dirty")

	widgets.widget_destroy(g)
}

@(test)
test_toggle_group_leave_state :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)
	g.hovered_index = 1
	g.dirty = false

	leave_event := core.Event{type = .Pointer_Leave}
	result := widgets.toggle_group_handle_event(g, &leave_event)

	testing.expect(t, result == true, "Toggle group should consume leave event")
	testing.expect(t, g.hovered_index == -1, "Hovered index should be -1")
	testing.expect(t, g.dirty == true, "Toggle group should be marked dirty")

	widgets.widget_destroy(g)
}

@(test)
test_toggle_group_disabled_ignores_events :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)
	g.enabled = false

	enter_event := core.Event{type = .Pointer_Enter, pointer_x = 10, pointer_y = 15}
	result := widgets.toggle_group_handle_event(g, &enter_event)

	testing.expect(t, result == false, "Disabled toggle group should not consume events")
	testing.expect(t, g.hovered_index == -1, "Disabled toggle group state should not change")

	widgets.widget_destroy(g)
}

@(test)
test_toggle_group_keyboard_navigation :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)
	g.focused = true
	testing.expect(t, g.selected_index == 0, "Initial selection should be 0")

	// Press Right
	right_event := core.Event{
		type = .Key_Press,
		keysym = u32(core.Keysym.Right),
	}
	result := widgets.toggle_group_handle_event(g, &right_event)
	testing.expect(t, result == true, "Toggle group should consume right key when focused")
	testing.expect(t, g.selected_index == 1, "Selection should move to 1")

	// Press Right again
	widgets.toggle_group_handle_event(g, &right_event)
	testing.expect(t, g.selected_index == 2, "Selection should move to 2")

	// Press Right at end - should stay at 2
	widgets.toggle_group_handle_event(g, &right_event)
	testing.expect(t, g.selected_index == 2, "Selection should stay at 2 (end)")

	// Press Left
	left_event := core.Event{
		type = .Key_Press,
		keysym = u32(core.Keysym.Left),
	}
	result = widgets.toggle_group_handle_event(g, &left_event)
	testing.expect(t, result == true, "Toggle group should consume left key when focused")
	testing.expect(t, g.selected_index == 1, "Selection should move to 1")

	widgets.widget_destroy(g)
}

@(test)
test_toggle_group_keyboard_at_bounds :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)
	g.focused = true
	g.selected_index = 0

	// Press Left at start - should stay at 0
	left_event := core.Event{
		type = .Key_Press,
		keysym = u32(core.Keysym.Left),
	}
	result := widgets.toggle_group_handle_event(g, &left_event)
	testing.expect(t, result == false, "Left at start should not be consumed")
	testing.expect(t, g.selected_index == 0, "Selection should stay at 0")

	widgets.widget_destroy(g)
}

// ============================================================================
// Toggle group in widget tree
// ============================================================================

@(test)
test_toggle_group_as_child :: proc(t: ^testing.T) {
	container := widgets.container_create()
	options := []string{"A", "B", "C"}
	toggle := widgets.toggle_group_create(options)

	widgets.widget_add_child(container, toggle)

	testing.expect(t, len(container.children) == 1, "Container should have 1 child")
	testing.expect(t, toggle.parent == container, "Toggle group's parent should be container")

	widgets.widget_destroy(container)
}

// ============================================================================
// Toggle group callback
// ============================================================================

@(test)
test_toggle_group_on_change_callback :: proc(t: ^testing.T) {
	options := []string{"One", "Two", "Three"}
	g := widgets.toggle_group_create(options)
	g.focused = true

	callback_called := false
	g.on_change = proc(group: ^widgets.Toggle_Group) {
		// Can't easily test this in Odin without context
	}

	// Just verify callback can be set without crash
	testing.expect(t, g.on_change != nil, "Callback should be set")

	widgets.widget_destroy(g)
}
