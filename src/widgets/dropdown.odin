package widgets

import "../core"
import "../render"

// Dropdown widget for selecting from a list of options
Dropdown :: struct {
    using base: Widget,

    options:        [dynamic]string,
    selected_index: i32,
    is_open:        bool,
    hovered_index:  i32,
    font:           ^render.Font,

    // Visual properties
    corner_radius:  i32,
    item_height:    i32,

    // Callback
    on_change:      proc(dropdown: ^Dropdown, index: i32),
}

// Shared vtable for all dropdowns
dropdown_vtable := Widget_VTable {
    draw           = dropdown_draw,
    handle_event   = dropdown_handle_event,
    layout         = dropdown_layout,
    destroy        = dropdown_destroy,
    measure        = dropdown_measure,
    contains_point = dropdown_contains_point,
    draw_overlay   = dropdown_draw_overlay,
}

// Create a new dropdown
dropdown_create :: proc(font: ^render.Font = nil) -> ^Dropdown {
    d := new(Dropdown)
    d.vtable = &dropdown_vtable
    d.visible = true
    d.enabled = true
    d.dirty = true
    d.focusable = true
    d.font = font
    d.selected_index = 0
    d.hovered_index = -1
    d.is_open = false
    d.corner_radius = 4
    d.item_height = 28
    d.padding = edges_symmetric(12, 6)

    return d
}

// Add an option to the dropdown
dropdown_add_option :: proc(d: ^Dropdown, option: string) {
    append(&d.options, option)
    widget_mark_dirty(d)
}

// Set all options at once
dropdown_set_options :: proc(d: ^Dropdown, options: []string) {
    clear(&d.options)
    for opt in options {
        append(&d.options, opt)
    }
    if d.selected_index >= i32(len(d.options)) {
        d.selected_index = max(0, i32(len(d.options)) - 1)
    }
    widget_mark_dirty(d)
}

// Set selected index
dropdown_set_selected :: proc(d: ^Dropdown, index: i32) {
    if index < 0 || index >= i32(len(d.options)) {
        return
    }
    if d.selected_index == index {
        return
    }
    d.selected_index = index
    widget_mark_dirty(d)
}

// Get selected index
dropdown_get_selected :: proc(d: ^Dropdown) -> i32 {
    return d.selected_index
}

// Get selected option text
dropdown_get_selected_text :: proc(d: ^Dropdown) -> string {
    if d.selected_index < 0 || d.selected_index >= i32(len(d.options)) {
        return ""
    }
    return d.options[d.selected_index]
}

// Draw dropdown
dropdown_draw :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
    d := cast(^Dropdown)w
    abs_rect := widget_get_absolute_rect(w)
    theme := theme_get()

    // Draw trigger button
    bg_color := theme.input_bg
    if w.hovered && !d.is_open {
        bg_color = theme.bg_hover
    }

    render.fill_rounded_rect(ctx, abs_rect, d.corner_radius, bg_color)
    render.draw_rounded_rect(ctx, abs_rect, d.corner_radius, theme.input_border)

    // Draw selected text
    if d.font != nil && len(d.options) > 0 {
        text := dropdown_get_selected_text(d)
        text_y := abs_rect.y + (abs_rect.height - render.font_get_logical_line_height(d.font)) / 2
        text_x := abs_rect.x + d.padding.left
        render.draw_text_top(ctx, d.font, text, text_x, text_y, theme.text_primary)
    }

    // Draw chevron on the right
    chevron_size: i32 = 6
    chevron_x := abs_rect.x + abs_rect.width - d.padding.right - chevron_size
    chevron_y := abs_rect.y + (abs_rect.height - chevron_size) / 2

    if d.is_open {
        // Up chevron
        render.fill_triangle(
            ctx,
            chevron_x, chevron_y + chevron_size,
            chevron_x + chevron_size, chevron_y + chevron_size,
            chevron_x + chevron_size / 2, chevron_y,
            theme.text_secondary,
        )
    } else {
        // Down chevron
        render.fill_triangle(
            ctx,
            chevron_x, chevron_y,
            chevron_x + chevron_size, chevron_y,
            chevron_x + chevron_size / 2, chevron_y + chevron_size,
            theme.text_secondary,
        )
    }

    // Draw focus indicator
    if w.focused {
        render.draw_rounded_rect(ctx, abs_rect, d.corner_radius, theme.border_focus)
    }
}

// Draw overlay (the dropdown panel when open)
dropdown_draw_overlay :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
    d := cast(^Dropdown)w
    if !d.is_open {
        return
    }
    abs_rect := widget_get_absolute_rect(w)
    dropdown_draw_panel(d, ctx, abs_rect)
}

// Draw the dropdown options panel
dropdown_draw_panel :: proc(d: ^Dropdown, ctx: ^render.Draw_Context, trigger_rect: core.Rect) {
    theme := theme_get()

    panel_height := i32(len(d.options)) * d.item_height
    panel_rect := core.Rect {
        x      = trigger_rect.x,
        y      = trigger_rect.y + trigger_rect.height + 2,
        width  = trigger_rect.width,
        height = panel_height,
    }

    // Panel background
    render.fill_rounded_rect(ctx, panel_rect, d.corner_radius, theme.bg_secondary)
    render.draw_rounded_rect(ctx, panel_rect, d.corner_radius, theme.border)

    // Draw each option
    for opt, i in d.options {
        idx := i32(i)
        item_rect := core.Rect {
            x      = panel_rect.x,
            y      = panel_rect.y + idx * d.item_height,
            width  = panel_rect.width,
            height = d.item_height,
        }

        // Highlight hovered item (very subtle)
        if idx == d.hovered_index {
            highlight := core.Color{255, 255, 255, 15}
            render.fill_rect_blend(ctx, item_rect, highlight)
        }

        // Draw option text
        if d.font != nil {
            text_y := item_rect.y + (item_rect.height - render.font_get_logical_line_height(d.font)) / 2
            text_x := item_rect.x + d.padding.left
            render.draw_text_top(ctx, d.font, opt, text_x, text_y, theme.text_primary)
        }
    }
}

// Get the full rect including dropdown panel when open
dropdown_get_full_rect :: proc(d: ^Dropdown) -> core.Rect {
    abs_rect := widget_get_absolute_rect(d)
    if !d.is_open {
        return abs_rect
    }

    panel_height := i32(len(d.options)) * d.item_height
    return core.Rect {
        x      = abs_rect.x,
        y      = abs_rect.y,
        width  = abs_rect.width,
        height = abs_rect.height + 2 + panel_height,
    }
}

// Handle dropdown events
dropdown_handle_event :: proc(w: ^Widget, event: ^core.Event) -> bool {
    d := cast(^Dropdown)w

    if !w.enabled {
        return false
    }

    abs_rect := widget_get_absolute_rect(w)

    #partial switch event.type {
    case .Pointer_Enter:
        widget_mark_dirty(d)
        return true

    case .Pointer_Leave:
        if !d.is_open {
            d.hovered_index = -1
        }
        widget_mark_dirty(d)
        return true

    case .Pointer_Motion:
        if d.is_open {
            // Check if hovering over an option
            panel_y := abs_rect.y + abs_rect.height + 2
            if event.pointer_y >= panel_y {
                idx := (event.pointer_y - panel_y) / d.item_height
                if idx >= 0 && idx < i32(len(d.options)) {
                    if d.hovered_index != idx {
                        d.hovered_index = idx
                        widget_mark_dirty(d)
                    }
                }
            } else {
                if d.hovered_index != -1 {
                    d.hovered_index = -1
                    widget_mark_dirty(d)
                }
            }
            return true
        }

    case .Pointer_Button_Press:
        if event.button != .Left {
            return false
        }

        if d.is_open {
            // Check if clicking on an option
            panel_y := abs_rect.y + abs_rect.height + 2
            panel_height := i32(len(d.options)) * d.item_height

            if event.pointer_y >= panel_y && event.pointer_y < panel_y + panel_height {
                idx := (event.pointer_y - panel_y) / d.item_height
                if idx >= 0 && idx < i32(len(d.options)) {
                    old_index := d.selected_index
                    d.selected_index = idx
                    d.is_open = false
                    d.hovered_index = -1
                    widget_mark_dirty(d)
                    if d.on_change != nil && old_index != idx {
                        d.on_change(d, idx)
                    }
                    return true
                }
            }
            // Click outside panel - close
            d.is_open = false
            d.hovered_index = -1
            widget_mark_dirty(d)
            return true
        } else {
            // Toggle open
            d.is_open = true
            widget_mark_dirty(d)
            return true
        }

    case .Key_Press:
        if !w.focused {
            return false
        }

        keycode := core.Keycode(event.keycode)

        if keycode == .Escape && d.is_open {
            d.is_open = false
            d.hovered_index = -1
            widget_mark_dirty(d)
            return true
        }

        if keycode == .Enter || keycode == .Space {
            if d.is_open {
                // Select hovered item
                if d.hovered_index >= 0 && d.hovered_index < i32(len(d.options)) {
                    old_index := d.selected_index
                    d.selected_index = d.hovered_index
                    d.is_open = false
                    d.hovered_index = -1
                    widget_mark_dirty(d)
                    if d.on_change != nil && old_index != d.selected_index {
                        d.on_change(d, d.selected_index)
                    }
                } else {
                    d.is_open = false
                    widget_mark_dirty(d)
                }
            } else {
                d.is_open = true
                d.hovered_index = d.selected_index
                widget_mark_dirty(d)
            }
            return true
        }

        if keycode == .Up {
            if d.is_open {
                if d.hovered_index > 0 {
                    d.hovered_index -= 1
                } else if d.hovered_index < 0 {
                    d.hovered_index = d.selected_index
                }
            } else {
                if d.selected_index > 0 {
                    d.selected_index -= 1
                    widget_mark_dirty(d)
                    if d.on_change != nil {
                        d.on_change(d, d.selected_index)
                    }
                }
            }
            widget_mark_dirty(d)
            return true
        }

        if keycode == .Down {
            if d.is_open {
                if d.hovered_index < i32(len(d.options)) - 1 {
                    d.hovered_index += 1
                } else if d.hovered_index < 0 {
                    d.hovered_index = d.selected_index
                }
            } else {
                if d.selected_index < i32(len(d.options)) - 1 {
                    d.selected_index += 1
                    widget_mark_dirty(d)
                    if d.on_change != nil {
                        d.on_change(d, d.selected_index)
                    }
                }
            }
            widget_mark_dirty(d)
            return true
        }
    }

    return false
}

// Dropdown layout (no children)
dropdown_layout :: proc(w: ^Widget) {
    // No children to layout
}

// Measure dropdown's preferred size
dropdown_measure :: proc(w: ^Widget, available_width: i32) -> core.Size {
    d := cast(^Dropdown)w

    // Find widest option
    max_width: i32 = 80  // Minimum width
    text_height: i32 = 16

    if d.font != nil {
        text_height = render.font_get_logical_line_height(d.font)
        for opt in d.options {
            opt_width := render.text_measure_logical(d.font, opt)
            max_width = max(max_width, opt_width)
        }
    }

    // Add padding and space for chevron
    chevron_space: i32 = 20
    width := max_width + w.padding.left + w.padding.right + chevron_space
    height := text_height + w.padding.top + w.padding.bottom

    return core.Size {
        width  = max(width, w.min_size.width),
        height = max(height, w.min_size.height),
    }
}

// Custom hit testing - include panel area when open
dropdown_contains_point :: proc(w: ^Widget, x, y: i32) -> bool {
    d := cast(^Dropdown)w
    if !w.visible {
        return false
    }
    full_rect := dropdown_get_full_rect(d)
    return core.rect_contains(full_rect, core.Point{x, y})
}

// Close dropdown
dropdown_close :: proc(d: ^Dropdown) {
    if d.is_open {
        d.is_open = false
        d.hovered_index = -1
        widget_mark_dirty(d)
    }
}

// Close any open dropdowns that don't contain the given point
// Call this on mouse click to implement "click outside to close"
close_dropdowns_outside :: proc(root: ^Widget, x, y: i32) {
    if root == nil {
        return
    }

    // Check if this is a dropdown
    if root.vtable == &dropdown_vtable {
        d := cast(^Dropdown)root
        if d.is_open {
            full_rect := dropdown_get_full_rect(d)
            if !core.rect_contains(full_rect, core.Point{x, y}) {
                dropdown_close(d)
            }
        }
    }

    // Check children
    for child in root.children {
        close_dropdowns_outside(child, x, y)
    }
}

// Destroy dropdown resources
dropdown_destroy :: proc(w: ^Widget) {
    d := cast(^Dropdown)w
    delete(d.options)
}
