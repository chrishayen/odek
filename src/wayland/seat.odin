package wayland

import "core:c"

// Seat capabilities (bitfield)
Wl_Seat_Capability :: enum u32 {
    POINTER = 1,
    KEYBOARD = 2,
    TOUCH = 4,
}

// Pointer button state
Wl_Pointer_Button_State :: enum u32 {
    RELEASED = 0,
    PRESSED = 1,
}

// Keyboard key state
Wl_Keyboard_Key_State :: enum u32 {
    RELEASED = 0,
    PRESSED = 1,
}

// Keyboard keymap format
Wl_Keyboard_Keymap_Format :: enum u32 {
    NO_KEYMAP = 0,
    XKB_V1 = 1,
}

// Axis source (for scroll)
Wl_Pointer_Axis :: enum u32 {
    VERTICAL_SCROLL = 0,
    HORIZONTAL_SCROLL = 1,
}

// Seat listener
Wl_Seat_Listener :: struct {
    capabilities: proc "c" (data: rawptr, seat: ^Wl_Seat, capabilities: u32),
    name: proc "c" (data: rawptr, seat: ^Wl_Seat, name: cstring),
}

// Pointer listener
Wl_Pointer_Listener :: struct {
    enter: proc "c" (
        data: rawptr,
        pointer: ^Wl_Pointer,
        serial: u32,
        surface: ^Wl_Surface,
        surface_x: i32, // wl_fixed_t
        surface_y: i32, // wl_fixed_t
    ),
    leave: proc "c" (
        data: rawptr,
        pointer: ^Wl_Pointer,
        serial: u32,
        surface: ^Wl_Surface,
    ),
    motion: proc "c" (
        data: rawptr,
        pointer: ^Wl_Pointer,
        time: u32,
        surface_x: i32, // wl_fixed_t
        surface_y: i32, // wl_fixed_t
    ),
    button: proc "c" (
        data: rawptr,
        pointer: ^Wl_Pointer,
        serial: u32,
        time: u32,
        button: u32,
        state: u32,
    ),
    axis: proc "c" (
        data: rawptr,
        pointer: ^Wl_Pointer,
        time: u32,
        axis: u32,
        value: i32, // wl_fixed_t
    ),
    frame: proc "c" (data: rawptr, pointer: ^Wl_Pointer),
    axis_source: proc "c" (data: rawptr, pointer: ^Wl_Pointer, axis_source: u32),
    axis_stop: proc "c" (data: rawptr, pointer: ^Wl_Pointer, time: u32, axis: u32),
    axis_discrete: proc "c" (data: rawptr, pointer: ^Wl_Pointer, axis: u32, discrete: i32),
    axis_value120: proc "c" (data: rawptr, pointer: ^Wl_Pointer, axis: u32, value120: i32),
    axis_relative_direction: proc "c" (data: rawptr, pointer: ^Wl_Pointer, axis: u32, direction: u32),
}

// Keyboard listener
Wl_Keyboard_Listener :: struct {
    keymap: proc "c" (
        data: rawptr,
        keyboard: ^Wl_Keyboard,
        format: u32,
        fd: i32,
        size: u32,
    ),
    enter: proc "c" (
        data: rawptr,
        keyboard: ^Wl_Keyboard,
        serial: u32,
        surface: ^Wl_Surface,
        keys: ^Wl_Array,
    ),
    leave: proc "c" (
        data: rawptr,
        keyboard: ^Wl_Keyboard,
        serial: u32,
        surface: ^Wl_Surface,
    ),
    key: proc "c" (
        data: rawptr,
        keyboard: ^Wl_Keyboard,
        serial: u32,
        time: u32,
        key: u32,
        state: u32,
    ),
    modifiers: proc "c" (
        data: rawptr,
        keyboard: ^Wl_Keyboard,
        serial: u32,
        mods_depressed: u32,
        mods_latched: u32,
        mods_locked: u32,
        group: u32,
    ),
    repeat_info: proc "c" (
        data: rawptr,
        keyboard: ^Wl_Keyboard,
        rate: i32,
        delay: i32,
    ),
}

// Seat operations
seat_add_listener :: proc(seat: ^Wl_Seat, listener: ^Wl_Seat_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(seat, listener, data)
}

// Get pointer from seat (opcode 0)
seat_get_pointer :: proc(seat: ^Wl_Seat) -> ^Wl_Pointer {
    return cast(^Wl_Pointer)wl_proxy_marshal_flags(
        seat, 0, &wl_pointer_interface, wl_proxy_get_version(seat), 0)
}

// Get keyboard from seat (opcode 1)
seat_get_keyboard :: proc(seat: ^Wl_Seat) -> ^Wl_Keyboard {
    return cast(^Wl_Keyboard)wl_proxy_marshal_flags(
        seat, 1, &wl_keyboard_interface, wl_proxy_get_version(seat), 0)
}

// Release seat (opcode 3)
seat_release :: proc(seat: ^Wl_Seat) {
    wl_proxy_marshal_flags(seat, 3, nil, wl_proxy_get_version(seat), WL_MARSHAL_FLAG_DESTROY)
}

// Pointer operations
pointer_add_listener :: proc(pointer: ^Wl_Pointer, listener: ^Wl_Pointer_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(pointer, listener, data)
}

// Set cursor (opcode 0)
pointer_set_cursor :: proc(pointer: ^Wl_Pointer, serial: u32, surface: ^Wl_Surface, hotspot_x, hotspot_y: i32) {
    wl_proxy_marshal_flags(pointer, 0, nil, wl_proxy_get_version(pointer), 0,
        serial, surface, hotspot_x, hotspot_y)
}

// Release pointer (opcode 1)
pointer_release :: proc(pointer: ^Wl_Pointer) {
    wl_proxy_marshal_flags(pointer, 1, nil, wl_proxy_get_version(pointer), WL_MARSHAL_FLAG_DESTROY)
}

// Keyboard operations
keyboard_add_listener :: proc(keyboard: ^Wl_Keyboard, listener: ^Wl_Keyboard_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(keyboard, listener, data)
}

// Release keyboard (opcode 0)
keyboard_release :: proc(keyboard: ^Wl_Keyboard) {
    wl_proxy_marshal_flags(keyboard, 0, nil, wl_proxy_get_version(keyboard), WL_MARSHAL_FLAG_DESTROY)
}

// Helper: convert wl_fixed_t to f64
wl_fixed_to_double :: proc(f: i32) -> f64 {
    return f64(f) / 256.0
}

// Helper: convert f64 to wl_fixed_t
wl_double_to_fixed :: proc(d: f64) -> i32 {
    return i32(d * 256.0)
}

// Mouse button codes (Linux input event codes)
BTN_LEFT :: 0x110
BTN_RIGHT :: 0x111
BTN_MIDDLE :: 0x112

// ============================================================================
// wp_cursor_shape_manager_v1 protocol
// ============================================================================

// Cursor shape device (per-pointer)
Wp_Cursor_Shape_Device :: struct {}

// Cursor shapes (wp_cursor_shape_device_v1.shape enum)
Cursor_Shape :: enum u32 {
    Default = 1,
    Context_Menu = 2,
    Help = 3,
    Pointer = 4,
    Progress = 5,
    Wait = 6,
    Cell = 7,
    Crosshair = 8,
    Text = 9,
    Vertical_Text = 10,
    Alias = 11,
    Copy = 12,
    Move = 13,
    No_Drop = 14,
    Not_Allowed = 15,
    Grab = 16,
    Grabbing = 17,
    E_Resize = 18,
    N_Resize = 19,
    NE_Resize = 20,
    NW_Resize = 21,
    S_Resize = 22,
    SE_Resize = 23,
    SW_Resize = 24,
    W_Resize = 25,
    EW_Resize = 26,
    NS_Resize = 27,
    NESW_Resize = 28,
    NWSE_Resize = 29,
    Col_Resize = 30,
    Row_Resize = 31,
    All_Scroll = 32,
    Zoom_In = 33,
    Zoom_Out = 34,
}

// wp_cursor_shape_manager_v1 requests: destroy, get_pointer
@(private)
cursor_shape_manager_requests := [2]Wl_Message{
    {name = "destroy", signature = "", types = nil},
    {name = "get_pointer", signature = "no", types = nil},
}

// wp_cursor_shape_device_v1 requests: set_shape, destroy
@(private)
cursor_shape_device_requests := [2]Wl_Message{
    {name = "set_shape", signature = "uu", types = nil},
    {name = "destroy", signature = "", types = nil},
}

// wp_cursor_shape_manager_v1 interface
wp_cursor_shape_manager_v1_interface := Wl_Interface{
    name = "wp_cursor_shape_manager_v1",
    version = 1,
    method_count = 2,
    methods = &cursor_shape_manager_requests[0],
    event_count = 0,
    events = nil,
}

// wp_cursor_shape_device_v1 interface
wp_cursor_shape_device_v1_interface := Wl_Interface{
    name = "wp_cursor_shape_device_v1",
    version = 1,
    method_count = 2,
    methods = &cursor_shape_device_requests[0],
    event_count = 0,
    events = nil,
}

// Get cursor shape device for a pointer (opcode 1)
cursor_shape_manager_get_pointer :: proc(manager: ^Wp_Cursor_Shape_Manager, pointer: ^Wl_Pointer) -> ^Wp_Cursor_Shape_Device {
    // Use array version for proper argument passing with new_id
    args: [2]Wl_Argument
    args[0].o = nil  // new_id placeholder
    args[1].o = pointer
    return cast(^Wp_Cursor_Shape_Device)wl_proxy_marshal_array_flags(
        manager, 1, &wp_cursor_shape_device_v1_interface, wl_proxy_get_version(manager), 0, &args[0])
}

// Destroy cursor shape manager (opcode 0)
cursor_shape_manager_destroy :: proc(manager: ^Wp_Cursor_Shape_Manager) {
    wl_proxy_marshal_flags(manager, 0, nil, wl_proxy_get_version(manager), WL_MARSHAL_FLAG_DESTROY)
}

// Set cursor shape (opcode 0)
cursor_shape_device_set_shape :: proc(device: ^Wp_Cursor_Shape_Device, serial: u32, shape: Cursor_Shape) {
    wl_proxy_marshal_flags(device, 0, nil, wl_proxy_get_version(device), 0, serial, u32(shape))
}

// Destroy cursor shape device (opcode 1)
cursor_shape_device_destroy :: proc(device: ^Wp_Cursor_Shape_Device) {
    wl_proxy_marshal_flags(device, 1, nil, wl_proxy_get_version(device), WL_MARSHAL_FLAG_DESTROY)
}
