package render

import "core:c"
import "../core"

// Cached glyph data
Glyph :: struct {
    bitmap:   []u8,        // Grayscale bitmap data (owned copy)
    width:    i32,
    height:   i32,
    bearing_x: i32,        // Offset from pen position to left edge
    bearing_y: i32,        // Offset from baseline to top edge
    advance:  i32,         // Horizontal advance in pixels
}

// Glyph cache key
Glyph_Key :: struct {
    codepoint: rune,
    size:      u32,
}

// Font instance with specific size
Font :: struct {
    face:        ^FT_FaceRec,
    size:        u32,         // Current pixel size (may be scaled)
    base_size:   u32,         // Base logical size (before scaling)
    ascender:    i32,         // Pixels above baseline
    descender:   i32,         // Pixels below baseline (negative)
    line_height: i32,         // Total line height in pixels
    cache:       map[Glyph_Key]Glyph,
}

// Text renderer manages FreeType and fonts
Text_Renderer :: struct {
    library: ^FT_Library,
}

// Initialize text rendering system
text_renderer_init :: proc() -> (Text_Renderer, bool) {
    renderer: Text_Renderer
    err := FT_Init_FreeType(&renderer.library)
    if err != 0 {
        return {}, false
    }
    return renderer, true
}

// Destroy text renderer
text_renderer_destroy :: proc(renderer: ^Text_Renderer) {
    if renderer.library != nil {
        FT_Done_FreeType(renderer.library)
        renderer.library = nil
    }
}

// Load a font from file
// size is the base logical size (will be scaled when font_set_scale is called)
font_load :: proc(renderer: ^Text_Renderer, path: string, size: u32) -> (Font, bool) {
    font: Font
    font.size = size
    font.base_size = size

    // Load face
    c_path := cstring(raw_data(path))
    err := FT_New_Face(renderer.library, c_path, 0, &font.face)
    if err != 0 {
        return {}, false
    }

    // Set pixel size
    err = FT_Set_Pixel_Sizes(font.face, 0, c.uint(size))
    if err != 0 {
        FT_Done_Face(font.face)
        return {}, false
    }

    // Calculate metrics (FreeType uses 26.6 fixed point for size metrics)
    font.ascender = i32(font.face.size.metrics.ascender >> 6)
    font.descender = i32(font.face.size.metrics.descender >> 6)
    font.line_height = i32(font.face.size.metrics.height >> 6)

    font.cache = make(map[Glyph_Key]Glyph)

    return font, true
}

// Load font from memory
// size is the base logical size (will be scaled when font_set_scale is called)
font_load_from_memory :: proc(renderer: ^Text_Renderer, data: []u8, size: u32) -> (Font, bool) {
    font: Font
    font.size = size
    font.base_size = size

    err := FT_New_Memory_Face(renderer.library, raw_data(data), c.long(len(data)), 0, &font.face)
    if err != 0 {
        return {}, false
    }

    err = FT_Set_Pixel_Sizes(font.face, 0, c.uint(size))
    if err != 0 {
        FT_Done_Face(font.face)
        return {}, false
    }

    font.ascender = i32(font.face.size.metrics.ascender >> 6)
    font.descender = i32(font.face.size.metrics.descender >> 6)
    font.line_height = i32(font.face.size.metrics.height >> 6)

    font.cache = make(map[Glyph_Key]Glyph)

    return font, true
}

// Set font scale factor
// Reloads the font at base_size * scale, clearing the glyph cache
// Returns true on success
font_set_scale :: proc(font: ^Font, scale: f64) -> bool {
    if font.face == nil {
        return false
    }

    // Calculate new pixel size
    new_size := u32(f64(font.base_size) * scale + 0.5)  // Round to nearest
    if new_size < 1 {
        new_size = 1
    }

    // Skip if size unchanged
    if new_size == font.size {
        return true
    }

    // Clear glyph cache (bitmaps are for old size)
    for _, glyph in font.cache {
        delete(glyph.bitmap)
    }
    delete(font.cache)
    font.cache = make(map[Glyph_Key]Glyph)

    // Set new pixel size
    err := FT_Set_Pixel_Sizes(font.face, 0, c.uint(new_size))
    if err != 0 {
        return false
    }

    // Update metrics
    font.size = new_size
    font.ascender = i32(font.face.size.metrics.ascender >> 6)
    font.descender = i32(font.face.size.metrics.descender >> 6)
    font.line_height = i32(font.face.size.metrics.height >> 6)

    return true
}

// Destroy font
font_destroy :: proc(font: ^Font) {
    // Free cached glyph bitmaps
    for _, glyph in font.cache {
        delete(glyph.bitmap)
    }
    delete(font.cache)

    if font.face != nil {
        FT_Done_Face(font.face)
        font.face = nil
    }
}

// Get or load a glyph
font_get_glyph :: proc(font: ^Font, codepoint: rune) -> (Glyph, bool) {
    key := Glyph_Key{codepoint, font.size}

    // Check cache
    if glyph, ok := font.cache[key]; ok {
        return glyph, true
    }

    // Load glyph
    err := FT_Load_Char(font.face, c.ulong(codepoint), FT_LOAD_RENDER)
    if err != 0 {
        return {}, false
    }

    slot := font.face.glyph
    if slot == nil {
        return {}, false
    }

    // Copy bitmap data
    bitmap_size := int(slot.bitmap.rows) * int(slot.bitmap.width)
    bitmap_copy: []u8
    if bitmap_size > 0 {
        bitmap_copy = make([]u8, bitmap_size)
        // Handle pitch (row stride) which might differ from width
        for row in 0 ..< int(slot.bitmap.rows) {
            src_offset := row * int(slot.bitmap.pitch)
            dst_offset := row * int(slot.bitmap.width)
            for col in 0 ..< int(slot.bitmap.width) {
                bitmap_copy[dst_offset + col] = slot.bitmap.buffer[src_offset + col]
            }
        }
    }

    glyph := Glyph{
        bitmap    = bitmap_copy,
        width     = i32(slot.bitmap.width),
        height    = i32(slot.bitmap.rows),
        bearing_x = i32(slot.bitmap_left),
        bearing_y = i32(slot.bitmap_top),
        advance   = i32(slot.advance.x >> 6), // Convert from 26.6 fixed point
    }

    font.cache[key] = glyph
    return glyph, true
}

// Get kerning between two characters
font_get_kerning :: proc(font: ^Font, left, right: rune) -> i32 {
    if !face_has_kerning(font.face) {
        return 0
    }

    left_idx := FT_Get_Char_Index(font.face, c.ulong(left))
    right_idx := FT_Get_Char_Index(font.face, c.ulong(right))

    kerning: FT_Vector
    err := FT_Get_Kerning(font.face, left_idx, right_idx, .DEFAULT, &kerning)
    if err != 0 {
        return 0
    }

    return i32(kerning.x >> 6) // Convert from 26.6 fixed point
}

// Measure text width
text_measure :: proc(font: ^Font, text: string) -> i32 {
    width: i32 = 0
    prev_char: rune = 0

    for ch in text {
        // Add kerning
        if prev_char != 0 {
            width += font_get_kerning(font, prev_char, ch)
        }

        glyph, ok := font_get_glyph(font, ch)
        if ok {
            width += glyph.advance
        }

        prev_char = ch
    }

    return width
}

// Get logical line height (converts from physical pixels)
font_get_logical_line_height :: proc(font: ^Font) -> i32 {
    if font.size == 0 || font.base_size == 0 {
        return font.line_height
    }
    // Convert physical line_height to logical using base_size/size ratio
    return i32(f32(font.line_height) * f32(font.base_size) / f32(font.size) + 0.5)
}

// Measure text width in logical pixels
text_measure_logical :: proc(font: ^Font, text: string) -> i32 {
    if font.size == 0 || font.base_size == 0 {
        return text_measure(font, text)
    }
    physical_width := text_measure(font, text)
    return i32(f32(physical_width) * f32(font.base_size) / f32(font.size) + 0.5)
}

// Measure text size (width and height) in logical pixels
text_measure_size :: proc(font: ^Font, text: string) -> core.Size {
    return core.Size{
        width  = text_measure_logical(font, text),
        height = font_get_logical_line_height(font),
    }
}

// Internal: draw text at physical (already-scaled) coordinates
@(private)
draw_text_phys :: proc(ctx: ^Draw_Context, font: ^Font, text: string, phys_x, phys_y: i32, color: core.Color) {
    pen_x := phys_x
    pen_y := phys_y
    prev_char: rune = 0

    for ch in text {
        // Add kerning (in physical units from scaled font)
        if prev_char != 0 {
            pen_x += font_get_kerning(font, prev_char, ch)
        }

        glyph, ok := font_get_glyph(font, ch)
        if !ok {
            prev_char = ch
            continue
        }

        // Calculate glyph position (all in physical coordinates)
        gx := pen_x + glyph.bearing_x
        gy := pen_y - glyph.bearing_y

        // Render glyph
        draw_glyph(ctx, &glyph, gx, gy, color)

        pen_x += glyph.advance
        prev_char = ch
    }
}

// Draw text at position
// x, y is the baseline position (left edge, baseline) in logical coordinates
// Coordinates are scaled to physical pixels based on ctx.scale
draw_text :: proc(ctx: ^Draw_Context, font: ^Font, text: string, x, y: i32, color: core.Color) {
    draw_text_phys(ctx, font, text, scale_coord(ctx, x), scale_coord(ctx, y), color)
}

// Draw text with top-left origin (more intuitive for UI)
// x, y are in logical coordinates
draw_text_top :: proc(ctx: ^Draw_Context, font: ^Font, text: string, x, y: i32, color: core.Color) {
    // Scale to physical, add ascender (physical units from scaled font) for baseline
    phys_baseline_y := scale_coord(ctx, y) + font.ascender
    draw_text_phys(ctx, font, text, scale_coord(ctx, x), phys_baseline_y, color)
}

// Draw a single glyph with alpha blending
draw_glyph :: proc(ctx: ^Draw_Context, glyph: ^Glyph, x, y: i32, color: core.Color) {
    if len(glyph.bitmap) == 0 {
        return
    }

    stride_pixels := ctx.stride / 4

    for row in 0 ..< glyph.height {
        py := y + row
        if py < ctx.clip.y || py >= ctx.clip.y + ctx.clip.height {
            continue
        }

        for col in 0 ..< glyph.width {
            px := x + col
            if px < ctx.clip.x || px >= ctx.clip.x + ctx.clip.width {
                continue
            }

            // Get glyph alpha
            alpha := glyph.bitmap[row * glyph.width + col]
            if alpha == 0 {
                continue
            }

            // Blend with text color
            // For premultiplied alpha: out = src * alpha + dst * (1 - alpha)
            text_alpha := (u32(color.a) * u32(alpha)) / 255

            if text_alpha == 0 {
                continue
            }

            idx := py * stride_pixels + px
            dst := ctx.pixels[idx]

            // Premultiply color components with combined alpha
            src_r := (u32(color.r) * u32(alpha)) / 255
            src_g := (u32(color.g) * u32(alpha)) / 255
            src_b := (u32(color.b) * u32(alpha)) / 255

            // Extract destination components
            dst_a := (dst >> 24) & 0xFF
            dst_r := (dst >> 16) & 0xFF
            dst_g := (dst >> 8) & 0xFF
            dst_b := dst & 0xFF

            // Blend
            inv_alpha := 255 - text_alpha
            out_a := text_alpha + (dst_a * inv_alpha) / 255
            out_r := src_r + (dst_r * inv_alpha) / 255
            out_g := src_g + (dst_g * inv_alpha) / 255
            out_b := src_b + (dst_b * inv_alpha) / 255

            ctx.pixels[idx] = (out_a << 24) | (out_r << 16) | (out_g << 8) | out_b
        }
    }
}
