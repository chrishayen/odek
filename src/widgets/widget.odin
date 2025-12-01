package widgets

import "../core"
import "../render"

// Widget vtable - function pointers for polymorphic dispatch
Widget_VTable :: struct {
    draw:         proc(w: ^Widget, ctx: ^render.Draw_Context),
    handle_event: proc(w: ^Widget, event: ^core.Event) -> bool,  // returns true if consumed
    layout:       proc(w: ^Widget),  // calculate size and position children
    destroy:      proc(w: ^Widget),  // cleanup widget-specific resources
    measure:      proc(w: ^Widget, available_width: i32) -> core.Size,  // preferred size given available width (-1 = unconstrained)
}

// Edge insets (padding, margin)
Edges :: struct {
    top, right, bottom, left: i32,
}

// Layout direction for containers
Direction :: enum {
    Row,
    Column,
}

// Alignment on cross axis
Align :: enum {
    Start,
    Center,
    End,
    Stretch,
}

// Cursor type for widget hover
Cursor :: enum {
    Arrow,      // Default arrow cursor
    Hand,       // Pointing hand for clickable items
    Text,       // I-beam for text input
    Wait,       // Busy/loading
    Crosshair,  // Crosshair
    Move,       // Move/drag
    Grab,       // Grabbing
}

// Base widget structure
Widget :: struct {
    vtable:   ^Widget_VTable,

    // Geometry
    rect:     core.Rect,      // position and size (relative to parent)
    padding:  Edges,          // internal padding
    margin:   Edges,          // external margin

    // Layout hints
    min_size: core.Size,
    max_size: core.Size,      // 0 = unlimited
    flex:     f32,            // flex grow factor (0 = fixed size)

    // Tree structure
    parent:   ^Widget,
    children: [dynamic]^Widget,

    // State
    visible:  bool,
    enabled:  bool,
    dirty:    bool,           // needs redraw
    focused:  bool,
    hovered:  bool,
    focusable: bool,          // can receive keyboard focus

    // Cursor to show when hovering
    cursor:   Cursor,

    // Widget-specific data
    user_data: rawptr,
}

// Create a new widget with the given vtable
widget_create :: proc(vtable: ^Widget_VTable) -> ^Widget {
    w := new(Widget)
    w.vtable = vtable
    w.visible = true
    w.enabled = true
    w.dirty = true
    return w
}

// Destroy a widget and all its children
widget_destroy :: proc(w: ^Widget) {
    if w == nil {
        return
    }

    // Clear hit state if this widget is hovered
    hit_state_clear_widget(w)

    // Destroy children first
    for child in w.children {
        widget_destroy(child)
    }
    delete(w.children)

    // Call widget-specific cleanup
    if w.vtable != nil && w.vtable.destroy != nil {
        w.vtable.destroy(w)
    }

    free(w)
}

// Add a child widget
widget_add_child :: proc(parent, child: ^Widget) {
    if parent == nil || child == nil {
        return
    }

    // Remove from old parent if any
    if child.parent != nil {
        widget_remove_child(child.parent, child)
    }

    child.parent = parent
    append(&parent.children, child)
    widget_mark_dirty(parent)
}

// Remove a child widget (does not destroy it)
widget_remove_child :: proc(parent, child: ^Widget) {
    if parent == nil || child == nil {
        return
    }

    for i := 0; i < len(parent.children); i += 1 {
        if parent.children[i] == child {
            ordered_remove(&parent.children, i)
            child.parent = nil
            widget_mark_dirty(parent)
            return
        }
    }
}

// Mark widget as needing redraw
widget_mark_dirty :: proc(w: ^Widget) {
    if w == nil {
        return
    }
    w.dirty = true
    // Propagate up to root so window knows to redraw
    if w.parent != nil {
        widget_mark_dirty(w.parent)
    }
}

// Get absolute rect in window coordinates
widget_get_absolute_rect :: proc(w: ^Widget) -> core.Rect {
    if w == nil {
        return {}
    }

    rect := w.rect

    // Walk up the tree adding parent offsets
    parent := w.parent
    for parent != nil {
        rect.x += parent.rect.x + parent.padding.left
        rect.y += parent.rect.y + parent.padding.top
        parent = parent.parent
    }

    return rect
}

// Draw a widget and its children
widget_draw :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
    if w == nil || !w.visible {
        return
    }

    // Get absolute position for clipping (in logical coordinates)
    abs_rect := widget_get_absolute_rect(w)

    // Save current logical clip and set widget clip
    // Use logical_clip for intersection since abs_rect is in logical coordinates
    old_logical_clip := ctx.logical_clip
    if clipped, ok := core.rect_intersection(abs_rect, old_logical_clip); ok {
        render.context_set_clip(ctx, clipped)
    } else {
        // Widget entirely outside clip region
        return
    }

    // Draw this widget
    if w.vtable != nil && w.vtable.draw != nil {
        w.vtable.draw(w, ctx)
    }

    // Draw debug border if enabled (skip root widget which has no parent)
    if debug_borders_enabled() && w.parent != nil {
        render.draw_rect(ctx, abs_rect, core.Color{255, 0, 0, 255})
    }

    // Draw children
    for child in w.children {
        widget_draw(child, ctx)
    }

    // Restore clip (pass logical clip to context_set_clip)
    render.context_set_clip(ctx, old_logical_clip)
    w.dirty = false
}

// Dispatch event to widget, returns true if consumed
widget_handle_event :: proc(w: ^Widget, event: ^core.Event) -> bool {
    if w == nil || !w.visible || !w.enabled {
        return false
    }

    if w.vtable != nil && w.vtable.handle_event != nil {
        return w.vtable.handle_event(w, event)
    }

    return false
}

// Perform layout on widget and children
widget_layout :: proc(w: ^Widget) {
    if w == nil {
        return
    }

    if w.vtable != nil && w.vtable.layout != nil {
        w.vtable.layout(w)
    }

    // Layout children recursively
    for child in w.children {
        widget_layout(child)
    }
}

// Measure preferred size given available width (-1 = unconstrained)
widget_measure :: proc(w: ^Widget, available_width: i32 = -1) -> core.Size {
    if w == nil {
        return {}
    }

    if w.vtable != nil && w.vtable.measure != nil {
        return w.vtable.measure(w, available_width)
    }

    // Default: use min_size or rect size
    return core.Size{
        width = max(w.min_size.width, w.rect.width),
        height = max(w.min_size.height, w.rect.height),
    }
}

// Set widget rect (position and size)
widget_set_rect :: proc(w: ^Widget, rect: core.Rect) {
    if w == nil {
        return
    }
    w.rect = rect
    widget_mark_dirty(w)
}

// Set widget position
widget_set_position :: proc(w: ^Widget, x, y: i32) {
    if w == nil {
        return
    }
    w.rect.x = x
    w.rect.y = y
    widget_mark_dirty(w)
}

// Set widget size
widget_set_size :: proc(w: ^Widget, width, height: i32) {
    if w == nil {
        return
    }
    w.rect.width = width
    w.rect.height = height
    widget_mark_dirty(w)
}

// Check if widget contains point (in absolute coordinates)
widget_contains_point :: proc(w: ^Widget, x, y: i32) -> bool {
    if w == nil || !w.visible {
        return false
    }
    abs_rect := widget_get_absolute_rect(w)
    return core.rect_contains(abs_rect, core.Point{x, y})
}

// Helper to create edges with same value on all sides
edges_all :: proc(value: i32) -> Edges {
    return Edges{value, value, value, value}
}

// Helper to create edges with horizontal and vertical values
edges_symmetric :: proc(horizontal, vertical: i32) -> Edges {
    return Edges{vertical, horizontal, vertical, horizontal}
}

// Get content rect (rect minus padding)
widget_get_content_rect :: proc(w: ^Widget) -> core.Rect {
    if w == nil {
        return {}
    }
    return core.Rect{
        x = w.padding.left,
        y = w.padding.top,
        width = w.rect.width - w.padding.left - w.padding.right,
        height = w.rect.height - w.padding.top - w.padding.bottom,
    }
}

// Default vtable for basic widgets (no-op implementations)
default_vtable := Widget_VTable{
    draw = nil,
    handle_event = nil,
    layout = nil,
    destroy = nil,
    measure = nil,
}
