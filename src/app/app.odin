package app

// High-level API for odek UI toolkit
// Provides GTK/Qt-style convenience - just create app, add widgets, run.

import "../core"
import "../widgets"
import "../render"

// App encapsulates all state - no global variables needed by user
App :: struct {
    // Core Wayland state
    core_app: ^core.App,
    window:   ^core.Window,

    // Text rendering (auto-initialized)
    text_renderer: render.Text_Renderer,
    font:          render.Font,
    font_loaded:   bool,

    // Widget tree
    root: ^widgets.Container,

    // State management (auto-handled)
    focus_manager: widgets.Focus_Manager,
    hit_state:     widgets.Hit_Test_State,
}

// Common font paths to try
DEFAULT_FONT_PATHS :: []string{
    // Arch Linux / TTF
    "/usr/share/fonts/TTF/DejaVuSans.ttf",
    "/usr/share/fonts/TTF/LiberationSans-Regular.ttf",
    "/usr/share/fonts/noto/NotoSans-Regular.ttf",
    // Debian/Ubuntu
    "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
    "/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
    "/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",
    // Fedora
    "/usr/share/fonts/dejavu-sans-fonts/DejaVuSans.ttf",
    "/usr/share/fonts/liberation-sans/LiberationSans-Regular.ttf",
    // Generic
    "/usr/share/fonts/TTF/FreeSans.ttf",
    "/usr/share/fonts/truetype/freefont/FreeSans.ttf",
}

// Global app pointer for callbacks (internal use only)
// This is necessary because Wayland callbacks can't capture state
@(private)
g_app: ^App

// Create a new application
create :: proc(title: string, width: i32 = 800, height: i32 = 600) -> ^App {
    a := new(App)

    // Initialize Wayland
    a.core_app = core.init()
    if a.core_app == nil {
        free(a)
        return nil
    }

    // Create window
    a.window = core.create_window(a.core_app, title, width, height)
    if a.window == nil {
        core.shutdown(a.core_app)
        free(a)
        return nil
    }

    // Initialize text renderer
    ok: bool
    a.text_renderer, ok = render.text_renderer_init()
    if !ok {
        core.shutdown(a.core_app)
        free(a)
        return nil
    }

    // Auto-load font using system-configured size from fontconfig
    font_size := render.fc_get_default_pixel_size(14)
    a.font_loaded = false
    for path in DEFAULT_FONT_PATHS {
        a.font, ok = render.font_load(&a.text_renderer, path, font_size)
        if ok {
            a.font_loaded = true
            break
        }
    }

    // Create root container
    theme := widgets.theme_get()
    a.root = widgets.container_create(.Column)
    a.root.background = theme.bg_primary

    // Initialize focus manager
    a.focus_manager = widgets.focus_manager_init(a.root)

    // Set up callbacks (store app pointer for callbacks)
    g_app = a
    a.window.on_draw = _draw_callback
    a.window.on_close = _close_callback
    a.window.on_pointer_enter = _pointer_enter_callback
    a.window.on_pointer_leave = _pointer_leave_callback
    a.window.on_pointer_motion = _pointer_motion_callback
    a.window.on_pointer_button = _pointer_button_callback
    a.window.on_scroll = _scroll_callback
    a.window.on_key = _key_callback
    a.window.on_scale_changed = _scale_changed_callback

    return a
}

// Destroy the application
destroy :: proc(a: ^App) {
    if a == nil {
        return
    }

    if a.root != nil {
        widgets.widget_destroy(a.root)
    }

    if a.font_loaded {
        render.font_destroy(&a.font)
    }

    render.text_renderer_destroy(&a.text_renderer)

    if a.core_app != nil {
        core.shutdown(a.core_app)
    }

    if g_app == a {
        g_app = nil
    }

    free(a)
}

// Run the application event loop
run :: proc(a: ^App) {
    if a == nil || a.core_app == nil {
        return
    }
    core.run(a.core_app)
}

// Get the app's font (for custom widgets)
get_font :: proc(a: ^App) -> ^render.Font {
    if a == nil || !a.font_loaded {
        return nil
    }
    return &a.font
}

// ============================================================================
// Widget Factory Functions
// ============================================================================

// Create a label with the app's font
label :: proc(a: ^App, text: string = "") -> ^widgets.Label {
    l := widgets.label_create(text, get_font(a))
    widgets.widget_add_child(a.root, l)
    return l
}

// Create a button with the app's font
button :: proc(a: ^App, text: string = "") -> ^widgets.Button {
    b := widgets.button_create(text, get_font(a))
    widgets.widget_add_child(a.root, b)
    return b
}

// Create a text input with the app's font
text_input :: proc(a: ^App) -> ^widgets.Text_Input {
    ti := widgets.text_input_create(get_font(a))
    ti.min_size = core.Size{0, 32}
    widgets.widget_add_child(a.root, ti)
    return ti
}

// Create a container
container :: proc(a: ^App, direction: widgets.Direction = .Column) -> ^widgets.Container {
    c := widgets.container_create(direction)
    widgets.widget_add_child(a.root, c)
    return c
}

// Create a scroll container
scroll_container :: proc(a: ^App, direction: widgets.Scroll_Direction = .Vertical) -> ^widgets.Scroll_Container {
    sc := widgets.scroll_container_create(direction)
    widgets.widget_add_child(a.root, sc)
    return sc
}

// ============================================================================
// Layout Helpers
// ============================================================================

// Arrange children in a column layout
column :: proc(a: ^App, children: []^widgets.Widget, spacing: i32 = 10, padding: i32 = 20) {
    a.root.spacing = spacing
    a.root.padding = widgets.edges_all(padding)
    widgets.container_set_direction(a.root, .Column)

    // Clear existing children and add new ones
    for len(a.root.children) > 0 {
        widgets.widget_remove_child(a.root, a.root.children[0])
    }

    for child in children {
        if child != nil {
            widgets.widget_add_child(a.root, child)
        }
    }

    // Re-init focus manager with new tree
    a.focus_manager = widgets.focus_manager_init(a.root)
}

// Arrange children in a row layout
row :: proc(a: ^App, children: []^widgets.Widget, spacing: i32 = 10, padding: i32 = 20) {
    a.root.spacing = spacing
    a.root.padding = widgets.edges_all(padding)
    widgets.container_set_direction(a.root, .Row)

    // Clear existing children and add new ones
    for len(a.root.children) > 0 {
        widgets.widget_remove_child(a.root, a.root.children[0])
    }

    for child in children {
        if child != nil {
            widgets.widget_add_child(a.root, child)
        }
    }

    // Re-init focus manager with new tree
    a.focus_manager = widgets.focus_manager_init(a.root)
}

// ============================================================================
// Internal Callbacks
// ============================================================================

@(private)
_request_redraw :: proc() {
    if g_app != nil && g_app.window != nil {
        core.window_request_redraw(g_app.window)
    }
}

@(private)
_draw_callback :: proc(win: ^core.Window, pixels: [^]u32, w, h, stride: i32) {
    if g_app == nil {
        return
    }

    ctx := render.context_create_scaled(
        pixels, w, h, stride,
        win.logical_width, win.logical_height, win.scale)

    theme := widgets.theme_get()
    render.clear(&ctx, theme.bg_primary)

    g_app.root.rect = core.Rect{0, 0, win.logical_width, win.logical_height}
    widgets.widget_layout(g_app.root)
    widgets.widget_draw(g_app.root, &ctx)
}

@(private)
_close_callback :: proc(win: ^core.Window) {
    if g_app != nil && g_app.core_app != nil {
        g_app.core_app.running = false
    }
}

@(private)
_pointer_enter_callback :: proc(win: ^core.Window, x, y: f64) {
    // Nothing special needed
}

@(private)
_pointer_leave_callback :: proc(win: ^core.Window) {
    if g_app == nil {
        return
    }
    // Clear hover state
    widgets.update_hover(&g_app.hit_state, g_app.root, -1000, -1000)
    core.window_request_redraw(win)
}

@(private)
_pointer_motion_callback :: proc(win: ^core.Window, x, y: f64) {
    if g_app == nil {
        return
    }

    hovered := widgets.update_hover(&g_app.hit_state, g_app.root, i32(x), i32(y))

    // Update cursor based on hovered widget
    if hovered != nil {
        #partial switch hovered.cursor {
        case .Hand:
            core.set_cursor(g_app.core_app, .Hand)
        case .Text:
            core.set_cursor(g_app.core_app, .Text)
        case:
            core.set_cursor(g_app.core_app, .Arrow)
        }
    } else {
        core.set_cursor(g_app.core_app, .Arrow)
    }

    core.window_request_redraw(win)
}

@(private)
_pointer_button_callback :: proc(win: ^core.Window, button: u32, pressed: bool) {
    if g_app == nil {
        return
    }

    x, y := core.get_pointer_pos(g_app.core_app)
    event := core.event_pointer_button(core.Mouse_Button(button), pressed, i32(x), i32(y), 0)
    widgets.dispatch_pointer_event(&g_app.hit_state, g_app.root, &event)

    // Focus clicked widget if focusable
    if pressed && button == u32(core.Mouse_Button.Left) {
        hovered := g_app.hit_state.hovered
        if hovered != nil && hovered.focusable {
            widgets.focus_set(&g_app.focus_manager, hovered)
        }
    }

    core.window_request_redraw(win)
}

@(private)
_scroll_callback :: proc(win: ^core.Window, delta: i32, axis: u32) {
    if g_app == nil {
        return
    }

    x, y := core.get_pointer_pos(g_app.core_app)
    event := core.event_scroll(delta, axis, i32(x), i32(y))
    widgets.dispatch_pointer_event(&g_app.hit_state, g_app.root, &event)

    core.window_request_redraw(win)
}

@(private)
_key_callback :: proc(win: ^core.Window, keycode: u32, pressed: bool, utf8: string) {
    if g_app == nil || !pressed {
        return
    }

    event := core.Event{
        type = .Key_Press,
        keycode = keycode,
    }

    // Handle Tab for focus navigation
    if widgets.focus_handle_tab(&g_app.focus_manager, &event) {
        core.window_request_redraw(win)
        return
    }

    // Send to focused widget
    focused := widgets.focus_get(&g_app.focus_manager)
    if focused != nil {
        // Add keysym for text input
        // Simple ASCII mapping for printable characters
        if keycode >= 2 && keycode <= 11 {
            // Number keys 1-0
            num := (keycode - 1) % 10
            event.keysym = u32('0') + u32(num)
        } else if keycode >= 16 && keycode <= 25 {
            // QWERTY row
            chars := "qwertyuiop"
            event.keysym = u32(chars[keycode - 16])
        } else if keycode >= 30 && keycode <= 38 {
            // ASDF row
            chars := "asdfghjkl"
            event.keysym = u32(chars[keycode - 30])
        } else if keycode >= 44 && keycode <= 50 {
            // ZXCV row
            chars := "zxcvbnm"
            event.keysym = u32(chars[keycode - 44])
        } else if keycode == u32(core.Keycode.Space) {
            event.keysym = u32(' ')
        }

        widgets.widget_handle_event(focused, &event)
        core.window_request_redraw(win)
    }
}

@(private)
_scale_changed_callback :: proc(win: ^core.Window, scale: f64) {
    if g_app == nil || !g_app.font_loaded {
        return
    }
    render.font_set_scale(&g_app.font, scale)
}
