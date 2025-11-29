package wayland

import "core:c"

// XDG shell interface definitions
// These need to be defined in Odin since xdg-shell is a protocol extension
// not part of libwayland-client itself

// Message type arrays for interfaces - these define the argument signatures
// xdg_wm_base requests: destroy, create_positioner, get_xdg_surface, pong
@(private)
xdg_wm_base_requests := [4]Wl_Message{
    {name = "destroy", signature = "", types = nil},
    {name = "create_positioner", signature = "n", types = nil},
    {name = "get_xdg_surface", signature = "no", types = &xdg_surface_types[0]},
    {name = "pong", signature = "u", types = nil},
}

@(private)
xdg_wm_base_events := [1]Wl_Message{
    {name = "ping", signature = "u", types = nil},
}

@(private)
xdg_surface_types := [2]rawptr{
    &xdg_surface_interface,
    &wl_surface_interface,
}

// xdg_surface requests: destroy, get_toplevel, get_popup, set_window_geometry, ack_configure
@(private)
xdg_surface_requests := [5]Wl_Message{
    {name = "destroy", signature = "", types = nil},
    {name = "get_toplevel", signature = "n", types = &xdg_toplevel_types[0]},
    {name = "get_popup", signature = "n?oo", types = nil},
    {name = "set_window_geometry", signature = "iiii", types = nil},
    {name = "ack_configure", signature = "u", types = nil},
}

@(private)
xdg_surface_events := [1]Wl_Message{
    {name = "configure", signature = "u", types = nil},
}

@(private)
xdg_toplevel_types := [1]rawptr{
    &xdg_toplevel_interface,
}

// xdg_toplevel requests
@(private)
xdg_toplevel_requests := [14]Wl_Message{
    {name = "destroy", signature = "", types = nil},
    {name = "set_parent", signature = "?o", types = nil},
    {name = "set_title", signature = "s", types = nil},
    {name = "set_app_id", signature = "s", types = nil},
    {name = "show_window_menu", signature = "ouii", types = nil},
    {name = "move", signature = "ou", types = nil},
    {name = "resize", signature = "ouu", types = nil},
    {name = "set_max_size", signature = "ii", types = nil},
    {name = "set_min_size", signature = "ii", types = nil},
    {name = "set_maximized", signature = "", types = nil},
    {name = "unset_maximized", signature = "", types = nil},
    {name = "set_fullscreen", signature = "?o", types = nil},
    {name = "unset_fullscreen", signature = "", types = nil},
    {name = "set_minimized", signature = "", types = nil},
}

@(private)
xdg_toplevel_events := [4]Wl_Message{
    {name = "configure", signature = "iia", types = nil},
    {name = "close", signature = "", types = nil},
    {name = "configure_bounds", signature = "ii", types = nil},
    {name = "wm_capabilities", signature = "a", types = nil},
}

// Interface definitions for xdg-shell protocol
xdg_wm_base_interface := Wl_Interface{
    name = "xdg_wm_base",
    version = 6,
    method_count = 4,
    methods = &xdg_wm_base_requests[0],
    event_count = 1,
    events = &xdg_wm_base_events[0],
}

xdg_surface_interface := Wl_Interface{
    name = "xdg_surface",
    version = 6,
    method_count = 5,
    methods = &xdg_surface_requests[0],
    event_count = 1,
    events = &xdg_surface_events[0],
}

xdg_toplevel_interface := Wl_Interface{
    name = "xdg_toplevel",
    version = 6,
    method_count = 14,
    methods = &xdg_toplevel_requests[0],
    event_count = 4,
    events = &xdg_toplevel_events[0],
}

// XDG WM Base listener (ping/pong for responsiveness check)
Xdg_Wm_Base_Listener :: struct {
    ping: proc "c" (data: rawptr, xdg_wm_base: ^Xdg_Wm_Base, serial: u32),
}

// XDG Surface listener
Xdg_Surface_Listener :: struct {
    configure: proc "c" (data: rawptr, xdg_surface: ^Xdg_Surface, serial: u32),
}

// XDG Toplevel listener
Xdg_Toplevel_Listener :: struct {
    configure: proc "c" (
        data: rawptr,
        xdg_toplevel: ^Xdg_Toplevel,
        width: i32,
        height: i32,
        states: ^Wl_Array,
    ),
    close: proc "c" (data: rawptr, xdg_toplevel: ^Xdg_Toplevel),
    configure_bounds: proc "c" (data: rawptr, xdg_toplevel: ^Xdg_Toplevel, width: i32, height: i32),
    wm_capabilities: proc "c" (data: rawptr, xdg_toplevel: ^Xdg_Toplevel, capabilities: ^Wl_Array),
}

// Wayland array type
Wl_Array :: struct {
    size: uint,
    alloc: uint,
    data: rawptr,
}

// XDG toplevel states
Xdg_Toplevel_State :: enum u32 {
    MAXIMIZED = 1,
    FULLSCREEN = 2,
    RESIZING = 3,
    ACTIVATED = 4,
    TILED_LEFT = 5,
    TILED_RIGHT = 6,
    TILED_TOP = 7,
    TILED_BOTTOM = 8,
    SUSPENDED = 9,
}

// XDG WM Base operations
xdg_wm_base_add_listener :: proc(wm_base: ^Xdg_Wm_Base, listener: ^Xdg_Wm_Base_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(wm_base, listener, data)
}

xdg_wm_base_pong :: proc(wm_base: ^Xdg_Wm_Base, serial: u32) {
    wl_proxy_marshal_flags(wm_base, 3, nil, wl_proxy_get_version(wm_base), 0, serial)
}

xdg_wm_base_get_xdg_surface :: proc(wm_base: ^Xdg_Wm_Base, surface: ^Wl_Surface) -> ^Xdg_Surface {
    // get_xdg_surface: opcode 2, creates xdg_surface, takes wl_surface
    // Use array version for proper argument passing
    args: [2]Wl_Argument
    args[0].o = nil  // new_id placeholder
    args[1].o = surface
    return cast(^Xdg_Surface)wl_proxy_marshal_array_flags(
        wm_base, 2, &xdg_surface_interface, wl_proxy_get_version(wm_base), 0, &args[0])
}

xdg_wm_base_destroy :: proc(wm_base: ^Xdg_Wm_Base) {
    // WL_MARSHAL_FLAG_DESTROY already destroys the proxy
    wl_proxy_marshal_flags(wm_base, 0, nil, wl_proxy_get_version(wm_base), WL_MARSHAL_FLAG_DESTROY)
}

// XDG Surface operations
xdg_surface_add_listener :: proc(surface: ^Xdg_Surface, listener: ^Xdg_Surface_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(surface, listener, data)
}

xdg_surface_get_toplevel :: proc(surface: ^Xdg_Surface) -> ^Xdg_Toplevel {
    // get_toplevel: opcode 1, creates xdg_toplevel
    return cast(^Xdg_Toplevel)wl_proxy_marshal_flags(
        surface, 1, &xdg_toplevel_interface, wl_proxy_get_version(surface), 0)
}

xdg_surface_ack_configure :: proc(surface: ^Xdg_Surface, serial: u32) {
    wl_proxy_marshal_flags(surface, 4, nil, wl_proxy_get_version(surface), 0, serial)
}

xdg_surface_destroy :: proc(surface: ^Xdg_Surface) {
    // WL_MARSHAL_FLAG_DESTROY already destroys the proxy
    wl_proxy_marshal_flags(surface, 0, nil, wl_proxy_get_version(surface), WL_MARSHAL_FLAG_DESTROY)
}

// XDG Toplevel operations
xdg_toplevel_add_listener :: proc(toplevel: ^Xdg_Toplevel, listener: ^Xdg_Toplevel_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(toplevel, listener, data)
}

xdg_toplevel_set_title :: proc(toplevel: ^Xdg_Toplevel, title: cstring) {
    wl_proxy_marshal_flags(toplevel, 2, nil, wl_proxy_get_version(toplevel), 0, title)
}

xdg_toplevel_set_app_id :: proc(toplevel: ^Xdg_Toplevel, app_id: cstring) {
    wl_proxy_marshal_flags(toplevel, 3, nil, wl_proxy_get_version(toplevel), 0, app_id)
}

xdg_toplevel_set_min_size :: proc(toplevel: ^Xdg_Toplevel, width, height: i32) {
    wl_proxy_marshal_flags(toplevel, 7, nil, wl_proxy_get_version(toplevel), 0, width, height)
}

xdg_toplevel_set_max_size :: proc(toplevel: ^Xdg_Toplevel, width, height: i32) {
    wl_proxy_marshal_flags(toplevel, 8, nil, wl_proxy_get_version(toplevel), 0, width, height)
}

xdg_toplevel_destroy :: proc(toplevel: ^Xdg_Toplevel) {
    // WL_MARSHAL_FLAG_DESTROY already destroys the proxy
    wl_proxy_marshal_flags(toplevel, 0, nil, wl_proxy_get_version(toplevel), WL_MARSHAL_FLAG_DESTROY)
}

// Helper to iterate over states array
xdg_toplevel_states_iter :: proc(states: ^Wl_Array) -> []Xdg_Toplevel_State {
    if states == nil || states.size == 0 {
        return nil
    }
    count := int(states.size / size_of(u32))
    return (cast([^]Xdg_Toplevel_State)states.data)[:count]
}
