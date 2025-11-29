package tests

import "../src/core"
import "../src/render"
import "core:testing"

@(test)
test_context_create :: proc(t: ^testing.T) {
    pixels: [100 * 100]u32
    ctx := render.context_create(&pixels[0], 100, 100, 400)

    testing.expect(t, ctx.width == 100, "Width should be 100")
    testing.expect(t, ctx.height == 100, "Height should be 100")
    testing.expect(t, ctx.stride == 400, "Stride should be 400")
    testing.expect(t, ctx.clip.x == 0, "Clip x should be 0")
    testing.expect(t, ctx.clip.y == 0, "Clip y should be 0")
    testing.expect(t, ctx.clip.width == 100, "Clip width should be 100")
    testing.expect(t, ctx.clip.height == 100, "Clip height should be 100")
}

@(test)
test_clear :: proc(t: ^testing.T) {
    pixels: [10 * 10]u32
    ctx := render.context_create(&pixels[0], 10, 10, 40)

    render.clear(&ctx, core.COLOR_WHITE)

    // Check all pixels are white
    for i in 0 ..< 100 {
        testing.expect(t, pixels[i] == 0xFFFFFFFF, "All pixels should be white after clear")
    }
}

@(test)
test_fill_rect :: proc(t: ^testing.T) {
    pixels: [10 * 10]u32
    ctx := render.context_create(&pixels[0], 10, 10, 40)

    // Clear to black first
    render.clear(&ctx, core.COLOR_BLACK)

    // Fill a small rect with red
    render.fill_rect(&ctx, core.Rect{2, 2, 3, 3}, core.COLOR_RED)

    // Check that the rect area is red
    red_argb := core.color_to_argb(core.COLOR_RED)
    black_argb := core.color_to_argb(core.COLOR_BLACK)

    for y in 0 ..< 10 {
        for x in 0 ..< 10 {
            idx := y * 10 + x
            if x >= 2 && x < 5 && y >= 2 && y < 5 {
                testing.expect(t, pixels[idx] == red_argb, "Pixel inside rect should be red")
            } else {
                testing.expect(t, pixels[idx] == black_argb, "Pixel outside rect should be black")
            }
        }
    }
}

@(test)
test_fill_rect_clipping :: proc(t: ^testing.T) {
    pixels: [10 * 10]u32
    ctx := render.context_create(&pixels[0], 10, 10, 40)

    // Clear to black
    render.clear(&ctx, core.COLOR_BLACK)

    // Try to fill a rect that extends beyond the buffer
    render.fill_rect(&ctx, core.Rect{8, 8, 10, 10}, core.COLOR_RED)

    // Only the portion inside the buffer should be red
    red_argb := core.color_to_argb(core.COLOR_RED)
    black_argb := core.color_to_argb(core.COLOR_BLACK)

    for y in 0 ..< 10 {
        for x in 0 ..< 10 {
            idx := y * 10 + x
            if x >= 8 && y >= 8 {
                testing.expect(t, pixels[idx] == red_argb, "Clipped rect portion should be red")
            } else {
                testing.expect(t, pixels[idx] == black_argb, "Outside area should be black")
            }
        }
    }
}

@(test)
test_set_pixel :: proc(t: ^testing.T) {
    pixels: [10 * 10]u32
    ctx := render.context_create(&pixels[0], 10, 10, 40)

    render.clear(&ctx, core.COLOR_BLACK)
    render.set_pixel(&ctx, 5, 5, core.COLOR_GREEN)

    green_argb := core.color_to_argb(core.COLOR_GREEN)
    testing.expect(t, pixels[5 * 10 + 5] == green_argb, "Pixel at 5,5 should be green")

    // Check that setting pixel outside bounds doesn't crash
    render.set_pixel(&ctx, -1, -1, core.COLOR_RED)
    render.set_pixel(&ctx, 100, 100, core.COLOR_RED)
}

@(test)
test_draw_rect :: proc(t: ^testing.T) {
    pixels: [20 * 20]u32
    ctx := render.context_create(&pixels[0], 20, 20, 80)

    render.clear(&ctx, core.COLOR_BLACK)
    render.draw_rect(&ctx, core.Rect{5, 5, 10, 10}, core.COLOR_WHITE, 1)

    white_argb := core.color_to_argb(core.COLOR_WHITE)
    black_argb := core.color_to_argb(core.COLOR_BLACK)

    // Top edge
    for x in 5 ..< 15 {
        testing.expect(t, pixels[5 * 20 + x] == white_argb, "Top edge should be white")
    }

    // Bottom edge
    for x in 5 ..< 15 {
        testing.expect(t, pixels[14 * 20 + x] == white_argb, "Bottom edge should be white")
    }

    // Left edge
    for y in 6 ..< 14 {
        testing.expect(t, pixels[y * 20 + 5] == white_argb, "Left edge should be white")
    }

    // Right edge
    for y in 6 ..< 14 {
        testing.expect(t, pixels[y * 20 + 14] == white_argb, "Right edge should be white")
    }

    // Interior should be black
    testing.expect(t, pixels[7 * 20 + 7] == black_argb, "Interior should be black")
}
