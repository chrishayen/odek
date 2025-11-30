package tests

import "../src/widgets"
import "../src/core"
import "core:testing"

// ============================================================================
// ScrollContainer creation and destruction
// ============================================================================

@(test)
test_scroll_container_create :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()
    testing.expect(t, sc != nil, "scroll_container_create should return non-nil")
    testing.expect(t, sc.visible == true, "ScrollContainer should be visible by default")
    testing.expect(t, sc.enabled == true, "ScrollContainer should be enabled by default")
    testing.expect(t, sc.direction == .Vertical, "Default direction should be Vertical")
    testing.expect(t, sc.scrollbar_width > 0, "Scrollbar width should be set")

    widgets.widget_destroy(sc)
}

@(test)
test_scroll_container_create_horizontal :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create(.Horizontal)
    testing.expect(t, sc.direction == .Horizontal, "Direction should be Horizontal")

    widgets.widget_destroy(sc)
}

// ============================================================================
// Content management
// ============================================================================

@(test)
test_scroll_container_set_content :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()
    content := widgets.container_create(.Column)

    widgets.scroll_container_set_content(sc, content)

    testing.expect(t, len(sc.children) == 1, "ScrollContainer should have 1 child")
    testing.expect(t, sc.children[0] == content, "Content should be the child")
    testing.expect(t, widgets.scroll_container_get_content(sc) == content, "get_content should return content")

    widgets.widget_destroy(sc)
}

@(test)
test_scroll_container_replace_content :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()
    content1 := widgets.container_create(.Column)
    content2 := widgets.container_create(.Row)

    widgets.scroll_container_set_content(sc, content1)
    widgets.scroll_container_set_content(sc, content2)

    testing.expect(t, len(sc.children) == 1, "ScrollContainer should still have 1 child")
    testing.expect(t, widgets.scroll_container_get_content(sc) == content2, "Content should be replaced")

    widgets.widget_destroy(sc)
    widgets.widget_destroy(content1)  // Not destroyed since it was removed
}

@(test)
test_scroll_container_no_content :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()

    testing.expect(t, widgets.scroll_container_get_content(sc) == nil, "get_content should return nil")

    widgets.widget_destroy(sc)
}

// ============================================================================
// Scroll operations
// ============================================================================

@(test)
test_scroll_container_scroll_to_top :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()
    sc.scroll.offset = 100

    widgets.scroll_container_scroll_to_top(sc)

    testing.expect(t, sc.scroll.offset == 0, "Offset should be 0 after scroll_to_top")

    widgets.widget_destroy(sc)
}

@(test)
test_scroll_container_scroll_to_bottom :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()
    sc.scroll.content_size = 500
    sc.scroll.viewport_size = 100

    widgets.scroll_container_scroll_to_bottom(sc)

    testing.expect(t, sc.scroll.offset == 400, "Offset should be at max after scroll_to_bottom")

    widgets.widget_destroy(sc)
}

// ============================================================================
// Viewport calculation
// ============================================================================

@(test)
test_scroll_container_viewport :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()
    sc.rect = core.Rect{0, 0, 200, 300}
    sc.padding = widgets.edges_all(10)

    viewport := widgets.scroll_container_get_viewport(sc)

    testing.expect(t, viewport.x == 10, "Viewport x should account for padding")
    testing.expect(t, viewport.y == 10, "Viewport y should account for padding")
    testing.expect(t, viewport.width <= 180, "Viewport width should account for padding")
    testing.expect(t, viewport.height == 280, "Viewport height should account for padding")

    widgets.widget_destroy(sc)
}

// ============================================================================
// Scrollbar rects
// ============================================================================

@(test)
test_scroll_container_track_rect :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()
    sc.rect = core.Rect{0, 0, 200, 300}
    sc.padding = widgets.edges_all(10)

    track := widgets.scroll_container_get_track_rect(sc)

    testing.expect(t, track.width == sc.scrollbar_width, "Track width should be scrollbar_width")
    testing.expect(t, track.height == 280, "Track height should match viewport")

    widgets.widget_destroy(sc)
}

@(test)
test_scroll_container_thumb_rect :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()
    sc.rect = core.Rect{0, 0, 200, 300}
    sc.padding = widgets.edges_all(10)
    sc.scroll.content_size = 600
    sc.scroll.viewport_size = 280
    sc.scroll.offset = 0

    thumb := widgets.scroll_container_get_thumb_rect(sc)

    testing.expect(t, thumb.width == sc.scrollbar_width, "Thumb width should be scrollbar_width")
    testing.expect(t, thumb.height > 0, "Thumb height should be > 0")
    testing.expect(t, thumb.height < 280, "Thumb height should be less than track")

    widgets.widget_destroy(sc)
}

// ============================================================================
// Point tests
// ============================================================================

@(test)
test_scroll_container_point_in_scrollbar :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()
    sc.rect = core.Rect{0, 0, 200, 300}
    sc.scrollbar_width = 10
    sc.scrollbar_padding = 2

    // Point in scrollbar area (near right edge)
    in_sb := widgets.scroll_container_point_in_scrollbar(sc, 195, 50)
    testing.expect(t, in_sb == true, "Point should be in scrollbar area")

    // Point outside scrollbar
    out_sb := widgets.scroll_container_point_in_scrollbar(sc, 50, 50)
    testing.expect(t, out_sb == false, "Point should not be in scrollbar area")

    widgets.widget_destroy(sc)
}

// ============================================================================
// Theme colors
// ============================================================================

@(test)
test_scroll_container_theme_colors :: proc(t: ^testing.T) {
    sc := widgets.scroll_container_create()

    // ScrollContainer should use theme colors
    theme := widgets.theme_get()
    testing.expect(t, sc.track_color == theme.scrollbar_track, "Track color should match theme")
    testing.expect(t, sc.thumb_color == theme.scrollbar_thumb, "Thumb color should match theme")
    testing.expect(t, sc.thumb_hover_color == theme.scrollbar_hover, "Hover color should match theme")

    widgets.widget_destroy(sc)
}
