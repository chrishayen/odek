package widgets

// Reusable scroll state for scrollable widgets
Scroll_State :: struct {
    offset:        i32,   // Current scroll offset (pixels from top)
    content_size:  i32,   // Total height of scrollable content
    viewport_size: i32,   // Visible area height
    scroll_speed:  i32,   // Pixels per scroll step
}

// Initialize scroll state with defaults
scroll_init :: proc(scroll_speed: i32 = 40) -> Scroll_State {
    return Scroll_State{
        offset = 0,
        content_size = 0,
        viewport_size = 0,
        scroll_speed = scroll_speed,
    }
}

// Set scroll offset, clamping to valid range
scroll_set_offset :: proc(s: ^Scroll_State, offset: i32) {
    max_offset := max(0, s.content_size - s.viewport_size)
    s.offset = clamp(offset, 0, max_offset)
}

// Scroll by delta pixels (negative = up, positive = down)
scroll_by :: proc(s: ^Scroll_State, delta: i32) {
    scroll_set_offset(s, s.offset + delta)
}

// Scroll by number of "steps" (e.g., mouse wheel ticks)
scroll_by_steps :: proc(s: ^Scroll_State, steps: i32) {
    scroll_by(s, steps * s.scroll_speed)
}

// Update content and viewport sizes
scroll_set_sizes :: proc(s: ^Scroll_State, content_size, viewport_size: i32) {
    s.content_size = content_size
    s.viewport_size = viewport_size
    // Re-clamp offset in case content shrunk
    scroll_set_offset(s, s.offset)
}

// Check if scrolling is needed
scroll_is_scrollable :: proc(s: ^Scroll_State) -> bool {
    return s.content_size > s.viewport_size
}

// Get scroll percentage (0.0 to 1.0)
scroll_get_fraction :: proc(s: ^Scroll_State) -> f32 {
    max_offset := s.content_size - s.viewport_size
    if max_offset <= 0 {
        return 0.0
    }
    return f32(s.offset) / f32(max_offset)
}

// Get maximum scroll offset
scroll_get_max_offset :: proc(s: ^Scroll_State) -> i32 {
    return max(0, s.content_size - s.viewport_size)
}

// Scroll to ensure a range is visible
scroll_ensure_visible :: proc(s: ^Scroll_State, item_start, item_height: i32) {
    item_end := item_start + item_height

    // If item is above viewport, scroll up
    if item_start < s.offset {
        scroll_set_offset(s, item_start)
        return
    }

    // If item is below viewport, scroll down
    if item_end > s.offset + s.viewport_size {
        scroll_set_offset(s, item_end - s.viewport_size)
    }
}
