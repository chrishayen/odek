package core

import wl "../wayland"
import "base:runtime"
import "core:c"
import "core:fmt"
import "core:strings"

// Application state
App :: struct {
    // Wayland core
    display:    ^wl.Wl_Display,
    registry:   ^wl.Wl_Registry,
    compositor: ^wl.Wl_Compositor,
    shm:        ^wl.Wl_Shm,
    xdg_wm_base: ^wl.Xdg_Wm_Base,

    // Listeners
    registry_listener: wl.Wl_Registry_Listener,
    xdg_wm_base_listener: wl.Xdg_Wm_Base_Listener,

    // State
    running: bool,
    shm_formats: [dynamic]wl.Wl_Shm_Format,

    // Windows
    windows: [dynamic]^Window,
}

// Window state
Window :: struct {
    app: ^App,

    // Wayland objects
    surface: ^wl.Wl_Surface,
    xdg_surface: ^wl.Xdg_Surface,
    xdg_toplevel: ^wl.Xdg_Toplevel,

    // Listeners
    xdg_surface_listener: wl.Xdg_Surface_Listener,
    xdg_toplevel_listener: wl.Xdg_Toplevel_Listener,
    frame_listener: wl.Wl_Callback_Listener,

    // Buffers
    pool: ^wl.Shm_Pool,
    buffers: [2]^wl.Buffer,
    current_buffer: int,

    // Dimensions
    width: i32,
    height: i32,
    configured: bool,
    pending_width: i32,
    pending_height: i32,

    // State
    closed: bool,
    needs_redraw: bool,

    // Callbacks
    on_draw: proc(win: ^Window, pixels: [^]u32, width, height, stride: i32),
    on_close: proc(win: ^Window),
    user_data: rawptr,
}

// Initialize the application
init :: proc() -> ^App {
    app := new(App)

    // Connect to Wayland display
    app.display = wl.wl_display_connect(nil)
    if app.display == nil {
        fmt.eprintln("Failed to connect to Wayland display")
        free(app)
        return nil
    }

    // Get registry
    app.registry = wl.display_get_registry(app.display)
    if app.registry == nil {
        fmt.eprintln("Failed to get Wayland registry")
        wl.wl_display_disconnect(app.display)
        free(app)
        return nil
    }

    // Set up registry listener
    app.registry_listener = wl.Wl_Registry_Listener{
        global = registry_global_handler,
        global_remove = registry_global_remove_handler,
    }
    wl.registry_add_listener(app.registry, &app.registry_listener, app)

    // Set up xdg_wm_base listener
    app.xdg_wm_base_listener = wl.Xdg_Wm_Base_Listener{
        ping = xdg_wm_base_ping_handler,
    }

    // Roundtrip to get globals
    wl.wl_display_roundtrip(app.display)


    // Check we have required globals
    if app.compositor == nil {
        fmt.eprintln("Compositor not available")
        shutdown(app)
        return nil
    }
    if app.shm == nil {
        fmt.eprintln("SHM not available")
        shutdown(app)
        return nil
    }
    if app.xdg_wm_base == nil {
        fmt.eprintln("XDG WM Base not available")
        shutdown(app)
        return nil
    }

    app.running = true
    return app
}

// Shutdown the application
shutdown :: proc(app: ^App) {
    if app == nil {
        return
    }

    // Destroy windows
    for win in app.windows {
        window_destroy(win)
    }
    delete(app.windows)

    // Destroy globals
    if app.xdg_wm_base != nil {
        wl.xdg_wm_base_destroy(app.xdg_wm_base)
    }

    delete(app.shm_formats)

    if app.registry != nil {
        wl.wl_proxy_destroy(app.registry)
    }

    if app.display != nil {
        wl.wl_display_disconnect(app.display)
    }

    free(app)
}

// Create a window
create_window :: proc(app: ^App, title: string, width, height: i32) -> ^Window {
    win := new(Window)
    win.app = app
    win.width = width
    win.height = height
    win.needs_redraw = true

    // Create surface
    win.surface = wl.compositor_create_surface(app.compositor)
    if win.surface == nil {
        fmt.eprintln("Failed to create surface")
        free(win)
        return nil
    }

    // Create xdg_surface
    win.xdg_surface = wl.xdg_wm_base_get_xdg_surface(app.xdg_wm_base, win.surface)
    if win.xdg_surface == nil {
        fmt.eprintln("Failed to create xdg_surface")
        wl.surface_destroy(win.surface)
        free(win)
        return nil
    }

    // Set up xdg_surface listener
    win.xdg_surface_listener = wl.Xdg_Surface_Listener{
        configure = xdg_surface_configure_handler,
    }
    wl.xdg_surface_add_listener(win.xdg_surface, &win.xdg_surface_listener, win)

    // Create xdg_toplevel
    win.xdg_toplevel = wl.xdg_surface_get_toplevel(win.xdg_surface)
    if win.xdg_toplevel == nil {
        fmt.eprintln("Failed to create xdg_toplevel")
        wl.xdg_surface_destroy(win.xdg_surface)
        wl.surface_destroy(win.surface)
        free(win)
        return nil
    }

    // Set up xdg_toplevel listener
    win.xdg_toplevel_listener = wl.Xdg_Toplevel_Listener{
        configure = xdg_toplevel_configure_handler,
        close = xdg_toplevel_close_handler,
        configure_bounds = xdg_toplevel_configure_bounds_handler,
        wm_capabilities = xdg_toplevel_wm_capabilities_handler,
    }
    wl.xdg_toplevel_add_listener(win.xdg_toplevel, &win.xdg_toplevel_listener, win)

    // Set title
    title_cstr := strings.clone_to_cstring(title)
    defer delete(title_cstr)
    wl.xdg_toplevel_set_title(win.xdg_toplevel, title_cstr)
    wl.xdg_toplevel_set_app_id(win.xdg_toplevel, "odek")

    // Frame callback listener
    win.frame_listener = wl.Wl_Callback_Listener{
        done = frame_done_handler,
    }

    // Commit to trigger configure
    wl.surface_commit(win.surface)

    // Wait for configure - this will update win.width/height to compositor's size
    wl.wl_display_roundtrip(app.display)

    // Create SHM pool and buffers using the configured size
    buffer_size := int(win.width * win.height * 4)
    pool, ok := wl.shm_pool_create(app.shm, buffer_size)
    if !ok {
        fmt.eprintln("Failed to create SHM pool")
        wl.xdg_toplevel_destroy(win.xdg_toplevel)
        wl.xdg_surface_destroy(win.xdg_surface)
        wl.surface_destroy(win.surface)
        free(win)
        return nil
    }
    win.pool = pool

    // Create buffer at configured size
    buf1, ok1 := wl.buffer_create(pool, win.width, win.height, .ARGB8888)
    if !ok1 {
        fmt.eprintln("Failed to create buffer 1")
        wl.shm_pool_destroy(pool)
        wl.xdg_toplevel_destroy(win.xdg_toplevel)
        wl.xdg_surface_destroy(win.xdg_surface)
        wl.surface_destroy(win.surface)
        free(win)
        return nil
    }
    win.buffers[0] = buf1

    append(&app.windows, win)
    return win
}

// Destroy a window
window_destroy :: proc(win: ^Window) {
    if win == nil {
        return
    }

    for buf in win.buffers {
        if buf != nil {
            wl.buffer_destroy_internal(buf)
        }
    }

    if win.pool != nil {
        wl.shm_pool_destroy(win.pool)
    }

    if win.xdg_toplevel != nil {
        wl.xdg_toplevel_destroy(win.xdg_toplevel)
    }

    if win.xdg_surface != nil {
        wl.xdg_surface_destroy(win.xdg_surface)
    }

    if win.surface != nil {
        wl.surface_destroy(win.surface)
    }

    free(win)
}

// Get the current buffer for drawing
window_get_buffer :: proc(win: ^Window) -> (pixels: [^]u32, width, height, stride: i32) {
    buf := win.buffers[win.current_buffer]
    if buf == nil {
        return nil, 0, 0, 0
    }
    return buf.data, win.width, win.height, win.width * 4
}

// Request a redraw
window_request_redraw :: proc(win: ^Window) {
    win.needs_redraw = true
}

// Draw and present the window
window_present :: proc(win: ^Window) {
    buf := win.buffers[win.current_buffer]
    if buf == nil || buf.busy {
        return
    }

    // Attach buffer
    wl.surface_attach(win.surface, buf.wl_buffer, 0, 0)
    wl.surface_damage(win.surface, 0, 0, win.width, win.height)

    // Request frame callback
    callback := wl.surface_frame(win.surface)
    wl.callback_add_listener(callback, &win.frame_listener, win)

    // Commit
    wl.surface_commit(win.surface)

    buf.busy = true
    win.needs_redraw = false
}

// Resize the window buffer
window_resize_buffer :: proc(win: ^Window, new_width, new_height: i32) {
    // Destroy old buffer
    if win.buffers[0] != nil {
        wl.buffer_destroy_internal(win.buffers[0])
        win.buffers[0] = nil
    }

    // Destroy old pool
    if win.pool != nil {
        wl.shm_pool_destroy(win.pool)
        win.pool = nil
    }

    // Update size
    win.width = new_width
    win.height = new_height

    // Create new pool and buffer at new size
    buffer_size := int(new_width * new_height * 4)
    pool, ok := wl.shm_pool_create(win.app.shm, buffer_size)
    if !ok {
        return
    }
    win.pool = pool

    buf, ok2 := wl.buffer_create(pool, new_width, new_height, .ARGB8888)
    if !ok2 {
        wl.shm_pool_destroy(pool)
        win.pool = nil
        return
    }
    win.buffers[0] = buf
}

// Run the main event loop
run :: proc(app: ^App) {
    for app.running {
        // Process Wayland events
        if wl.wl_display_dispatch(app.display) < 0 {
            break
        }

        // Check for closed windows
        all_closed := true
        for win in app.windows {
            if !win.closed {
                all_closed = false
                break
            }
        }
        if all_closed && len(app.windows) > 0 {
            app.running = false
        }
    }
}

// Poll events without blocking
poll :: proc(app: ^App) -> bool {
    wl.wl_display_flush(app.display)
    return wl.wl_display_dispatch_pending(app.display) >= 0
}

// Stop the application
quit :: proc(app: ^App) {
    app.running = false
}

// ============================================================================
// Internal handlers
// ============================================================================

registry_global_handler :: proc "c" (
    data: rawptr,
    registry: ^wl.Wl_Registry,
    name: u32,
    interface: cstring,
    version: u32,
) {
    context = runtime.default_context()
    app := cast(^App)data
    iface := string(interface)

    if iface == "wl_compositor" {
        app.compositor = cast(^wl.Wl_Compositor)wl.registry_bind(
            registry, name, &wl.wl_compositor_interface, min(version, 6))
    } else if iface == "wl_shm" {
        app.shm = cast(^wl.Wl_Shm)wl.registry_bind(
            registry, name, &wl.wl_shm_interface, min(version, 1))
    } else if iface == "xdg_wm_base" {
        app.xdg_wm_base = cast(^wl.Xdg_Wm_Base)wl.registry_bind(
            registry, name, &wl.xdg_wm_base_interface, min(version, 6))
        wl.xdg_wm_base_add_listener(app.xdg_wm_base, &app.xdg_wm_base_listener, app)
    }
}

registry_global_remove_handler :: proc "c" (data: rawptr, registry: ^wl.Wl_Registry, name: u32) {
    // Handle global removal if needed
}

xdg_wm_base_ping_handler :: proc "c" (data: rawptr, xdg_wm_base: ^wl.Xdg_Wm_Base, serial: u32) {
    context = runtime.default_context()
    wl.xdg_wm_base_pong(xdg_wm_base, serial)
}

xdg_surface_configure_handler :: proc "c" (data: rawptr, xdg_surface: ^wl.Xdg_Surface, serial: u32) {
    context = runtime.default_context()
    win := cast(^Window)data

    // Check if we need to resize
    new_width := win.pending_width if win.pending_width > 0 else win.width
    new_height := win.pending_height if win.pending_height > 0 else win.height

    // Only resize if buffer exists and size actually changed
    if win.buffers[0] != nil && (new_width != win.width || new_height != win.height) {
        // Can't resize while buffer is busy - store pending size for later
        buf := win.buffers[0]
        if buf.busy {
            win.pending_width = new_width
            win.pending_height = new_height
        } else {
            window_resize_buffer(win, new_width, new_height)
        }
    } else if win.buffers[0] == nil {
        // First configure - just update the size, buffer will be created in create_window
        win.width = new_width
        win.height = new_height
    }

    win.configured = true
    wl.xdg_surface_ack_configure(xdg_surface, serial)

    // Trigger redraw if callback is set
    if win.on_draw != nil && win.buffers[0] != nil {
        pixels, width, height, stride := window_get_buffer(win)
        if pixels != nil {
            win.on_draw(win, pixels, width, height, stride)
            window_present(win)
        }
    }
}

xdg_toplevel_configure_handler :: proc "c" (
    data: rawptr,
    xdg_toplevel: ^wl.Xdg_Toplevel,
    width: i32,
    height: i32,
    states: ^wl.Wl_Array,
) {
    win := cast(^Window)data

    // Store compositor's suggested size (0 means "you choose")
    if width > 0 && height > 0 {
        win.pending_width = width
        win.pending_height = height
    }
}

xdg_toplevel_close_handler :: proc "c" (data: rawptr, xdg_toplevel: ^wl.Xdg_Toplevel) {
    context = runtime.default_context()
    win := cast(^Window)data
    win.closed = true

    if win.on_close != nil {
        win.on_close(win)
    }
}

xdg_toplevel_configure_bounds_handler :: proc "c" (data: rawptr, xdg_toplevel: ^wl.Xdg_Toplevel, width: i32, height: i32) {
    // Optional: could store max suggested size
}

xdg_toplevel_wm_capabilities_handler :: proc "c" (data: rawptr, xdg_toplevel: ^wl.Xdg_Toplevel, capabilities: ^wl.Wl_Array) {
    // Optional: could check what window manager capabilities are available
}

frame_done_handler :: proc "c" (data: rawptr, callback: ^wl.Wl_Callback, time: u32) {
    context = runtime.default_context()
    win := cast(^Window)data
    wl.callback_destroy(callback)

    // Mark buffer as not busy
    buf := win.buffers[win.current_buffer]
    if buf != nil {
        buf.busy = false
    }

    // Check if we have a deferred resize
    if win.pending_width > 0 && win.pending_height > 0 &&
       (win.pending_width != win.width || win.pending_height != win.height) {
        window_resize_buffer(win, win.pending_width, win.pending_height)
        win.needs_redraw = true
    }

    // Redraw if needed
    if win.needs_redraw && win.on_draw != nil {
        pixels, width, height, stride := window_get_buffer(win)
        if pixels != nil {
            win.on_draw(win, pixels, width, height, stride)
            window_present(win)
        }
    }
}
