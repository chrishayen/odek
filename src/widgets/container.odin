package widgets

import "../core"
import "../render"

// Container widget with flexbox-lite layout
Container :: struct {
    using base: Widget,

    direction:   Direction,   // Row or Column
    align_items: Align,       // Cross-axis alignment
    spacing:     i32,         // Gap between children
    background:  core.Color,  // Background color (transparent = no fill)
}

// Shared vtable for all containers
container_vtable := Widget_VTable{
    draw         = container_draw,
    handle_event = container_handle_event,
    layout       = container_layout,
    destroy      = container_destroy,
    measure      = container_measure,
}

// Create a new container
container_create :: proc(direction: Direction = .Column) -> ^Container {
    c := new(Container)
    c.vtable = &container_vtable
    c.visible = true
    c.enabled = true
    c.dirty = true
    c.direction = direction
    c.align_items = .Start
    c.spacing = 0
    c.background = core.COLOR_TRANSPARENT
    return c
}

// Draw container background and children
container_draw :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
    c := cast(^Container)w
    abs_rect := widget_get_absolute_rect(w)

    // Draw background if not transparent
    if c.background.a > 0 {
        render.fill_rect(ctx, abs_rect, c.background)
    }

    // Children are drawn by widget_draw after this returns
}

// Container doesn't handle events itself, just passes to children
container_handle_event :: proc(w: ^Widget, event: ^core.Event) -> bool {
    return false
}

// Flexbox-lite layout algorithm
container_layout :: proc(w: ^Widget) {
    c := cast(^Container)w

    visible_children := make([dynamic]^Widget, context.temp_allocator)
    for child in c.children {
        if child.visible {
            append(&visible_children, child)
        }
    }

    if len(visible_children) == 0 {
        return
    }

    // Calculate available space
    content := widget_get_content_rect(w)
    num_gaps := max(0, len(visible_children) - 1)
    total_spacing := c.spacing * i32(num_gaps)

    is_row := c.direction == .Row
    main_size := content.width if is_row else content.height
    cross_size := content.height if is_row else content.width
    available_main := main_size - total_spacing

    // First pass: measure all children and calculate flex totals
    child_sizes := make([]core.Size, len(visible_children), context.temp_allocator)
    total_flex: f32 = 0
    fixed_main: i32 = 0

    for i := 0; i < len(visible_children); i += 1 {
        child := visible_children[i]
        child_sizes[i] = widget_measure(child)

        if child.flex > 0 {
            total_flex += child.flex
        } else {
            // Fixed size child
            child_main := child_sizes[i].width if is_row else child_sizes[i].height
            fixed_main += child_main
        }
    }

    // Calculate space for flex children
    flex_space := available_main - fixed_main
    if flex_space < 0 {
        flex_space = 0
    }

    // Second pass: position children
    main_pos: i32 = 0

    for i := 0; i < len(visible_children); i += 1 {
        child := visible_children[i]
        child_size := child_sizes[i]

        // Calculate main axis size
        child_main: i32
        if child.flex > 0 && total_flex > 0 {
            // Flex child gets proportional share
            child_main = i32(f32(flex_space) * child.flex / total_flex)
        } else {
            child_main = child_size.width if is_row else child_size.height
        }

        // Apply min/max constraints
        if child.min_size.width > 0 && is_row {
            child_main = max(child_main, child.min_size.width)
        }
        if child.min_size.height > 0 && !is_row {
            child_main = max(child_main, child.min_size.height)
        }
        if child.max_size.width > 0 && is_row {
            child_main = min(child_main, child.max_size.width)
        }
        if child.max_size.height > 0 && !is_row {
            child_main = min(child_main, child.max_size.height)
        }

        // Calculate cross axis size and position
        child_cross := child_size.height if is_row else child_size.width
        cross_pos: i32 = 0

        switch c.align_items {
        case .Start:
            cross_pos = 0
        case .Center:
            cross_pos = (cross_size - child_cross) / 2
        case .End:
            cross_pos = cross_size - child_cross
        case .Stretch:
            child_cross = cross_size
            cross_pos = 0
        }

        // Set child rect
        if is_row {
            child.rect = core.Rect{
                x = main_pos,
                y = cross_pos,
                width = child_main,
                height = child_cross,
            }
        } else {
            child.rect = core.Rect{
                x = cross_pos,
                y = main_pos,
                width = child_cross,
                height = child_main,
            }
        }

        main_pos += child_main + c.spacing
    }
}

// Measure container's preferred size based on children
container_measure :: proc(w: ^Widget) -> core.Size {
    c := cast(^Container)w

    visible_children := make([dynamic]^Widget, context.temp_allocator)
    for child in c.children {
        if child.visible {
            append(&visible_children, child)
        }
    }

    if len(visible_children) == 0 {
        return core.Size{
            width = w.padding.left + w.padding.right,
            height = w.padding.top + w.padding.bottom,
        }
    }

    is_row := c.direction == .Row
    num_gaps := max(0, len(visible_children) - 1)
    total_spacing := c.spacing * i32(num_gaps)

    main_total: i32 = total_spacing
    cross_max: i32 = 0

    for child in visible_children {
        child_size := widget_measure(child)

        if is_row {
            main_total += child_size.width + child.margin.left + child.margin.right
            cross_max = max(cross_max, child_size.height + child.margin.top + child.margin.bottom)
        } else {
            main_total += child_size.height + child.margin.top + child.margin.bottom
            cross_max = max(cross_max, child_size.width + child.margin.left + child.margin.right)
        }
    }

    if is_row {
        return core.Size{
            width = main_total + w.padding.left + w.padding.right,
            height = cross_max + w.padding.top + w.padding.bottom,
        }
    } else {
        return core.Size{
            width = cross_max + w.padding.left + w.padding.right,
            height = main_total + w.padding.top + w.padding.bottom,
        }
    }
}

// Container cleanup (children destroyed by widget_destroy)
container_destroy :: proc(w: ^Widget) {
    // No container-specific cleanup needed
    // Children are destroyed by widget_destroy
}

// Helper to set container properties
container_set_spacing :: proc(c: ^Container, spacing: i32) {
    c.spacing = spacing
    widget_mark_dirty(c)
}

container_set_align :: proc(c: ^Container, align: Align) {
    c.align_items = align
    widget_mark_dirty(c)
}

container_set_direction :: proc(c: ^Container, direction: Direction) {
    c.direction = direction
    widget_mark_dirty(c)
}

container_set_background :: proc(c: ^Container, color: core.Color) {
    c.background = color
    widget_mark_dirty(c)
}
