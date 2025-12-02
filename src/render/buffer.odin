package render

import "../core"
import "base:intrinsics"
import "core:slice"

// Draw context for software rendering
Draw_Context :: struct {
    pixels: [^]u32,
    width: i32,           // Physical buffer width
    height: i32,          // Physical buffer height
    stride: i32,
    clip: core.Rect,      // Clip rect in physical coordinates
    logical_clip: core.Rect,  // Clip rect in logical coordinates (for widgets)
    scale: f64,           // Scale factor (logical -> physical)
    logical_width: i32,   // Logical width (for widgets)
    logical_height: i32,  // Logical height (for widgets)
}

// Create a draw context
context_create :: proc(pixels: [^]u32, width, height, stride: i32) -> Draw_Context {
    return Draw_Context{
        pixels = pixels,
        width = width,
        height = height,
        stride = stride,
        clip = core.Rect{0, 0, width, height},
        logical_clip = core.Rect{0, 0, width, height},
        scale = 1.0,
        logical_width = width,
        logical_height = height,
    }
}

// Create a draw context with scale factor
context_create_scaled :: proc(pixels: [^]u32, phys_width, phys_height, stride: i32, logical_width, logical_height: i32, scale: f64) -> Draw_Context {
    return Draw_Context{
        pixels = pixels,
        width = phys_width,
        height = phys_height,
        stride = stride,
        clip = core.Rect{0, 0, phys_width, phys_height},
        logical_clip = core.Rect{0, 0, logical_width, logical_height},
        scale = scale,
        logical_width = logical_width,
        logical_height = logical_height,
    }
}

// Scale a logical coordinate to physical
scale_coord :: proc(ctx: ^Draw_Context, val: i32) -> i32 {
    return i32(f64(val) * ctx.scale)
}

// Scale a logical rect to physical
scale_rect :: proc(ctx: ^Draw_Context, rect: core.Rect) -> core.Rect {
    return core.Rect{
        x = i32(f64(rect.x) * ctx.scale),
        y = i32(f64(rect.y) * ctx.scale),
        width = i32(f64(rect.width) * ctx.scale),
        height = i32(f64(rect.height) * ctx.scale),
    }
}

// Set clipping rectangle
// clip is in logical coordinates, will be scaled to physical
context_set_clip :: proc(ctx: ^Draw_Context, clip: core.Rect) {
    // Store logical clip (for widget clipping calculations)
    logical_bounds := core.Rect{0, 0, ctx.logical_width, ctx.logical_height}
    if logical_clipped, ok := core.rect_intersection(clip, logical_bounds); ok {
        ctx.logical_clip = logical_clipped
    } else {
        ctx.logical_clip = core.Rect{}
    }

    // Scale logical clip to physical
    phys_clip := scale_rect(ctx, clip)
    // Intersect with buffer bounds
    bounds := core.Rect{0, 0, ctx.width, ctx.height}
    if clipped, ok := core.rect_intersection(phys_clip, bounds); ok {
        ctx.clip = clipped
    } else {
        ctx.clip = core.Rect{}
    }
}

// Clear the entire buffer
clear :: proc(ctx: ^Draw_Context, color: core.Color) {
    fill_rect(ctx, core.Rect{0, 0, ctx.logical_width, ctx.logical_height}, color)
}

// Fill a rectangle with a solid color (optimized)
// rect is in logical coordinates, will be scaled to physical
fill_rect :: proc(ctx: ^Draw_Context, rect: core.Rect, color: core.Color) {
    // Scale logical rect to physical
    phys_rect := scale_rect(ctx, rect)
    clipped, ok := core.rect_intersection(phys_rect, ctx.clip)
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
// rect is in logical coordinates, will be scaled to physical
fill_rect_blend :: proc(ctx: ^Draw_Context, rect: core.Rect, color: core.Color) {
    if color.a == 255 {
        fill_rect(ctx, rect, color)
        return
    }
    if color.a == 0 {
        return
    }

    // Scale logical rect to physical
    phys_rect := scale_rect(ctx, rect)
    clipped, ok := core.rect_intersection(phys_rect, ctx.clip)
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
// cx, cy, r are in logical coordinates
draw_corner :: proc(ctx: ^Draw_Context, cx, cy: i32, r: i32, color: core.Color, corner: Corner) {
    pixel := core.color_to_argb(color)
    stride_pixels := ctx.stride / 4

    // Scale to physical coordinates
    phys_cx := scale_coord(ctx, cx)
    phys_cy := scale_coord(ctx, cy)
    phys_r := scale_coord(ctx, r)
    r_sq := phys_r * phys_r

    for y in 0 ..= phys_r {
        for x in 0 ..= phys_r {
            if x * x + y * y <= r_sq {
                px, py: i32
                switch corner {
                case .TopLeft:
                    px = phys_cx - x
                    py = phys_cy - y
                case .TopRight:
                    px = phys_cx + x
                    py = phys_cy - y
                case .BottomLeft:
                    px = phys_cx - x
                    py = phys_cy + y
                case .BottomRight:
                    px = phys_cx + x
                    py = phys_cy + y
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

// Draw a horizontal line (coordinates in logical)
draw_hline :: proc(ctx: ^Draw_Context, x1, x2, y: i32, color: core.Color) {
    // Scale to physical
    phys_y := scale_coord(ctx, y)
    phys_x1 := scale_coord(ctx, x1)
    phys_x2 := scale_coord(ctx, x2)

    if phys_y < ctx.clip.y || phys_y >= ctx.clip.y + ctx.clip.height {
        return
    }

    start := max(min(phys_x1, phys_x2), ctx.clip.x)
    end := min(max(phys_x1, phys_x2), ctx.clip.x + ctx.clip.width - 1)

    if start > end {
        return
    }

    pixel := core.color_to_argb(color)
    stride_pixels := ctx.stride / 4
    row_start := phys_y * stride_pixels

    for x in start ..= end {
        ctx.pixels[row_start + x] = pixel
    }
}

// Draw a vertical line (coordinates in logical)
draw_vline :: proc(ctx: ^Draw_Context, x, y1, y2: i32, color: core.Color) {
    // Scale to physical
    phys_x := scale_coord(ctx, x)
    phys_y1 := scale_coord(ctx, y1)
    phys_y2 := scale_coord(ctx, y2)

    if phys_x < ctx.clip.x || phys_x >= ctx.clip.x + ctx.clip.width {
        return
    }

    start := max(min(phys_y1, phys_y2), ctx.clip.y)
    end := min(max(phys_y1, phys_y2), ctx.clip.y + ctx.clip.height - 1)

    if start > end {
        return
    }

    pixel := core.color_to_argb(color)
    stride_pixels := ctx.stride / 4

    for y in start ..= end {
        ctx.pixels[y * stride_pixels + phys_x] = pixel
    }
}

// Set a single pixel (coordinates in logical)
set_pixel :: proc(ctx: ^Draw_Context, x, y: i32, color: core.Color) {
    phys_x := scale_coord(ctx, x)
    phys_y := scale_coord(ctx, y)

    if phys_x < ctx.clip.x || phys_x >= ctx.clip.x + ctx.clip.width ||
       phys_y < ctx.clip.y || phys_y >= ctx.clip.y + ctx.clip.height {
        return
    }

    stride_pixels := ctx.stride / 4
    ctx.pixels[phys_y * stride_pixels + phys_x] = core.color_to_argb(color)
}

// Set a single pixel with blending (coordinates in logical)
set_pixel_blend :: proc(ctx: ^Draw_Context, x, y: i32, color: core.Color) {
    phys_x := scale_coord(ctx, x)
    phys_y := scale_coord(ctx, y)

    if phys_x < ctx.clip.x || phys_x >= ctx.clip.x + ctx.clip.width ||
       phys_y < ctx.clip.y || phys_y >= ctx.clip.y + ctx.clip.height {
        return
    }

    stride_pixels := ctx.stride / 4
    idx := phys_y * stride_pixels + phys_x
    ctx.pixels[idx] = blend_pixel(ctx.pixels[idx], color)
}

// Draw a rounded rectangle outline
draw_rounded_rect :: proc(ctx: ^Draw_Context, rect: core.Rect, radius: i32, color: core.Color) {
    if radius <= 0 {
        draw_rect(ctx, rect, color)
        return
    }

    r := min(radius, rect.width / 2, rect.height / 2)

    // Top edge
    draw_hline(ctx, rect.x + r, rect.x + rect.width - r - 1, rect.y, color)
    // Bottom edge
    draw_hline(ctx, rect.x + r, rect.x + rect.width - r - 1, rect.y + rect.height - 1, color)
    // Left edge
    draw_vline(ctx, rect.x, rect.y + r, rect.y + rect.height - r - 1, color)
    // Right edge
    draw_vline(ctx, rect.x + rect.width - 1, rect.y + r, rect.y + rect.height - r - 1, color)

    // Corner arcs
    draw_corner_arc(ctx, rect.x + r, rect.y + r, r, color, .TopLeft)
    draw_corner_arc(ctx, rect.x + rect.width - r - 1, rect.y + r, r, color, .TopRight)
    draw_corner_arc(ctx, rect.x + r, rect.y + rect.height - r - 1, r, color, .BottomLeft)
    draw_corner_arc(ctx, rect.x + rect.width - r - 1, rect.y + rect.height - r - 1, r, color, .BottomRight)
}

// Draw a quarter circle arc (outline) for rounded corners
draw_corner_arc :: proc(ctx: ^Draw_Context, cx, cy: i32, r: i32, color: core.Color, corner: Corner) {
    pixel := core.color_to_argb(color)
    stride_pixels := ctx.stride / 4

    // Scale to physical coordinates
    phys_cx := scale_coord(ctx, cx)
    phys_cy := scale_coord(ctx, cy)
    phys_r := scale_coord(ctx, r)

    // Midpoint circle algorithm
    x := phys_r
    y: i32 = 0
    err := 1 - phys_r

    for x >= y {
        // Draw the 2 points for this quarter
        px1, py1, px2, py2: i32
        switch corner {
        case .TopLeft:
            px1 = phys_cx - x
            py1 = phys_cy - y
            px2 = phys_cx - y
            py2 = phys_cy - x
        case .TopRight:
            px1 = phys_cx + x
            py1 = phys_cy - y
            px2 = phys_cx + y
            py2 = phys_cy - x
        case .BottomLeft:
            px1 = phys_cx - x
            py1 = phys_cy + y
            px2 = phys_cx - y
            py2 = phys_cy + x
        case .BottomRight:
            px1 = phys_cx + x
            py1 = phys_cy + y
            px2 = phys_cx + y
            py2 = phys_cy + x
        }

        if px1 >= ctx.clip.x && px1 < ctx.clip.x + ctx.clip.width &&
           py1 >= ctx.clip.y && py1 < ctx.clip.y + ctx.clip.height {
            ctx.pixels[py1 * stride_pixels + px1] = pixel
        }
        if px2 >= ctx.clip.x && px2 < ctx.clip.x + ctx.clip.width &&
           py2 >= ctx.clip.y && py2 < ctx.clip.y + ctx.clip.height {
            ctx.pixels[py2 * stride_pixels + px2] = pixel
        }

        y += 1
        if err < 0 {
            err += 2 * y + 1
        } else {
            x -= 1
            err += 2 * (y - x) + 1
        }
    }
}

// Fill a triangle with a solid color
// Coordinates are in logical pixels
fill_triangle :: proc(ctx: ^Draw_Context, x1, y1, x2, y2, x3, y3: i32, color: core.Color) {
    // Sort vertices by y coordinate
    v1_x, v1_y := x1, y1
    v2_x, v2_y := x2, y2
    v3_x, v3_y := x3, y3

    if v1_y > v2_y {
        v1_x, v2_x = v2_x, v1_x
        v1_y, v2_y = v2_y, v1_y
    }
    if v2_y > v3_y {
        v2_x, v3_x = v3_x, v2_x
        v2_y, v3_y = v3_y, v2_y
    }
    if v1_y > v2_y {
        v1_x, v2_x = v2_x, v1_x
        v1_y, v2_y = v2_y, v1_y
    }

    // Now v1_y <= v2_y <= v3_y

    pixel := core.color_to_argb(color)
    stride_pixels := ctx.stride / 4

    // Helper to draw a horizontal span
    draw_span :: proc(ctx: ^Draw_Context, y, x_start, x_end: i32, pixel: u32, stride_pixels: i32) {
        phys_y := scale_coord(ctx, y)
        if phys_y < ctx.clip.y || phys_y >= ctx.clip.y + ctx.clip.height {
            return
        }

        phys_x1 := scale_coord(ctx, min(x_start, x_end))
        phys_x2 := scale_coord(ctx, max(x_start, x_end))

        start := max(phys_x1, ctx.clip.x)
        end := min(phys_x2, ctx.clip.x + ctx.clip.width - 1)

        if start > end {
            return
        }

        row_start := phys_y * stride_pixels
        for x in start ..= end {
            ctx.pixels[row_start + x] = pixel
        }
    }

    // Fill the triangle using scanline
    if v2_y == v1_y {
        // Flat top triangle
        for y in v1_y ..= v3_y {
            if v3_y == v1_y {
                break
            }
            t := f64(y - v1_y) / f64(v3_y - v1_y)
            x_left := v1_x + i32(t * f64(v3_x - v1_x))
            x_right := v2_x + i32(t * f64(v3_x - v2_x))
            draw_span(ctx, y, x_left, x_right, pixel, stride_pixels)
        }
    } else if v3_y == v2_y {
        // Flat bottom triangle
        for y in v1_y ..= v2_y {
            if v2_y == v1_y {
                break
            }
            t := f64(y - v1_y) / f64(v2_y - v1_y)
            x_left := v1_x + i32(t * f64(v2_x - v1_x))
            x_right := v1_x + i32(t * f64(v3_x - v1_x))
            draw_span(ctx, y, x_left, x_right, pixel, stride_pixels)
        }
    } else {
        // General case - split into flat bottom and flat top
        // Calculate x coordinate at the split point
        t := f64(v2_y - v1_y) / f64(v3_y - v1_y)
        v4_x := v1_x + i32(t * f64(v3_x - v1_x))

        // Upper part (flat bottom)
        for y in v1_y ..= v2_y {
            if v2_y == v1_y {
                break
            }
            t2 := f64(y - v1_y) / f64(v2_y - v1_y)
            x_left := v1_x + i32(t2 * f64(v2_x - v1_x))
            x_right := v1_x + i32(t2 * f64(v4_x - v1_x))
            draw_span(ctx, y, x_left, x_right, pixel, stride_pixels)
        }

        // Lower part (flat top)
        for y in v2_y ..= v3_y {
            if v3_y == v2_y {
                break
            }
            t2 := f64(y - v2_y) / f64(v3_y - v2_y)
            x_left := v2_x + i32(t2 * f64(v3_x - v2_x))
            x_right := v4_x + i32(t2 * f64(v3_x - v4_x))
            draw_span(ctx, y, x_left, x_right, pixel, stride_pixels)
        }
    }
}
