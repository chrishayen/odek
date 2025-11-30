package widgets

import "../core"
import "../render"

// Scroll direction
Scroll_Direction :: enum {
    Vertical,
    Horizontal,
    Both,
}

// ScrollContainer wraps a content widget and provides scrolling
Scroll_Container :: struct {
    using base: Widget,

    scroll:     Scroll_State,
    direction:  Scroll_Direction,

    // Scrollbar appearance
    scrollbar_width:    i32,
    scrollbar_padding:  i32,  // padding from edge

    // Scrollbar state
    scrollbar_dragging: bool,
    drag_start_y:       i32,
    drag_start_offset:  i32,

    // Colors (use theme if not set)
    track_color:        core.Color,
    thumb_color:        core.Color,
    thumb_hover_color:  core.Color,

    // State
    thumb_hovered:      bool,
}

// Shared vtable for scroll containers
scroll_container_vtable := Widget_VTable{
    draw         = scroll_container_draw,
    handle_event = scroll_container_handle_event,
    layout       = scroll_container_layout,
    destroy      = scroll_container_destroy,
    measure      = scroll_container_measure,
}

// Create a new scroll container
scroll_container_create :: proc(direction: Scroll_Direction = .Vertical) -> ^Scroll_Container {
    sc := new(Scroll_Container)
    sc.vtable = &scroll_container_vtable
    sc.visible = true
    sc.enabled = true
    sc.dirty = true

    sc.scroll = scroll_init()
    sc.direction = direction
    sc.scrollbar_width = 8
    sc.scrollbar_padding = 2

    // Default colors from theme
    theme := theme_get()
    sc.track_color = theme.scrollbar_track
    sc.thumb_color = theme.scrollbar_thumb
    sc.thumb_hover_color = theme.scrollbar_hover

    return sc
}

// Set content widget (the widget to scroll)
scroll_container_set_content :: proc(sc: ^Scroll_Container, content: ^Widget) {
    // Remove existing children
    for len(sc.children) > 0 {
        widget_remove_child(sc, sc.children[0])
    }

    // Add new content
    if content != nil {
        widget_add_child(sc, content)
    }

    widget_mark_dirty(sc)
}

// Get the content widget
scroll_container_get_content :: proc(sc: ^Scroll_Container) -> ^Widget {
    if len(sc.children) > 0 {
        return sc.children[0]
    }
    return nil
}

// Scroll to top
scroll_container_scroll_to_top :: proc(sc: ^Scroll_Container) {
    scroll_set_offset(&sc.scroll, 0)
    widget_mark_dirty(sc)
}

// Scroll to bottom
scroll_container_scroll_to_bottom :: proc(sc: ^Scroll_Container) {
    max_offset := scroll_get_max_offset(&sc.scroll)
    scroll_set_offset(&sc.scroll, max_offset)
    widget_mark_dirty(sc)
}

// Get viewport rect (area where content is visible)
scroll_container_get_viewport :: proc(sc: ^Scroll_Container) -> core.Rect {
    width := sc.rect.width - sc.padding.left - sc.padding.right
    height := sc.rect.height - sc.padding.top - sc.padding.bottom

    // Account for scrollbar if scrollable
    if scroll_is_scrollable(&sc.scroll) {
        if sc.direction == .Vertical || sc.direction == .Both {
            width -= sc.scrollbar_width + sc.scrollbar_padding * 2
        }
    }

    return core.Rect{
        x = sc.padding.left,
        y = sc.padding.top,
        width = width,
        height = height,
    }
}

// Get scrollbar track rect
scroll_container_get_track_rect :: proc(sc: ^Scroll_Container) -> core.Rect {
    viewport := scroll_container_get_viewport(sc)
    return core.Rect{
        x = sc.rect.width - sc.scrollbar_width - sc.scrollbar_padding,
        y = sc.padding.top,
        width = sc.scrollbar_width,
        height = viewport.height,
    }
}

// Get scrollbar thumb rect
scroll_container_get_thumb_rect :: proc(sc: ^Scroll_Container) -> core.Rect {
    track := scroll_container_get_track_rect(sc)
    max_offset := scroll_get_max_offset(&sc.scroll)

    if max_offset <= 0 {
        return track
    }

    // Calculate thumb size proportional to visible content
    visible_fraction := f32(sc.scroll.viewport_size) / f32(sc.scroll.content_size)
    thumb_height := max(i32(f32(track.height) * visible_fraction), 20)

    // Calculate thumb position
    scroll_fraction := scroll_get_fraction(&sc.scroll)
    available_space := track.height - thumb_height
    thumb_y := track.y + i32(f32(available_space) * scroll_fraction)

    return core.Rect{
        x = track.x,
        y = thumb_y,
        width = track.width,
        height = thumb_height,
    }
}

// Check if point is in scrollbar area
scroll_container_point_in_scrollbar :: proc(sc: ^Scroll_Container, x, y: i32) -> bool {
    track := scroll_container_get_track_rect(sc)
    return x >= track.x && x < track.x + track.width &&
           y >= track.y && y < track.y + track.height
}

// Check if point is in thumb
scroll_container_point_in_thumb :: proc(sc: ^Scroll_Container, x, y: i32) -> bool {
    thumb := scroll_container_get_thumb_rect(sc)
    return x >= thumb.x && x < thumb.x + thumb.width &&
           y >= thumb.y && y < thumb.y + thumb.height
}

// Layout implementation
scroll_container_layout :: proc(w: ^Widget) {
    sc := cast(^Scroll_Container)w

    if len(sc.children) == 0 {
        return
    }

    content := sc.children[0]
    viewport := scroll_container_get_viewport(sc)

    // Measure content to get its natural size
    content_size := widget_measure(content)

    // Position content at top-left of viewport, offset by scroll
    content.rect = core.Rect{
        x = viewport.x,
        y = viewport.y - sc.scroll.offset,
        width = viewport.width,
        height = max(content_size.height, viewport.height),
    }

    // Update scroll state
    scroll_set_sizes(&sc.scroll, content_size.height, viewport.height)
}

// Draw implementation
scroll_container_draw :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
    sc := cast(^Scroll_Container)w
    abs_rect := widget_get_absolute_rect(w)

    // Draw background if set (optional - containers often transparent)
    theme := theme_get()

    // Set up clipping for viewport
    viewport := scroll_container_get_viewport(sc)
    viewport_abs := core.Rect{
        x = abs_rect.x + viewport.x,
        y = abs_rect.y + viewport.y,
        width = viewport.width,
        height = viewport.height,
    }

    // Draw content with clipping
    old_clip := ctx.logical_clip
    if clipped, ok := core.rect_intersection(viewport_abs, old_clip); ok {
        render.context_set_clip(ctx, clipped)

        // Draw children (content)
        for child in sc.children {
            widget_draw(child, ctx)
        }

        render.context_set_clip(ctx, old_clip)
    }

    // Draw scrollbar if scrollable
    if scroll_is_scrollable(&sc.scroll) {
        track := scroll_container_get_track_rect(sc)
        thumb := scroll_container_get_thumb_rect(sc)

        // Transform to absolute coordinates
        track_abs := core.Rect{
            x = abs_rect.x + track.x,
            y = abs_rect.y + track.y,
            width = track.width,
            height = track.height,
        }
        thumb_abs := core.Rect{
            x = abs_rect.x + thumb.x,
            y = abs_rect.y + thumb.y,
            width = thumb.width,
            height = thumb.height,
        }

        // Draw track
        render.fill_rounded_rect(ctx, track_abs, 4, sc.track_color)

        // Draw thumb
        thumb_color := sc.thumb_hover_color if (sc.thumb_hovered || sc.scrollbar_dragging) else sc.thumb_color
        render.fill_rounded_rect(ctx, thumb_abs, 4, thumb_color)
    }
}

// Handle events
scroll_container_handle_event :: proc(w: ^Widget, event: ^core.Event) -> bool {
    sc := cast(^Scroll_Container)w
    abs_rect := widget_get_absolute_rect(w)

    #partial switch event.type {
    case .Pointer_Motion:
        local_x := event.pointer_x - abs_rect.x
        local_y := event.pointer_y - abs_rect.y

        // Handle scrollbar dragging
        if sc.scrollbar_dragging {
            track := scroll_container_get_track_rect(sc)
            thumb := scroll_container_get_thumb_rect(sc)
            max_offset := scroll_get_max_offset(&sc.scroll)

            if max_offset > 0 {
                drag_delta := local_y - sc.drag_start_y
                scroll_range := track.height - thumb.height
                if scroll_range > 0 {
                    offset_delta := (drag_delta * max_offset) / scroll_range
                    scroll_set_offset(&sc.scroll, sc.drag_start_offset + offset_delta)
                    widget_mark_dirty(sc)
                }
            }
            return true
        }

        // Update thumb hover state
        if scroll_is_scrollable(&sc.scroll) {
            was_hovered := sc.thumb_hovered
            sc.thumb_hovered = scroll_container_point_in_thumb(sc, local_x, local_y)
            if was_hovered != sc.thumb_hovered {
                widget_mark_dirty(sc)
            }
        }

    case .Pointer_Button_Press:
        if event.button == .Left {
            local_x := event.pointer_x - abs_rect.x
            local_y := event.pointer_y - abs_rect.y

            // Check if clicking scrollbar
            if scroll_is_scrollable(&sc.scroll) && scroll_container_point_in_scrollbar(sc, local_x, local_y) {
                thumb := scroll_container_get_thumb_rect(sc)

                // Check if clicking on thumb
                if scroll_container_point_in_thumb(sc, local_x, local_y) {
                    // Start dragging
                    sc.scrollbar_dragging = true
                    sc.drag_start_y = local_y
                    sc.drag_start_offset = sc.scroll.offset
                    widget_mark_dirty(sc)
                    return true
                }

                // Clicked on track - page up/down
                track := scroll_container_get_track_rect(sc)
                if local_y < thumb.y {
                    // Page up
                    scroll_by(&sc.scroll, -sc.scroll.viewport_size)
                } else {
                    // Page down
                    scroll_by(&sc.scroll, sc.scroll.viewport_size)
                }
                widget_mark_dirty(sc)
                return true
            }
        }

    case .Pointer_Button_Release:
        if event.button == .Left && sc.scrollbar_dragging {
            sc.scrollbar_dragging = false
            widget_mark_dirty(sc)
            return true
        }

    case .Pointer_Leave:
        sc.thumb_hovered = false
        if sc.scrollbar_dragging {
            sc.scrollbar_dragging = false
        }
        widget_mark_dirty(sc)

    case .Scroll:
        // Handle mouse wheel scroll
        if scroll_is_scrollable(&sc.scroll) {
            scroll_by_steps(&sc.scroll, event.scroll_delta)
            widget_mark_dirty(sc)
            return true
        }
    }

    // Pass events to content if not handled by scrollbar
    if len(sc.children) > 0 {
        return widget_handle_event(sc.children[0], event)
    }

    return false
}

// Measure preferred size
scroll_container_measure :: proc(w: ^Widget) -> core.Size {
    sc := cast(^Scroll_Container)w

    // Return min_size if set, otherwise measure content
    if sc.min_size.width > 0 && sc.min_size.height > 0 {
        return sc.min_size
    }

    if len(sc.children) > 0 {
        content_size := widget_measure(sc.children[0])
        return core.Size{
            width = max(content_size.width + sc.scrollbar_width + sc.scrollbar_padding * 2, sc.min_size.width),
            height = max(content_size.height, sc.min_size.height),
        }
    }

    return sc.min_size
}

// Destroy scroll container
scroll_container_destroy :: proc(w: ^Widget) {
    // Children are destroyed by widget_destroy in widget.odin
}
