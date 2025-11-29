package main

import "core"
import "render"
import "core:fmt"

main :: proc() {
    fmt.println("Odek - Odin Native Wayland Toolkit")

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

    // Set up draw callback
    window.on_draw = draw_callback
    window.on_close = close_callback

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

    // Draw some example content
    center_x := width / 2
    center_y := height / 2

    // Draw a rounded button
    button_rect := core.Rect{center_x - 100, center_y - 25, 200, 50}
    render.fill_rounded_rect(&ctx, button_rect, 8, core.color_hex(0x4A90D9))

    // Draw a border around the button
    render.draw_rect(&ctx, button_rect, core.color_hex(0x3A7BC8), 2)

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
