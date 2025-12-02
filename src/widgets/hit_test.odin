package widgets

import "../core"

// Hit test state for tracking hover changes
Hit_Test_State :: struct {
    hovered:  ^Widget,  // currently hovered widget
    captured: ^Widget,  // widget that has captured pointer (for drag operations)
}

// Global hit state pointer (set by app package)
// Used to clear hovered widget when it's destroyed
g_hit_state: ^Hit_Test_State

// Set the global hit state pointer
hit_state_set_global :: proc(state: ^Hit_Test_State) {
    g_hit_state = state
}

// Clear hit state if it references the given widget
hit_state_clear_widget :: proc(w: ^Widget) {
    if g_hit_state != nil {
        if g_hit_state.hovered == w {
            g_hit_state.hovered = nil
        }
        if g_hit_state.captured == w {
            g_hit_state.captured = nil
        }
    }
}

// Capture pointer events to a widget (for drag operations)
pointer_capture :: proc(state: ^Hit_Test_State, w: ^Widget) {
    state.captured = w
}

// Release pointer capture
pointer_release :: proc(state: ^Hit_Test_State) {
    state.captured = nil
}

// Find the deepest widget at the given point (window coordinates)
// Returns nil if no widget contains the point
hit_test :: proc(root: ^Widget, x, y: i32) -> ^Widget {
    if root == nil || !root.visible {
        return nil
    }

    // First check for widgets with active overlays (popups, dropdowns, etc.)
    // These take priority over normal hit testing
    if overlay_widget := hit_test_overlays(root, x, y); overlay_widget != nil {
        return overlay_widget
    }

    return hit_test_normal(root, x, y)
}

// Check for widgets with active overlays that contain the point
hit_test_overlays :: proc(w: ^Widget, x, y: i32) -> ^Widget {
    if w == nil || !w.visible {
        return nil
    }

    // Check children first (depth-first, reverse order)
    for i := len(w.children) - 1; i >= 0; i -= 1 {
        if hit := hit_test_overlays(w.children[i], x, y); hit != nil {
            return hit
        }
    }

    // Check if this widget has a custom contains_point (indicates overlay capability)
    // and if the point is in the extended area (not just the base rect)
    if w.vtable != nil && w.vtable.contains_point != nil {
        abs_rect := widget_get_absolute_rect(w)
        base_contains := core.rect_contains(abs_rect, core.Point{x, y})
        custom_contains := w.vtable.contains_point(w, x, y)
        // If custom check passes but base rect doesn't, this is an overlay hit
        if custom_contains && !base_contains {
            return w
        }
    }

    return nil
}

// Normal hit testing within widget bounds
hit_test_normal :: proc(root: ^Widget, x, y: i32) -> ^Widget {
    if root == nil || !root.visible {
        return nil
    }

    // Check if point is in this widget
    if !widget_contains_point(root, x, y) {
        return nil
    }

    // Check children in reverse order (top-most first)
    // Later children are drawn on top, so they should be hit first
    for i := len(root.children) - 1; i >= 0; i -= 1 {
        child := root.children[i]
        if hit := hit_test_normal(child, x, y); hit != nil {
            return hit
        }
    }

    // No child hit, return this widget
    return root
}

// Update hover state and send enter/leave events
// Returns the newly hovered widget
update_hover :: proc(state: ^Hit_Test_State, root: ^Widget, x, y: i32) -> ^Widget {
    new_hover := hit_test(root, x, y)

    if new_hover != state.hovered {
        // Send leave event to old widget
        if state.hovered != nil {
            state.hovered.hovered = false
            leave_event := core.Event{type = .Pointer_Leave}
            widget_handle_event(state.hovered, &leave_event)
        }

        // Send enter event to new widget
        if new_hover != nil {
            new_hover.hovered = true
            enter_event := core.Event{
                type = .Pointer_Enter,
                pointer_x = x,
                pointer_y = y,
            }
            widget_handle_event(new_hover, &enter_event)
        }

        state.hovered = new_hover
    }

    return new_hover
}

// Dispatch a pointer event to the appropriate widget
// Handles motion, button press/release, scroll
dispatch_pointer_event :: proc(state: ^Hit_Test_State, root: ^Widget, event: ^core.Event) -> bool {
    if root == nil {
        return false
    }

    #partial switch event.type {
    case .Pointer_Motion:
        // If pointer is captured, send motion to captured widget
        if state.captured != nil {
            return widget_handle_event(state.captured, event)
        }

        // Update hover state first
        update_hover(state, root, event.pointer_x, event.pointer_y)

        // Dispatch motion to hovered widget
        if state.hovered != nil {
            return widget_handle_event(state.hovered, event)
        }

    case .Pointer_Button_Press:
        // Find widget at click position
        target := hit_test(root, event.pointer_x, event.pointer_y)
        if target != nil {
            return widget_handle_event(target, event)
        }

    case .Pointer_Button_Release:
        // If pointer is captured, send release to captured widget
        if state.captured != nil {
            result := widget_handle_event(state.captured, event)
            return result
        }
        // Otherwise find widget at position
        target := hit_test(root, event.pointer_x, event.pointer_y)
        if target != nil {
            return widget_handle_event(target, event)
        }

    case .Scroll:
        // Find widget at scroll position and bubble up until handled
        target := hit_test(root, event.pointer_x, event.pointer_y)
        for target != nil {
            if widget_handle_event(target, event) {
                return true
            }
            target = target.parent
        }

    case .Pointer_Enter, .Pointer_Leave:
        // These are generated internally by update_hover
        return false
    }

    return false
}

// Dispatch a keyboard event to the focused widget
// If not consumed, bubbles up to parent
dispatch_key_event :: proc(focused: ^Widget, event: ^core.Event) -> bool {
    if focused == nil {
        return false
    }

    // Try focused widget first
    if widget_handle_event(focused, event) {
        return true
    }

    // Bubble up to parent
    return dispatch_key_event(focused.parent, event)
}

// Find next focusable widget in tree order (for Tab navigation)
find_next_focusable :: proc(root: ^Widget, current: ^Widget) -> ^Widget {
    if root == nil {
        return nil
    }

    // Collect all focusable widgets in tree order
    focusables := make([dynamic]^Widget, context.temp_allocator)
    collect_focusables(root, &focusables)

    if len(focusables) == 0 {
        return nil
    }

    // Find current index
    current_idx := -1
    for i := 0; i < len(focusables); i += 1 {
        if focusables[i] == current {
            current_idx = i
            break
        }
    }

    // Return next (wrap around)
    next_idx := (current_idx + 1) % len(focusables)
    return focusables[next_idx]
}

// Find previous focusable widget (for Shift+Tab)
find_prev_focusable :: proc(root: ^Widget, current: ^Widget) -> ^Widget {
    if root == nil {
        return nil
    }

    focusables := make([dynamic]^Widget, context.temp_allocator)
    collect_focusables(root, &focusables)

    if len(focusables) == 0 {
        return nil
    }

    // Find current index
    current_idx := -1
    for i := 0; i < len(focusables); i += 1 {
        if focusables[i] == current {
            current_idx = i
            break
        }
    }

    // Return previous (wrap around)
    prev_idx := current_idx - 1
    if prev_idx < 0 {
        prev_idx = len(focusables) - 1
    }
    return focusables[prev_idx]
}

// Collect all focusable widgets in tree order
collect_focusables :: proc(w: ^Widget, out: ^[dynamic]^Widget) {
    if w == nil || !w.visible || !w.enabled {
        return
    }

    if w.focusable {
        append(out, w)
    }

    for child in w.children {
        collect_focusables(child, out)
    }
}
