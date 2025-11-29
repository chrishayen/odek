package render

import "../core"
import "base:intrinsics"
import "core:slice"

// Draw context for software rendering
Draw_Context :: struct {
    pixels: [^]u32,
    width: i32,
    height: i32,
    stride: i32,
    clip: core.Rect,
}

// Create a draw context
context_create :: proc(pixels: [^]u32, width, height, stride: i32) -> Draw_Context {
    return Draw_Context{
        pixels = pixels,
        width = width,
        height = height,
        stride = stride,
        clip = core.Rect{0, 0, width, height},
    }
}

// Set clipping rectangle
context_set_clip :: proc(ctx: ^Draw_Context, clip: core.Rect) {
    // Intersect with buffer bounds
    bounds := core.Rect{0, 0, ctx.width, ctx.height}
    if clipped, ok := core.rect_intersection(clip, bounds); ok {
        ctx.clip = clipped
    } else {
        ctx.clip = core.Rect{}
    }
}

// Clear the entire buffer
clear :: proc(ctx: ^Draw_Context, color: core.Color) {
    fill_rect(ctx, core.Rect{0, 0, ctx.width, ctx.height}, color)
}

// Fill a rectangle with a solid color (optimized)
fill_rect :: proc(ctx: ^Draw_Context, rect: core.Rect, color: core.Color) {
    clipped, ok := core.rect_intersection(rect, ctx.clip)
    if !ok {
        return
    }

    pixel := core.color_to_argb(color)
    stride_pixels := ctx.stride / 4
    w := int(clipped.width)
    h := clipped.height

    if w <= 0 || h <= 0 {
        return
    }

    // Fill first row using slice.fill
    first_row := clipped.y * stride_pixels + clipped.x
    first_row_slice := ctx.pixels[first_row:][:w]
    slice.fill(first_row_slice, pixel)

    // Copy first row to remaining rows
    row_bytes := w * size_of(u32)
    src := rawptr(&ctx.pixels[first_row])
    for y in 1 ..< h {
        dst_offset := (clipped.y + y) * stride_pixels + clipped.x
        intrinsics.mem_copy(rawptr(&ctx.pixels[dst_offset]), src, row_bytes)
    }
}

// Fill a rectangle with alpha blending
fill_rect_blend :: proc(ctx: ^Draw_Context, rect: core.Rect, color: core.Color) {
    if color.a == 255 {
        fill_rect(ctx, rect, color)
        return
    }
    if color.a == 0 {
        return
    }

    clipped, ok := core.rect_intersection(rect, ctx.clip)
    if !ok {
        return
    }

    stride_pixels := ctx.stride / 4

    for y in clipped.y ..< clipped.y + clipped.height {
        row_start := y * stride_pixels + clipped.x
        for x in 0 ..< clipped.width {
            dst := ctx.pixels[row_start + x]
            ctx.pixels[row_start + x] = blend_pixel(dst, color)
        }
    }
}

// Draw a rectangle outline
draw_rect :: proc(ctx: ^Draw_Context, rect: core.Rect, color: core.Color, thickness: i32 = 1) {
    // Top
    fill_rect(ctx, core.Rect{rect.x, rect.y, rect.width, thickness}, color)
    // Bottom
    fill_rect(ctx, core.Rect{rect.x, rect.y + rect.height - thickness, rect.width, thickness}, color)
    // Left
    fill_rect(ctx, core.Rect{rect.x, rect.y + thickness, thickness, rect.height - 2 * thickness}, color)
    // Right
    fill_rect(ctx, core.Rect{rect.x + rect.width - thickness, rect.y + thickness, thickness, rect.height - 2 * thickness}, color)
}

// Draw a rounded rectangle (filled)
fill_rounded_rect :: proc(ctx: ^Draw_Context, rect: core.Rect, radius: i32, color: core.Color) {
    if radius <= 0 {
        fill_rect(ctx, rect, color)
        return
    }

    r := min(radius, rect.width / 2, rect.height / 2)

    // Center rectangle
    fill_rect(ctx, core.Rect{rect.x + r, rect.y, rect.width - 2 * r, rect.height}, color)
    // Left strip
    fill_rect(ctx, core.Rect{rect.x, rect.y + r, r, rect.height - 2 * r}, color)
    // Right strip
    fill_rect(ctx, core.Rect{rect.x + rect.width - r, rect.y + r, r, rect.height - 2 * r}, color)

    // Corners
    draw_corner(ctx, rect.x + r, rect.y + r, r, color, .TopLeft)
    draw_corner(ctx, rect.x + rect.width - r - 1, rect.y + r, r, color, .TopRight)
    draw_corner(ctx, rect.x + r, rect.y + rect.height - r - 1, r, color, .BottomLeft)
    draw_corner(ctx, rect.x + rect.width - r - 1, rect.y + rect.height - r - 1, r, color, .BottomRight)
}

Corner :: enum {
    TopLeft,
    TopRight,
    BottomLeft,
    BottomRight,
}

// Draw a filled quarter circle (for rounded corners)
draw_corner :: proc(ctx: ^Draw_Context, cx, cy: i32, r: i32, color: core.Color, corner: Corner) {
    pixel := core.color_to_argb(color)
    stride_pixels := ctx.stride / 4
    r_sq := r * r

    for y in 0 ..= r {
        for x in 0 ..= r {
            if x * x + y * y <= r_sq {
                px, py: i32
                switch corner {
                case .TopLeft:
                    px = cx - x
                    py = cy - y
                case .TopRight:
                    px = cx + x
                    py = cy - y
                case .BottomLeft:
                    px = cx - x
                    py = cy + y
                case .BottomRight:
                    px = cx + x
                    py = cy + y
                }

                if px >= ctx.clip.x && px < ctx.clip.x + ctx.clip.width &&
                   py >= ctx.clip.y && py < ctx.clip.y + ctx.clip.height {
                    ctx.pixels[py * stride_pixels + px] = pixel
                }
            }
        }
    }
}

// Blend source color over destination pixel (premultiplied alpha)
blend_pixel :: proc(dst: u32, src: core.Color) -> u32 {
    // Extract destination components
    da := u8(dst >> 24)
    dr := u8(dst >> 16)
    dg := u8(dst >> 8)
    db := u8(dst)

    // Source is premultiplied, destination should be too
    // out = src + dst * (1 - src_alpha)
    inv_alpha := 255 - src.a

    ra := src.a + u8((u32(da) * u32(inv_alpha)) / 255)
    rr := src.r + u8((u32(dr) * u32(inv_alpha)) / 255)
    rg := src.g + u8((u32(dg) * u32(inv_alpha)) / 255)
    rb := src.b + u8((u32(db) * u32(inv_alpha)) / 255)

    return (u32(ra) << 24) | (u32(rr) << 16) | (u32(rg) << 8) | u32(rb)
}

// Draw a horizontal line
draw_hline :: proc(ctx: ^Draw_Context, x1, x2, y: i32, color: core.Color) {
    if y < ctx.clip.y || y >= ctx.clip.y + ctx.clip.height {
        return
    }

    start := max(min(x1, x2), ctx.clip.x)
    end := min(max(x1, x2), ctx.clip.x + ctx.clip.width - 1)

    if start > end {
        return
    }

    pixel := core.color_to_argb(color)
    stride_pixels := ctx.stride / 4
    row_start := y * stride_pixels

    for x in start ..= end {
        ctx.pixels[row_start + x] = pixel
    }
}

// Draw a vertical line
draw_vline :: proc(ctx: ^Draw_Context, x, y1, y2: i32, color: core.Color) {
    if x < ctx.clip.x || x >= ctx.clip.x + ctx.clip.width {
        return
    }

    start := max(min(y1, y2), ctx.clip.y)
    end := min(max(y1, y2), ctx.clip.y + ctx.clip.height - 1)

    if start > end {
        return
    }

    pixel := core.color_to_argb(color)
    stride_pixels := ctx.stride / 4

    for y in start ..= end {
        ctx.pixels[y * stride_pixels + x] = pixel
    }
}

// Set a single pixel
set_pixel :: proc(ctx: ^Draw_Context, x, y: i32, color: core.Color) {
    if x < ctx.clip.x || x >= ctx.clip.x + ctx.clip.width ||
       y < ctx.clip.y || y >= ctx.clip.y + ctx.clip.height {
        return
    }

    stride_pixels := ctx.stride / 4
    ctx.pixels[y * stride_pixels + x] = core.color_to_argb(color)
}

// Set a single pixel with blending
set_pixel_blend :: proc(ctx: ^Draw_Context, x, y: i32, color: core.Color) {
    if x < ctx.clip.x || x >= ctx.clip.x + ctx.clip.width ||
       y < ctx.clip.y || y >= ctx.clip.y + ctx.clip.height {
        return
    }

    stride_pixels := ctx.stride / 4
    idx := y * stride_pixels + x
    ctx.pixels[idx] = blend_pixel(ctx.pixels[idx], color)
}
