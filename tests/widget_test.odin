package tests

import "../src/widgets"
import "../src/core"
import "core:testing"

// ============================================================================
// Widget creation and destruction
// ============================================================================

@(test)
test_widget_create_destroy :: proc(t: ^testing.T) {
    w := widgets.widget_create(&widgets.default_vtable)
    testing.expect(t, w != nil, "widget_create should return non-nil")
    testing.expect(t, w.visible == true, "Widget should be visible by default")
    testing.expect(t, w.enabled == true, "Widget should be enabled by default")
    testing.expect(t, w.dirty == true, "Widget should be dirty by default")

    widgets.widget_destroy(w)
    // No crash = success
}

@(test)
test_widget_destroy_nil :: proc(t: ^testing.T) {
    // Should not crash
    widgets.widget_destroy(nil)
}

// ============================================================================
// Parent/child relationships
// ============================================================================

@(test)
test_widget_add_child :: proc(t: ^testing.T) {
    parent := widgets.widget_create(&widgets.default_vtable)
    child := widgets.widget_create(&widgets.default_vtable)

    widgets.widget_add_child(parent, child)

    testing.expect(t, len(parent.children) == 1, "Parent should have 1 child")
    testing.expect(t, parent.children[0] == child, "Child should be in parent's children")
    testing.expect(t, child.parent == parent, "Child's parent should be set")

    widgets.widget_destroy(parent)
}

@(test)
test_widget_remove_child :: proc(t: ^testing.T) {
    parent := widgets.widget_create(&widgets.default_vtable)
    child := widgets.widget_create(&widgets.default_vtable)

    widgets.widget_add_child(parent, child)
    widgets.widget_remove_child(parent, child)

    testing.expect(t, len(parent.children) == 0, "Parent should have 0 children")
    testing.expect(t, child.parent == nil, "Child's parent should be nil")

    widgets.widget_destroy(parent)
    widgets.widget_destroy(child)
}

@(test)
test_widget_reparent :: proc(t: ^testing.T) {
    parent1 := widgets.widget_create(&widgets.default_vtable)
    parent2 := widgets.widget_create(&widgets.default_vtable)
    child := widgets.widget_create(&widgets.default_vtable)

    widgets.widget_add_child(parent1, child)
    widgets.widget_add_child(parent2, child)

    testing.expect(t, len(parent1.children) == 0, "Old parent should have 0 children")
    testing.expect(t, len(parent2.children) == 1, "New parent should have 1 child")
    testing.expect(t, child.parent == parent2, "Child's parent should be new parent")

    widgets.widget_destroy(parent1)
    widgets.widget_destroy(parent2)
}

@(test)
test_widget_destroy_with_children :: proc(t: ^testing.T) {
    parent := widgets.widget_create(&widgets.default_vtable)
    child1 := widgets.widget_create(&widgets.default_vtable)
    child2 := widgets.widget_create(&widgets.default_vtable)

    widgets.widget_add_child(parent, child1)
    widgets.widget_add_child(parent, child2)

    // Destroying parent should destroy children too
    widgets.widget_destroy(parent)
    // No crash = success
}

// ============================================================================
// Hit testing
// ============================================================================

@(test)
test_hit_test_simple :: proc(t: ^testing.T) {
    root := widgets.widget_create(&widgets.default_vtable)
    root.rect = core.Rect{0, 0, 100, 100}

    // Point inside
    hit := widgets.hit_test(root, 50, 50)
    testing.expect(t, hit == root, "Hit test should find root widget")

    // Point outside
    hit = widgets.hit_test(root, 150, 50)
    testing.expect(t, hit == nil, "Hit test should return nil for point outside")

    widgets.widget_destroy(root)
}

@(test)
test_hit_test_nested :: proc(t: ^testing.T) {
    root := widgets.widget_create(&widgets.default_vtable)
    root.rect = core.Rect{0, 0, 200, 200}

    child := widgets.widget_create(&widgets.default_vtable)
    child.rect = core.Rect{50, 50, 100, 100}
    widgets.widget_add_child(root, child)

    // Point in child
    hit := widgets.hit_test(root, 75, 75)
    testing.expect(t, hit == child, "Hit test should find child widget")

    // Point in root but not child
    hit = widgets.hit_test(root, 25, 25)
    testing.expect(t, hit == root, "Hit test should find root when not in child")

    widgets.widget_destroy(root)
}

@(test)
test_hit_test_overlapping :: proc(t: ^testing.T) {
    root := widgets.widget_create(&widgets.default_vtable)
    root.rect = core.Rect{0, 0, 200, 200}

    // Two overlapping children - second should be hit first (on top)
    child1 := widgets.widget_create(&widgets.default_vtable)
    child1.rect = core.Rect{25, 25, 100, 100}
    widgets.widget_add_child(root, child1)

    child2 := widgets.widget_create(&widgets.default_vtable)
    child2.rect = core.Rect{50, 50, 100, 100}
    widgets.widget_add_child(root, child2)

    // Point in both children - should hit child2 (later in list = on top)
    hit := widgets.hit_test(root, 75, 75)
    testing.expect(t, hit == child2, "Hit test should find top-most (later) child")

    // Point only in child1
    hit = widgets.hit_test(root, 30, 30)
    testing.expect(t, hit == child1, "Hit test should find child1 in its area")

    widgets.widget_destroy(root)
}

@(test)
test_hit_test_invisible :: proc(t: ^testing.T) {
    root := widgets.widget_create(&widgets.default_vtable)
    root.rect = core.Rect{0, 0, 100, 100}
    root.visible = false

    hit := widgets.hit_test(root, 50, 50)
    testing.expect(t, hit == nil, "Hit test should not find invisible widget")

    widgets.widget_destroy(root)
}

// ============================================================================
// Focus management
// ============================================================================

@(test)
test_focus_set_clear :: proc(t: ^testing.T) {
    root := widgets.widget_create(&widgets.default_vtable)
    child := widgets.widget_create(&widgets.default_vtable)
    child.focusable = true
    widgets.widget_add_child(root, child)

    fm := widgets.focus_manager_init(root)

    widgets.focus_set(&fm, child)
    testing.expect(t, widgets.focus_get(&fm) == child, "Focus should be set to child")
    testing.expect(t, child.focused == true, "Child's focused flag should be true")

    widgets.focus_clear(&fm)
    testing.expect(t, widgets.focus_get(&fm) == nil, "Focus should be cleared")
    testing.expect(t, child.focused == false, "Child's focused flag should be false")

    widgets.widget_destroy(root)
}

@(test)
test_focus_next_prev :: proc(t: ^testing.T) {
    root := widgets.widget_create(&widgets.default_vtable)

    child1 := widgets.widget_create(&widgets.default_vtable)
    child1.focusable = true
    widgets.widget_add_child(root, child1)

    child2 := widgets.widget_create(&widgets.default_vtable)
    child2.focusable = true
    widgets.widget_add_child(root, child2)

    child3 := widgets.widget_create(&widgets.default_vtable)
    child3.focusable = true
    widgets.widget_add_child(root, child3)

    fm := widgets.focus_manager_init(root)

    // Start with no focus, focus_next should go to first
    widgets.focus_next(&fm)
    testing.expect(t, widgets.focus_get(&fm) == child1, "First focus_next should go to child1")

    widgets.focus_next(&fm)
    testing.expect(t, widgets.focus_get(&fm) == child2, "Second focus_next should go to child2")

    widgets.focus_next(&fm)
    testing.expect(t, widgets.focus_get(&fm) == child3, "Third focus_next should go to child3")

    widgets.focus_next(&fm)
    testing.expect(t, widgets.focus_get(&fm) == child1, "Fourth focus_next should wrap to child1")

    widgets.focus_prev(&fm)
    testing.expect(t, widgets.focus_get(&fm) == child3, "focus_prev should go to child3")

    widgets.widget_destroy(root)
}

@(test)
test_focus_first_last :: proc(t: ^testing.T) {
    root := widgets.widget_create(&widgets.default_vtable)

    child1 := widgets.widget_create(&widgets.default_vtable)
    child1.focusable = true
    widgets.widget_add_child(root, child1)

    child2 := widgets.widget_create(&widgets.default_vtable)
    child2.focusable = true
    widgets.widget_add_child(root, child2)

    fm := widgets.focus_manager_init(root)

    widgets.focus_first(&fm)
    testing.expect(t, widgets.focus_get(&fm) == child1, "focus_first should go to child1")

    widgets.focus_last(&fm)
    testing.expect(t, widgets.focus_get(&fm) == child2, "focus_last should go to child2")

    widgets.widget_destroy(root)
}

// ============================================================================
// Container creation
// ============================================================================

@(test)
test_container_create :: proc(t: ^testing.T) {
    c := widgets.container_create(.Row)
    testing.expect(t, c != nil, "container_create should return non-nil")
    testing.expect(t, c.direction == .Row, "Container direction should be Row")
    testing.expect(t, c.vtable == &widgets.container_vtable, "Container should have container vtable")

    widgets.widget_destroy(c)
}

@(test)
test_container_default_column :: proc(t: ^testing.T) {
    c := widgets.container_create()
    testing.expect(t, c.direction == .Column, "Default container direction should be Column")

    widgets.widget_destroy(c)
}

// ============================================================================
// Container layout - Row
// ============================================================================

@(test)
test_container_layout_row :: proc(t: ^testing.T) {
    c := widgets.container_create(.Row)
    c.rect = core.Rect{0, 0, 300, 100}

    child1 := widgets.widget_create(&widgets.default_vtable)
    child1.min_size = core.Size{50, 50}
    widgets.widget_add_child(c, child1)

    child2 := widgets.widget_create(&widgets.default_vtable)
    child2.min_size = core.Size{50, 50}
    widgets.widget_add_child(c, child2)

    widgets.container_layout(c)

    // Children should be positioned side by side
    testing.expect(t, child1.rect.x == 0, "Child1 x should be 0")
    testing.expect(t, child2.rect.x == 50, "Child2 x should be 50")
    testing.expect(t, child1.rect.width == 50, "Child1 width should be 50")
    testing.expect(t, child2.rect.width == 50, "Child2 width should be 50")

    widgets.widget_destroy(c)
}

@(test)
test_container_layout_row_with_spacing :: proc(t: ^testing.T) {
    c := widgets.container_create(.Row)
    c.rect = core.Rect{0, 0, 300, 100}
    c.spacing = 10

    child1 := widgets.widget_create(&widgets.default_vtable)
    child1.min_size = core.Size{50, 50}
    widgets.widget_add_child(c, child1)

    child2 := widgets.widget_create(&widgets.default_vtable)
    child2.min_size = core.Size{50, 50}
    widgets.widget_add_child(c, child2)

    widgets.container_layout(c)

    testing.expect(t, child1.rect.x == 0, "Child1 x should be 0")
    testing.expect(t, child2.rect.x == 60, "Child2 x should be 60 (50 + 10 spacing)")

    widgets.widget_destroy(c)
}

@(test)
test_container_layout_row_flex :: proc(t: ^testing.T) {
    c := widgets.container_create(.Row)
    c.rect = core.Rect{0, 0, 300, 100}

    // Fixed child
    child1 := widgets.widget_create(&widgets.default_vtable)
    child1.min_size = core.Size{50, 50}
    child1.flex = 0
    widgets.widget_add_child(c, child1)

    // Flex child should take remaining space
    child2 := widgets.widget_create(&widgets.default_vtable)
    child2.flex = 1
    widgets.widget_add_child(c, child2)

    widgets.container_layout(c)

    testing.expect(t, child1.rect.width == 50, "Fixed child width should be 50")
    testing.expect(t, child2.rect.width == 250, "Flex child width should be 250")
    testing.expect(t, child2.rect.x == 50, "Flex child x should be 50")

    widgets.widget_destroy(c)
}

// ============================================================================
// Container layout - Column
// ============================================================================

@(test)
test_container_layout_column :: proc(t: ^testing.T) {
    c := widgets.container_create(.Column)
    c.rect = core.Rect{0, 0, 100, 300}

    child1 := widgets.widget_create(&widgets.default_vtable)
    child1.min_size = core.Size{50, 50}
    widgets.widget_add_child(c, child1)

    child2 := widgets.widget_create(&widgets.default_vtable)
    child2.min_size = core.Size{50, 50}
    widgets.widget_add_child(c, child2)

    widgets.container_layout(c)

    // Children should be stacked vertically
    testing.expect(t, child1.rect.y == 0, "Child1 y should be 0")
    testing.expect(t, child2.rect.y == 50, "Child2 y should be 50")
    testing.expect(t, child1.rect.height == 50, "Child1 height should be 50")
    testing.expect(t, child2.rect.height == 50, "Child2 height should be 50")

    widgets.widget_destroy(c)
}

// ============================================================================
// Container layout - Alignment
// ============================================================================

@(test)
test_container_align_center :: proc(t: ^testing.T) {
    c := widgets.container_create(.Row)
    c.rect = core.Rect{0, 0, 300, 100}
    c.align_items = .Center

    child := widgets.widget_create(&widgets.default_vtable)
    child.min_size = core.Size{50, 40}
    widgets.widget_add_child(c, child)

    widgets.container_layout(c)

    // Child should be centered on cross axis
    expected_y: i32 = (100 - 40) / 2
    testing.expect(t, child.rect.y == expected_y, "Child should be vertically centered")

    widgets.widget_destroy(c)
}

@(test)
test_container_align_end :: proc(t: ^testing.T) {
    c := widgets.container_create(.Row)
    c.rect = core.Rect{0, 0, 300, 100}
    c.align_items = .End

    child := widgets.widget_create(&widgets.default_vtable)
    child.min_size = core.Size{50, 40}
    widgets.widget_add_child(c, child)

    widgets.container_layout(c)

    // Child should be at the end (bottom)
    expected_y: i32 = 100 - 40
    testing.expect(t, child.rect.y == expected_y, "Child should be at bottom")

    widgets.widget_destroy(c)
}

@(test)
test_container_align_stretch :: proc(t: ^testing.T) {
    c := widgets.container_create(.Row)
    c.rect = core.Rect{0, 0, 300, 100}
    c.align_items = .Stretch

    child := widgets.widget_create(&widgets.default_vtable)
    child.min_size = core.Size{50, 40}
    widgets.widget_add_child(c, child)

    widgets.container_layout(c)

    // Child should stretch to fill cross axis
    testing.expect(t, child.rect.y == 0, "Child y should be 0")
    testing.expect(t, child.rect.height == 100, "Child should stretch to container height")

    widgets.widget_destroy(c)
}

// ============================================================================
// Dirty marking
// ============================================================================

@(test)
test_widget_mark_dirty_propagates :: proc(t: ^testing.T) {
    root := widgets.widget_create(&widgets.default_vtable)
    root.dirty = false

    child := widgets.widget_create(&widgets.default_vtable)
    child.dirty = false
    widgets.widget_add_child(root, child)

    grandchild := widgets.widget_create(&widgets.default_vtable)
    grandchild.dirty = false
    widgets.widget_add_child(child, grandchild)

    // Marking grandchild dirty should propagate up
    widgets.widget_mark_dirty(grandchild)

    testing.expect(t, grandchild.dirty == true, "Grandchild should be dirty")
    testing.expect(t, child.dirty == true, "Child should be dirty")
    testing.expect(t, root.dirty == true, "Root should be dirty")

    widgets.widget_destroy(root)
}

// ============================================================================
// Absolute rect calculation
// ============================================================================

@(test)
test_widget_absolute_rect :: proc(t: ^testing.T) {
    root := widgets.widget_create(&widgets.default_vtable)
    root.rect = core.Rect{10, 10, 200, 200}
    root.padding = widgets.Edges{5, 5, 5, 5}

    child := widgets.widget_create(&widgets.default_vtable)
    child.rect = core.Rect{20, 20, 50, 50}
    widgets.widget_add_child(root, child)

    abs := widgets.widget_get_absolute_rect(child)

    // Child at (20, 20) relative to parent content area
    // Parent at (10, 10) with 5px padding
    // Child absolute: (10 + 5 + 20, 10 + 5 + 20) = (35, 35)
    testing.expect(t, abs.x == 35, "Absolute x should be 35")
    testing.expect(t, abs.y == 35, "Absolute y should be 35")
    testing.expect(t, abs.width == 50, "Width should be preserved")
    testing.expect(t, abs.height == 50, "Height should be preserved")

    widgets.widget_destroy(root)
}

// ============================================================================
// Helper functions
// ============================================================================

@(test)
test_edges_all :: proc(t: ^testing.T) {
    e := widgets.edges_all(10)
    testing.expect(t, e.top == 10, "Top should be 10")
    testing.expect(t, e.right == 10, "Right should be 10")
    testing.expect(t, e.bottom == 10, "Bottom should be 10")
    testing.expect(t, e.left == 10, "Left should be 10")
}

@(test)
test_edges_symmetric :: proc(t: ^testing.T) {
    e := widgets.edges_symmetric(10, 20)
    testing.expect(t, e.top == 20, "Top should be 20")
    testing.expect(t, e.right == 10, "Right should be 10")
    testing.expect(t, e.bottom == 20, "Bottom should be 20")
    testing.expect(t, e.left == 10, "Left should be 10")
}
