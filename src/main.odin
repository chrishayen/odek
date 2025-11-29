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

    // Labels
    header_label:  ^widgets.Label,
    footer_label:  ^widgets.Label,
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
    g_state.root.align_items = .Stretch  // Children fill width

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
    g_state.main_area.align_items = .Stretch  // Labels fill width for word wrapping
    widgets.widget_add_child(content, g_state.main_area)

    // Footer
    g_state.footer = widgets.container_create(.Row)
    g_state.footer.min_size = core.Size{0, 30}
    g_state.footer.background = core.color_hex(0x404040)
    g_state.footer.padding = widgets.edges_symmetric(15, 5)
    g_state.footer.align_items = .Center
    widgets.widget_add_child(g_state.root, g_state.footer)

    // Perform initial layout to set container rects before adding labels
    widgets.widget_layout(g_state.root)

    // Add label widgets if font is loaded
    if g_state.font_loaded {
        // Header label
        g_state.header_label = widgets.label_create("Odek Widget Demo", &g_state.font)
        widgets.label_set_color(g_state.header_label, core.COLOR_WHITE)
        widgets.widget_add_child(g_state.header, g_state.header_label)

        // Sidebar labels
        sidebar_items := []string{"Dashboard", "Widgets", "Settings", "About"}
        for item in sidebar_items {
            label := widgets.label_create(item, &g_state.font)
            widgets.label_set_color(label, core.color_hex(0xCCCCCC))
            widgets.widget_add_child(g_state.sidebar, label)
        }

        // Main area labels
        title := widgets.label_create("Phase 5: MVP Widgets", &g_state.font)
        widgets.label_set_color(title, core.COLOR_WHITE)
        widgets.widget_add_child(g_state.main_area, title)

        subtitle := widgets.label_create("Labels and Buttons with interactive states", &g_state.font)
        widgets.label_set_color(subtitle, core.color_hex(0x88CC88))
        widgets.widget_add_child(g_state.main_area, subtitle)

        // Button row
        button_row := widgets.container_create(.Row)
        button_row.spacing = 10
        button_row.align_items = .Center
        widgets.widget_add_child(g_state.main_area, button_row)

        // Create some demo buttons
        btn1 := widgets.button_create("Primary", &g_state.font)
        widgets.button_set_on_click(btn1, demo_button_click)
        widgets.widget_add_child(button_row, btn1)

        btn2 := widgets.button_create("Secondary", &g_state.font)
        widgets.button_set_colors(btn2, core.color_hex(0x606060), core.color_hex(0x707070), core.color_hex(0x505050))
        widgets.button_set_on_click(btn2, demo_button_click)
        widgets.widget_add_child(button_row, btn2)

        btn3 := widgets.button_create("Success", &g_state.font)
        widgets.button_set_colors(btn3, core.color_hex(0x28A745), core.color_hex(0x38B755), core.color_hex(0x189735))
        widgets.button_set_on_click(btn3, demo_button_click)
        widgets.widget_add_child(button_row, btn3)

        btn4 := widgets.button_create("Danger", &g_state.font)
        widgets.button_set_colors(btn4, core.color_hex(0xDC3545), core.color_hex(0xEC4555), core.color_hex(0xCC2535))
        widgets.button_set_on_click(btn4, demo_button_click)
        widgets.widget_add_child(button_row, btn4)

        // Info text
        info := widgets.label_create("Buttons support hover and pressed states. Click to see callback.", &g_state.font)
        widgets.label_set_color(info, core.color_hex(0xAAAAAA))
        widgets.widget_add_child(g_state.main_area, info)

        // Footer label
        g_state.footer_label = widgets.label_create("Press any key or move mouse to see events", &g_state.font)
        widgets.label_set_color(g_state.footer_label, core.color_hex(0x888888))
        widgets.widget_add_child(g_state.footer, g_state.footer_label)
    }

    // Initialize focus manager
    g_state.focus_manager = widgets.focus_manager_init(g_state.root)

    // Perform initial layout
    widgets.widget_layout(g_state.root)
}

// Demo button click callback
demo_button_click :: proc(button: ^widgets.Button) {
    fmt.printf("Button clicked: '%s'\n", button.text)
}

draw_callback :: proc(win: ^core.Window, pixels: [^]u32, width, height, stride: i32) {
    ctx := render.context_create(pixels, width, height, stride)

    // Update root size if window resized
    if g_state.root.rect.width != width || g_state.root.rect.height != height {
        g_state.root.rect = core.Rect{0, 0, width, height}
        widgets.widget_layout(g_state.root)
    }

    // Draw widget tree (containers and labels)
    widgets.widget_draw(g_state.root, &ctx)
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
    x, y := core.get_pointer_pos(win.app)
    event := core.Event{
        type = pressed ? .Pointer_Button_Press : .Pointer_Button_Release,
        button = core.Mouse_Button(button),
        pointer_x = i32(x),
        pointer_y = i32(y),
    }
    widgets.dispatch_pointer_event(&g_state.hit_state, g_state.root, &event)
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
