package wayland

import "core:c"
import "core:sys/posix"

// XKB opaque types
Xkb_Context :: struct {}
Xkb_Keymap :: struct {}
Xkb_State :: struct {}

// XKB context flags
Xkb_Context_Flags :: enum c.int {
    NO_FLAGS = 0,
    NO_DEFAULT_INCLUDES = 1,
    NO_ENVIRONMENT_NAMES = 2,
}

// XKB keymap compile flags
Xkb_Keymap_Compile_Flags :: enum c.int {
    NO_FLAGS = 0,
}

// XKB keymap format
Xkb_Keymap_Format :: enum c.int {
    TEXT_V1 = 1,
}

// XKB key direction (for state update)
Xkb_Key_Direction :: enum c.int {
    UP = 0,
    DOWN = 1,
}

// XKB state component (for state update)
Xkb_State_Component :: enum c.int {
    MODS_DEPRESSED = 1,
    MODS_LATCHED = 2,
    MODS_LOCKED = 4,
    MODS_EFFECTIVE = 8,
    LAYOUT_DEPRESSED = 16,
    LAYOUT_LATCHED = 32,
    LAYOUT_LOCKED = 64,
    LAYOUT_EFFECTIVE = 128,
    LEDS = 256,
}

// XKB keysym type
Xkb_Keysym :: u32

// Common keysyms
XKB_KEY_BackSpace :: 0xff08
XKB_KEY_Tab :: 0xff09
XKB_KEY_Return :: 0xff0d
XKB_KEY_Escape :: 0xff1b
XKB_KEY_Delete :: 0xffff
XKB_KEY_Home :: 0xff50
XKB_KEY_Left :: 0xff51
XKB_KEY_Up :: 0xff52
XKB_KEY_Right :: 0xff53
XKB_KEY_Down :: 0xff54
XKB_KEY_End :: 0xff57
XKB_KEY_Page_Up :: 0xff55
XKB_KEY_Page_Down :: 0xff56

// Modifier indices
XKB_MOD_NAME_SHIFT :: "Shift"
XKB_MOD_NAME_CAPS :: "Lock"
XKB_MOD_NAME_CTRL :: "Control"
XKB_MOD_NAME_ALT :: "Mod1"
XKB_MOD_NAME_LOGO :: "Mod4"

// Foreign bindings to libxkbcommon
foreign import xkbcommon "system:xkbcommon"

@(default_calling_convention = "c")
foreign xkbcommon {
    // Context
    xkb_context_new :: proc(flags: Xkb_Context_Flags) -> ^Xkb_Context ---
    xkb_context_ref :: proc(ctx: ^Xkb_Context) -> ^Xkb_Context ---
    xkb_context_unref :: proc(ctx: ^Xkb_Context) ---

    // Keymap
    xkb_keymap_new_from_string :: proc(
        ctx: ^Xkb_Context,
        str: cstring,
        format: Xkb_Keymap_Format,
        flags: Xkb_Keymap_Compile_Flags,
    ) -> ^Xkb_Keymap ---

    xkb_keymap_ref :: proc(keymap: ^Xkb_Keymap) -> ^Xkb_Keymap ---
    xkb_keymap_unref :: proc(keymap: ^Xkb_Keymap) ---

    // State
    xkb_state_new :: proc(keymap: ^Xkb_Keymap) -> ^Xkb_State ---
    xkb_state_ref :: proc(state: ^Xkb_State) -> ^Xkb_State ---
    xkb_state_unref :: proc(state: ^Xkb_State) ---

    // Get keysym from keycode
    xkb_state_key_get_one_sym :: proc(state: ^Xkb_State, key: u32) -> Xkb_Keysym ---

    // Get UTF-8 string for keycode
    xkb_state_key_get_utf8 :: proc(state: ^Xkb_State, key: u32, buffer: [^]u8, size: c.size_t) -> c.int ---

    // Update state with modifiers
    xkb_state_update_mask :: proc(
        state: ^Xkb_State,
        depressed_mods: u32,
        latched_mods: u32,
        locked_mods: u32,
        depressed_layout: u32,
        latched_layout: u32,
        locked_layout: u32,
    ) -> Xkb_State_Component ---

    // Update state with key event
    xkb_state_update_key :: proc(state: ^Xkb_State, key: u32, direction: Xkb_Key_Direction) -> Xkb_State_Component ---

    // Check if modifier is active
    xkb_state_mod_name_is_active :: proc(
        state: ^Xkb_State,
        name: cstring,
        type: Xkb_State_Component,
    ) -> c.int ---

    // Get modifier index
    xkb_keymap_mod_get_index :: proc(keymap: ^Xkb_Keymap, name: cstring) -> u32 ---
}

// Helper struct for managing XKB state
Xkb_Handler :: struct {
    ctx: ^Xkb_Context,
    keymap: ^Xkb_Keymap,
    state: ^Xkb_State,
}

// Initialize XKB handler
xkb_handler_init :: proc() -> (Xkb_Handler, bool) {
    ctx := xkb_context_new(.NO_FLAGS)
    if ctx == nil {
        return {}, false
    }
    return Xkb_Handler{ctx = ctx}, true
}

// Destroy XKB handler
xkb_handler_destroy :: proc(handler: ^Xkb_Handler) {
    if handler.state != nil {
        xkb_state_unref(handler.state)
    }
    if handler.keymap != nil {
        xkb_keymap_unref(handler.keymap)
    }
    if handler.ctx != nil {
        xkb_context_unref(handler.ctx)
    }
    handler^ = {}
}

// Load keymap from file descriptor (as received from Wayland keyboard event)
xkb_handler_load_keymap_from_fd :: proc(handler: ^Xkb_Handler, fd: i32, size: u32) -> bool {
    // mmap the keymap string
    data := posix.mmap(nil, uint(size), {.READ}, {.PRIVATE}, posix.FD(fd), 0)
    if data == posix.MAP_FAILED {
        return false
    }
    defer posix.munmap(data, uint(size))

    // Unref old state and keymap
    if handler.state != nil {
        xkb_state_unref(handler.state)
        handler.state = nil
    }
    if handler.keymap != nil {
        xkb_keymap_unref(handler.keymap)
        handler.keymap = nil
    }

    // Create new keymap from string
    keymap_str := cast(cstring)data
    handler.keymap = xkb_keymap_new_from_string(handler.ctx, keymap_str, .TEXT_V1, .NO_FLAGS)
    if handler.keymap == nil {
        return false
    }

    // Create new state
    handler.state = xkb_state_new(handler.keymap)
    if handler.state == nil {
        xkb_keymap_unref(handler.keymap)
        handler.keymap = nil
        return false
    }

    return true
}

// Update modifier state (from Wayland keyboard.modifiers event)
xkb_handler_update_modifiers :: proc(
    handler: ^Xkb_Handler,
    mods_depressed, mods_latched, mods_locked: u32,
    group: u32,
) {
    if handler.state == nil {
        return
    }
    xkb_state_update_mask(handler.state, mods_depressed, mods_latched, mods_locked, 0, 0, group)
}

// Get UTF-8 string for a key event
// Returns number of bytes written (excluding null terminator), 0 if no output
xkb_handler_get_utf8 :: proc(handler: ^Xkb_Handler, keycode: u32, buffer: []u8) -> int {
    if handler.state == nil || len(buffer) == 0 {
        return 0
    }
    // Wayland keycodes are evdev keycodes + 8
    xkb_keycode := keycode + 8
    size := xkb_state_key_get_utf8(handler.state, xkb_keycode, raw_data(buffer), c.size_t(len(buffer)))
    return int(size)
}

// Get keysym for a key event
xkb_handler_get_keysym :: proc(handler: ^Xkb_Handler, keycode: u32) -> Xkb_Keysym {
    if handler.state == nil {
        return 0
    }
    // Wayland keycodes are evdev keycodes + 8
    xkb_keycode := keycode + 8
    return xkb_state_key_get_one_sym(handler.state, xkb_keycode)
}

// Check modifier state
xkb_handler_is_shift_active :: proc(handler: ^Xkb_Handler) -> bool {
    if handler.state == nil {
        return false
    }
    return xkb_state_mod_name_is_active(handler.state, XKB_MOD_NAME_SHIFT, .MODS_EFFECTIVE) > 0
}

xkb_handler_is_ctrl_active :: proc(handler: ^Xkb_Handler) -> bool {
    if handler.state == nil {
        return false
    }
    return xkb_state_mod_name_is_active(handler.state, XKB_MOD_NAME_CTRL, .MODS_EFFECTIVE) > 0
}

xkb_handler_is_alt_active :: proc(handler: ^Xkb_Handler) -> bool {
    if handler.state == nil {
        return false
    }
    return xkb_state_mod_name_is_active(handler.state, XKB_MOD_NAME_ALT, .MODS_EFFECTIVE) > 0
}

xkb_handler_is_super_active :: proc(handler: ^Xkb_Handler) -> bool {
    if handler.state == nil {
        return false
    }
    return xkb_state_mod_name_is_active(handler.state, XKB_MOD_NAME_LOGO, .MODS_EFFECTIVE) > 0
}
