package tests

import "../src/core"
import "../src/render"
import "core:testing"

@(test)
test_image_load_invalid :: proc(t: ^testing.T) {
    img, ok := render.image_load("nonexistent_file.png")
    testing.expect(t, !ok, "Loading nonexistent image should fail")
    testing.expect(t, img.pixels == nil, "Failed load should have nil pixels")
    testing.expect(t, img.width == 0, "Failed load should have zero width")
    testing.expect(t, img.height == 0, "Failed load should have zero height")
}

@(test)
test_image_is_valid :: proc(t: ^testing.T) {
    // Empty image
    empty := render.Image{}
    testing.expect(t, !render.image_is_valid(&empty), "Empty image should not be valid")

    // Image with nil pixels
    nil_pixels := render.Image{pixels = nil, width = 10, height = 10}
    testing.expect(t, !render.image_is_valid(&nil_pixels), "Image with nil pixels should not be valid")

    // Image with zero dimensions
    zero_dim := render.Image{pixels = make([]u32, 1), width = 0, height = 10}
    testing.expect(t, !render.image_is_valid(&zero_dim), "Image with zero width should not be valid")
    delete(zero_dim.pixels)

    // Valid image
    valid := render.Image{pixels = make([]u32, 4), width = 2, height = 2}
    testing.expect(t, render.image_is_valid(&valid), "Image with pixels and dimensions should be valid")
    delete(valid.pixels)
}

@(test)
test_image_destroy :: proc(t: ^testing.T) {
    // Create a valid image manually
    img := render.Image{
        pixels = make([]u32, 16),
        width = 4,
        height = 4,
    }

    testing.expect(t, render.image_is_valid(&img), "Image should be valid before destroy")

    render.image_destroy(&img)

    testing.expect(t, img.pixels == nil, "Pixels should be nil after destroy")
    testing.expect(t, img.width == 0, "Width should be 0 after destroy")
    testing.expect(t, img.height == 0, "Height should be 0 after destroy")
    testing.expect(t, !render.image_is_valid(&img), "Image should not be valid after destroy")
}

@(test)
test_draw_image :: proc(t: ^testing.T) {
    // Create a 2x2 test image with known pixel values (opaque red)
    img := render.Image{
        pixels = make([]u32, 4),
        width = 2,
        height = 2,
    }
    defer render.image_destroy(&img)

    // Fill with opaque red (ARGB: 0xFFFF0000)
    for i in 0 ..< 4 {
        img.pixels[i] = 0xFFFF0000
    }

    // Create a 10x10 buffer
    pixels: [100]u32
    ctx := render.context_create(&pixels[0], 10, 10, 40)

    // Clear to black
    render.clear(&ctx, core.COLOR_BLACK)

    // Draw the image at (3, 3)
    render.draw_image(&ctx, &img, 3, 3)

    // Verify the pixels
    for y in 0 ..< 10 {
        for x in 0 ..< 10 {
            idx := y * 10 + x
            if x >= 3 && x < 5 && y >= 3 && y < 5 {
                testing.expect(t, pixels[idx] == 0xFFFF0000, "Pixel inside image should be red")
            } else {
                testing.expect(t, pixels[idx] == 0xFF000000, "Pixel outside image should be black")
            }
        }
    }
}

@(test)
test_draw_image_clipping :: proc(t: ^testing.T) {
    // Create a 4x4 test image
    img := render.Image{
        pixels = make([]u32, 16),
        width = 4,
        height = 4,
    }
    defer render.image_destroy(&img)

    // Fill with opaque green
    for i in 0 ..< 16 {
        img.pixels[i] = 0xFF00FF00
    }

    // Create a 10x10 buffer
    pixels: [100]u32
    ctx := render.context_create(&pixels[0], 10, 10, 40)
    render.clear(&ctx, core.COLOR_BLACK)

    // Set clip region
    render.context_set_clip(&ctx, core.Rect{5, 5, 3, 3})

    // Draw image at (4, 4) - only partially in clip region
    render.draw_image(&ctx, &img, 4, 4)

    // Verify only clipped region was drawn
    for y in 0 ..< 10 {
        for x in 0 ..< 10 {
            idx := y * 10 + x
            // Image covers 4-7 in x and y, clip covers 5-7
            // So green should only be at intersection: 5-7 x 5-7
            if x >= 5 && x < 8 && y >= 5 && y < 8 {
                testing.expect(t, pixels[idx] == 0xFF00FF00, "Pixel in clipped image area should be green")
            } else {
                testing.expect(t, pixels[idx] == 0xFF000000, "Pixel outside clipped area should be black")
            }
        }
    }
}

@(test)
test_draw_image_alpha_blending :: proc(t: ^testing.T) {
    // Create a 2x2 test image with 50% alpha red
    img := render.Image{
        pixels = make([]u32, 4),
        width = 2,
        height = 2,
    }
    defer render.image_destroy(&img)

    // 50% alpha red, premultiplied: alpha=128, r=128 (255*0.5), g=0, b=0
    // ARGB: 0x80800000
    for i in 0 ..< 4 {
        img.pixels[i] = 0x80800000
    }

    // Create buffer with white background
    pixels: [100]u32
    ctx := render.context_create(&pixels[0], 10, 10, 40)
    render.clear(&ctx, core.COLOR_WHITE)

    // Draw at (0, 0)
    render.draw_image(&ctx, &img, 0, 0)

    // Check that blending occurred (result should not be pure white or pure red)
    result := pixels[0]
    testing.expect(t, result != 0xFFFFFFFF, "Blended pixel should not be pure white")
    testing.expect(t, result != 0xFFFF0000, "Blended pixel should not be pure red")

    // The alpha channel should be 255 (fully opaque result)
    result_a := (result >> 24) & 0xFF
    testing.expect(t, result_a == 255, "Result alpha should be fully opaque")
}

@(test)
test_draw_image_scaled :: proc(t: ^testing.T) {
    // Create a 2x2 test image
    img := render.Image{
        pixels = make([]u32, 4),
        width = 2,
        height = 2,
    }
    defer render.image_destroy(&img)

    // Fill with opaque blue
    for i in 0 ..< 4 {
        img.pixels[i] = 0xFF0000FF
    }

    // Create a 20x20 buffer
    pixels: [400]u32
    ctx := render.context_create(&pixels[0], 20, 20, 80)
    render.clear(&ctx, core.COLOR_BLACK)

    // Draw scaled to 10x10 rect at (5, 5)
    render.draw_image_scaled(&ctx, &img, core.Rect{5, 5, 10, 10})

    // Verify some pixels in the scaled area are blue
    center := pixels[10 * 20 + 10]
    testing.expect(t, center == 0xFF0000FF, "Center of scaled image should be blue")

    // Verify pixels outside are black
    corner := pixels[0]
    testing.expect(t, corner == 0xFF000000, "Corner should be black")
}

@(test)
test_draw_image_scaled_maintains_aspect :: proc(t: ^testing.T) {
    // Create a 4x2 test image (wide)
    img := render.Image{
        pixels = make([]u32, 8),
        width = 4,
        height = 2,
    }
    defer render.image_destroy(&img)

    // Fill with magenta
    for i in 0 ..< 8 {
        img.pixels[i] = 0xFFFF00FF
    }

    // Create a 20x20 buffer
    pixels: [400]u32
    ctx := render.context_create(&pixels[0], 20, 20, 80)
    render.clear(&ctx, core.COLOR_BLACK)

    // Draw to 10x10 rect - should maintain 2:1 aspect ratio
    // Result should be 10x5 centered in the 10x10 area
    render.draw_image_scaled(&ctx, &img, core.Rect{5, 5, 10, 10})

    // Center should be magenta (image centered vertically with black bars)
    center := pixels[10 * 20 + 10]
    testing.expect(t, center == 0xFFFF00FF, "Center should be magenta")

    // Top of rect area should be black (letterboxing)
    top := pixels[5 * 20 + 10]  // y=5 should be in black letterbox
    testing.expect(t, top == 0xFF000000, "Top letterbox should be black")
}

@(test)
test_draw_image_invalid :: proc(t: ^testing.T) {
    // Should not crash when drawing invalid image
    pixels: [100]u32
    ctx := render.context_create(&pixels[0], 10, 10, 40)

    empty := render.Image{}
    render.draw_image(&ctx, &empty, 0, 0)
    render.draw_image_scaled(&ctx, &empty, core.Rect{0, 0, 10, 10})

    // Should complete without crashing - test passes if we get here
    testing.expect(t, true, "Drawing invalid image should not crash")
}

@(test)
test_image_create_thumbnail_invalid :: proc(t: ^testing.T) {
    empty := render.Image{}
    thumb, ok := render.image_create_thumbnail(&empty, 100, 100)
    testing.expect(t, !ok, "Creating thumbnail from invalid image should fail")
    testing.expect(t, thumb.pixels == nil, "Thumbnail should have nil pixels on failure")
}

@(test)
test_image_create_thumbnail_invalid_size :: proc(t: ^testing.T) {
    img := render.Image{
        pixels = make([]u32, 4),
        width = 2,
        height = 2,
    }
    defer render.image_destroy(&img)

    // Zero max width
    _, ok := render.image_create_thumbnail(&img, 0, 100)
    testing.expect(t, !ok, "Creating thumbnail with zero max width should fail")

    // Negative max height
    _, ok2 := render.image_create_thumbnail(&img, 100, -1)
    testing.expect(t, !ok2, "Creating thumbnail with negative max height should fail")
}

@(test)
test_image_create_thumbnail_maintains_aspect :: proc(t: ^testing.T) {
    // Create a 100x50 image (2:1 aspect ratio)
    img := render.Image{
        pixels = make([]u32, 100 * 50),
        width = 100,
        height = 50,
    }
    // Fill with opaque blue
    for i in 0 ..< 100 * 50 {
        img.pixels[i] = 0xFF0000FF
    }
    defer render.image_destroy(&img)

    // Create thumbnail to fit in 40x40 - should be 40x20 (maintaining 2:1)
    thumb, ok := render.image_create_thumbnail(&img, 40, 40)
    testing.expect(t, ok, "Creating thumbnail should succeed")
    defer render.image_destroy(&thumb)

    testing.expect(t, thumb.width == 40, "Thumbnail width should be 40")
    testing.expect(t, thumb.height == 20, "Thumbnail height should be 20")

    // Check pixels are valid
    testing.expect(t, render.image_is_valid(&thumb), "Thumbnail should be valid")
}

@(test)
test_image_create_thumbnail_no_upscale :: proc(t: ^testing.T) {
    // Create a small 10x10 image
    img := render.Image{
        pixels = make([]u32, 100),
        width = 10,
        height = 10,
    }
    for i in 0 ..< 100 {
        img.pixels[i] = 0xFFFF0000
    }
    defer render.image_destroy(&img)

    // Request 100x100 thumbnail - should return copy of original size
    thumb, ok := render.image_create_thumbnail(&img, 100, 100)
    testing.expect(t, ok, "Creating thumbnail should succeed")
    defer render.image_destroy(&thumb)

    testing.expect(t, thumb.width == 10, "Thumbnail should not upscale - width should be 10")
    testing.expect(t, thumb.height == 10, "Thumbnail should not upscale - height should be 10")
}

@(test)
test_image_create_thumbnail_downscale :: proc(t: ^testing.T) {
    // Create a 200x200 image
    img := render.Image{
        pixels = make([]u32, 200 * 200),
        width = 200,
        height = 200,
    }
    // Fill with opaque green
    for i in 0 ..< 200 * 200 {
        img.pixels[i] = 0xFF00FF00
    }
    defer render.image_destroy(&img)

    // Create 50x50 thumbnail
    thumb, ok := render.image_create_thumbnail(&img, 50, 50)
    testing.expect(t, ok, "Creating thumbnail should succeed")
    defer render.image_destroy(&thumb)

    testing.expect(t, thumb.width == 50, "Thumbnail width should be 50")
    testing.expect(t, thumb.height == 50, "Thumbnail height should be 50")

    // Check that thumbnail has pixel data
    testing.expect(t, len(thumb.pixels) == 2500, "Thumbnail should have 2500 pixels")

    // Check pixels are greenish (may not be exactly 0xFF00FF00 due to filtering)
    center := thumb.pixels[25 * 50 + 25]
    center_g := (center >> 8) & 0xFF
    testing.expect(t, center_g > 200, "Center pixel should be mostly green")
}

@(test)
test_image_copy :: proc(t: ^testing.T) {
    img := render.Image{
        pixels = make([]u32, 4),
        width = 2,
        height = 2,
    }
    img.pixels[0] = 0xFF112233
    img.pixels[1] = 0xFF445566
    img.pixels[2] = 0xFF778899
    img.pixels[3] = 0xFFAABBCC
    defer render.image_destroy(&img)

    copy_img, ok := render.image_copy(&img)
    testing.expect(t, ok, "Image copy should succeed")
    defer render.image_destroy(&copy_img)

    testing.expect(t, copy_img.width == 2, "Copy width should match")
    testing.expect(t, copy_img.height == 2, "Copy height should match")
    testing.expect(t, copy_img.pixels[0] == 0xFF112233, "Pixel 0 should match")
    testing.expect(t, copy_img.pixels[1] == 0xFF445566, "Pixel 1 should match")
    testing.expect(t, copy_img.pixels[2] == 0xFF778899, "Pixel 2 should match")
    testing.expect(t, copy_img.pixels[3] == 0xFFAABBCC, "Pixel 3 should match")

    // Verify it's a true copy (modifying copy doesn't affect original)
    copy_img.pixels[0] = 0x00000000
    testing.expect(t, img.pixels[0] == 0xFF112233, "Original should be unchanged")
}

@(test)
test_image_copy_invalid :: proc(t: ^testing.T) {
    empty := render.Image{}
    copy_img, ok := render.image_copy(&empty)
    testing.expect(t, !ok, "Copying invalid image should fail")
    testing.expect(t, copy_img.pixels == nil, "Copy should have nil pixels")
}
