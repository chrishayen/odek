package widgets

import "../core"
import "../render"
import "core:strings"

// Label widget for displaying text with optional word wrapping
Label :: struct {
    using base: Widget,

    text:       string,
    font:       ^render.Font,
    color:      core.Color,
    h_align:    Align,          // Horizontal text alignment (Start, Center, End)
    wrap:       bool,           // Enable word wrapping

    // Cached layout data
    lines:        [dynamic]string,  // Wrapped lines (slices into original text)
    cached_width: i32,              // Width used for wrapping calculation
}

// Shared vtable for all labels
label_vtable := Widget_VTable{
    draw         = label_draw,
    handle_event = label_handle_event,
    layout       = label_layout,
    destroy      = label_destroy,
    measure      = label_measure,
}

// Create a new label
label_create :: proc(text: string = "", font: ^render.Font = nil) -> ^Label {
    l := new(Label)
    l.vtable = &label_vtable
    l.visible = true
    l.enabled = true
    l.dirty = true
    l.text = text
    l.font = font
    l.color = theme_get().text_primary
    l.h_align = .Start
    l.wrap = true
    l.cached_width = -1  // Force recalculation
    return l
}

// Set label text
label_set_text :: proc(l: ^Label, text: string) {
    if l.text == text {
        return
    }
    l.text = text
    l.cached_width = -1  // Invalidate cache
    widget_mark_dirty(l)
}

// Set label font
label_set_font :: proc(l: ^Label, font: ^render.Font) {
    if l.font == font {
        return
    }
    l.font = font
    l.cached_width = -1  // Invalidate cache
    widget_mark_dirty(l)
}

// Set label color
label_set_color :: proc(l: ^Label, color: core.Color) {
    l.color = color
    widget_mark_dirty(l)
}

// Set horizontal alignment
label_set_align :: proc(l: ^Label, align: Align) {
    l.h_align = align
    widget_mark_dirty(l)
}

// Set word wrap enabled
label_set_wrap :: proc(l: ^Label, wrap: bool) {
    if l.wrap == wrap {
        return
    }
    l.wrap = wrap
    l.cached_width = -1  // Invalidate cache
    widget_mark_dirty(l)
}

// Update wrapped lines cache
label_update_lines :: proc(l: ^Label, available_width: i32) {
    // Only skip if cache is valid AND we have lines
    if l.cached_width == available_width && len(l.lines) > 0 {
        return  // Cache is valid
    }

    clear(&l.lines)
    l.cached_width = available_width

    if l.font == nil || l.text == "" {
        return
    }

    if !l.wrap || available_width <= 0 {
        // No wrapping - single line
        append(&l.lines, l.text)
        return
    }

    // Word wrap algorithm
    label_wrap_text(l.font, l.text, available_width, &l.lines)
}

// Wrap text to fit within max_width, breaking at word boundaries
label_wrap_text :: proc(font: ^render.Font, text: string, max_width: i32, lines: ^[dynamic]string) {
    // Process text line by line (respecting existing newlines)
    remaining := text

    for len(remaining) > 0 {
        // Find next newline
        newline_idx := strings.index_byte(remaining, '\n')

        line: string
        if newline_idx >= 0 {
            line = remaining[:newline_idx]
            remaining = remaining[newline_idx + 1:]
        } else {
            line = remaining
            remaining = ""
        }

        // Wrap this line
        wrap_single_line(font, line, max_width, lines)
    }
}

// Wrap a single line (no embedded newlines) to fit within max_width
wrap_single_line :: proc(font: ^render.Font, line: string, max_width: i32, lines: ^[dynamic]string) {
    if line == "" {
        append(lines, "")
        return
    }

    // Check if line fits without wrapping
    if render.text_measure(font, line) <= max_width {
        append(lines, line)
        return
    }

    // Need to wrap - iterate word by word
    line_start := 0
    word_start := 0
    last_break := 0  // Last valid break point

    i := 0
    for i <= len(line) {
        // Find end of current word
        is_end := i >= len(line)
        is_space := !is_end && line[i] == ' '

        if is_space || is_end {
            // We have a complete word from word_start to i
            current_segment := line[line_start:i]
            width := render.text_measure(font, current_segment)

            if width > max_width {
                // Line is too long
                if last_break > line_start {
                    // Break at last valid position
                    append(lines, line[line_start:last_break])
                    // Skip space after break
                    line_start = last_break
                    for line_start < len(line) && line[line_start] == ' ' {
                        line_start += 1
                    }
                    last_break = line_start
                    word_start = line_start
                    // Re-check from new position
                    i = line_start
                    continue
                } else {
                    // Single word is too long - just add it (overflow)
                    append(lines, line[line_start:i])
                    line_start = i
                    for line_start < len(line) && line[line_start] == ' ' {
                        line_start += 1
                    }
                    last_break = line_start
                    word_start = line_start
                    i = line_start
                    continue
                }
            }

            if is_space {
                last_break = i
                // Skip consecutive spaces
                for i < len(line) && line[i] == ' ' {
                    i += 1
                }
                word_start = i
                continue
            }
        }

        i += 1
    }

    // Add remaining text
    if line_start < len(line) {
        append(lines, line[line_start:])
    }
}

// Draw label
label_draw :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
    l := cast(^Label)w

    if l.font == nil || l.text == "" {
        return
    }

    abs_rect := widget_get_absolute_rect(w)
    content_width := abs_rect.width - w.padding.left - w.padding.right

    // Ensure lines are calculated
    label_update_lines(l, content_width)

    if len(l.lines) == 0 {
        return
    }

    // Draw each line
    y := abs_rect.y + w.padding.top

    for line in l.lines {
        line_width := render.text_measure(l.font, line)

        // Calculate x position based on alignment
        x: i32
        switch l.h_align {
        case .Start:
            x = abs_rect.x + w.padding.left
        case .Center:
            x = abs_rect.x + w.padding.left + (content_width - line_width) / 2
        case .End:
            x = abs_rect.x + abs_rect.width - w.padding.right - line_width
        case .Stretch:
            x = abs_rect.x + w.padding.left  // Same as Start for text
        }

        render.draw_text_top(ctx, l.font, line, x, y, l.color)
        y += l.font.line_height
    }
}

// Label doesn't handle events
label_handle_event :: proc(w: ^Widget, event: ^core.Event) -> bool {
    return false
}

// Layout - recalculate lines if width changed
label_layout :: proc(w: ^Widget) {
    l := cast(^Label)w
    content_width := w.rect.width - w.padding.left - w.padding.right
    label_update_lines(l, content_width)
}

// Measure label's preferred size
label_measure :: proc(w: ^Widget) -> core.Size {
    l := cast(^Label)w

    if l.font == nil || l.text == "" {
        return core.Size{
            width = w.padding.left + w.padding.right,
            height = w.padding.top + w.padding.bottom,
        }
    }

    // Try to get width constraint from rect, or from parent's content area
    available_width := w.rect.width - w.padding.left - w.padding.right
    if available_width <= 0 && l.wrap && w.parent != nil {
        // Use parent's content width as hint for wrapping calculation
        parent_content := widget_get_content_rect(w.parent)
        available_width = parent_content.width - w.margin.left - w.margin.right - w.padding.left - w.padding.right
    }

    if available_width <= 0 || !l.wrap {
        // No width constraint, measure single line
        text_size := render.text_measure_size(l.font, l.text)
        return core.Size{
            width = text_size.width + w.padding.left + w.padding.right,
            height = text_size.height + w.padding.top + w.padding.bottom,
        }
    }

    // Wrap and calculate height based on available width
    label_update_lines(l, available_width)
    height := i32(len(l.lines)) * l.font.line_height

    // For wrapping labels, prefer parent's width
    preferred_width := w.rect.width
    if preferred_width <= 0 && w.parent != nil {
        parent_content := widget_get_content_rect(w.parent)
        preferred_width = parent_content.width - w.margin.left - w.margin.right
    }

    return core.Size{
        width = preferred_width,
        height = max(height + w.padding.top + w.padding.bottom, w.min_size.height),
    }
}

// Destroy label resources
label_destroy :: proc(w: ^Widget) {
    l := cast(^Label)w
    delete(l.lines)
}
