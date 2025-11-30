package widgets

import "../core"
import "../render"
import "core:fmt"
import "core:path/filepath"

// Grid item type
Grid_Item_Type :: enum {
    Image,
    Folder,
    File,
}

// Grid item
Grid_Item :: struct {
    type:      Grid_Item_Type,
    image:     ^render.Image,  // Reference to image (not owned)
    thumbnail: ^render.Image,  // Reference to thumbnail for display (not owned)
    path:      string,         // File path for identification
    name:      string,         // Display name (for folders)
    loading:   bool,           // True while async loading in progress
    user_data: rawptr,         // User-defined data
}

// Image grid widget with scrolling and selection
Image_Grid :: struct {
    using base: Widget,

    // Items
    items: [dynamic]Grid_Item,

    // Font for labels
    font: ^render.Font,

    // Layout configuration
    cell_width:  i32,     // Width of each cell
    cell_height: i32,     // Height of each cell
    spacing:     i32,     // Gap between cells
    columns:     i32,     // Number of columns (0 = auto-calculate)

    // Visual
    background:      core.Color,
    selection_color: core.Color,
    hover_color:     core.Color,

    // State
    scroll:       Scroll_State,
    selected_idx: i32,    // -1 = none selected
    hovered_idx:  i32,    // -1 = none hovered

    // Scrollbar state
    scrollbar_width:    i32,
    scrollbar_dragging: bool,
    scrollbar_drag_start_y: i32,
    scrollbar_drag_start_offset: i32,

    // Callbacks
    on_click:            proc(grid: ^Image_Grid, index: i32, item: ^Grid_Item),
    on_folder_click:     proc(grid: ^Image_Grid, path: string),
    on_selection_change: proc(grid: ^Image_Grid, old_idx, new_idx: i32),
}

// Shared vtable
image_grid_vtable := Widget_VTable{
    draw         = image_grid_draw,
    handle_event = image_grid_handle_event,
    layout       = image_grid_layout,
    destroy      = image_grid_destroy,
    measure      = image_grid_measure,
}

// Create image grid
image_grid_create :: proc() -> ^Image_Grid {
    g := new(Image_Grid)
    g.vtable = &image_grid_vtable
    g.visible = true
    g.enabled = true
    g.dirty = true
    g.focusable = true

    // Defaults
    g.cell_width = 150
    g.cell_height = 150
    g.spacing = 10
    g.columns = 0  // Auto-calculate

    g.background = core.color_hex(0x2D2D2D)
    g.selection_color = core.color_rgba(74, 144, 217, 200)
    g.hover_color = core.color_rgba(255, 255, 255, 40)

    g.scroll = scroll_init()
    g.selected_idx = -1
    g.hovered_idx = -1
    g.scrollbar_width = 8

    return g
}

// Add an image to the grid
image_grid_add_item :: proc(g: ^Image_Grid, img: ^render.Image, thumbnail: ^render.Image = nil, path: string = "", user_data: rawptr = nil) {
    item := Grid_Item{
        type = .Image,
        image = img,
        thumbnail = thumbnail if thumbnail != nil else img,
        path = path,
        user_data = user_data,
    }
    append(&g.items, item)
    widget_mark_dirty(g)
}

// Add a folder to the grid
image_grid_add_folder :: proc(g: ^Image_Grid, name: string, path: string) {
    item := Grid_Item{
        type = .Folder,
        path = path,
        name = name,
    }
    append(&g.items, item)
    widget_mark_dirty(g)
}

// Add a file to the grid
image_grid_add_file :: proc(g: ^Image_Grid, name: string, path: string) {
    item := Grid_Item{
        type = .File,
        path = path,
        name = name,
    }
    append(&g.items, item)
    widget_mark_dirty(g)
}

// Add a placeholder image item (for async loading)
// Returns the index of the added item
image_grid_add_placeholder :: proc(g: ^Image_Grid, name: string, path: string) -> i32 {
    item := Grid_Item{
        type = .Image,
        path = path,
        name = name,
        loading = true,
    }
    append(&g.items, item)
    widget_mark_dirty(g)
    return i32(len(g.items) - 1)
}

// Set image for a grid item (call when async load completes)
image_grid_set_image :: proc(g: ^Image_Grid, index: i32, img: ^render.Image, thumbnail: ^render.Image) {
    if index < 0 || index >= i32(len(g.items)) {
        return
    }
    g.items[index].image = img
    g.items[index].thumbnail = thumbnail
    g.items[index].loading = false
    widget_mark_dirty(g)
}

// Remove item at index
image_grid_remove_item :: proc(g: ^Image_Grid, index: i32) {
    if index < 0 || index >= i32(len(g.items)) {
        return
    }
    ordered_remove(&g.items, int(index))

    // Adjust selection
    if g.selected_idx == index {
        g.selected_idx = -1
    } else if g.selected_idx > index {
        g.selected_idx -= 1
    }

    widget_mark_dirty(g)
}

// Clear all items
image_grid_clear :: proc(g: ^Image_Grid) {
    clear(&g.items)
    g.selected_idx = -1
    g.hovered_idx = -1
    g.scroll.offset = 0
    widget_mark_dirty(g)
}

// Set selected item
image_grid_set_selected :: proc(g: ^Image_Grid, index: i32) {
    if index == g.selected_idx {
        return
    }

    old_idx := g.selected_idx
    g.selected_idx = index

    if g.on_selection_change != nil {
        g.on_selection_change(g, old_idx, index)
    }

    widget_mark_dirty(g)
}

// Get number of columns based on widget width
image_grid_get_columns :: proc(g: ^Image_Grid) -> i32 {
    if g.columns > 0 {
        return g.columns
    }

    content_width := g.rect.width - g.padding.left - g.padding.right
    cell_total := g.cell_width + g.spacing
    if cell_total <= 0 {
        return 1
    }

    cols := (content_width + g.spacing) / cell_total
    return max(1, cols)
}

// Get grid position from item index (relative to widget, accounting for scroll)
image_grid_get_item_rect :: proc(g: ^Image_Grid, index: i32) -> core.Rect {
    cols := image_grid_get_columns(g)
    row := index / cols
    col := index % cols

    x := g.padding.left + col * (g.cell_width + g.spacing)
    y := g.padding.top + row * (g.cell_height + g.spacing) - g.scroll.offset

    return core.Rect{x, y, g.cell_width, g.cell_height}
}

// Get item index from point (relative to widget)
image_grid_get_item_at :: proc(g: ^Image_Grid, x, y: i32) -> i32 {
    // Adjust for scroll
    adjusted_y := y + g.scroll.offset

    cols := image_grid_get_columns(g)
    cell_total_w := g.cell_width + g.spacing
    cell_total_h := g.cell_height + g.spacing

    rel_x := x - g.padding.left
    rel_y := adjusted_y - g.padding.top

    if rel_x < 0 || rel_y < 0 {
        return -1
    }

    col := rel_x / cell_total_w
    row := rel_y / cell_total_h

    if col >= cols {
        return -1
    }

    // Check if within cell bounds (not in spacing)
    local_x := rel_x % cell_total_w
    local_y := rel_y % cell_total_h

    if local_x >= g.cell_width || local_y >= g.cell_height {
        return -1  // In spacing
    }

    index := row * cols + col
    if index < 0 || index >= i32(len(g.items)) {
        return -1
    }

    return index
}

// Calculate total content height
image_grid_get_content_height :: proc(g: ^Image_Grid) -> i32 {
    if len(g.items) == 0 {
        return g.padding.top + g.padding.bottom
    }

    cols := image_grid_get_columns(g)
    rows := (i32(len(g.items)) + cols - 1) / cols

    return g.padding.top + rows * (g.cell_height + g.spacing) - g.spacing + g.padding.bottom
}

// Get scrollbar track rect (relative to widget)
image_grid_get_scrollbar_rect :: proc(g: ^Image_Grid) -> core.Rect {
    viewport_h := g.rect.height - g.padding.top - g.padding.bottom
    return core.Rect{
        x = g.rect.width - g.scrollbar_width - 2,
        y = g.padding.top,
        width = g.scrollbar_width,
        height = viewport_h,
    }
}

// Get scrollbar thumb rect (relative to widget)
image_grid_get_thumb_rect :: proc(g: ^Image_Grid) -> core.Rect {
    viewport_h := g.rect.height - g.padding.top - g.padding.bottom
    max_offset := scroll_get_max_offset(&g.scroll)

    if max_offset <= 0 {
        return core.Rect{}
    }

    thumb_height := max(20, (viewport_h * viewport_h) / g.scroll.content_size)
    thumb_y := g.padding.top + ((viewport_h - thumb_height) * g.scroll.offset) / max_offset

    return core.Rect{
        x = g.rect.width - g.scrollbar_width - 2,
        y = thumb_y,
        width = g.scrollbar_width,
        height = thumb_height,
    }
}

// Check if point is in scrollbar area
image_grid_point_in_scrollbar :: proc(g: ^Image_Grid, x, y: i32) -> bool {
    sb := image_grid_get_scrollbar_rect(g)
    return x >= sb.x && x < sb.x + sb.width && y >= sb.y && y < sb.y + sb.height
}

// Draw a loading placeholder in the given rect
draw_loading_placeholder :: proc(ctx: ^render.Draw_Context, rect: core.Rect) {
    // Dark background
    bg_color := core.color_hex(0x3A3A3A)
    render.fill_rect(ctx, rect, bg_color)

    // Draw a simple loading indicator (gray box with lighter center)
    center_size := min(rect.width, rect.height) * 30 / 100
    center_x := rect.x + (rect.width - center_size) / 2
    center_y := rect.y + (rect.height - center_size) / 2

    indicator_rect := core.Rect{center_x, center_y, center_size, center_size}
    indicator_color := core.color_hex(0x555555)
    render.fill_rect(ctx, indicator_rect, indicator_color)
}

// Draw a folder icon in the given rect
draw_folder_icon :: proc(ctx: ^render.Draw_Context, rect: core.Rect) {
    // Folder colors
    folder_body := core.color_hex(0x5B9BD5)     // Blue folder body
    folder_tab := core.color_hex(0x4A8AC4)      // Slightly darker tab
    folder_shadow := core.color_hex(0x3A7AB4)   // Shadow/depth

    // Calculate icon dimensions (centered, 60% of cell size)
    icon_w := rect.width * 60 / 100
    icon_h := rect.height * 50 / 100
    icon_x := rect.x + (rect.width - icon_w) / 2
    icon_y := rect.y + (rect.height - icon_h) / 2

    // Draw tab (top-left portion)
    tab_w := icon_w * 35 / 100
    tab_h := icon_h * 15 / 100
    tab_rect := core.Rect{icon_x, icon_y, tab_w, tab_h}
    render.fill_rect(ctx, tab_rect, folder_tab)

    // Draw main folder body
    body_rect := core.Rect{icon_x, icon_y + tab_h, icon_w, icon_h - tab_h}
    render.fill_rect(ctx, body_rect, folder_body)

    // Draw bottom edge for depth
    edge_rect := core.Rect{icon_x, icon_y + icon_h - 4, icon_w, 4}
    render.fill_rect(ctx, edge_rect, folder_shadow)
}

// Draw a blank file icon in the given rect
draw_file_icon :: proc(ctx: ^render.Draw_Context, rect: core.Rect) {
    // File colors
    file_body := core.color_hex(0x888888)       // Gray file body
    file_corner := core.color_hex(0x666666)     // Darker folded corner
    file_lines := core.color_hex(0x777777)      // Lines on file

    // Calculate icon dimensions (centered, 45% width, 55% height)
    icon_w := rect.width * 45 / 100
    icon_h := rect.height * 55 / 100
    icon_x := rect.x + (rect.width - icon_w) / 2
    icon_y := rect.y + (rect.height - icon_h) / 2

    // Draw main file body
    render.fill_rect(ctx, core.Rect{icon_x, icon_y, icon_w, icon_h}, file_body)

    // Draw folded corner (top-right)
    corner_size := icon_w * 25 / 100
    corner_x := icon_x + icon_w - corner_size
    corner_y := icon_y
    render.fill_rect(ctx, core.Rect{corner_x, corner_y, corner_size, corner_size}, file_corner)

    // Draw a few lines to suggest text
    line_margin := icon_w * 15 / 100
    line_y := icon_y + icon_h * 40 / 100
    line_w := icon_w - line_margin * 2
    line_h: i32 = 2
    for i in 0..<3 {
        render.fill_rect(ctx, core.Rect{icon_x + line_margin, line_y + i32(i) * 8, line_w, line_h}, file_lines)
    }
}

// Draw function
image_grid_draw :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
    g := cast(^Image_Grid)w
    abs_rect := widget_get_absolute_rect(w)

    // Draw background
    if g.background.a > 0 {
        render.fill_rect(ctx, abs_rect, g.background)
    }

    // Set up clipping for content area (use logical_clip since content_rect is logical)
    old_logical_clip := ctx.logical_clip
    content_rect := core.Rect{
        x = abs_rect.x + w.padding.left,
        y = abs_rect.y + w.padding.top,
        width = abs_rect.width - w.padding.left - w.padding.right,
        height = abs_rect.height - w.padding.top - w.padding.bottom,
    }

    if clipped, ok := core.rect_intersection(content_rect, old_logical_clip); ok {
        render.context_set_clip(ctx, clipped)
    } else {
        return  // Nothing visible
    }

    // Draw each visible item
    for i in 0 ..< i32(len(g.items)) {
        item_rect := image_grid_get_item_rect(g, i)

        // Transform to absolute coordinates
        item_abs := core.Rect{
            x = abs_rect.x + item_rect.x,
            y = abs_rect.y + item_rect.y,
            width = item_rect.width,
            height = item_rect.height,
        }

        // Skip if completely outside visible area
        if item_abs.y + item_abs.height < content_rect.y ||
           item_abs.y > content_rect.y + content_rect.height {
            continue
        }

        // Draw selection/hover highlight
        if i == g.selected_idx {
            render.fill_rect_blend(ctx, item_abs, g.selection_color)
        } else if i == g.hovered_idx {
            render.fill_rect_blend(ctx, item_abs, g.hover_color)
        }

        item := &g.items[i]

        // Reserve space for label at bottom
        label_height: i32 = 20 if g.font != nil else 0
        icon_rect := core.Rect{
            x = item_abs.x,
            y = item_abs.y,
            width = item_abs.width,
            height = item_abs.height - label_height,
        }

        if item.type == .Folder {
            // Draw folder icon
            draw_folder_icon(ctx, icon_rect)
        } else if item.type == .File {
            // Draw file icon
            draw_file_icon(ctx, icon_rect)
        } else if item.loading {
            // Draw loading placeholder
            draw_loading_placeholder(ctx, icon_rect)
        } else {
            // Draw image (use thumbnail if available)
            display_img := item.thumbnail if item.thumbnail != nil else item.image
            if display_img != nil {
                render.draw_image_scaled(ctx, display_img, icon_rect)
            }
        }

        // Draw selection border
        if i == g.selected_idx {
            render.draw_rect(ctx, item_abs, g.selection_color, 2)
        }

        // Draw label
        if g.font != nil {
            label_text: string
            if item.type == .Folder || item.type == .File {
                label_text = item.name
            } else {
                // Extract filename from path
                label_text = filepath.base(item.path)
            }

            if len(label_text) > 0 {
                // Convert max_width to physical for comparison with text measurements
                // (font is loaded at scaled size, so text_measure returns physical width)
                max_width_phys := i32(f64(g.cell_width - 10) * ctx.scale)
                text_width_phys := render.text_measure(g.font, label_text)

                // Truncate with ellipsis if too wide
                display_text := label_text
                truncated_buf: [256]u8
                if text_width_phys > max_width_phys {
                    ellipsis :: "..."
                    ellipsis_width := render.text_measure(g.font, ellipsis)
                    target_width := max_width_phys - ellipsis_width

                    // Find truncation point
                    for i := len(label_text) - 1; i > 0; i -= 1 {
                        truncated := label_text[:i]
                        width := render.text_measure(g.font, truncated)
                        if width <= target_width {
                            // Build truncated string with ellipsis
                            display_text = fmt.bprintf(truncated_buf[:], "%s...", truncated)
                            break
                        }
                    }
                    text_width_phys = render.text_measure(g.font, display_text)
                }

                // Convert text width back to logical for centering calculation
                text_width_logical := i32(f64(text_width_phys) / ctx.scale)

                // Center the text (all in logical coordinates)
                text_x := item_abs.x + (item_abs.width - text_width_logical) / 2
                text_y := item_abs.y + item_abs.height - label_height - 12

                label_color := core.color_hex(0xCCCCCC)
                render.draw_text_top(ctx, g.font, display_text, text_x, text_y, label_color)
            }
        }
    }

    // Restore clip
    render.context_set_clip(ctx, old_logical_clip)

    // Draw scroll indicator if scrollable
    if scroll_is_scrollable(&g.scroll) {
        sb_rect := image_grid_get_scrollbar_rect(g)
        thumb_rect := image_grid_get_thumb_rect(g)

        // Transform to absolute coordinates
        sb_abs := core.Rect{
            x = abs_rect.x + sb_rect.x,
            y = abs_rect.y + sb_rect.y,
            width = sb_rect.width,
            height = sb_rect.height,
        }
        thumb_abs := core.Rect{
            x = abs_rect.x + thumb_rect.x,
            y = abs_rect.y + thumb_rect.y,
            width = thumb_rect.width,
            height = thumb_rect.height,
        }

        // Draw track
        track_color := core.color_rgba(255, 255, 255, 20)
        render.fill_rect(ctx, sb_abs, track_color)

        // Draw thumb (brighter when dragging)
        thumb_color := core.color_rgba(255, 255, 255, 120) if g.scrollbar_dragging else core.color_rgba(255, 255, 255, 80)
        render.fill_rounded_rect(ctx, thumb_abs, 3, thumb_color)
    }
}

// Event handling
image_grid_handle_event :: proc(w: ^Widget, event: ^core.Event) -> bool {
    g := cast(^Image_Grid)w

    if !w.enabled {
        return false
    }

    #partial switch event.type {
    case .Pointer_Motion:
        abs_rect := widget_get_absolute_rect(w)
        local_x := event.pointer_x - abs_rect.x
        local_y := event.pointer_y - abs_rect.y

        // Handle scrollbar dragging
        if g.scrollbar_dragging {
            viewport_h := g.rect.height - g.padding.top - g.padding.bottom
            thumb_rect := image_grid_get_thumb_rect(g)
            max_offset := scroll_get_max_offset(&g.scroll)

            if max_offset > 0 && viewport_h > thumb_rect.height {
                // Calculate new scroll offset based on drag delta
                drag_delta := local_y - g.scrollbar_drag_start_y
                scroll_range := viewport_h - thumb_rect.height
                offset_delta := (drag_delta * max_offset) / scroll_range
                scroll_set_offset(&g.scroll, g.scrollbar_drag_start_offset + offset_delta)
                widget_mark_dirty(g)
            }
            return true
        }

        // Update hover state
        new_hover := image_grid_get_item_at(g, local_x, local_y)
        if new_hover != g.hovered_idx {
            g.hovered_idx = new_hover
            widget_mark_dirty(g)
        }
        return true

    case .Pointer_Leave:
        if g.hovered_idx != -1 {
            g.hovered_idx = -1
            widget_mark_dirty(g)
        }
        return true

    case .Pointer_Button_Press:
        if event.button == .Left {
            abs_rect := widget_get_absolute_rect(w)
            local_x := event.pointer_x - abs_rect.x
            local_y := event.pointer_y - abs_rect.y

            // Check if clicking scrollbar
            if scroll_is_scrollable(&g.scroll) && image_grid_point_in_scrollbar(g, local_x, local_y) {
                thumb_rect := image_grid_get_thumb_rect(g)

                // Check if clicking on thumb
                if local_y >= thumb_rect.y && local_y < thumb_rect.y + thumb_rect.height {
                    // Start dragging thumb
                    g.scrollbar_dragging = true
                    g.scrollbar_drag_start_y = local_y
                    g.scrollbar_drag_start_offset = g.scroll.offset
                } else {
                    // Click on track - jump to position
                    viewport_h := g.rect.height - g.padding.top - g.padding.bottom
                    max_offset := scroll_get_max_offset(&g.scroll)
                    track_y := local_y - g.padding.top
                    new_offset := (track_y * max_offset) / viewport_h
                    scroll_set_offset(&g.scroll, new_offset)
                    widget_mark_dirty(g)
                }
                return true
            }

            // Check if clicking on item
            clicked_idx := image_grid_get_item_at(g, local_x, local_y)
            if clicked_idx >= 0 {
                image_grid_set_selected(g, clicked_idx)
                item := &g.items[clicked_idx]
                if item.type == .Folder {
                    if g.on_folder_click != nil {
                        g.on_folder_click(g, item.path)
                    }
                } else {
                    if g.on_click != nil {
                        g.on_click(g, clicked_idx, item)
                    }
                }
            }
            return true
        }

    case .Pointer_Button_Release:
        if event.button == .Left && g.scrollbar_dragging {
            g.scrollbar_dragging = false
            return true
        }

    case .Scroll:
        // Scroll by event amount
        scroll_by(&g.scroll, event.scroll_delta)
        widget_mark_dirty(g)
        return true

    case .Key_Press:
        // Keyboard navigation
        if w.focused {
            cols := image_grid_get_columns(g)
            new_sel := g.selected_idx

            // Initialize selection if none
            if new_sel < 0 && len(g.items) > 0 {
                new_sel = 0
            }

            // Arrow key navigation
            if event.keycode == 105 {  // Left
                if new_sel > 0 {
                    new_sel -= 1
                }
            } else if event.keycode == 106 {  // Right
                if new_sel < i32(len(g.items)) - 1 {
                    new_sel += 1
                }
            } else if event.keycode == 103 {  // Up
                if new_sel >= cols {
                    new_sel -= cols
                }
            } else if event.keycode == 108 {  // Down
                if new_sel + cols < i32(len(g.items)) {
                    new_sel += cols
                }
            } else if event.keycode == 28 {  // Enter
                if g.selected_idx >= 0 && g.on_click != nil {
                    g.on_click(g, g.selected_idx, &g.items[g.selected_idx])
                }
                return true
            }

            if new_sel != g.selected_idx && new_sel >= 0 {
                image_grid_set_selected(g, new_sel)
                image_grid_ensure_visible(g, new_sel)
                return true
            }
        }
    }

    return false
}

// Ensure item at index is visible (scroll if needed)
image_grid_ensure_visible :: proc(g: ^Image_Grid, index: i32) {
    if index < 0 || index >= i32(len(g.items)) {
        return
    }

    cols := image_grid_get_columns(g)
    row := index / cols

    item_top := g.padding.top + row * (g.cell_height + g.spacing)
    item_bottom := item_top + g.cell_height

    scroll_ensure_visible(&g.scroll, item_top, g.cell_height)
}

// Layout function
image_grid_layout :: proc(w: ^Widget) {
    g := cast(^Image_Grid)w

    // Update scroll sizes
    content_h := image_grid_get_content_height(g)
    viewport_h := w.rect.height - w.padding.top - w.padding.bottom
    scroll_set_sizes(&g.scroll, content_h, viewport_h)
}

// Measure function
image_grid_measure :: proc(w: ^Widget) -> core.Size {
    g := cast(^Image_Grid)w

    // Preferred size showing at least 2 rows
    min_rows: i32 = 2
    min_height := g.padding.top + min_rows * (g.cell_height + g.spacing) - g.spacing + g.padding.bottom

    return core.Size{
        width = max(w.min_size.width, w.rect.width),
        height = max(w.min_size.height, min_height),
    }
}

// Destroy function
image_grid_destroy :: proc(w: ^Widget) {
    g := cast(^Image_Grid)w
    delete(g.items)
    // Note: Images are not owned by grid, caller must destroy them
}

// Get item count
image_grid_count :: proc(g: ^Image_Grid) -> i32 {
    return i32(len(g.items))
}

// Get selected item
image_grid_get_selected :: proc(g: ^Image_Grid) -> (^Grid_Item, bool) {
    if g.selected_idx < 0 || g.selected_idx >= i32(len(g.items)) {
        return nil, false
    }
    return &g.items[g.selected_idx], true
}
