package widgets

import "../core"
import "../render"
import "core:strings"
import "core:unicode/utf8"

// TextInput is a single-line text entry widget
Text_Input :: struct {
    using base: Widget,

    // Text content
    text:        [dynamic]u8,
    cursor_pos:  int,  // Byte position in text

    // Display
    font:        ^render.Font,
    placeholder: string,

    // Scroll offset for long text
    scroll_offset: i32,

    // Colors (use theme defaults if zero)
    bg_color:     core.Color,
    text_color:   core.Color,
    cursor_color: core.Color,
    border_color: core.Color,
    focus_color:  core.Color,
    placeholder_color: core.Color,

    // Visual properties
    border_width:  i32,
    corner_radius: i32,

    // Cursor blink state
    cursor_visible: bool,
    cursor_blink_time: f64,

    // Callbacks
    on_change:   proc(input: ^Text_Input),
    on_submit:   proc(input: ^Text_Input),  // Called on Enter
}

// Shared vtable for text inputs
text_input_vtable := Widget_VTable{
    draw         = text_input_draw,
    handle_event = text_input_handle_event,
    layout       = text_input_layout,
    destroy      = text_input_destroy,
    measure      = text_input_measure,
}

// Create a new text input
text_input_create :: proc(font: ^render.Font = nil) -> ^Text_Input {
    ti := new(Text_Input)
    ti.vtable = &text_input_vtable
    ti.visible = true
    ti.enabled = true
    ti.dirty = true
    ti.focusable = true
    ti.cursor = .Text

    ti.font = font
    ti.cursor_pos = 0
    ti.scroll_offset = 0
    ti.cursor_visible = true

    // Default appearance from theme
    theme := theme_get()
    ti.bg_color = theme.input_bg
    ti.text_color = theme.text_primary
    ti.cursor_color = theme.text_primary
    ti.border_color = theme.input_border
    ti.focus_color = theme.input_focus
    ti.placeholder_color = theme.text_secondary

    ti.border_width = 1
    ti.corner_radius = 4
    ti.padding = edges_symmetric(8, 6)

    return ti
}

// Set the text content
text_input_set_text :: proc(ti: ^Text_Input, text: string) {
    clear(&ti.text)
    append(&ti.text, ..transmute([]u8)text)
    ti.cursor_pos = len(ti.text)
    text_input_ensure_cursor_visible(ti)
    widget_mark_dirty(ti)
}

// Get the text content as a string
text_input_get_text :: proc(ti: ^Text_Input) -> string {
    return string(ti.text[:])
}

// Set placeholder text
text_input_set_placeholder :: proc(ti: ^Text_Input, placeholder: string) {
    ti.placeholder = placeholder
    widget_mark_dirty(ti)
}

// Set font
text_input_set_font :: proc(ti: ^Text_Input, font: ^render.Font) {
    ti.font = font
    widget_mark_dirty(ti)
}

// Clear the text
text_input_clear :: proc(ti: ^Text_Input) {
    clear(&ti.text)
    ti.cursor_pos = 0
    ti.scroll_offset = 0
    widget_mark_dirty(ti)
}

// Get cursor position in pixels relative to text start
text_input_get_cursor_x :: proc(ti: ^Text_Input) -> i32 {
    if ti.font == nil || ti.cursor_pos == 0 {
        return 0
    }

    text_before_cursor := string(ti.text[:ti.cursor_pos])
    return render.text_measure_logical(ti.font, text_before_cursor)
}

// Ensure cursor is visible by adjusting scroll
text_input_ensure_cursor_visible :: proc(ti: ^Text_Input) {
    cursor_x := text_input_get_cursor_x(ti)
    content_width := ti.rect.width - ti.padding.left - ti.padding.right

    // Cursor is to the left of visible area
    if cursor_x < ti.scroll_offset {
        ti.scroll_offset = max(0, cursor_x - 10)
    }

    // Cursor is to the right of visible area
    if cursor_x > ti.scroll_offset + content_width {
        ti.scroll_offset = cursor_x - content_width + 10
    }
}

// Move cursor left by one character
text_input_cursor_left :: proc(ti: ^Text_Input) {
    if ti.cursor_pos > 0 {
        // Find start of previous UTF-8 character
        ti.cursor_pos -= 1
        for ti.cursor_pos > 0 && (ti.text[ti.cursor_pos] & 0xC0) == 0x80 {
            ti.cursor_pos -= 1
        }
        text_input_ensure_cursor_visible(ti)
        widget_mark_dirty(ti)
    }
}

// Move cursor right by one character
text_input_cursor_right :: proc(ti: ^Text_Input) {
    if ti.cursor_pos < len(ti.text) {
        // Skip to next UTF-8 character
        ti.cursor_pos += 1
        for ti.cursor_pos < len(ti.text) && (ti.text[ti.cursor_pos] & 0xC0) == 0x80 {
            ti.cursor_pos += 1
        }
        text_input_ensure_cursor_visible(ti)
        widget_mark_dirty(ti)
    }
}

// Move cursor to start
text_input_cursor_home :: proc(ti: ^Text_Input) {
    ti.cursor_pos = 0
    text_input_ensure_cursor_visible(ti)
    widget_mark_dirty(ti)
}

// Move cursor to end
text_input_cursor_end :: proc(ti: ^Text_Input) {
    ti.cursor_pos = len(ti.text)
    text_input_ensure_cursor_visible(ti)
    widget_mark_dirty(ti)
}

// Insert text at cursor position
text_input_insert :: proc(ti: ^Text_Input, text: string) {
    if len(text) == 0 {
        return
    }

    bytes := transmute([]u8)text

    // Insert bytes at cursor position
    if ti.cursor_pos >= len(ti.text) {
        append(&ti.text, ..bytes)
    } else {
        // Make room and insert
        old_len := len(ti.text)
        resize(&ti.text, old_len + len(bytes))
        // Shift existing bytes right
        for i := old_len - 1; i >= ti.cursor_pos; i -= 1 {
            ti.text[i + len(bytes)] = ti.text[i]
        }
        // Insert new bytes
        for i := 0; i < len(bytes); i += 1 {
            ti.text[ti.cursor_pos + i] = bytes[i]
        }
    }

    ti.cursor_pos += len(bytes)
    text_input_ensure_cursor_visible(ti)
    widget_mark_dirty(ti)

    if ti.on_change != nil {
        ti.on_change(ti)
    }
}

// Delete character before cursor (backspace)
text_input_backspace :: proc(ti: ^Text_Input) {
    if ti.cursor_pos == 0 {
        return
    }

    // Find start of previous UTF-8 character
    char_start := ti.cursor_pos - 1
    for char_start > 0 && (ti.text[char_start] & 0xC0) == 0x80 {
        char_start -= 1
    }

    // Remove bytes from char_start to cursor_pos
    bytes_to_remove := ti.cursor_pos - char_start
    ordered_remove_range(&ti.text, char_start, ti.cursor_pos)
    ti.cursor_pos = char_start

    text_input_ensure_cursor_visible(ti)
    widget_mark_dirty(ti)

    if ti.on_change != nil {
        ti.on_change(ti)
    }
}

// Delete character after cursor (delete key)
text_input_delete :: proc(ti: ^Text_Input) {
    if ti.cursor_pos >= len(ti.text) {
        return
    }

    // Find end of current UTF-8 character
    char_end := ti.cursor_pos + 1
    for char_end < len(ti.text) && (ti.text[char_end] & 0xC0) == 0x80 {
        char_end += 1
    }

    // Remove bytes from cursor_pos to char_end
    ordered_remove_range(&ti.text, ti.cursor_pos, char_end)

    widget_mark_dirty(ti)

    if ti.on_change != nil {
        ti.on_change(ti)
    }
}

// Helper to remove a range from dynamic array
@(private)
ordered_remove_range :: proc(arr: ^[dynamic]u8, start, end: int) {
    if start >= end || start >= len(arr^) {
        return
    }
    end := min(end, len(arr^))
    count := end - start

    // Shift elements left
    for i := start; i < len(arr^) - count; i += 1 {
        arr[i] = arr[i + count]
    }
    resize(arr, len(arr^) - count)
}

// Layout implementation
text_input_layout :: proc(w: ^Widget) {
    // Text input has no children to layout
}

// Draw implementation
text_input_draw :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
    ti := cast(^Text_Input)w
    abs_rect := widget_get_absolute_rect(w)

    // Determine colors
    bg := ti.bg_color
    border := ti.focus_color if w.focused else ti.border_color

    // Draw background
    if ti.corner_radius > 0 {
        render.fill_rounded_rect(ctx, abs_rect, ti.corner_radius, bg)
    } else {
        render.fill_rect(ctx, abs_rect, bg)
    }

    // Draw border
    if ti.border_width > 0 {
        render.draw_rect(ctx, abs_rect, border, ti.border_width)
    }

    // Calculate text area
    text_x := abs_rect.x + ti.padding.left - ti.scroll_offset
    text_y := abs_rect.y + ti.padding.top
    content_width := abs_rect.width - ti.padding.left - ti.padding.right

    // Set up clipping for text area
    clip_rect := core.Rect{
        x = abs_rect.x + ti.padding.left,
        y = abs_rect.y,
        width = content_width,
        height = abs_rect.height,
    }

    old_clip := ctx.logical_clip
    if clipped, ok := core.rect_intersection(clip_rect, old_clip); ok {
        render.context_set_clip(ctx, clipped)

        if ti.font != nil {
            // Draw text or placeholder
            if len(ti.text) > 0 {
                text := string(ti.text[:])
                render.draw_text_top(ctx, ti.font, text, text_x, text_y, ti.text_color)
            } else if ti.placeholder != "" {
                render.draw_text_top(ctx, ti.font, ti.placeholder, text_x, text_y, ti.placeholder_color)
            }

            // Draw cursor when focused
            if w.focused && ti.cursor_visible {
                cursor_x := text_x + text_input_get_cursor_x(ti)
                cursor_rect := core.Rect{
                    x = cursor_x,
                    y = abs_rect.y + ti.padding.top,
                    width = 2,
                    height = render.font_get_logical_line_height(ti.font),
                }
                render.fill_rect(ctx, cursor_rect, ti.cursor_color)
            }
        }

        render.context_set_clip(ctx, old_clip)
    }
}

// Handle events
text_input_handle_event :: proc(w: ^Widget, event: ^core.Event) -> bool {
    ti := cast(^Text_Input)w

    if !w.enabled {
        return false
    }

    #partial switch event.type {
    case .Text_Input:
        // Handle text input from keyboard
        // The utf8 text comes from the event's text_input field
        // For now, we'll handle this through Key_Press with character data
        return false

    case .Key_Press:
        if !w.focused {
            return false
        }

        // Handle special keys
        keycode := core.Keycode(event.keycode)

        #partial switch keycode {
        case .Backspace:
            text_input_backspace(ti)
            return true

        case .Delete:
            text_input_delete(ti)
            return true

        case .Left:
            text_input_cursor_left(ti)
            return true

        case .Right:
            text_input_cursor_right(ti)
            return true

        case .Home:
            text_input_cursor_home(ti)
            return true

        case .End:
            text_input_cursor_end(ti)
            return true

        case .Enter:
            if ti.on_submit != nil {
                ti.on_submit(ti)
            }
            return true
        }

        // Check for printable character via keysym
        // XKB keysyms for printable ASCII are the same as ASCII codes
        if event.keysym >= 0x20 && event.keysym <= 0x7E {
            // Simple ASCII character
            char_buf: [1]u8
            char_buf[0] = u8(event.keysym)
            text_input_insert(ti, string(char_buf[:]))
            return true
        }

        return false

    case .Pointer_Button_Press:
        if event.button == .Left {
            // Click to position cursor
            abs_rect := widget_get_absolute_rect(w)
            click_x := event.pointer_x - abs_rect.x - ti.padding.left + ti.scroll_offset

            if ti.font != nil && len(ti.text) > 0 {
                // Find character position at click location
                ti.cursor_pos = 0
                text_width: i32 = 0

                for i := 0; i < len(ti.text); {
                    // Get next character
                    _, char_len := utf8.decode_rune(ti.text[i:])
                    char_text := string(ti.text[i:i+char_len])
                    char_width := render.text_measure_logical(ti.font, char_text)

                    // Check if click is before midpoint of this character
                    if click_x < text_width + char_width / 2 {
                        break
                    }

                    text_width += char_width
                    ti.cursor_pos = i + char_len
                    i += char_len
                }

                text_input_ensure_cursor_visible(ti)
                widget_mark_dirty(ti)
            }
            return true
        }

    case .Window_Focus:
        ti.cursor_visible = true
        widget_mark_dirty(ti)

    case .Window_Unfocus:
        widget_mark_dirty(ti)
    }

    return false
}

// Measure preferred size
text_input_measure :: proc(w: ^Widget, available_width: i32) -> core.Size {
    ti := cast(^Text_Input)w

    height: i32 = 32  // Default height
    if ti.font != nil {
        height = render.font_get_logical_line_height(ti.font) + ti.padding.top + ti.padding.bottom
    }

    width := ti.min_size.width
    if width == 0 {
        width = 200  // Default width
    }

    return core.Size{
        width = width,
        height = max(height, ti.min_size.height),
    }
}

// Destroy text input
text_input_destroy :: proc(w: ^Widget) {
    ti := cast(^Text_Input)w
    delete(ti.text)
}
