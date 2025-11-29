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
