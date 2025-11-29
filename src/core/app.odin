package core

import wl "../wayland"
import "base:runtime"
import "core:c"
import "core:fmt"
import "core:strings"
import "core:sys/posix"

// Application state
App :: struct {
    // Wayland core
    display:    ^wl.Wl_Display,
    registry:   ^wl.Wl_Registry,
    compositor: ^wl.Wl_Compositor,
    shm:        ^wl.Wl_Shm,
    xdg_wm_base: ^wl.Xdg_Wm_Base,

    // Input devices
    seat:     ^wl.Wl_Seat,
    pointer:  ^wl.Wl_Pointer,
    keyboard: ^wl.Wl_Keyboard,

    // Cursor
    cursor_shape_manager: ^wl.Wp_Cursor_Shape_Manager,
    cursor_shape_device:  ^wl.Wp_Cursor_Shape_Device,

    // XKB state
    xkb: wl.Xkb_Handler,

    // Listeners
    registry_listener: wl.Wl_Registry_Listener,
    xdg_wm_base_listener: wl.Xdg_Wm_Base_Listener,
    seat_listener: wl.Wl_Seat_Listener,
    pointer_listener: wl.Wl_Pointer_Listener,
    keyboard_listener: wl.Wl_Keyboard_Listener,

    // Input state
    pointer_x: f64,
    pointer_y: f64,
    pointer_surface: ^wl.Wl_Surface,  // Surface pointer is over
    pointer_buttons: u32,              // Button state bitmask
    pointer_serial: u32,               // Last pointer enter serial

    keyboard_surface: ^wl.Wl_Surface,  // Surface with keyboard focus
    keyboard_serial: u32,              // Last keyboard enter serial
    key_repeat_rate: i32,              // Keys per second
    key_repeat_delay: i32,             // Delay in ms before repeat

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
    frame_pending: bool,  // True if waiting for frame callback

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
    on_pointer_enter: proc(win: ^Window, x, y: f64),
    on_pointer_leave: proc(win: ^Window),
    on_pointer_motion: proc(win: ^Window, x, y: f64),
    on_pointer_button: proc(win: ^Window, button: u32, pressed: bool),
    on_scroll: proc(win: ^Window, delta: i32, axis: u32),
    on_key: proc(win: ^Window, key: u32, pressed: bool, utf8: string),
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

    // Set up seat listener
    app.seat_listener = wl.Wl_Seat_Listener{
        capabilities = seat_capabilities_handler,
        name = seat_name_handler,
    }

    // Set up pointer listener
    app.pointer_listener = wl.Wl_Pointer_Listener{
        enter = pointer_enter_handler,
        leave = pointer_leave_handler,
        motion = pointer_motion_handler,
        button = pointer_button_handler,
        axis = pointer_axis_handler,
        frame = pointer_frame_handler,
        axis_source = pointer_axis_source_handler,
        axis_stop = pointer_axis_stop_handler,
        axis_discrete = pointer_axis_discrete_handler,
        axis_value120 = pointer_axis_value120_handler,
        axis_relative_direction = pointer_axis_relative_direction_handler,
    }

    // Set up keyboard listener
    app.keyboard_listener = wl.Wl_Keyboard_Listener{
        keymap = keyboard_keymap_handler,
        enter = keyboard_enter_handler,
        leave = keyboard_leave_handler,
        key = keyboard_key_handler,
        modifiers = keyboard_modifiers_handler,
        repeat_info = keyboard_repeat_info_handler,
    }

    // Initialize XKB
    xkb, xkb_ok := wl.xkb_handler_init()
    if !xkb_ok {
        fmt.eprintln("Failed to initialize XKB")
        wl.wl_display_disconnect(app.display)
        free(app)
        return nil
    }
    app.xkb = xkb

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

    // Destroy input devices
    if app.keyboard != nil {
        wl.keyboard_release(app.keyboard)
    }
    if app.pointer != nil {
        wl.pointer_release(app.pointer)
    }
    if app.seat != nil {
        wl.seat_release(app.seat)
    }

    // Destroy XKB
    wl.xkb_handler_destroy(&app.xkb)

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

    // Create SHM pool for double buffers
    buffer_size := int(win.width * win.height * 4 * 2)  // Space for 2 buffers
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

    // Create double buffers at configured size
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

    buf2, ok2 := wl.buffer_create(pool, win.width, win.height, .ARGB8888)
    if !ok2 {
        fmt.eprintln("Failed to create buffer 2")
        wl.buffer_destroy_internal(buf1)
        wl.shm_pool_destroy(pool)
        wl.xdg_toplevel_destroy(win.xdg_toplevel)
        wl.xdg_surface_destroy(win.xdg_surface)
        wl.surface_destroy(win.surface)
        free(win)
        return nil
    }
    win.buffers[1] = buf2

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

// Get a free buffer for drawing (double buffering)
window_get_buffer :: proc(win: ^Window) -> (pixels: [^]u32, width, height, stride: i32) {
    // Try current buffer first
    buf := win.buffers[win.current_buffer]
    if buf != nil && !buf.busy {
        return buf.data, win.width, win.height, win.width * 4
    }

    // Try the other buffer
    other := 1 - win.current_buffer
    buf = win.buffers[other]
    if buf != nil && !buf.busy {
        win.current_buffer = other
        return buf.data, win.width, win.height, win.width * 4
    }

    // Both buffers busy
    return nil, 0, 0, 0
}

// Check if window can be drawn to (at least one buffer not busy)
window_can_draw :: proc(win: ^Window) -> bool {
    for buf in win.buffers {
        if buf != nil && !buf.busy {
            return true
        }
    }
    return false
}

// Request a redraw - will draw on next available frame
window_request_redraw :: proc(win: ^Window) {
    win.needs_redraw = true

    // If no frame callback is pending, kick off the render loop
    if !win.frame_pending {
        pixels, width, height, stride := window_get_buffer(win)
        if pixels != nil && win.on_draw != nil {
            win.on_draw(win, pixels, width, height, stride)
            window_present(win)
        }
    }
    // Otherwise, the pending frame callback will handle the redraw
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
    win.frame_pending = true
    win.needs_redraw = false
}

// Resize the window buffers
window_resize_buffer :: proc(win: ^Window, new_width, new_height: i32) {
    // Destroy old buffers
    for i in 0..<len(win.buffers) {
        if win.buffers[i] != nil {
            wl.buffer_destroy_internal(win.buffers[i])
            win.buffers[i] = nil
        }
    }

    // Destroy old pool
    if win.pool != nil {
        wl.shm_pool_destroy(win.pool)
        win.pool = nil
    }

    // Update size
    win.width = new_width
    win.height = new_height
    win.current_buffer = 0

    // Create new pool for double buffers
    buffer_size := int(new_width * new_height * 4 * 2)  // Space for 2 buffers
    pool, ok := wl.shm_pool_create(win.app.shm, buffer_size)
    if !ok {
        return
    }
    win.pool = pool

    // Create both buffers
    buf1, ok1 := wl.buffer_create(pool, new_width, new_height, .ARGB8888)
    if !ok1 {
        wl.shm_pool_destroy(pool)
        win.pool = nil
        return
    }
    win.buffers[0] = buf1

    buf2, ok2 := wl.buffer_create(pool, new_width, new_height, .ARGB8888)
    if !ok2 {
        wl.buffer_destroy_internal(buf1)
        win.buffers[0] = nil
        wl.shm_pool_destroy(pool)
        win.pool = nil
        return
    }
    win.buffers[1] = buf2
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
// Input state queries
// ============================================================================

// Get current pointer position
get_pointer_pos :: proc(app: ^App) -> (x, y: f64) {
    return app.pointer_x, app.pointer_y
}

// Check if a mouse button is pressed (1=left, 2=right, 4=middle)
is_button_pressed :: proc(app: ^App, button: u32) -> bool {
    return (app.pointer_buttons & button) != 0
}

// Check if left mouse button is pressed
is_left_button_pressed :: proc(app: ^App) -> bool {
    return (app.pointer_buttons & 1) != 0
}

// Check if right mouse button is pressed
is_right_button_pressed :: proc(app: ^App) -> bool {
    return (app.pointer_buttons & 2) != 0
}

// Check if middle mouse button is pressed
is_middle_button_pressed :: proc(app: ^App) -> bool {
    return (app.pointer_buttons & 4) != 0
}

// Check if Shift is held
is_shift_pressed :: proc(app: ^App) -> bool {
    return wl.xkb_handler_is_shift_active(&app.xkb)
}

// Check if Ctrl is held
is_ctrl_pressed :: proc(app: ^App) -> bool {
    return wl.xkb_handler_is_ctrl_active(&app.xkb)
}

// Check if Alt is held
is_alt_pressed :: proc(app: ^App) -> bool {
    return wl.xkb_handler_is_alt_active(&app.xkb)
}

// Check if Super/Logo is held
is_super_pressed :: proc(app: ^App) -> bool {
    return wl.xkb_handler_is_super_active(&app.xkb)
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
    } else if iface == "wl_seat" {
        app.seat = cast(^wl.Wl_Seat)wl.registry_bind(
            registry, name, &wl.wl_seat_interface, min(version, 8))
        wl.seat_add_listener(app.seat, &app.seat_listener, app)
    }
    // Note: wp_cursor_shape_manager_v1 requires proper protocol bindings
    // TODO: Implement cursor support via libwayland-cursor instead
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

    win.frame_pending = false

    // Mark all buffers as not busy (frame callback means compositor is done with them)
    for &buf in win.buffers {
        if buf != nil {
            buf.busy = false
        }
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

// ============================================================================
// Seat handlers
// ============================================================================

seat_capabilities_handler :: proc "c" (data: rawptr, seat: ^wl.Wl_Seat, capabilities: u32) {
    context = runtime.default_context()
    app := cast(^App)data

    has_pointer := (capabilities & u32(wl.Wl_Seat_Capability.POINTER)) != 0
    has_keyboard := (capabilities & u32(wl.Wl_Seat_Capability.KEYBOARD)) != 0

    // Handle pointer capability
    if has_pointer && app.pointer == nil {
        app.pointer = wl.seat_get_pointer(seat)
        wl.pointer_add_listener(app.pointer, &app.pointer_listener, app)
    } else if !has_pointer && app.pointer != nil {
        if app.cursor_shape_device != nil {
            wl.cursor_shape_device_destroy(app.cursor_shape_device)
            app.cursor_shape_device = nil
        }
        wl.pointer_release(app.pointer)
        app.pointer = nil
    }

    // Handle keyboard capability
    if has_keyboard && app.keyboard == nil {
        app.keyboard = wl.seat_get_keyboard(seat)
        wl.keyboard_add_listener(app.keyboard, &app.keyboard_listener, app)
    } else if !has_keyboard && app.keyboard != nil {
        wl.keyboard_release(app.keyboard)
        app.keyboard = nil
    }
}

seat_name_handler :: proc "c" (data: rawptr, seat: ^wl.Wl_Seat, name: cstring) {
    // Optional: could store seat name
}

// ============================================================================
// Pointer handlers
// ============================================================================

// Find window by surface
find_window_by_surface :: proc(app: ^App, surface: ^wl.Wl_Surface) -> ^Window {
    for win in app.windows {
        if win.surface == surface {
            return win
        }
    }
    return nil
}

pointer_enter_handler :: proc "c" (
    data: rawptr,
    pointer: ^wl.Wl_Pointer,
    serial: u32,
    surface: ^wl.Wl_Surface,
    surface_x: i32,
    surface_y: i32,
) {
    context = runtime.default_context()
    app := cast(^App)data

    app.pointer_surface = surface
    app.pointer_serial = serial
    app.pointer_x = wl.wl_fixed_to_double(surface_x)
    app.pointer_y = wl.wl_fixed_to_double(surface_y)

    // Note: Cursor shape requires either wp_cursor_shape_manager_v1 protocol
    // or loading cursors via libwayland-cursor. For now, we keep the system default.
    // TODO: Implement proper cursor support via libwayland-cursor

    win := find_window_by_surface(app, surface)
    if win != nil && win.on_pointer_enter != nil {
        win.on_pointer_enter(win, app.pointer_x, app.pointer_y)
    }
}

pointer_leave_handler :: proc "c" (
    data: rawptr,
    pointer: ^wl.Wl_Pointer,
    serial: u32,
    surface: ^wl.Wl_Surface,
) {
    context = runtime.default_context()
    app := cast(^App)data

    win := find_window_by_surface(app, surface)
    if win != nil && win.on_pointer_leave != nil {
        win.on_pointer_leave(win)
    }

    app.pointer_surface = nil
}

pointer_motion_handler :: proc "c" (
    data: rawptr,
    pointer: ^wl.Wl_Pointer,
    time: u32,
    surface_x: i32,
    surface_y: i32,
) {
    context = runtime.default_context()
    app := cast(^App)data

    app.pointer_x = wl.wl_fixed_to_double(surface_x)
    app.pointer_y = wl.wl_fixed_to_double(surface_y)

    win := find_window_by_surface(app, app.pointer_surface)
    if win != nil && win.on_pointer_motion != nil {
        win.on_pointer_motion(win, app.pointer_x, app.pointer_y)
    }
}

pointer_button_handler :: proc "c" (
    data: rawptr,
    pointer: ^wl.Wl_Pointer,
    serial: u32,
    time: u32,
    button: u32,
    state: u32,
) {
    context = runtime.default_context()
    app := cast(^App)data

    pressed := state == u32(wl.Wl_Pointer_Button_State.PRESSED)

    // Update button state bitmask
    button_bit: u32
    switch button {
    case wl.BTN_LEFT:
        button_bit = 1
    case wl.BTN_RIGHT:
        button_bit = 2
    case wl.BTN_MIDDLE:
        button_bit = 4
    case:
        button_bit = 1 << ((button - wl.BTN_LEFT) & 31)
    }

    if pressed {
        app.pointer_buttons |= button_bit
    } else {
        app.pointer_buttons &= ~button_bit
    }

    win := find_window_by_surface(app, app.pointer_surface)
    if win != nil && win.on_pointer_button != nil {
        win.on_pointer_button(win, button, pressed)
    }
}

pointer_axis_handler :: proc "c" (
    data: rawptr,
    pointer: ^wl.Wl_Pointer,
    time: u32,
    axis: u32,
    value: i32,
) {
    context = runtime.default_context()
    app := cast(^App)data
    if app == nil {
        return
    }

    // Wayland sends axis values in fixed-point (8.24 format)
    // Divide by 256 to get reasonable pixel scroll amounts
    delta := value / 256

    // Call callback on window under pointer
    win := find_window_by_surface(app, app.pointer_surface)
    if win != nil && win.on_scroll != nil {
        win.on_scroll(win, delta, axis)
    }
}

pointer_frame_handler :: proc "c" (data: rawptr, pointer: ^wl.Wl_Pointer) {
    // Frame marks the end of a set of pointer events
    // Can be used to batch updates
}

pointer_axis_source_handler :: proc "c" (data: rawptr, pointer: ^wl.Wl_Pointer, axis_source: u32) {
    // Axis source (wheel, finger, continuous, etc.)
}

pointer_axis_stop_handler :: proc "c" (data: rawptr, pointer: ^wl.Wl_Pointer, time: u32, axis: u32) {
    // Axis movement stopped
}

pointer_axis_discrete_handler :: proc "c" (data: rawptr, pointer: ^wl.Wl_Pointer, axis: u32, discrete: i32) {
    // Discrete axis step (deprecated in favor of axis_value120)
    // Use 40 pixels per discrete step
    context = runtime.default_context()
    app := cast(^App)data
    if app == nil {
        return
    }

    delta := discrete * 40

    win := find_window_by_surface(app, app.pointer_surface)
    if win != nil && win.on_scroll != nil {
        win.on_scroll(win, delta, axis)
    }
}

pointer_axis_value120_handler :: proc "c" (data: rawptr, pointer: ^wl.Wl_Pointer, axis: u32, value120: i32) {
    // High-resolution axis value (120 units per wheel notch)
    // Convert to pixels: 120 units = 1 notch = ~40 pixels
    context = runtime.default_context()
    app := cast(^App)data
    if app == nil {
        return
    }

    delta := (value120 * 40) / 120

    win := find_window_by_surface(app, app.pointer_surface)
    if win != nil && win.on_scroll != nil {
        win.on_scroll(win, delta, axis)
    }
}

pointer_axis_relative_direction_handler :: proc "c" (data: rawptr, pointer: ^wl.Wl_Pointer, axis: u32, direction: u32) {
    // Relative direction for natural scrolling
}

// ============================================================================
// Keyboard handlers
// ============================================================================

keyboard_keymap_handler :: proc "c" (
    data: rawptr,
    keyboard: ^wl.Wl_Keyboard,
    format: u32,
    fd: i32,
    size: u32,
) {
    context = runtime.default_context()
    app := cast(^App)data

    if format != u32(wl.Wl_Keyboard_Keymap_Format.XKB_V1) {
        posix.close(posix.FD(fd))
        return
    }

    if !wl.xkb_handler_load_keymap_from_fd(&app.xkb, fd, size) {
        fmt.eprintln("Failed to load XKB keymap")
    }

    posix.close(posix.FD(fd))
}

keyboard_enter_handler :: proc "c" (
    data: rawptr,
    keyboard: ^wl.Wl_Keyboard,
    serial: u32,
    surface: ^wl.Wl_Surface,
    keys: ^wl.Wl_Array,
) {
    context = runtime.default_context()
    app := cast(^App)data

    app.keyboard_surface = surface
    app.keyboard_serial = serial
}

keyboard_leave_handler :: proc "c" (
    data: rawptr,
    keyboard: ^wl.Wl_Keyboard,
    serial: u32,
    surface: ^wl.Wl_Surface,
) {
    context = runtime.default_context()
    app := cast(^App)data

    app.keyboard_surface = nil
}

keyboard_key_handler :: proc "c" (
    data: rawptr,
    keyboard: ^wl.Wl_Keyboard,
    serial: u32,
    time: u32,
    key: u32,
    state: u32,
) {
    context = runtime.default_context()
    app := cast(^App)data

    pressed := state == u32(wl.Wl_Keyboard_Key_State.PRESSED)

    // Get UTF-8 representation of the key
    utf8_buf: [8]u8
    utf8_len := wl.xkb_handler_get_utf8(&app.xkb, key, utf8_buf[:])
    utf8_str := string(utf8_buf[:utf8_len])

    win := find_window_by_surface(app, app.keyboard_surface)
    if win != nil && win.on_key != nil {
        win.on_key(win, key, pressed, utf8_str)
    }
}

keyboard_modifiers_handler :: proc "c" (
    data: rawptr,
    keyboard: ^wl.Wl_Keyboard,
    serial: u32,
    mods_depressed: u32,
    mods_latched: u32,
    mods_locked: u32,
    group: u32,
) {
    context = runtime.default_context()
    app := cast(^App)data

    wl.xkb_handler_update_modifiers(&app.xkb, mods_depressed, mods_latched, mods_locked, group)
}

keyboard_repeat_info_handler :: proc "c" (
    data: rawptr,
    keyboard: ^wl.Wl_Keyboard,
    rate: i32,
    delay: i32,
) {
    context = runtime.default_context()
    app := cast(^App)data

    app.key_repeat_rate = rate
    app.key_repeat_delay = delay
}
