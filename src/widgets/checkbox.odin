package widgets

import "../core"
import "../render"

// Checkbox widget with toggle state
Checkbox :: struct {
    using base: Widget,

    // State
    checked: bool,

    // Visual properties
    box_size:      i32,
    corner_radius: i32,

    // Colors
    box_color:       core.Color,
    box_hover_color: core.Color,
    check_color:     core.Color,
    border_color:    core.Color,

    // Interaction state
    hovered_internal: bool,
    pressed:          bool,

    // Callback
    on_change: proc(checkbox: ^Checkbox),
}

// Shared vtable for all checkboxes
checkbox_vtable := Widget_VTable{
    draw         = checkbox_draw,
    handle_event = checkbox_handle_event,
    layout       = checkbox_layout,
    destroy      = checkbox_destroy,
    measure      = checkbox_measure,
}

// Create a new checkbox
checkbox_create :: proc() -> ^Checkbox {
    cb := new(Checkbox)
    cb.vtable = &checkbox_vtable
    cb.visible = true
    cb.enabled = true
    cb.dirty = true
    cb.focusable = true
    cb.checked = false

    // Default appearance
    cb.box_size = 18
    cb.corner_radius = 3

    // Colors from theme
    theme := theme_get()
    cb.box_color = theme.input_bg
    cb.box_hover_color = theme.bg_hover
    cb.check_color = theme.accent
    cb.border_color = theme.input_border

    cb.padding = edges_all(2)

    return cb
}

// Set checked state
checkbox_set_checked :: proc(cb: ^Checkbox, checked: bool) {
    if cb.checked == checked {
        return
    }
    cb.checked = checked
    widget_mark_dirty(cb)
}

// Toggle checkbox state
checkbox_toggle :: proc(cb: ^Checkbox) {
    cb.checked = !cb.checked
    widget_mark_dirty(cb)
    if cb.on_change != nil {
        cb.on_change(cb)
    }
}

// Get checked state
checkbox_is_checked :: proc(cb: ^Checkbox) -> bool {
    return cb.checked
}

// Set on_change callback
checkbox_set_on_change :: proc(cb: ^Checkbox, callback: proc(checkbox: ^Checkbox)) {
    cb.on_change = callback
}

// Draw checkbox
checkbox_draw :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
    cb := cast(^Checkbox)w
    abs_rect := widget_get_absolute_rect(w)

    // Center the box in the widget area
    box_x := abs_rect.x + (abs_rect.width - cb.box_size) / 2
    box_y := abs_rect.y + (abs_rect.height - cb.box_size) / 2
    box_rect := core.Rect{box_x, box_y, cb.box_size, cb.box_size}

    // Background color based on hover state
    bg_color := cb.box_hover_color if cb.hovered_internal else cb.box_color

    // Draw box background
    if cb.corner_radius > 0 {
        render.fill_rounded_rect(ctx, box_rect, cb.corner_radius, bg_color)
    } else {
        render.fill_rect(ctx, box_rect, bg_color)
    }

    // Draw border
    render.draw_rect(ctx, box_rect, cb.border_color)

    // Draw check indicator when checked
    if cb.checked {
        // Inner filled square/rounded rect as check indicator
        inner_margin: i32 = 4
        inner_rect := core.Rect{
            box_x + inner_margin,
            box_y + inner_margin,
            cb.box_size - inner_margin * 2,
            cb.box_size - inner_margin * 2,
        }
        inner_radius := max(cb.corner_radius - 1, 0)
        if inner_radius > 0 {
            render.fill_rounded_rect(ctx, inner_rect, inner_radius, cb.check_color)
        } else {
            render.fill_rect(ctx, inner_rect, cb.check_color)
        }
    }

    // Draw focus indicator
    if w.focused {
        focus_rect := core.Rect{box_x - 2, box_y - 2, cb.box_size + 4, cb.box_size + 4}
        theme := theme_get()
        render.draw_rect(ctx, focus_rect, theme.border_focus)
    }
}

// Handle checkbox events
checkbox_handle_event :: proc(w: ^Widget, event: ^core.Event) -> bool {
    cb := cast(^Checkbox)w

    if !w.enabled {
        return false
    }

    #partial switch event.type {
    case .Pointer_Enter:
        cb.hovered_internal = true
        widget_mark_dirty(cb)
        return true

    case .Pointer_Leave:
        cb.hovered_internal = false
        cb.pressed = false
        widget_mark_dirty(cb)
        return true

    case .Pointer_Button_Press:
        if event.button == .Left {
            cb.pressed = true
            widget_mark_dirty(cb)
            return true
        }

    case .Pointer_Button_Release:
        if event.button == .Left && cb.pressed {
            cb.pressed = false
            checkbox_toggle(cb)
            return true
        }

    case .Key_Press:
        // Space toggles checkbox when focused
        if w.focused && event.keycode == u32(core.Keycode.Space) {
            checkbox_toggle(cb)
            return true
        }
    }

    return false
}

// Checkbox layout (no children)
checkbox_layout :: proc(w: ^Widget) {
    // No children to layout
}

// Measure checkbox preferred size
checkbox_measure :: proc(w: ^Widget, available_width: i32) -> core.Size {
    cb := cast(^Checkbox)w
    size := cb.box_size + w.padding.left + w.padding.right
    return core.Size{
        width = max(size, w.min_size.width),
        height = max(size, w.min_size.height),
    }
}

// Destroy checkbox resources
checkbox_destroy :: proc(w: ^Widget) {
    // No checkbox-specific cleanup needed
}
