package tests

import "../src/core"
import "core:testing"

@(test)
test_rect_contains :: proc(t: ^testing.T) {
    r := core.Rect{10, 10, 100, 100}

    // Point inside
    testing.expect(t, core.rect_contains(r, core.Point{50, 50}), "Point 50,50 should be inside rect")
    testing.expect(t, core.rect_contains(r, core.Point{10, 10}), "Point at top-left corner should be inside")

    // Point outside
    testing.expect(t, !core.rect_contains(r, core.Point{0, 0}), "Point 0,0 should be outside rect")
    testing.expect(t, !core.rect_contains(r, core.Point{110, 110}), "Point 110,110 should be outside rect")
    testing.expect(t, !core.rect_contains(r, core.Point{200, 50}), "Point 200,50 should be outside rect")
}

@(test)
test_rect_intersects :: proc(t: ^testing.T) {
    r1 := core.Rect{0, 0, 100, 100}
    r2 := core.Rect{50, 50, 100, 100}
    r3 := core.Rect{200, 200, 50, 50}

    testing.expect(t, core.rect_intersects(r1, r2), "Overlapping rects should intersect")
    testing.expect(t, !core.rect_intersects(r1, r3), "Non-overlapping rects should not intersect")
}

@(test)
test_rect_intersection :: proc(t: ^testing.T) {
    r1 := core.Rect{0, 0, 100, 100}
    r2 := core.Rect{50, 50, 100, 100}

    result, ok := core.rect_intersection(r1, r2)
    testing.expect(t, ok, "Intersection should succeed")
    testing.expect(t, result.x == 50, "Intersection x should be 50")
    testing.expect(t, result.y == 50, "Intersection y should be 50")
    testing.expect(t, result.width == 50, "Intersection width should be 50")
    testing.expect(t, result.height == 50, "Intersection height should be 50")

    // Non-intersecting
    r3 := core.Rect{200, 200, 50, 50}
    _, ok2 := core.rect_intersection(r1, r3)
    testing.expect(t, !ok2, "Non-overlapping rects should return false")
}

@(test)
test_rect_union :: proc(t: ^testing.T) {
    r1 := core.Rect{0, 0, 50, 50}
    r2 := core.Rect{50, 50, 50, 50}

    result := core.rect_union(r1, r2)
    testing.expect(t, result.x == 0, "Union x should be 0")
    testing.expect(t, result.y == 0, "Union y should be 0")
    testing.expect(t, result.width == 100, "Union width should be 100")
    testing.expect(t, result.height == 100, "Union height should be 100")
}

@(test)
test_rect_is_empty :: proc(t: ^testing.T) {
    empty1 := core.Rect{0, 0, 0, 0}
    empty2 := core.Rect{10, 10, 0, 50}
    empty3 := core.Rect{10, 10, 50, 0}
    not_empty := core.Rect{0, 0, 50, 50}

    testing.expect(t, core.rect_is_empty(empty1), "0x0 rect should be empty")
    testing.expect(t, core.rect_is_empty(empty2), "0-width rect should be empty")
    testing.expect(t, core.rect_is_empty(empty3), "0-height rect should be empty")
    testing.expect(t, !core.rect_is_empty(not_empty), "50x50 rect should not be empty")
}

@(test)
test_color_to_argb :: proc(t: ^testing.T) {
    white := core.COLOR_WHITE
    argb := core.color_to_argb(white)
    testing.expect(t, argb == 0xFFFFFFFF, "White should be 0xFFFFFFFF")

    black := core.COLOR_BLACK
    argb2 := core.color_to_argb(black)
    testing.expect(t, argb2 == 0xFF000000, "Black should be 0xFF000000")

    transparent := core.COLOR_TRANSPARENT
    argb3 := core.color_to_argb(transparent)
    testing.expect(t, argb3 == 0x00000000, "Transparent should be 0x00000000")
}

@(test)
test_color_hex :: proc(t: ^testing.T) {
    red := core.color_hex(0xFF0000)
    testing.expect(t, red.r == 255, "Red component should be 255")
    testing.expect(t, red.g == 0, "Green component should be 0")
    testing.expect(t, red.b == 0, "Blue component should be 0")
    testing.expect(t, red.a == 255, "Alpha should be 255")
}

@(test)
test_color_rgb :: proc(t: ^testing.T) {
    c := core.color_rgb(128, 64, 32)
    testing.expect(t, c.r == 128, "Red should be 128")
    testing.expect(t, c.g == 64, "Green should be 64")
    testing.expect(t, c.b == 32, "Blue should be 32")
    testing.expect(t, c.a == 255, "Alpha should be 255")
}
