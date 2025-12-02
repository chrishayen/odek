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

    // Selection (anchor is where selection started, cursor_pos is other end)
    selection_anchor: int,  // -1 means no selection
    mouse_selecting:  bool, // True while mouse drag selecting

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
    selection_color: core.Color,

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
    ti.selection_anchor = -1
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
    ti.selection_color = theme.accent

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

// Check if there's an active selection
text_input_has_selection :: proc(ti: ^Text_Input) -> bool {
    return ti.selection_anchor >= 0 && ti.selection_anchor != ti.cursor_pos
}

// Get selection bounds (start, end) in byte positions
text_input_get_selection :: proc(ti: ^Text_Input) -> (start: int, end: int) {
    if ti.selection_anchor < 0 {
        return ti.cursor_pos, ti.cursor_pos
    }
    if ti.selection_anchor < ti.cursor_pos {
        return ti.selection_anchor, ti.cursor_pos
    }
    return ti.cursor_pos, ti.selection_anchor
}

// Clear selection
text_input_clear_selection :: proc(ti: ^Text_Input) {
    ti.selection_anchor = -1
}

// Start selection at current cursor position
text_input_start_selection :: proc(ti: ^Text_Input) {
    if ti.selection_anchor < 0 {
        ti.selection_anchor = ti.cursor_pos
    }
}

// Select all text
text_input_select_all :: proc(ti: ^Text_Input) {
    ti.selection_anchor = 0
    ti.cursor_pos = len(ti.text)
    text_input_ensure_cursor_visible(ti)
    widget_mark_dirty(ti)
}

// Global clipboard callbacks (set by app package)
Clipboard_Copy_Proc :: #type proc(text: string)
Clipboard_Paste_Proc :: #type proc() -> string

@(private)
g_clipboard_copy: Clipboard_Copy_Proc
@(private)
g_clipboard_paste: Clipboard_Paste_Proc

// Set clipboard callbacks (called by app package)
clipboard_set_handlers :: proc(copy_fn: Clipboard_Copy_Proc, paste_fn: Clipboard_Paste_Proc) {
    g_clipboard_copy = copy_fn
    g_clipboard_paste = paste_fn
}

// Copy selected text to clipboard
text_input_copy :: proc(ti: ^Text_Input) {
    if !text_input_has_selection(ti) {
        return
    }
    start, end := text_input_get_selection(ti)
    text := string(ti.text[start:end])
    if g_clipboard_copy != nil {
        g_clipboard_copy(text)
    }
}

// Cut selected text to clipboard
text_input_cut :: proc(ti: ^Text_Input) {
    if !text_input_has_selection(ti) {
        return
    }
    text_input_copy(ti)
    text_input_delete_selection(ti)
}

// Paste from clipboard
text_input_paste :: proc(ti: ^Text_Input) {
    if g_clipboard_paste == nil {
        return
    }
    text := g_clipboard_paste()
    if len(text) > 0 {
        text_input_insert(ti, text)
        delete(text)  // clipboard_paste allocates
    }
}

// Delete selected text
text_input_delete_selection :: proc(ti: ^Text_Input) {
    if !text_input_has_selection(ti) {
        return
    }

    start, end := text_input_get_selection(ti)
    ordered_remove_range(&ti.text, start, end)
    ti.cursor_pos = start
    ti.selection_anchor = -1
    text_input_ensure_cursor_visible(ti)
    widget_mark_dirty(ti)

    if ti.on_change != nil {
        ti.on_change(ti)
    }
}

// Get X position for a byte offset
text_input_get_x_for_pos :: proc(ti: ^Text_Input, pos: int) -> i32 {
    if ti.font == nil || pos == 0 {
        return 0
    }
    pos := min(pos, len(ti.text))
    text_before := string(ti.text[:pos])
    return render.text_measure_logical(ti.font, text_before)
}

// Get byte position from X coordinate (for click handling)
text_input_pos_from_x :: proc(ti: ^Text_Input, click_x: i32) -> int {
    if ti.font == nil || len(ti.text) == 0 {
        return 0
    }

    pos := 0
    text_width: i32 = 0

    for i := 0; i < len(ti.text); {
        _, char_len := utf8.decode_rune(ti.text[i:])
        char_text := string(ti.text[i:i+char_len])
        char_width := render.text_measure_logical(ti.font, char_text)

        // Check if click is before midpoint of this character
        if click_x < text_width + char_width / 2 {
            break
        }

        text_width += char_width
        pos = i + char_len
        i += char_len
    }

    return pos
}

// Move cursor left by one character
text_input_cursor_left :: proc(ti: ^Text_Input, extend_selection: bool = false) {
    if extend_selection {
        text_input_start_selection(ti)
    } else if text_input_has_selection(ti) {
        // Move cursor to start of selection
        start, _ := text_input_get_selection(ti)
        ti.cursor_pos = start
        text_input_clear_selection(ti)
        text_input_ensure_cursor_visible(ti)
        widget_mark_dirty(ti)
        return
    } else {
        text_input_clear_selection(ti)
    }
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
text_input_cursor_right :: proc(ti: ^Text_Input, extend_selection: bool = false) {
    if extend_selection {
        text_input_start_selection(ti)
    } else if text_input_has_selection(ti) {
        // Move cursor to end of selection
        _, end := text_input_get_selection(ti)
        ti.cursor_pos = end
        text_input_clear_selection(ti)
        text_input_ensure_cursor_visible(ti)
        widget_mark_dirty(ti)
        return
    } else {
        text_input_clear_selection(ti)
    }

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
text_input_cursor_home :: proc(ti: ^Text_Input, extend_selection: bool = false) {
    if extend_selection {
        text_input_start_selection(ti)
    } else {
        text_input_clear_selection(ti)
    }
    ti.cursor_pos = 0
    text_input_ensure_cursor_visible(ti)
    widget_mark_dirty(ti)
}

// Move cursor to end
text_input_cursor_end :: proc(ti: ^Text_Input, extend_selection: bool = false) {
    if extend_selection {
        text_input_start_selection(ti)
    } else {
        text_input_clear_selection(ti)
    }
    ti.cursor_pos = len(ti.text)
    text_input_ensure_cursor_visible(ti)
    widget_mark_dirty(ti)
}

// Insert text at cursor position
text_input_insert :: proc(ti: ^Text_Input, text: string) {
    if len(text) == 0 {
        return
    }

    // Delete selection first if any
    if text_input_has_selection(ti) {
        text_input_delete_selection(ti)
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
    // Delete selection if any
    if text_input_has_selection(ti) {
        text_input_delete_selection(ti)
        return
    }

    if ti.cursor_pos == 0 {
        return
    }

    // Find start of previous UTF-8 character
    char_start := ti.cursor_pos - 1
    for char_start > 0 && (ti.text[char_start] & 0xC0) == 0x80 {
        char_start -= 1
    }

    // Remove bytes from char_start to cursor_pos
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
    // Delete selection if any
    if text_input_has_selection(ti) {
        text_input_delete_selection(ti)
        return
    }

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
            line_height := render.font_get_logical_line_height(ti.font)

            // Draw selection highlight
            if text_input_has_selection(ti) {
                sel_start, sel_end := text_input_get_selection(ti)
                start_x := text_x + text_input_get_x_for_pos(ti, sel_start)
                end_x := text_x + text_input_get_x_for_pos(ti, sel_end)
                sel_rect := core.Rect{
                    x = start_x,
                    y = abs_rect.y + ti.padding.top,
                    width = end_x - start_x,
                    height = line_height,
                }
                // Use semi-transparent selection color with blending
                sel_color := ti.selection_color
                sel_color.a = 100
                render.fill_rect_blend(ctx, sel_rect, sel_color)
            }

            // Draw text or placeholder
            if len(ti.text) > 0 {
                text := string(ti.text[:])
                render.draw_text_top(ctx, ti.font, text, text_x, text_y, ti.text_color)
            } else if ti.placeholder != "" {
                render.draw_text_top(ctx, ti.font, ti.placeholder, text_x, text_y, ti.placeholder_color)
            }

            // Draw cursor when focused (not when there's a selection)
            if w.focused && ti.cursor_visible && !text_input_has_selection(ti) {
                cursor_x := text_x + text_input_get_cursor_x(ti)
                cursor_rect := core.Rect{
                    x = cursor_x,
                    y = abs_rect.y + ti.padding.top,
                    width = 2,
                    height = line_height,
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

        shift := .Shift in event.modifiers
        ctrl := .Ctrl in event.modifiers

        #partial switch keycode {
        case .Backspace:
            text_input_backspace(ti)
            return true

        case .Delete:
            text_input_delete(ti)
            return true

        case .Left:
            text_input_cursor_left(ti, shift)
            return true

        case .Right:
            text_input_cursor_right(ti, shift)
            return true

        case .Home:
            text_input_cursor_home(ti, shift)
            return true

        case .End:
            text_input_cursor_end(ti, shift)
            return true

        case .Enter:
            if ti.on_submit != nil {
                ti.on_submit(ti)
            }
            return true

        case .A:
            // Ctrl+A to select all
            if ctrl {
                text_input_select_all(ti)
                return true
            }

        case .C:
            // Ctrl+C to copy
            if ctrl {
                text_input_copy(ti)
                return true
            }

        case .X:
            // Ctrl+X to cut
            if ctrl {
                text_input_cut(ti)
                return true
            }

        case .V:
            // Ctrl+V to paste
            if ctrl {
                text_input_paste(ti)
                return true
            }

        case .Insert:
            // CUA standard: Ctrl+Insert to copy, Shift+Insert to paste
            if ctrl {
                text_input_copy(ti)
                return true
            }
            if shift {
                text_input_paste(ti)
                return true
            }
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
            new_pos := text_input_pos_from_x(ti, event.pointer_x - abs_rect.x - ti.padding.left + ti.scroll_offset)

            shift := .Shift in event.modifiers
            if shift {
                // Shift+click extends selection
                text_input_start_selection(ti)
            } else {
                // Regular click clears selection
                text_input_clear_selection(ti)
            }

            ti.cursor_pos = new_pos
            ti.mouse_selecting = true

            // Capture pointer for drag selection
            if g_hit_state != nil {
                pointer_capture(g_hit_state, w)
            }

            text_input_ensure_cursor_visible(ti)
            widget_mark_dirty(ti)
            return true
        }

    case .Pointer_Button_Release:
        if event.button == .Left && ti.mouse_selecting {
            ti.mouse_selecting = false
            if g_hit_state != nil {
                pointer_release(g_hit_state)
            }
            return true
        }

    case .Pointer_Motion:
        if ti.mouse_selecting {
            abs_rect := widget_get_absolute_rect(w)
            new_pos := text_input_pos_from_x(ti, event.pointer_x - abs_rect.x - ti.padding.left + ti.scroll_offset)

            // Start selection if not already started
            if ti.selection_anchor < 0 {
                ti.selection_anchor = ti.cursor_pos
            }

            ti.cursor_pos = new_pos
            text_input_ensure_cursor_visible(ti)
            widget_mark_dirty(ti)
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
