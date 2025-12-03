package widgets

import "../core"
import "../render"
import "core:strings"

// Label widget for displaying text with optional word wrapping
Label :: struct {
    using base: Widget,

    text:         string,
    font:         ^render.Font,
    font_bold:    ^render.Font,    // Bold font variant (optional)
    bold:         bool,            // Use bold font
    color:        core.Color,
    h_align:      Align,          // Horizontal text alignment (Start, Center, End)
    wrap:         bool,           // Enable word wrapping
    strikethrough: bool,          // Draw line through text

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

// Set strikethrough enabled
label_set_strikethrough :: proc(l: ^Label, strikethrough: bool) {
    if l.strikethrough == strikethrough {
        return
    }
    l.strikethrough = strikethrough
    widget_mark_dirty(l)
}

// Set bold enabled
label_set_bold :: proc(l: ^Label, bold: bool) {
    if l.bold == bold {
        return
    }
    l.bold = bold
    l.cached_width = -1  // Invalidate cache (bold may have different metrics)
    widget_mark_dirty(l)
}

// Set bold font
label_set_font_bold :: proc(l: ^Label, font: ^render.Font) {
    if l.font_bold == font {
        return
    }
    l.font_bold = font
    if l.bold {
        l.cached_width = -1  // Invalidate cache if currently using bold
    }
    widget_mark_dirty(l)
}

// Get the active font (bold or regular)
label_get_active_font :: proc(l: ^Label) -> ^render.Font {
    if l.bold && l.font_bold != nil {
        return l.font_bold
    }
    return l.font
}

// Update wrapped lines cache
label_update_lines :: proc(l: ^Label, available_width: i32) {
    // Only skip if cache is valid AND we have lines
    if l.cached_width == available_width && len(l.lines) > 0 {
        return  // Cache is valid
    }

    clear(&l.lines)
    l.cached_width = available_width

    font := label_get_active_font(l)
    if font == nil || l.text == "" {
        return
    }

    if !l.wrap || available_width <= 0 {
        // No wrapping - single line
        append(&l.lines, l.text)
        return
    }

    // Word wrap algorithm
    label_wrap_text(font, l.text, available_width, &l.lines)
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

// Wrap a single line (no embedded newlines) to fit within max_width (logical pixels)
wrap_single_line :: proc(font: ^render.Font, line: string, max_width: i32, lines: ^[dynamic]string) {
    if line == "" {
        append(lines, "")
        return
    }

    // Check if line fits without wrapping (use logical pixels)
    if render.text_measure_logical(font, line) <= max_width {
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
            width := render.text_measure_logical(font, current_segment)

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

    font := label_get_active_font(l)
    if font == nil || l.text == "" {
        return
    }

    abs_rect := widget_get_absolute_rect(w)
    content_width := abs_rect.width - w.padding.left - w.padding.right

    // Ensure lines are calculated
    label_update_lines(l, content_width)

    if len(l.lines) == 0 {
        return
    }

    // Draw each line (all coordinates in logical pixels)
    line_height := render.font_get_logical_line_height(font)
    y := abs_rect.y + w.padding.top

    for line in l.lines {
        line_width := render.text_measure_logical(font, line)

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

        render.draw_text_top(ctx, font, line, x, y, l.color)

        // Draw strikethrough line if enabled
        if l.strikethrough && line_width > 0 {
            // Draw line through vertical center of text
            strike_y := y + line_height / 2
            render.draw_hline(ctx, x, x + line_width, strike_y, l.color)
        }

        y += line_height
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

// Measure label's preferred size given available width (-1 = unconstrained)
label_measure :: proc(w: ^Widget, available_width: i32) -> core.Size {
    l := cast(^Label)w

    font := label_get_active_font(l)
    if font == nil || l.text == "" {
        return core.Size{
            width = w.padding.left + w.padding.right,
            height = w.padding.top + w.padding.bottom,
        }
    }

    // Calculate content width for wrapping
    content_width := available_width - w.padding.left - w.padding.right if available_width > 0 else -1

    if content_width <= 0 || !l.wrap {
        // No width constraint or wrapping disabled - measure single line
        text_size := render.text_measure_size(font, l.text)
        return core.Size{
            width = text_size.width + w.padding.left + w.padding.right,
            height = text_size.height + w.padding.top + w.padding.bottom,
        }
    }

    // Wrap and calculate height based on available width
    label_update_lines(l, content_width)
    line_height := render.font_get_logical_line_height(font)
    num_lines := max(1, i32(len(l.lines)))  // At least 1 line height
    height := num_lines * line_height

    // Calculate width: use actual text width for proper centering
    max_line_width: i32 = 0
    for line in l.lines {
        line_width := render.text_measure_logical(font, line)
        max_line_width = max(max_line_width, line_width)
    }
    // If no lines yet, use full text width
    if max_line_width == 0 {
        max_line_width = render.text_measure_logical(font, l.text)
    }

    return core.Size{
        width = max_line_width + w.padding.left + w.padding.right,
        height = max(height + w.padding.top + w.padding.bottom, w.min_size.height),
    }
}

// Destroy label resources
label_destroy :: proc(w: ^Widget) {
    l := cast(^Label)w
    delete(l.lines)
}
