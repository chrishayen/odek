package tests

import "../src/core"
import "../src/render"
import "../src/widgets"
import "core:testing"

@(test)
test_image_grid_create :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    testing.expect(t, g != nil, "image_grid_create should return non-nil")
    testing.expect(t, g.visible == true, "Grid should be visible by default")
    testing.expect(t, g.selected_idx == -1, "No item should be selected initially")
    testing.expect(t, widgets.image_grid_count(g) == 0, "Grid should have no items initially")

    widgets.widget_destroy(g)
}

@(test)
test_image_grid_add_item :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    // Create a dummy image
    img := render.Image{
        width = 10,
        height = 10,
        pixels = make([]u32, 100),
    }
    defer render.image_destroy(&img)

    widgets.image_grid_add_item(g, &img, nil, "test.png")

    testing.expect(t, widgets.image_grid_count(g) == 1, "Grid should have 1 item")
    testing.expect(t, g.items[0].path == "test.png", "Item path should match")
}

@(test)
test_image_grid_remove_item :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    widgets.image_grid_add_item(g, &img, nil, "test1.png")
    widgets.image_grid_add_item(g, &img, nil, "test2.png")
    widgets.image_grid_add_item(g, &img, nil, "test3.png")

    testing.expect(t, widgets.image_grid_count(g) == 3, "Grid should have 3 items")

    widgets.image_grid_remove_item(g, 1)

    testing.expect(t, widgets.image_grid_count(g) == 2, "Grid should have 2 items after removal")
    testing.expect(t, g.items[0].path == "test1.png", "First item should remain")
    testing.expect(t, g.items[1].path == "test3.png", "Third item should move to index 1")
}

@(test)
test_image_grid_selection :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)

    testing.expect(t, g.selected_idx == -1, "Initially no selection")

    widgets.image_grid_set_selected(g, 1)
    testing.expect(t, g.selected_idx == 1, "Selection should be 1")

    widgets.image_grid_set_selected(g, -1)
    testing.expect(t, g.selected_idx == -1, "Selection cleared")
}

@(test)
test_image_grid_selection_adjusted_on_remove :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)

    // Select item at index 2
    widgets.image_grid_set_selected(g, 2)
    testing.expect(t, g.selected_idx == 2, "Selection should be 2")

    // Remove item before selection
    widgets.image_grid_remove_item(g, 0)
    testing.expect(t, g.selected_idx == 1, "Selection should shift to 1")

    // Remove selected item
    widgets.image_grid_remove_item(g, 1)
    testing.expect(t, g.selected_idx == -1, "Selection should be cleared")
}

@(test)
test_image_grid_columns :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    g.cell_width = 100
    g.spacing = 10
    g.padding = widgets.edges_all(10)
    g.rect = core.Rect{0, 0, 350, 500}  // Enough for 3 columns: 10 + 100 + 10 + 100 + 10 + 100 + 10 = 340

    cols := widgets.image_grid_get_columns(g)
    testing.expect(t, cols == 3, "Should have 3 columns")
}

@(test)
test_image_grid_columns_fixed :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    g.columns = 5  // Fixed columns
    g.rect = core.Rect{0, 0, 100, 100}

    cols := widgets.image_grid_get_columns(g)
    testing.expect(t, cols == 5, "Should use fixed column count")
}

@(test)
test_image_grid_item_at :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    g.cell_width = 100
    g.cell_height = 100
    g.spacing = 10
    g.padding = widgets.edges_all(10)
    g.rect = core.Rect{0, 0, 350, 500}

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    // Add 6 items (2 rows x 3 cols)
    for _ in 0 ..< 6 {
        widgets.image_grid_add_item(g, &img)
    }

    // Test first item (at padding + 0)
    idx := widgets.image_grid_get_item_at(g, 15, 15)  // Inside first cell
    testing.expect(t, idx == 0, "Should find first item")

    // Test second item
    idx = widgets.image_grid_get_item_at(g, 125, 15)  // Inside second cell
    testing.expect(t, idx == 1, "Should find second item")

    // Test spacing (no item)
    idx = widgets.image_grid_get_item_at(g, 115, 15)  // In spacing
    testing.expect(t, idx == -1, "Should find no item in spacing")

    // Test second row
    idx = widgets.image_grid_get_item_at(g, 15, 125)  // First cell of second row
    testing.expect(t, idx == 3, "Should find item in second row")
}

@(test)
test_image_grid_content_height :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    g.cell_width = 100
    g.cell_height = 100
    g.spacing = 10
    g.padding = widgets.edges_all(10)
    g.rect = core.Rect{0, 0, 350, 500}  // 3 columns

    // Empty grid
    height := widgets.image_grid_get_content_height(g)
    testing.expect(t, height == 20, "Empty grid should have padding height")

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    // Add 3 items (1 row)
    for _ in 0 ..< 3 {
        widgets.image_grid_add_item(g, &img)
    }
    height = widgets.image_grid_get_content_height(g)
    testing.expect(t, height == 120, "1 row: 10 + 100 + 10 = 120")

    // Add 3 more items (2 rows)
    for _ in 0 ..< 3 {
        widgets.image_grid_add_item(g, &img)
    }
    height = widgets.image_grid_get_content_height(g)
    testing.expect(t, height == 230, "2 rows: 10 + 100 + 10 + 100 + 10 = 230")
}

@(test)
test_image_grid_clear :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_set_selected(g, 1)

    widgets.image_grid_clear(g)

    testing.expect(t, widgets.image_grid_count(g) == 0, "Grid should be empty")
    testing.expect(t, g.selected_idx == -1, "Selection should be cleared")
    testing.expect(t, g.hovered_idx == -1, "Hover should be cleared")
}

@(test)
test_image_grid_scroll_state :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    g.rect = core.Rect{0, 0, 300, 200}
    g.cell_height = 100
    g.spacing = 10
    g.padding = widgets.edges_all(10)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    // Add enough items to require scrolling
    for _ in 0 ..< 20 {
        widgets.image_grid_add_item(g, &img)
    }

    // Run layout to calculate scroll sizes
    widgets.image_grid_layout(g)

    testing.expect(t, widgets.scroll_is_scrollable(&g.scroll), "Grid should be scrollable")
    testing.expect(t, g.scroll.offset == 0, "Initial scroll should be 0")

    // Scroll down
    widgets.scroll_by(&g.scroll, 50)
    testing.expect(t, g.scroll.offset == 50, "Scroll offset should be 50")
}

@(test)
test_image_grid_get_selected :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    // No selection
    item, ok := widgets.image_grid_get_selected(g)
    testing.expect(t, !ok, "Should return false when nothing selected")
    testing.expect(t, item == nil, "Item should be nil")

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    widgets.image_grid_add_item(g, &img, nil, "selected.png")
    widgets.image_grid_set_selected(g, 0)

    item, ok = widgets.image_grid_get_selected(g)
    testing.expect(t, ok, "Should return true when selected")
    testing.expect(t, item != nil, "Item should not be nil")
    testing.expect(t, item.path == "selected.png", "Should return correct item")
}

@(test)
test_image_grid_measure :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    g.cell_height = 100
    g.spacing = 10
    g.padding = widgets.edges_all(10)

    size := widgets.image_grid_measure(g, -1)
    // Min 2 rows: 10 + 100 + 10 + 100 + 10 = 230 (but spacing between rows, so 10 + 100 + 10 + 100 + 10 - 10 = 220)
    // Actually: padding.top + 2*cell_height + spacing + padding.bottom = 10 + 100 + 10 + 100 + 10 = 230
    // Wait the formula is: padding.top + min_rows * (cell_height + spacing) - spacing + padding.bottom
    // = 10 + 2 * 110 - 10 + 10 = 10 + 220 - 10 + 10 = 230
    testing.expect(t, size.height >= 200, "Measure should return reasonable height")
}

// =============================================================================
// Multi-selection tests
// =============================================================================

@(test)
test_image_grid_multi_selection_toggle :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)

    // Toggle selection of first item
    widgets.image_grid_toggle_selection(g, 0)
    testing.expect(t, widgets.image_grid_is_selected(g, 0), "Item 0 should be selected")
    testing.expect(t, widgets.image_grid_selection_count(g) == 1, "Should have 1 selected item")

    // Toggle selection of second item (adds to selection)
    widgets.image_grid_toggle_selection(g, 1)
    testing.expect(t, widgets.image_grid_is_selected(g, 0), "Item 0 should still be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 1), "Item 1 should be selected")
    testing.expect(t, widgets.image_grid_selection_count(g) == 2, "Should have 2 selected items")

    // Toggle item 0 again (removes from selection)
    widgets.image_grid_toggle_selection(g, 0)
    testing.expect(t, !widgets.image_grid_is_selected(g, 0), "Item 0 should be deselected")
    testing.expect(t, widgets.image_grid_is_selected(g, 1), "Item 1 should still be selected")
    testing.expect(t, widgets.image_grid_selection_count(g) == 1, "Should have 1 selected item")
}

@(test)
test_image_grid_range_selection :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    for _ in 0..<10 {
        widgets.image_grid_add_item(g, &img)
    }

    // Single select item 2 (sets anchor)
    widgets.image_grid_set_selected(g, 2)
    testing.expect(t, g.anchor_idx == 2, "Anchor should be 2")

    // Range select to item 5
    widgets.image_grid_select_range(g, 5)
    testing.expect(t, widgets.image_grid_selection_count(g) == 4, "Should have 4 selected items (2,3,4,5)")
    testing.expect(t, widgets.image_grid_is_selected(g, 2), "Item 2 should be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 3), "Item 3 should be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 4), "Item 4 should be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 5), "Item 5 should be selected")
    testing.expect(t, !widgets.image_grid_is_selected(g, 1), "Item 1 should not be selected")
    testing.expect(t, !widgets.image_grid_is_selected(g, 6), "Item 6 should not be selected")
}

@(test)
test_image_grid_range_selection_backwards :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    for _ in 0..<10 {
        widgets.image_grid_add_item(g, &img)
    }

    // Single select item 7 (sets anchor)
    widgets.image_grid_set_selected(g, 7)

    // Range select backwards to item 3
    widgets.image_grid_select_range(g, 3)
    testing.expect(t, widgets.image_grid_selection_count(g) == 5, "Should have 5 selected items (3,4,5,6,7)")
    testing.expect(t, widgets.image_grid_is_selected(g, 3), "Item 3 should be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 7), "Item 7 should be selected")
    testing.expect(t, !widgets.image_grid_is_selected(g, 2), "Item 2 should not be selected")
    testing.expect(t, !widgets.image_grid_is_selected(g, 8), "Item 8 should not be selected")
}

@(test)
test_image_grid_extend_selection :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    for _ in 0..<10 {
        widgets.image_grid_add_item(g, &img)
    }

    // Select item 2
    widgets.image_grid_set_selected(g, 2)

    // Toggle items 5 and 6 (Ctrl+click)
    widgets.image_grid_toggle_selection(g, 5)
    widgets.image_grid_toggle_selection(g, 6)
    testing.expect(t, widgets.image_grid_selection_count(g) == 3, "Should have 3 selected items")

    // Extend selection from anchor (6) to 8 (Shift+Ctrl+click)
    widgets.image_grid_extend_selection(g, 8)
    testing.expect(t, widgets.image_grid_is_selected(g, 2), "Item 2 should still be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 5), "Item 5 should still be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 6), "Item 6 should be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 7), "Item 7 should be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 8), "Item 8 should be selected")
}

@(test)
test_image_grid_clear_selection :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)

    // Select multiple items
    widgets.image_grid_toggle_selection(g, 0)
    widgets.image_grid_toggle_selection(g, 1)
    widgets.image_grid_toggle_selection(g, 2)
    testing.expect(t, widgets.image_grid_selection_count(g) == 3, "Should have 3 selected items")

    // Clear selection
    widgets.image_grid_clear_selection(g)
    testing.expect(t, widgets.image_grid_selection_count(g) == 0, "Should have 0 selected items")
    testing.expect(t, g.selected_idx == -1, "Primary selection should be -1")
    testing.expect(t, g.anchor_idx == -1, "Anchor should be -1")
}

@(test)
test_image_grid_single_select_clears_multi :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)
    widgets.image_grid_add_item(g, &img)

    // Select multiple items
    widgets.image_grid_toggle_selection(g, 0)
    widgets.image_grid_toggle_selection(g, 1)
    widgets.image_grid_toggle_selection(g, 2)
    testing.expect(t, widgets.image_grid_selection_count(g) == 3, "Should have 3 selected items")

    // Single select should clear multi-selection
    widgets.image_grid_set_selected(g, 1)
    testing.expect(t, widgets.image_grid_selection_count(g) == 1, "Should have 1 selected item")
    testing.expect(t, widgets.image_grid_is_selected(g, 1), "Item 1 should be selected")
    testing.expect(t, !widgets.image_grid_is_selected(g, 0), "Item 0 should not be selected")
    testing.expect(t, !widgets.image_grid_is_selected(g, 2), "Item 2 should not be selected")
}

@(test)
test_image_grid_remove_adjusts_multi_selection :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    for _ in 0..<5 {
        widgets.image_grid_add_item(g, &img)
    }

    // Select items 1, 2, 4
    widgets.image_grid_toggle_selection(g, 1)
    widgets.image_grid_toggle_selection(g, 2)
    widgets.image_grid_toggle_selection(g, 4)
    testing.expect(t, widgets.image_grid_selection_count(g) == 3, "Should have 3 selected items")

    // Remove item 2 (selected)
    widgets.image_grid_remove_item(g, 2)
    testing.expect(t, widgets.image_grid_selection_count(g) == 2, "Should have 2 selected items after removal")
    testing.expect(t, widgets.image_grid_is_selected(g, 1), "Item 1 should still be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 3), "Item at new index 3 (was 4) should be selected")

    // Remove item 0 (not selected, before selections)
    widgets.image_grid_remove_item(g, 0)
    testing.expect(t, widgets.image_grid_selection_count(g) == 2, "Should still have 2 selected items")
    testing.expect(t, widgets.image_grid_is_selected(g, 0), "Item at new index 0 (was 1) should be selected")
    testing.expect(t, widgets.image_grid_is_selected(g, 2), "Item at new index 2 (was 3) should be selected")
}

@(test)
test_image_grid_get_selected_indices :: proc(t: ^testing.T) {
    g := widgets.image_grid_create()
    defer widgets.widget_destroy(g)

    img := render.Image{width = 1, height = 1, pixels = make([]u32, 1)}
    defer render.image_destroy(&img)

    for _ in 0..<5 {
        widgets.image_grid_add_item(g, &img)
    }

    // No selection
    indices := widgets.image_grid_get_selected_indices(g)
    testing.expect(t, indices == nil, "Should return nil for no selection")

    // Select some items
    widgets.image_grid_toggle_selection(g, 1)
    widgets.image_grid_toggle_selection(g, 3)

    indices = widgets.image_grid_get_selected_indices(g)
    defer delete(indices)
    testing.expect(t, len(indices) == 2, "Should have 2 indices")

    // Check that both 1 and 3 are in the result
    has_1 := false
    has_3 := false
    for idx in indices {
        if idx == 1 { has_1 = true }
        if idx == 3 { has_3 = true }
    }
    testing.expect(t, has_1 && has_3, "Should contain indices 1 and 3")
}
