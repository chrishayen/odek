package tests

import "../src/core"
import "../src/render"
import "core:testing"

// Font path for testing (may vary by system)
TEST_FONT_PATH :: "/usr/share/fonts/noto/NotoSans-Regular.ttf\x00"

@(test)
test_text_renderer_init :: proc(t: ^testing.T) {
    renderer, ok := render.text_renderer_init()
    testing.expect(t, ok, "Text renderer should initialize successfully")

    if ok {
        render.text_renderer_destroy(&renderer)
        testing.expect(t, renderer.library == nil, "Library should be nil after destroy")
    }
}

@(test)
test_font_load :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        // Font might not be available on all systems, skip test
        return
    }
    defer render.font_destroy(&font)

    testing.expect(t, font.face != nil, "Font face should not be nil")
    testing.expect(t, font.size == 16, "Font size should be 16")
    testing.expect(t, font.ascender > 0, "Ascender should be positive")
    testing.expect(t, font.descender < 0, "Descender should be negative")
    testing.expect(t, font.line_height > 0, "Line height should be positive")
}

@(test)
test_font_metrics :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 24)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    // Check that metrics are reasonable for a 24px font
    testing.expect(t, font.line_height > 0, "Line height should be positive")
    testing.expect(t, font.ascender > 0, "Ascender should be positive")
    testing.expect(t, font.descender <= 0, "Descender should be zero or negative")
    // Line height should be larger for larger font size
    testing.expect(t, font.line_height >= 20, "Line height should be at least 20px for 24pt font")
}

@(test)
test_glyph_loading :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    // Load a simple character
    glyph, glyph_ok := render.font_get_glyph(&font, 'A')
    testing.expect(t, glyph_ok, "Should be able to load glyph for 'A'")

    if glyph_ok {
        testing.expect(t, glyph.width > 0, "Glyph width should be positive")
        testing.expect(t, glyph.height > 0, "Glyph height should be positive")
        testing.expect(t, glyph.advance > 0, "Glyph advance should be positive")
        testing.expect(t, len(glyph.bitmap) > 0, "Glyph should have bitmap data")
    }
}

@(test)
test_glyph_caching :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    // Load same glyph twice - should get from cache
    glyph1, ok1 := render.font_get_glyph(&font, 'X')
    glyph2, ok2 := render.font_get_glyph(&font, 'X')

    testing.expect(t, ok1 && ok2, "Both glyph loads should succeed")

    if ok1 && ok2 {
        // Cached glyph should have same properties
        testing.expect(t, glyph1.width == glyph2.width, "Cached glyph width should match")
        testing.expect(t, glyph1.height == glyph2.height, "Cached glyph height should match")
        testing.expect(t, glyph1.advance == glyph2.advance, "Cached glyph advance should match")
    }
}

@(test)
test_text_measure :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    // Empty string should have zero width
    empty_width := render.text_measure(&font, "")
    testing.expect(t, empty_width == 0, "Empty string should have zero width")

    // Single character
    single_width := render.text_measure(&font, "A")
    testing.expect(t, single_width > 0, "Single character should have positive width")

    // Longer string should be wider
    longer_width := render.text_measure(&font, "AAAA")
    testing.expect(t, longer_width > single_width, "Longer string should be wider")

    // Width should be roughly proportional to character count for same characters
    testing.expect(t, longer_width <= single_width * 5, "Width should be roughly proportional")
}

@(test)
test_text_measure_size :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    size := render.text_measure_size(&font, "Hello")

    testing.expect(t, size.width > 0, "Text width should be positive")
    testing.expect(t, size.height == font.line_height, "Text height should equal line height")
}

@(test)
test_draw_text :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    // Create a small buffer
    pixels: [100 * 50]u32
    ctx := render.context_create(&pixels[0], 100, 50, 400)

    // Clear to black
    render.clear(&ctx, core.COLOR_BLACK)

    // Draw some text
    render.draw_text_top(&ctx, &font, "Hi", 10, 10, core.COLOR_WHITE)

    // Check that some pixels were modified
    non_black_count := 0
    for i in 0 ..< 100 * 50 {
        if pixels[i] != 0xFF000000 {
            non_black_count += 1
        }
    }

    testing.expect(t, non_black_count > 0, "Drawing text should modify some pixels")
}

@(test)
test_draw_text_clipping :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    // Create a small buffer and pre-fill with black
    pixels: [50 * 30]u32
    for i in 0 ..< 50 * 30 {
        pixels[i] = 0xFF000000  // Black with full alpha
    }

    ctx := render.context_create(&pixels[0], 50, 30, 200)

    // Set a clip region - only this area should be modified
    render.context_set_clip(&ctx, core.Rect{10, 10, 30, 15})

    // Draw text that would extend beyond clip
    render.draw_text_top(&ctx, &font, "Hello World", 0, 5, core.COLOR_WHITE)

    // Verify pixels above clip region are still black
    for y in 0 ..< 10 {
        for x in 0 ..< 50 {
            idx := y * 50 + x
            testing.expect(t, pixels[idx] == 0xFF000000, "Pixels above clip should be black")
        }
    }

    // Verify pixels left of clip region are still black
    for y in 10 ..< 25 {
        for x in 0 ..< 10 {
            idx := y * 50 + x
            testing.expect(t, pixels[idx] == 0xFF000000, "Pixels left of clip should be black")
        }
    }
}

@(test)
test_space_character :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    font, font_ok := render.font_load(&renderer, TEST_FONT_PATH, 16)
    if !font_ok {
        return
    }
    defer render.font_destroy(&font)

    // Space should have an advance but minimal/no bitmap
    glyph, ok := render.font_get_glyph(&font, ' ')
    testing.expect(t, ok, "Should be able to load space glyph")

    if ok {
        testing.expect(t, glyph.advance > 0, "Space should have positive advance")
    }
}

// ============================================================================
// Fontconfig tests
// ============================================================================

@(test)
test_fc_get_font_path :: proc(t: ^testing.T) {
    // Get sans font path
    path := render.fc_get_font_path("sans", false)
    defer delete(path)

    testing.expect(t, len(path) > 0, "Should find a sans font")

    // Path should end in .ttf or .otf
    if len(path) > 4 {
        ext := path[len(path)-4:]
        has_valid_ext := ext == ".ttf" || ext == ".otf"
        testing.expect(t, has_valid_ext, "Font path should end in .ttf or .otf")
    }
}

@(test)
test_fc_get_font_path_bold :: proc(t: ^testing.T) {
    // Get bold sans font path
    path := render.fc_get_font_path("sans", true)
    defer delete(path)

    testing.expect(t, len(path) > 0, "Should find a bold sans font")
}

@(test)
test_fc_get_font_path_loads :: proc(t: ^testing.T) {
    renderer, renderer_ok := render.text_renderer_init()
    if !renderer_ok {
        return
    }
    defer render.text_renderer_destroy(&renderer)

    // Get font path from fontconfig
    path := render.fc_get_font_path("sans", false)
    defer delete(path)

    if len(path) == 0 {
        return
    }

    // Should be able to load font at this path
    font, font_ok := render.font_load(&renderer, path, 16)
    testing.expect(t, font_ok, "Should be able to load font from fontconfig path")

    if font_ok {
        render.font_destroy(&font)
    }
}
