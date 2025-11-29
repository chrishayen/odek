package main

import "core"
import "render"
import "core:fmt"

// Global state for text rendering
g_text_renderer: render.Text_Renderer
g_font: render.Font
g_font_loaded: bool

main :: proc() {
    fmt.println("Odek - Odin Native Wayland Toolkit")

    // Initialize text rendering
    text_renderer, text_ok := render.text_renderer_init()
    if !text_ok {
        fmt.eprintln("Failed to initialize text renderer")
        return
    }
    g_text_renderer = text_renderer
    defer render.text_renderer_destroy(&g_text_renderer)

    // Try to load a font
    font, font_ok := render.font_load(&g_text_renderer, "/usr/share/fonts/noto/NotoSans-Regular.ttf\x00", 16)
    if font_ok {
        g_font = font
        g_font_loaded = true
        fmt.println("Font loaded successfully")
    } else {
        fmt.eprintln("Warning: Could not load font, text will not be displayed")
    }
    defer if g_font_loaded {
        render.font_destroy(&g_font)
    }

    // Initialize application
    app := core.init()
    if app == nil {
        fmt.eprintln("Failed to initialize application")
        return
    }
    defer core.shutdown(app)

    // Create a window
    window := core.create_window(app, "Odek Demo", 800, 600)
    if window == nil {
        fmt.eprintln("Failed to create window")
        return
    }

    // Set up callbacks
    window.on_draw = draw_callback
    window.on_close = close_callback
    window.on_pointer_enter = pointer_enter_callback
    window.on_pointer_leave = pointer_leave_callback
    window.on_pointer_motion = pointer_motion_callback
    window.on_pointer_button = pointer_button_callback
    window.on_key = key_callback

    // Trigger initial draw now that callback is set
    pixels, width, height, stride := core.window_get_buffer(window)
    if pixels != nil {
        draw_callback(window, pixels, width, height, stride)
        core.window_present(window)
    }

    fmt.println("Window created, entering event loop...")

    // Run the event loop
    core.run(app)

    fmt.println("Application closed")
}

draw_callback :: proc(win: ^core.Window, pixels: [^]u32, width, height, stride: i32) {
    ctx := render.context_create(pixels, width, height, stride)

    // Background - dark gray
    render.clear(&ctx, core.color_hex(0x2D2D2D))

    // Draw a title bar area
    render.fill_rect(&ctx, core.Rect{0, 0, width, 40}, core.color_hex(0x404040))

    // Draw title text
    if g_font_loaded {
        render.draw_text_top(&ctx, &g_font, "Odek Demo", 12, 12, core.COLOR_WHITE)
    }

    // Draw some example content
    center_x := width / 2
    center_y := height / 2

    // Draw a rounded button
    button_rect := core.Rect{center_x - 100, center_y - 25, 200, 50}
    render.fill_rounded_rect(&ctx, button_rect, 8, core.color_hex(0x4A90D9))

    // Draw a border around the button
    render.draw_rect(&ctx, button_rect, core.color_hex(0x3A7BC8), 2)

    // Draw button text (centered)
    if g_font_loaded {
        button_text := "Click Me!"
        text_width := render.text_measure(&g_font, button_text)
        text_x := center_x - text_width / 2
        text_y := center_y + g_font.line_height / 4  // Approximate vertical centering
        render.draw_text(&ctx, &g_font, button_text, text_x, text_y, core.COLOR_WHITE)
    }

    // Draw info text
    if g_font_loaded {
        render.draw_text_top(&ctx, &g_font, "Press any key or move the mouse to see events in console", 12, 60, core.color_hex(0xAAAAAA))
        render.draw_text_top(&ctx, &g_font, "Text rendering with FreeType", 12, 85, core.color_hex(0x88CC88))
    }

    // Draw some decorative elements
    for i in 0 ..< 5 {
        x := 50 + i32(i) * 60
        y := height - 100
        render.fill_rounded_rect(&ctx, core.Rect{x, y, 40, 40}, 5, core.color_hex(0x5A5A5A))
    }
}

close_callback :: proc(win: ^core.Window) {
    fmt.println("Window close requested")
}

pointer_enter_callback :: proc(win: ^core.Window, x, y: f64) {
    fmt.printf("Pointer entered at (%.1f, %.1f)\n", x, y)
}

pointer_leave_callback :: proc(win: ^core.Window) {
    fmt.println("Pointer left window")
}

pointer_motion_callback :: proc(win: ^core.Window, x, y: f64) {
    // Uncomment to see motion events (very verbose)
    // fmt.printf("Motion: (%.1f, %.1f)\n", x, y)
}

pointer_button_callback :: proc(win: ^core.Window, button: u32, pressed: bool) {
    action := "pressed" if pressed else "released"
    button_name: string
    switch button {
    case 0x110: button_name = "Left"
    case 0x111: button_name = "Right"
    case 0x112: button_name = "Middle"
    case: button_name = "Unknown"
    }
    fmt.printf("%s button %s\n", button_name, action)
}

key_callback :: proc(win: ^core.Window, key: u32, pressed: bool, utf8: string) {
    if pressed {
        if len(utf8) > 0 {
            fmt.printf("Key %d pressed: '%s'\n", key, utf8)
        } else {
            fmt.printf("Key %d pressed (no printable char)\n", key)
        }
    }
}
