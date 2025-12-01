package widgets

import "../core"

// Focus manager tracks keyboard focus within a widget tree
Focus_Manager :: struct {
    root:    ^Widget,
    focused: ^Widget,
}

// Global focus manager pointer (set by app package)
// Used to clear focused widget when it's destroyed
g_focus_manager: ^Focus_Manager

// Set the global focus manager pointer
focus_manager_set_global :: proc(fm: ^Focus_Manager) {
    g_focus_manager = fm
}

// Clear focus if it references the given widget
focus_manager_clear_widget :: proc(w: ^Widget) {
    if g_focus_manager == nil || g_focus_manager.focused != w {
        return
    }
    g_focus_manager.focused.focused = false
    g_focus_manager.focused = nil
}

// Initialize focus manager for a widget tree
focus_manager_init :: proc(root: ^Widget) -> Focus_Manager {
    return Focus_Manager{
        root = root,
        focused = nil,
    }
}

// Set focus to a specific widget
// Handles unfocusing previous widget and marking dirty
focus_set :: proc(fm: ^Focus_Manager, widget: ^Widget) {
    if fm.focused == widget {
        return
    }

    // Unfocus current
    if fm.focused != nil {
        fm.focused.focused = false
        widget_mark_dirty(fm.focused)

        // Send unfocus event
        unfocus_event := core.Event{type = .Window_Unfocus}
        widget_handle_event(fm.focused, &unfocus_event)
    }

    // Focus new widget
    fm.focused = widget

    if widget != nil {
        widget.focused = true
        widget_mark_dirty(widget)

        // Send focus event
        focus_event := core.Event{type = .Window_Focus}
        widget_handle_event(widget, &focus_event)
    }
}

// Clear focus (no widget focused)
focus_clear :: proc(fm: ^Focus_Manager) {
    focus_set(fm, nil)
}

// Move focus to next focusable widget (Tab)
focus_next :: proc(fm: ^Focus_Manager) {
    next := find_next_focusable(fm.root, fm.focused)
    if next != nil {
        focus_set(fm, next)
    }
}

// Move focus to previous focusable widget (Shift+Tab)
focus_prev :: proc(fm: ^Focus_Manager) {
    prev := find_prev_focusable(fm.root, fm.focused)
    if prev != nil {
        focus_set(fm, prev)
    }
}

// Get currently focused widget
focus_get :: proc(fm: ^Focus_Manager) -> ^Widget {
    return fm.focused
}

// Check if a widget has focus
focus_has :: proc(fm: ^Focus_Manager, widget: ^Widget) -> bool {
    return fm.focused == widget
}

// Handle Tab/Shift+Tab key events for focus navigation
// Returns true if event was consumed
focus_handle_tab :: proc(fm: ^Focus_Manager, event: ^core.Event) -> bool {
    if event.type != .Key_Press {
        return false
    }

    // Check for Tab key
    if event.keysym == u32(core.Keysym.Tab) || event.keycode == u32(core.Keycode.Tab) {
        if .Shift in event.modifiers {
            focus_prev(fm)
        } else {
            focus_next(fm)
        }
        return true
    }

    return false
}

// Focus first focusable widget in tree
focus_first :: proc(fm: ^Focus_Manager) {
    if fm.root == nil {
        return
    }

    focusables := make([dynamic]^Widget, context.temp_allocator)
    collect_focusables(fm.root, &focusables)

    if len(focusables) > 0 {
        focus_set(fm, focusables[0])
    }
}

// Focus last focusable widget in tree
focus_last :: proc(fm: ^Focus_Manager) {
    if fm.root == nil {
        return
    }

    focusables := make([dynamic]^Widget, context.temp_allocator)
    collect_focusables(fm.root, &focusables)

    if len(focusables) > 0 {
        focus_set(fm, focusables[len(focusables) - 1])
    }
}
