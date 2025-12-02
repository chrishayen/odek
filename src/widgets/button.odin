package widgets

import "../core"
import "../render"

// Button state
Button_State :: enum {
    Normal,
    Hovered,
    Pressed,
}

// Button widget with click handling
Button :: struct {
    using base: Widget,

    text:       string,
    font:       ^render.Font,

    // Colors for each state
    bg_normal:  core.Color,
    bg_hover:   core.Color,
    bg_pressed: core.Color,
    text_color: core.Color,

    // Visual properties
    corner_radius: i32,

    // State
    state:      Button_State,

    // Callback
    on_click:   proc(button: ^Button),
}

// Shared vtable for all buttons
button_vtable := Widget_VTable{
    draw         = button_draw,
    handle_event = button_handle_event,
    layout       = button_layout,
    destroy      = button_destroy,
    measure      = button_measure,
}

// Create a new button
button_create :: proc(text: string = "", font: ^render.Font = nil) -> ^Button {
    b := new(Button)
    b.vtable = &button_vtable
    b.visible = true
    b.enabled = true
    b.dirty = true
    b.focusable = true  // Buttons can receive focus
    b.text = text
    b.font = font

    // Default colors from theme
    theme := theme_get()
    b.bg_normal = theme.accent
    b.bg_hover = theme.accent_hover
    b.bg_pressed = theme.accent_pressed
    b.text_color = theme.text_on_accent

    b.corner_radius = 4
    b.state = .Normal
    b.padding = edges_symmetric(16, 8)  // Horizontal 16, vertical 8

    return b
}

// Set button text
button_set_text :: proc(b: ^Button, text: string) {
    if b.text == text {
        return
    }
    b.text = text
    widget_mark_dirty(b)
}

// Set button font
button_set_font :: proc(b: ^Button, font: ^render.Font) {
    if b.font == font {
        return
    }
    b.font = font
    widget_mark_dirty(b)
}

// Set button colors
button_set_colors :: proc(b: ^Button, normal, hover, pressed: core.Color) {
    b.bg_normal = normal
    b.bg_hover = hover
    b.bg_pressed = pressed
    widget_mark_dirty(b)
}

// Set text color
button_set_text_color :: proc(b: ^Button, color: core.Color) {
    b.text_color = color
    widget_mark_dirty(b)
}

// Set click callback
button_set_on_click :: proc(b: ^Button, callback: proc(button: ^Button)) {
    b.on_click = callback
}

// Draw button
button_draw :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
    b := cast(^Button)w
    abs_rect := widget_get_absolute_rect(w)

    // Choose background color based on state
    bg_color: core.Color
    text_color := b.text_color

    if !w.enabled {
        // Disabled: use gray colors with readable text
        bg_color = core.Color{80, 80, 80, 255}
        text_color = core.Color{140, 140, 140, 255}
    } else {
        switch b.state {
        case .Normal:
            bg_color = b.bg_normal
        case .Hovered:
            bg_color = b.bg_hover
        case .Pressed:
            bg_color = b.bg_pressed
        }
    }

    // Draw background
    if b.corner_radius > 0 {
        render.fill_rounded_rect(ctx, abs_rect, b.corner_radius, bg_color)
    } else {
        render.fill_rect(ctx, abs_rect, bg_color)
    }

    // Draw text centered (use logical pixels for layout)
    if b.font != nil && b.text != "" {
        text_width := render.text_measure_logical(b.font, b.text)
        text_x := abs_rect.x + (abs_rect.width - text_width) / 2
        text_y := abs_rect.y + (abs_rect.height - render.font_get_logical_line_height(b.font)) / 2
        render.draw_text_top(ctx, b.font, b.text, text_x, text_y, text_color)
    }

    // Draw focus indicator if focused and enabled
    if w.focused && w.enabled {
        // Draw a subtle border
        focus_color := core.color_rgba(255, 255, 255, 100)
        render.draw_rounded_rect(ctx, abs_rect, b.corner_radius, focus_color)
    }
}

// Handle button events
button_handle_event :: proc(w: ^Widget, event: ^core.Event) -> bool {
    b := cast(^Button)w

    if !w.enabled {
        return false
    }

    #partial switch event.type {
    case .Pointer_Enter:
        b.state = .Hovered
        widget_mark_dirty(b)
        return true

    case .Pointer_Leave:
        b.state = .Normal
        widget_mark_dirty(b)
        return true

    case .Pointer_Button_Press:
        if event.button == .Left {
            b.state = .Pressed
            widget_mark_dirty(b)
            return true
        }

    case .Pointer_Button_Release:
        if event.button == .Left && b.state == .Pressed {
            b.state = .Hovered
            widget_mark_dirty(b)
            // Fire callback
            if b.on_click != nil {
                b.on_click(b)
            }
            return true
        }

    case .Key_Press:
        // Enter or Space activates the button when focused
        if w.focused && (event.keycode == u32(core.Keycode.Enter) || event.keycode == u32(core.Keycode.Space)) {
            if b.on_click != nil {
                // Visual feedback
                b.state = .Pressed
                widget_mark_dirty(b)
                b.on_click(b)
                b.state = .Hovered if w.hovered else .Normal
                widget_mark_dirty(b)
            }
            return true
        }
    }

    return false
}

// Button layout (no children)
button_layout :: proc(w: ^Widget) {
    // No children to layout
}

// Measure button's preferred size
button_measure :: proc(w: ^Widget, available_width: i32) -> core.Size {
    b := cast(^Button)w

    text_width: i32 = 0
    text_height: i32 = 0

    if b.font != nil && b.text != "" {
        text_size := render.text_measure_size(b.font, b.text)
        text_width = text_size.width
        text_height = text_size.height
    } else {
        text_height = 16  // Default height if no font
    }

    return core.Size{
        width = max(text_width + w.padding.left + w.padding.right, w.min_size.width),
        height = max(text_height + w.padding.top + w.padding.bottom, w.min_size.height),
    }
}

// Destroy button resources
button_destroy :: proc(w: ^Widget) {
    // No button-specific cleanup needed
}
