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
