package core

// Basic geometric types
Rect :: struct {
    x, y: i32,
    width, height: i32,
}

Point :: struct {
    x, y: i32,
}

Size :: struct {
    width, height: i32,
}

// Color in ARGB8888 format (Wayland standard)
// Stored as premultiplied alpha for correct blending
Color :: struct {
    r, g, b, a: u8,
}

// Common color constructors
color_rgba :: proc(r, g, b, a: u8) -> Color {
    // Premultiply alpha
    alpha := f32(a) / 255.0
    return Color{
        r = u8(f32(r) * alpha),
        g = u8(f32(g) * alpha),
        b = u8(f32(b) * alpha),
        a = a,
    }
}

color_rgb :: proc(r, g, b: u8) -> Color {
    return Color{r = r, g = g, b = b, a = 255}
}

color_hex :: proc(hex: u32) -> Color {
    return Color{
        r = u8((hex >> 16) & 0xFF),
        g = u8((hex >> 8) & 0xFF),
        b = u8(hex & 0xFF),
        a = 255,
    }
}

// Convert to ARGB8888 u32 for pixel buffer
color_to_argb :: proc(c: Color) -> u32 {
    return (u32(c.a) << 24) | (u32(c.r) << 16) | (u32(c.g) << 8) | u32(c.b)
}

// Common colors
COLOR_BLACK :: Color{0, 0, 0, 255}
COLOR_WHITE :: Color{255, 255, 255, 255}
COLOR_RED :: Color{255, 0, 0, 255}
COLOR_GREEN :: Color{0, 255, 0, 255}
COLOR_BLUE :: Color{0, 0, 255, 255}
COLOR_TRANSPARENT :: Color{0, 0, 0, 0}

// Rect operations
rect_contains :: proc(r: Rect, p: Point) -> bool {
    return p.x >= r.x && p.x < r.x + r.width &&
           p.y >= r.y && p.y < r.y + r.height
}

rect_intersects :: proc(a, b: Rect) -> bool {
    return a.x < b.x + b.width && a.x + a.width > b.x &&
           a.y < b.y + b.height && a.y + a.height > b.y
}

rect_intersection :: proc(a, b: Rect) -> (Rect, bool) {
    x1 := max(a.x, b.x)
    y1 := max(a.y, b.y)
    x2 := min(a.x + a.width, b.x + b.width)
    y2 := min(a.y + a.height, b.y + b.height)

    if x2 <= x1 || y2 <= y1 {
        return Rect{}, false
    }

    return Rect{x1, y1, x2 - x1, y2 - y1}, true
}

rect_union :: proc(a, b: Rect) -> Rect {
    x1 := min(a.x, b.x)
    y1 := min(a.y, b.y)
    x2 := max(a.x + a.width, b.x + b.width)
    y2 := max(a.y + a.height, b.y + b.height)

    return Rect{x1, y1, x2 - x1, y2 - y1}
}

rect_is_empty :: proc(r: Rect) -> bool {
    return r.width <= 0 || r.height <= 0
}
