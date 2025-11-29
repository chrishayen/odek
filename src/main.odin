package main

import "core"
import "render"
import "widgets"
import "core:fmt"

// Application state
App_State :: struct {
    text_renderer: render.Text_Renderer,
    font:          render.Font,
    font_loaded:   bool,

    // Widget system
    root:          ^widgets.Container,
    header:        ^widgets.Container,
    sidebar:       ^widgets.Container,
    main_area:     ^widgets.Container,
    footer:        ^widgets.Container,
    focus_manager: widgets.Focus_Manager,
    hit_state:     widgets.Hit_Test_State,
}

g_state: App_State

main :: proc() {
    fmt.println("Odek - Odin Native Wayland Toolkit")

    // Initialize text rendering
    text_renderer, text_ok := render.text_renderer_init()
    if !text_ok {
        fmt.eprintln("Failed to initialize text renderer")
        return
    }
    g_state.text_renderer = text_renderer
    defer render.text_renderer_destroy(&g_state.text_renderer)

    // Try to load a font
    font, font_ok := render.font_load(&g_state.text_renderer, "/usr/share/fonts/noto/NotoSans-Regular.ttf\x00", 16)
    if font_ok {
        g_state.font = font
        g_state.font_loaded = true
        fmt.println("Font loaded successfully")
    } else {
        fmt.eprintln("Warning: Could not load font, text will not be displayed")
    }
    defer if g_state.font_loaded {
        render.font_destroy(&g_state.font)
    }

    // Initialize application
    app := core.init()
    if app == nil {
        fmt.eprintln("Failed to initialize application")
        return
    }
    defer core.shutdown(app)

    // Create a window
    window := core.create_window(app, "Odek Widget Demo", 800, 600)
    if window == nil {
        fmt.eprintln("Failed to create window")
        return
    }

    // Build widget tree
    build_ui(window.width, window.height)
    defer widgets.widget_destroy(g_state.root)

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
    fmt.println("Widget system active with flexbox layout")

    // Run the event loop
    core.run(app)

    fmt.println("Application closed")
}

// Build the widget UI
build_ui :: proc(width, height: i32) {
    // Create root container (fills window)
    g_state.root = widgets.container_create(.Column)
    g_state.root.rect = core.Rect{0, 0, width, height}
    g_state.root.padding = widgets.edges_all(10)
    g_state.root.background = core.color_hex(0x2D2D2D)
    g_state.root.spacing = 10

    // Header container (row)
    g_state.header = widgets.container_create(.Row)
    g_state.header.min_size = core.Size{0, 50}
    g_state.header.background = core.color_hex(0x404040)
    g_state.header.padding = widgets.edges_symmetric(15, 10)
    g_state.header.align_items = .Center
    widgets.widget_add_child(g_state.root, g_state.header)

    // Content area (row with sidebar and main)
    content := widgets.container_create(.Row)
    content.flex = 1  // Take remaining space
    content.spacing = 10
    widgets.widget_add_child(g_state.root, content)

    // Sidebar (column)
    g_state.sidebar = widgets.container_create(.Column)
    g_state.sidebar.min_size = core.Size{150, 0}
    g_state.sidebar.background = core.color_hex(0x383838)
    g_state.sidebar.padding = widgets.edges_all(10)
    g_state.sidebar.spacing = 8
    widgets.widget_add_child(content, g_state.sidebar)

    // Main area (flex)
    g_state.main_area = widgets.container_create(.Column)
    g_state.main_area.flex = 1
    g_state.main_area.background = core.color_hex(0x333333)
    g_state.main_area.padding = widgets.edges_all(20)
    g_state.main_area.spacing = 15
    widgets.widget_add_child(content, g_state.main_area)

    // Footer
    g_state.footer = widgets.container_create(.Row)
    g_state.footer.min_size = core.Size{0, 30}
    g_state.footer.background = core.color_hex(0x404040)
    g_state.footer.padding = widgets.edges_symmetric(15, 5)
    g_state.footer.align_items = .Center
    widgets.widget_add_child(g_state.root, g_state.footer)

    // Initialize focus manager
    g_state.focus_manager = widgets.focus_manager_init(g_state.root)

    // Perform initial layout
    widgets.widget_layout(g_state.root)
}

draw_callback :: proc(win: ^core.Window, pixels: [^]u32, width, height, stride: i32) {
    ctx := render.context_create(pixels, width, height, stride)

    // Update root size if window resized
    if g_state.root.rect.width != width || g_state.root.rect.height != height {
        g_state.root.rect = core.Rect{0, 0, width, height}
        widgets.widget_layout(g_state.root)
    }

    // Draw widget tree (this draws backgrounds)
    widgets.widget_draw(g_state.root, &ctx)

    // Draw text overlays (since we don't have Label widgets yet)
    if g_state.font_loaded {
        // Get absolute positions of widgets
        header_rect := widgets.widget_get_absolute_rect(g_state.header)
        sidebar_rect := widgets.widget_get_absolute_rect(g_state.sidebar)
        main_rect := widgets.widget_get_absolute_rect(g_state.main_area)
        footer_rect := widgets.widget_get_absolute_rect(g_state.footer)

        // Header text (centered vertically in header)
        header_text_y := header_rect.y + (header_rect.height - g_state.font.line_height) / 2
        render.draw_text_top(&ctx, &g_state.font, "Odek Widget Demo",
            header_rect.x + g_state.header.padding.left,
            header_text_y,
            core.COLOR_WHITE)

        // Sidebar items (relative to sidebar)
        sidebar_text_x := sidebar_rect.x + g_state.sidebar.padding.left
        sidebar_text_y := sidebar_rect.y + g_state.sidebar.padding.top
        items := []string{"Dashboard", "Widgets", "Settings", "About"}
        for item in items {
            render.draw_text_top(&ctx, &g_state.font, item, sidebar_text_x, sidebar_text_y, core.color_hex(0xCCCCCC))
            sidebar_text_y += g_state.font.line_height + g_state.sidebar.spacing
        }

        // Main area text (relative to main area)
        main_text_x := main_rect.x + g_state.main_area.padding.left
        main_text_y := main_rect.y + g_state.main_area.padding.top
        line_spacing := g_state.font.line_height + 8

        render.draw_text_top(&ctx, &g_state.font, "Phase 4: Widget System", main_text_x, main_text_y, core.COLOR_WHITE)
        main_text_y += line_spacing
        render.draw_text_top(&ctx, &g_state.font, "Flexbox-lite layout with containers", main_text_x, main_text_y, core.color_hex(0x88CC88))
        main_text_y += line_spacing
        render.draw_text_top(&ctx, &g_state.font, "- Row and Column directions", main_text_x, main_text_y, core.color_hex(0xAAAAAA))
        main_text_y += line_spacing - 4
        render.draw_text_top(&ctx, &g_state.font, "- Flex grow for dynamic sizing", main_text_x, main_text_y, core.color_hex(0xAAAAAA))
        main_text_y += line_spacing - 4
        render.draw_text_top(&ctx, &g_state.font, "- Cross-axis alignment", main_text_x, main_text_y, core.color_hex(0xAAAAAA))
        main_text_y += line_spacing - 4
        render.draw_text_top(&ctx, &g_state.font, "- Padding and spacing", main_text_x, main_text_y, core.color_hex(0xAAAAAA))
        main_text_y += line_spacing + 10

        // Draw a demo button (centered in remaining main area space)
        button_width: i32 = 180
        button_height: i32 = 45
        button_x := main_rect.x + (main_rect.width - button_width) / 2
        button_y := main_text_y + 10
        button_rect := core.Rect{button_x, button_y, button_width, button_height}
        render.fill_rounded_rect(&ctx, button_rect, 6, core.color_hex(0x4A90D9))
        button_text := "Demo Button"
        text_width := render.text_measure(&g_state.font, button_text)
        text_x := button_rect.x + (button_rect.width - text_width) / 2
        text_y := button_rect.y + (button_rect.height + g_state.font.line_height) / 2 - 2
        render.draw_text_top(&ctx, &g_state.font, button_text, text_x, text_y - g_state.font.line_height + 5, core.COLOR_WHITE)

        // Footer text (centered vertically in footer)
        footer_text_y := footer_rect.y + (footer_rect.height - g_state.font.line_height) / 2
        render.draw_text_top(&ctx, &g_state.font, "Press any key or move mouse to see events",
            footer_rect.x + g_state.footer.padding.left,
            footer_text_y,
            core.color_hex(0x888888))
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
    // Update hover state in widget tree
    widgets.update_hover(&g_state.hit_state, g_state.root, i32(x), i32(y))
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

    // Dispatch to widget system
    if pressed {
        x, y := core.get_pointer_pos(win.app)
        event := core.Event{
            type = .Pointer_Button_Press,
            pointer_x = i32(x),
            pointer_y = i32(y),
        }
        widgets.dispatch_pointer_event(&g_state.hit_state, g_state.root, &event)
    }
}

key_callback :: proc(win: ^core.Window, key: u32, pressed: bool, utf8: string) {
    if pressed {
        if len(utf8) > 0 {
            fmt.printf("Key %d pressed: '%s'\n", key, utf8)
        } else {
            fmt.printf("Key %d pressed (no printable char)\n", key)
        }

        // Handle Tab for focus navigation
        event := core.Event{
            type = .Key_Press,
            keycode = key,
        }
        widgets.focus_handle_tab(&g_state.focus_manager, &event)
    }
}
