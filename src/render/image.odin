package render

import "../core"
import stbi "vendor:stb/image"

// Image data structure (owns pixel data)
Image :: struct {
    pixels: []u32,      // ARGB8888 format (premultiplied alpha), owned
    width:  i32,
    height: i32,
}

// Load image from file path
// Returns (Image, true) on success, ({}, false) on failure
image_load :: proc(path: string) -> (Image, bool) {
    width, height, channels: i32

    // Request 4 channels (RGBA)
    data := stbi.load(cstring(raw_data(path)), &width, &height, &channels, 4)
    if data == nil {
        return {}, false
    }
    defer stbi.image_free(data)

    // Convert RGBA to ARGB8888 with premultiplied alpha
    pixel_count := int(width * height)
    pixels := make([]u32, pixel_count)

    for i in 0 ..< pixel_count {
        r := data[i * 4 + 0]
        g := data[i * 4 + 1]
        b := data[i * 4 + 2]
        a := data[i * 4 + 3]

        // Premultiply alpha
        if a == 0 {
            pixels[i] = 0
        } else if a == 255 {
            pixels[i] = (u32(a) << 24) | (u32(r) << 16) | (u32(g) << 8) | u32(b)
        } else {
            alpha := f32(a) / 255.0
            pr := u8(f32(r) * alpha)
            pg := u8(f32(g) * alpha)
            pb := u8(f32(b) * alpha)
            pixels[i] = (u32(a) << 24) | (u32(pr) << 16) | (u32(pg) << 8) | u32(pb)
        }
    }

    return Image{
        pixels = pixels,
        width  = width,
        height = height,
    }, true
}

// Destroy image and free memory
image_destroy :: proc(img: ^Image) {
    if img.pixels != nil {
        delete(img.pixels)
    }
    img.pixels = nil
    img.width = 0
    img.height = 0
}

// Check if image is valid
image_is_valid :: proc(img: ^Image) -> bool {
    return img != nil && img.pixels != nil && img.width > 0 && img.height > 0
}

// Blit image to draw context at position (x, y)
// Respects clipping, performs alpha blending
draw_image :: proc(ctx: ^Draw_Context, img: ^Image, x, y: i32) {
    if !image_is_valid(img) {
        return
    }

    stride_pixels := ctx.stride / 4

    for row in 0 ..< img.height {
        py := y + row
        if py < ctx.clip.y || py >= ctx.clip.y + ctx.clip.height {
            continue
        }

        for col in 0 ..< img.width {
            px := x + col
            if px < ctx.clip.x || px >= ctx.clip.x + ctx.clip.width {
                continue
            }

            src := img.pixels[row * img.width + col]
            src_a := (src >> 24) & 0xFF

            if src_a == 0 {
                continue  // Fully transparent
            }

            idx := py * stride_pixels + px

            if src_a == 255 {
                // Fully opaque - direct copy
                ctx.pixels[idx] = src
            } else {
                // Alpha blend (src is premultiplied)
                dst := ctx.pixels[idx]
                inv_alpha := 255 - src_a

                dst_a := (dst >> 24) & 0xFF
                dst_r := (dst >> 16) & 0xFF
                dst_g := (dst >> 8) & 0xFF
                dst_b := dst & 0xFF

                src_r := (src >> 16) & 0xFF
                src_g := (src >> 8) & 0xFF
                src_b := src & 0xFF

                out_a := src_a + (dst_a * inv_alpha) / 255
                out_r := src_r + (dst_r * inv_alpha) / 255
                out_g := src_g + (dst_g * inv_alpha) / 255
                out_b := src_b + (dst_b * inv_alpha) / 255

                ctx.pixels[idx] = (out_a << 24) | (out_r << 16) | (out_g << 8) | out_b
            }
        }
    }
}

// Blit image scaled to fit within dst_rect (maintains aspect ratio)
// Uses nearest-neighbor scaling
draw_image_scaled :: proc(ctx: ^Draw_Context, img: ^Image, dst_rect: core.Rect) {
    if !image_is_valid(img) {
        return
    }

    // Calculate scale to fit while maintaining aspect ratio
    scale_x := f32(dst_rect.width) / f32(img.width)
    scale_y := f32(dst_rect.height) / f32(img.height)
    scale := min(scale_x, scale_y)

    scaled_w := i32(f32(img.width) * scale)
    scaled_h := i32(f32(img.height) * scale)

    // Center within dst_rect
    offset_x := (dst_rect.width - scaled_w) / 2
    offset_y := (dst_rect.height - scaled_h) / 2

    stride_pixels := ctx.stride / 4

    for dy in 0 ..< scaled_h {
        py := dst_rect.y + offset_y + dy
        if py < ctx.clip.y || py >= ctx.clip.y + ctx.clip.height {
            continue
        }

        // Source y (nearest neighbor)
        sy := i32(f32(dy) / scale)
        if sy >= img.height {
            sy = img.height - 1
        }

        for dx in 0 ..< scaled_w {
            px := dst_rect.x + offset_x + dx
            if px < ctx.clip.x || px >= ctx.clip.x + ctx.clip.width {
                continue
            }

            // Source x (nearest neighbor)
            sx := i32(f32(dx) / scale)
            if sx >= img.width {
                sx = img.width - 1
            }

            src := img.pixels[sy * img.width + sx]
            src_a := (src >> 24) & 0xFF

            if src_a == 0 {
                continue
            }

            idx := py * stride_pixels + px

            if src_a == 255 {
                ctx.pixels[idx] = src
            } else {
                dst := ctx.pixels[idx]
                inv_alpha := 255 - src_a

                dst_a := (dst >> 24) & 0xFF
                dst_r := (dst >> 16) & 0xFF
                dst_g := (dst >> 8) & 0xFF
                dst_b := dst & 0xFF

                src_r := (src >> 16) & 0xFF
                src_g := (src >> 8) & 0xFF
                src_b := src & 0xFF

                out_a := src_a + (dst_a * inv_alpha) / 255
                out_r := src_r + (dst_r * inv_alpha) / 255
                out_g := src_g + (dst_g * inv_alpha) / 255
                out_b := src_b + (dst_b * inv_alpha) / 255

                ctx.pixels[idx] = (out_a << 24) | (out_r << 16) | (out_g << 8) | out_b
            }
        }
    }
}

// Create a thumbnail that fits within max_width x max_height
// Maintains aspect ratio, uses high-quality box filter for downscaling
// Returns new Image that caller owns
image_create_thumbnail :: proc(img: ^Image, max_width, max_height: i32) -> (Image, bool) {
    if !image_is_valid(img) {
        return {}, false
    }

    if max_width <= 0 || max_height <= 0 {
        return {}, false
    }

    // Calculate thumbnail size maintaining aspect ratio
    scale_x := f32(max_width) / f32(img.width)
    scale_y := f32(max_height) / f32(img.height)
    scale := min(scale_x, scale_y)

    // Don't upscale
    if scale >= 1.0 {
        // Just copy the image
        return image_copy(img)
    }

    thumb_w := max(1, i32(f32(img.width) * scale))
    thumb_h := max(1, i32(f32(img.height) * scale))

    // Convert source ARGB to RGBA bytes for stb_image_resize
    src_rgba := make([]u8, img.width * img.height * 4)
    defer delete(src_rgba)

    for i in 0 ..< int(img.width * img.height) {
        pixel := img.pixels[i]
        // Extract components (premultiplied ARGB)
        a := u8((pixel >> 24) & 0xFF)
        r := u8((pixel >> 16) & 0xFF)
        g := u8((pixel >> 8) & 0xFF)
        b := u8(pixel & 0xFF)

        // Un-premultiply for resize (stb_image_resize handles alpha better this way)
        if a > 0 && a < 255 {
            factor := 255.0 / f32(a)
            r = u8(min(255.0, f32(r) * factor))
            g = u8(min(255.0, f32(g) * factor))
            b = u8(min(255.0, f32(b) * factor))
        }

        // Store as RGBA
        src_rgba[i * 4 + 0] = r
        src_rgba[i * 4 + 1] = g
        src_rgba[i * 4 + 2] = b
        src_rgba[i * 4 + 3] = a
    }

    // Allocate destination RGBA buffer
    dst_rgba := make([]u8, thumb_w * thumb_h * 4)
    defer delete(dst_rgba)

    // Use stb_image_resize
    result := stbi.resize_uint8(
        raw_data(src_rgba), img.width, img.height, img.width * 4,
        raw_data(dst_rgba), thumb_w, thumb_h, thumb_w * 4,
        4,
    )

    if result == 0 {
        return {}, false
    }

    // Convert RGBA back to premultiplied ARGB
    thumb_pixels := make([]u32, thumb_w * thumb_h)

    for i in 0 ..< int(thumb_w * thumb_h) {
        r := dst_rgba[i * 4 + 0]
        g := dst_rgba[i * 4 + 1]
        b := dst_rgba[i * 4 + 2]
        a := dst_rgba[i * 4 + 3]

        // Premultiply alpha
        if a == 0 {
            thumb_pixels[i] = 0
        } else if a == 255 {
            thumb_pixels[i] = (u32(a) << 24) | (u32(r) << 16) | (u32(g) << 8) | u32(b)
        } else {
            alpha := f32(a) / 255.0
            pr := u8(f32(r) * alpha)
            pg := u8(f32(g) * alpha)
            pb := u8(f32(b) * alpha)
            thumb_pixels[i] = (u32(a) << 24) | (u32(pr) << 16) | (u32(pg) << 8) | u32(pb)
        }
    }

    return Image{
        pixels = thumb_pixels,
        width  = thumb_w,
        height = thumb_h,
    }, true
}

// Copy an image
image_copy :: proc(img: ^Image) -> (Image, bool) {
    if !image_is_valid(img) {
        return {}, false
    }

    pixels := make([]u32, img.width * img.height)
    copy(pixels, img.pixels)

    return Image{
        pixels = pixels,
        width  = img.width,
        height = img.height,
    }, true
}
