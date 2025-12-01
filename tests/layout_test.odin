package tests

import "../src/core"
import "../src/render"
import "../src/widgets"
import "core:fmt"
import "core:testing"

// ============================================================================
// Todo App Footer Layout Test
// ============================================================================
// This test replicates the exact widget hierarchy of the todo app to verify
// that the footer layout is calculated correctly.

@(test)
test_todo_app_footer_layout :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 14)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    // Create root container (Column)
    root := widgets.container_create(.Column)
    root.rect = core.Rect{0, 0, 300, 500}
    root.padding = widgets.edges_all(15)
    root.spacing = 15
    root.align_items = .Stretch

    // Create header (Row, min_size={0, 40})
    header := widgets.container_create(.Row)
    header.min_size = core.Size{0, 40}
    widgets.widget_add_child(root, header)

    // Create input_row (Row, min_size={0, 36})
    input_row := widgets.container_create(.Row)
    input_row.min_size = core.Size{0, 36}
    widgets.widget_add_child(root, input_row)

    // Create scroll container (flex=1)
    scroll := widgets.scroll_container_create()
    scroll.flex = 1
    widgets.widget_add_child(root, scroll)

    // Create footer (Column, spacing=10)
    footer := widgets.container_create(.Column)
    footer.spacing = 10
    widgets.widget_add_child(root, footer)

    // Create filter_row inside footer (Row)
    filter_row := widgets.container_create(.Row)
    filter_row.spacing = 5
    widgets.widget_add_child(footer, filter_row)

    // Create 4 buttons in filter_row (min_size={70, 28})
    for i := 0; i < 4; i += 1 {
        btn := widgets.button_create("Btn")
        btn.min_size = core.Size{70, 28}
        widgets.widget_add_child(filter_row, btn)
    }

    // Create count_row inside footer (Row)
    count_row := widgets.container_create(.Row)
    widgets.widget_add_child(footer, count_row)

    // Create label in count_row
    label := widgets.label_create("10 items left", &font)
    label.wrap = false
    widgets.widget_add_child(count_row, label)

    // Run layout
    widgets.container_layout(root)

    // Assertions
    // 1. Header should get its min height
    testing.expect(t, header.rect.height >= 40,
        fmt.tprintf("Header height (%d) should be >= 40", header.rect.height))

    // 2. Input row should get its min height
    testing.expect(t, input_row.rect.height >= 36,
        fmt.tprintf("Input row height (%d) should be >= 36", input_row.rect.height))

    // 3. Footer should fit within window (critical assertion)
    footer_bottom := footer.rect.y + footer.rect.height
    content_bottom: i32 = 500 - 15  // window height - bottom padding
    testing.expect(t, footer_bottom <= content_bottom,
        fmt.tprintf("Footer bottom (%d) should be <= content bottom (%d)", footer_bottom, content_bottom))

    // 4. Count row should fit within footer
    count_row_bottom := count_row.rect.y + count_row.rect.height
    testing.expect(t, count_row_bottom <= footer.rect.height,
        fmt.tprintf("Count row bottom (%d) should be <= footer height (%d)", count_row_bottom, footer.rect.height))

    // 5. Scroll should have gotten the flex space (positive height)
    testing.expect(t, scroll.rect.height > 0,
        fmt.tprintf("Scroll height (%d) should be > 0", scroll.rect.height))

    widgets.widget_destroy(root)
}

// ============================================================================
// Simpler test - just fixed children (no flex)
// ============================================================================

@(test)
test_column_layout_fixed_children :: proc(t: ^testing.T) {
    root := widgets.container_create(.Column)
    root.rect = core.Rect{0, 0, 300, 500}
    root.padding = widgets.edges_all(15)
    root.spacing = 15
    root.align_items = .Stretch

    child1 := widgets.container_create(.Row)
    child1.min_size = core.Size{0, 40}
    widgets.widget_add_child(root, child1)

    child2 := widgets.container_create(.Row)
    child2.min_size = core.Size{0, 36}
    widgets.widget_add_child(root, child2)

    child3 := widgets.container_create(.Row)
    child3.min_size = core.Size{0, 100}
    widgets.widget_add_child(root, child3)

    child4 := widgets.container_create(.Row)
    child4.min_size = core.Size{0, 56}
    widgets.widget_add_child(root, child4)

    widgets.container_layout(root)

    // All children should get their min heights
    testing.expect(t, child1.rect.height == 40,
        fmt.tprintf("Child1 height (%d) should be 40", child1.rect.height))
    testing.expect(t, child2.rect.height == 36,
        fmt.tprintf("Child2 height (%d) should be 36", child2.rect.height))
    testing.expect(t, child3.rect.height == 100,
        fmt.tprintf("Child3 height (%d) should be 100", child3.rect.height))
    testing.expect(t, child4.rect.height == 56,
        fmt.tprintf("Child4 height (%d) should be 56", child4.rect.height))

    // Child4 should be positioned correctly
    expected_y4: i32 = 40 + 15 + 36 + 15 + 100 + 15  // 221
    testing.expect(t, child4.rect.y == expected_y4,
        fmt.tprintf("Child4 Y (%d) should be %d", child4.rect.y, expected_y4))

    // Child4 bottom should be within content area
    child4_bottom := child4.rect.y + child4.rect.height
    content_bottom: i32 = 470  // 500 - 30 padding
    testing.expect(t, child4_bottom <= content_bottom,
        fmt.tprintf("Child4 bottom (%d) should be <= content bottom (%d)", child4_bottom, content_bottom))

    widgets.widget_destroy(root)
}

// ============================================================================
// Test with one flex child
// ============================================================================

@(test)
test_column_layout_with_flex :: proc(t: ^testing.T) {
    root := widgets.container_create(.Column)
    root.rect = core.Rect{0, 0, 300, 500}
    root.padding = widgets.edges_all(15)
    root.spacing = 15
    root.align_items = .Stretch

    top := widgets.container_create(.Row)
    top.min_size = core.Size{0, 40}
    widgets.widget_add_child(root, top)

    middle := widgets.container_create(.Column)
    middle.flex = 1
    widgets.widget_add_child(root, middle)

    bottom := widgets.container_create(.Row)
    bottom.min_size = core.Size{0, 56}
    widgets.widget_add_child(root, bottom)

    widgets.container_layout(root)

    // Content height: 470, Spacing: 30, Available: 440
    // Fixed: 40 + 56 = 96, Flex space: 344
    expected_flex_height: i32 = 470 - 30 - 40 - 56  // 344
    testing.expect(t, middle.rect.height == expected_flex_height,
        fmt.tprintf("Middle height (%d) should be %d", middle.rect.height, expected_flex_height))

    // Bottom should be at correct position
    expected_bottom_y: i32 = 40 + 15 + expected_flex_height + 15  // 414
    testing.expect(t, bottom.rect.y == expected_bottom_y,
        fmt.tprintf("Bottom Y (%d) should be %d", bottom.rect.y, expected_bottom_y))

    // Bottom's bottom edge should be at content edge
    bottom_bottom := bottom.rect.y + bottom.rect.height
    content_bottom: i32 = 470
    testing.expect(t, bottom_bottom == content_bottom,
        fmt.tprintf("Bottom's bottom edge (%d) should be %d", bottom_bottom, content_bottom))

    widgets.widget_destroy(root)
}

// ============================================================================
// Test nested containers (closer to real app)
// ============================================================================

@(test)
test_nested_container_measurement :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 14)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    footer := widgets.container_create(.Column)
    footer.spacing = 10

    filter_row := widgets.container_create(.Row)
    filter_row.spacing = 5
    widgets.widget_add_child(footer, filter_row)

    for i := 0; i < 4; i += 1 {
        btn := widgets.button_create("Btn")
        btn.min_size = core.Size{70, 28}
        widgets.widget_add_child(filter_row, btn)
    }

    count_row := widgets.container_create(.Row)
    widgets.widget_add_child(footer, count_row)

    label := widgets.label_create("10 items left", &font)
    label.wrap = false
    widgets.widget_add_child(count_row, label)

    footer_size := widgets.widget_measure(footer, 270)
    filter_row_size := widgets.widget_measure(filter_row, 270)
    count_row_size := widgets.widget_measure(count_row, 270)

    // Buttons have default padding of 8 top + 8 bottom, text height 16
    // So button content height = 16 + 16 = 32 > min_size.height of 28
    testing.expect(t, filter_row_size.height == 32,
        fmt.tprintf("Filter row height (%d) should be 32", filter_row_size.height))

    expected_footer_height := filter_row_size.height + 10 + count_row_size.height
    testing.expect(t, footer_size.height == expected_footer_height,
        fmt.tprintf("Footer height (%d) should be %d", footer_size.height, expected_footer_height))

    widgets.widget_destroy(footer)
}
