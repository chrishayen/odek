package main

import "core"
import "render"
import "widgets"
import "core:fmt"
import "core:os"
import "core:strings"
import "core:path/filepath"

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

    // Header widgets
    back_button:   ^widgets.Button,
    header_label:  ^widgets.Label,

    // Footer
    footer_label:  ^widgets.Label,

    // Image grid
    image_cache:   ^render.Image_Cache,
    image_grid:    ^widgets.Image_Grid,

    // Navigation
    current_directory:  string,
    directory_history:  [dynamic]string,
    window:             ^core.Window,  // Store window reference for callbacks
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

    // Initialize image cache
    g_state.image_cache = render.image_cache_create(50, 150)
    defer render.image_cache_destroy(g_state.image_cache)

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
    g_state.window = window

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
    window.on_scroll = scroll_callback
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
        // Back button
        g_state.back_button = widgets.button_create("<", &g_state.font)
        g_state.back_button.min_size = core.Size{36, 30}
        g_state.back_button.on_click = back_button_callback
        widgets.button_set_colors(g_state.back_button,
            core.color_hex(0x505050),  // Normal - subtle gray
            core.color_hex(0x606060),  // Hover
            core.color_hex(0x404040))  // Pressed
        widgets.widget_add_child(g_state.header, g_state.back_button)

        // Header label
        g_state.header_label = widgets.label_create("Odek Image Grid Demo", &g_state.font)
        widgets.label_set_color(g_state.header_label, core.COLOR_WHITE)
        widgets.widget_add_child(g_state.header, g_state.header_label)

        // Sidebar bookmark buttons
        Bookmark :: struct { name: string, path: cstring }
        bookmarks := []Bookmark{
            {"Home", "/home/chris"},
            {"Pictures", "/home/chris/Pictures"},
            {"Downloads", "/home/chris/Downloads"},
            {"Code", "/home/chris/Code"},
        }
        for bm in bookmarks {
            btn := widgets.button_create(bm.name, &g_state.font)
            btn.user_data = rawptr(bm.path)
            btn.on_click = bookmark_callback
            btn.min_size = core.Size{0, 32}
            widgets.button_set_colors(btn,
                core.color_hex(0x444444),  // Normal - match sidebar
                core.color_hex(0x555555),  // Hover
                core.color_hex(0x333333))  // Pressed
            widgets.widget_add_child(g_state.sidebar, btn)
        }

        // Footer label
        g_state.footer_label = widgets.label_create("Click an image to select, scroll to navigate", &g_state.font)
        widgets.label_set_color(g_state.footer_label, core.color_hex(0x888888))
        widgets.widget_add_child(g_state.footer, g_state.footer_label)
    }

    // Create image grid in main area
    g_state.image_grid = widgets.image_grid_create()
    g_state.image_grid.cell_width = 150
    g_state.image_grid.cell_height = 150
    g_state.image_grid.spacing = 10
    g_state.image_grid.padding = widgets.edges_all(10)
    g_state.image_grid.on_click = image_grid_click_callback
    g_state.image_grid.on_folder_click = folder_click_callback
    if g_state.font_loaded {
        g_state.image_grid.font = &g_state.font
    }
    widgets.widget_add_child(g_state.main_area, g_state.image_grid)

    // Load initial directory (home has folders to test navigation)
    navigate_to_directory("/home/chris", add_to_history = false)

    // Initialize focus manager
    g_state.focus_manager = widgets.focus_manager_init(g_state.root)

    // Perform initial layout
    widgets.widget_layout(g_state.root)
}

// Navigate to a directory - clears grid and loads contents
navigate_to_directory :: proc(dir_path: string, add_to_history: bool = true) {
    // Push current directory to history before navigating
    if add_to_history && len(g_state.current_directory) > 0 {
        append(&g_state.directory_history, strings.clone(g_state.current_directory))
    }

    // Update current directory
    if len(g_state.current_directory) > 0 {
        delete(g_state.current_directory)
    }
    g_state.current_directory = strings.clone(dir_path)

    // Clear the grid
    widgets.image_grid_clear(g_state.image_grid)

    handle, err := os.open(dir_path)
    if err != nil {
        fmt.eprintln("Failed to open directory:", dir_path)
        return
    }
    defer os.close(handle)

    entries, read_err := os.read_dir(handle, -1)
    if read_err != nil {
        fmt.eprintln("Failed to read directory")
        return
    }
    defer delete(entries)

    folder_count := 0
    image_count := 0

    // First pass: add folders (skip hidden)
    for entry in entries {
        if !entry.is_dir {
            continue
        }
        if strings.has_prefix(entry.name, ".") {
            continue  // Skip hidden folders
        }

        full_path := filepath.join({dir_path, entry.name})
        name_clone := strings.clone(entry.name)
        path_clone := strings.clone(full_path)
        delete(full_path)

        widgets.image_grid_add_folder(g_state.image_grid, name_clone, path_clone)
        folder_count += 1
    }

    // Second pass: add images
    for entry in entries {
        if entry.is_dir {
            continue
        }

        // Check for image extensions
        ext := strings.to_lower(filepath.ext(entry.name))
        defer delete(ext)

        if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
            full_path := filepath.join({dir_path, entry.name})
            defer delete(full_path)

            img, thumb, ok := render.image_cache_load(g_state.image_cache, full_path)
            if ok {
                path_clone := strings.clone(full_path)
                widgets.image_grid_add_item(g_state.image_grid, img, thumb, path_clone)
                image_count += 1
            }
        }
    }

    fmt.printf("Loaded %d folders, %d images from %s\n", folder_count, image_count, dir_path)

    // Update header to show current path
    if g_state.header_label != nil {
        widgets.label_set_text(g_state.header_label, dir_path)
    }

    // Request redraw
    if g_state.window != nil {
        core.window_request_redraw(g_state.window)
    }
}

// Navigate back to previous directory
navigate_back :: proc() {
    if len(g_state.directory_history) == 0 {
        return
    }

    prev := pop(&g_state.directory_history)
    navigate_to_directory(prev, add_to_history = false)
    delete(prev)
}

// Folder click callback - navigate into folder
folder_click_callback :: proc(grid: ^widgets.Image_Grid, path: string) {
    fmt.printf("Navigating to folder: %s\n", path)
    navigate_to_directory(path)
}

// Back button callback
back_button_callback :: proc(button: ^widgets.Button) {
    navigate_back()
}

// Bookmark button callback - navigate to stored path
bookmark_callback :: proc(button: ^widgets.Button) {
    path := cast(cstring)button.user_data
    if path != nil {
        navigate_to_directory(string(path))
    }
}

// Image grid click callback
image_grid_click_callback :: proc(grid: ^widgets.Image_Grid, index: i32, item: ^widgets.Grid_Item) {
    fmt.printf("Image clicked: index=%d, path=%s\n", index, item.path)
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
    // Dispatch motion to image grid (for scrollbar dragging)
    if g_state.image_grid != nil && g_state.image_grid.scrollbar_dragging {
        event := core.Event{
            type = .Pointer_Motion,
            pointer_x = i32(x),
            pointer_y = i32(y),
        }
        if widgets.widget_handle_event(g_state.image_grid, &event) {
            // Just request redraw - let frame callback handle it
            core.window_request_redraw(win)
            return
        }
    }

    // Update hover state in widget tree
    widgets.update_hover(&g_state.hit_state, g_state.root, i32(x), i32(y))
}

pointer_button_callback :: proc(win: ^core.Window, button: u32, pressed: bool) {
    x, y := core.get_pointer_pos(win.app)
    event := core.Event{
        type = pressed ? .Pointer_Button_Press : .Pointer_Button_Release,
        button = core.Mouse_Button(button),
        pointer_x = i32(x),
        pointer_y = i32(y),
    }

    // Handle scrollbar drag release
    if !pressed && g_state.image_grid != nil && g_state.image_grid.scrollbar_dragging {
        widgets.widget_handle_event(g_state.image_grid, &event)
        return
    }

    // Dispatch to widget system
    widgets.dispatch_pointer_event(&g_state.hit_state, g_state.root, &event)
}

scroll_callback :: proc(win: ^core.Window, delta: i32, axis: u32) {
    // Dispatch scroll event to image grid
    if g_state.image_grid != nil {
        x, y := core.get_pointer_pos(win.app)
        event := core.event_scroll(delta, axis, i32(x), i32(y))
        if widgets.widget_handle_event(g_state.image_grid, &event) {
            core.window_request_redraw(win)
        }
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
