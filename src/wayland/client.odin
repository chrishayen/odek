package wayland

import "core:c"

// Opaque Wayland types
Wl_Display :: struct {}
Wl_Registry :: struct {}
Wl_Compositor :: struct {}
Wl_Surface :: struct {}
Wl_Callback :: struct {}
Wl_Shm :: struct {}
Wl_Shm_Pool :: struct {}
Wl_Buffer :: struct {}
Wl_Seat :: struct {}
Wl_Pointer :: struct {}
Wl_Keyboard :: struct {}

// XDG shell types
Xdg_Wm_Base :: struct {}
Xdg_Surface :: struct {}
Xdg_Toplevel :: struct {}

// Cursor shape protocol
Wp_Cursor_Shape_Manager :: struct {}

// Wayland proxy (base type for all objects)
Wl_Proxy :: struct {}

// Wayland interface descriptor
Wl_Interface :: struct {
    name: cstring,
    version: c.int,
    method_count: c.int,
    methods: rawptr,
    event_count: c.int,
    events: rawptr,
}

// Wayland message for dispatcher
Wl_Message :: struct {
    name: cstring,
    signature: cstring,
    types: ^rawptr,
}

// Wayland argument union
Wl_Argument :: struct #raw_union {
    i: i32,
    u: u32,
    f: i32, // wl_fixed_t
    s: cstring,
    o: rawptr,
    n: u32,
    a: rawptr,
    h: i32, // fd
}

// Listener callback type
Wl_Dispatcher_Func :: #type proc "c" (
    data: rawptr,
    target: rawptr,
    opcode: u32,
    msg: ^Wl_Message,
    args: [^]Wl_Argument,
) -> c.int

// Foreign bindings to libwayland-client
@(default_calling_convention = "c")
foreign wayland_client {
    // Display
    wl_display_connect :: proc(name: cstring) -> ^Wl_Display ---
    wl_display_disconnect :: proc(display: ^Wl_Display) ---
    wl_display_dispatch :: proc(display: ^Wl_Display) -> c.int ---
    wl_display_dispatch_pending :: proc(display: ^Wl_Display) -> c.int ---
    wl_display_roundtrip :: proc(display: ^Wl_Display) -> c.int ---
    wl_display_flush :: proc(display: ^Wl_Display) -> c.int ---
    wl_display_get_fd :: proc(display: ^Wl_Display) -> c.int ---

    // Proxy (generic object operations)
    wl_proxy_marshal_flags :: proc(
        proxy: rawptr,
        opcode: u32,
        interface: ^Wl_Interface,
        version: u32,
        flags: u32,
        #c_vararg args: ..any,
    ) -> rawptr ---

    wl_proxy_add_listener :: proc(
        proxy: rawptr,
        listener: rawptr,
        data: rawptr,
    ) -> c.int ---

    wl_proxy_destroy :: proc(proxy: rawptr) ---
    wl_proxy_get_version :: proc(proxy: rawptr) -> u32 ---

    wl_proxy_marshal_array_flags :: proc(
        proxy: rawptr,
        opcode: u32,
        interface: ^Wl_Interface,
        version: u32,
        flags: u32,
        args: [^]Wl_Argument,
    ) -> rawptr ---
}

@(default_calling_convention = "c")
foreign wayland_client {
    // Interface descriptors (exported by libwayland-client)
    wl_registry_interface: Wl_Interface
    wl_compositor_interface: Wl_Interface
    wl_surface_interface: Wl_Interface
    wl_callback_interface: Wl_Interface
    wl_shm_interface: Wl_Interface
    wl_shm_pool_interface: Wl_Interface
    wl_buffer_interface: Wl_Interface
    wl_seat_interface: Wl_Interface
    wl_pointer_interface: Wl_Interface
    wl_keyboard_interface: Wl_Interface
}

foreign import wayland_client "system:wayland-client"

// Get registry from display (wl_display_get_registry is inlined in Wayland)
display_get_registry :: proc(display: ^Wl_Display) -> ^Wl_Registry {
    return cast(^Wl_Registry)wl_proxy_marshal_flags(
        display, 1, &wl_registry_interface, 1, 0)
}

// Marshal flag for creating new object
WL_MARSHAL_FLAG_DESTROY :: 1

// SHM formats
Wl_Shm_Format :: enum u32 {
    ARGB8888 = 0,
    XRGB8888 = 1,
}

// Registry listener
Wl_Registry_Listener :: struct {
    global: proc "c" (
        data: rawptr,
        registry: ^Wl_Registry,
        name: u32,
        interface: cstring,
        version: u32,
    ),
    global_remove: proc "c" (
        data: rawptr,
        registry: ^Wl_Registry,
        name: u32,
    ),
}

// Callback listener (for frame callbacks)
Wl_Callback_Listener :: struct {
    done: proc "c" (data: rawptr, callback: ^Wl_Callback, callback_data: u32),
}

// Buffer listener
Wl_Buffer_Listener :: struct {
    release: proc "c" (data: rawptr, buffer: ^Wl_Buffer),
}

// SHM listener
Wl_Shm_Listener :: struct {
    format: proc "c" (data: rawptr, shm: ^Wl_Shm, format: u32),
}

// Registry operations
registry_add_listener :: proc(registry: ^Wl_Registry, listener: ^Wl_Registry_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(registry, listener, data)
}

// WL_REGISTRY_BIND opcode is 0
// Arguments: name (uint), interface_name (string), version (uint), id (new_id)
registry_bind :: proc(registry: ^Wl_Registry, name: u32, interface: ^Wl_Interface, version: u32) -> rawptr {
    return wl_proxy_marshal_flags(
        registry, 0, interface, version, 0,
        name, interface.name, version, nil)
}

// Compositor operations
compositor_create_surface :: proc(compositor: ^Wl_Compositor) -> ^Wl_Surface {
    return cast(^Wl_Surface)wl_proxy_marshal_flags(compositor, 0, &wl_surface_interface, wl_proxy_get_version(compositor), 0)
}

// Surface operations
surface_attach :: proc(surface: ^Wl_Surface, buffer: ^Wl_Buffer, x, y: i32) {
    wl_proxy_marshal_flags(surface, 1, nil, wl_proxy_get_version(surface), 0, buffer, x, y)
}

surface_damage :: proc(surface: ^Wl_Surface, x, y, width, height: i32) {
    wl_proxy_marshal_flags(surface, 2, nil, wl_proxy_get_version(surface), 0, x, y, width, height)
}

surface_commit :: proc(surface: ^Wl_Surface) {
    wl_proxy_marshal_flags(surface, 6, nil, wl_proxy_get_version(surface), 0)
}

surface_frame :: proc(surface: ^Wl_Surface) -> ^Wl_Callback {
    return cast(^Wl_Callback)wl_proxy_marshal_flags(surface, 3, &wl_callback_interface, wl_proxy_get_version(surface), 0)
}

surface_destroy :: proc(surface: ^Wl_Surface) {
    // WL_MARSHAL_FLAG_DESTROY already destroys the proxy
    wl_proxy_marshal_flags(surface, 0, nil, wl_proxy_get_version(surface), WL_MARSHAL_FLAG_DESTROY)
}

// Callback operations
callback_add_listener :: proc(callback: ^Wl_Callback, listener: ^Wl_Callback_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(callback, listener, data)
}

callback_destroy :: proc(callback: ^Wl_Callback) {
    wl_proxy_destroy(callback)
}

// Buffer operations
buffer_add_listener :: proc(buffer: ^Wl_Buffer, listener: ^Wl_Buffer_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(buffer, listener, data)
}

buffer_destroy :: proc(buffer: ^Wl_Buffer) {
    // WL_MARSHAL_FLAG_DESTROY already destroys the proxy, don't call wl_proxy_destroy
    wl_proxy_marshal_flags(buffer, 0, nil, wl_proxy_get_version(buffer), WL_MARSHAL_FLAG_DESTROY)
}
