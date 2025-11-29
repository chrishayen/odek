package tests

import "../src/wayland"
import "core:testing"

@(test)
test_wl_fixed_to_double :: proc(t: ^testing.T) {
    // wl_fixed_t is a signed 24.8 fixed point number
    // 1.0 = 256
    testing.expect(t, wayland.wl_fixed_to_double(256) == 1.0, "256 should be 1.0")
    testing.expect(t, wayland.wl_fixed_to_double(512) == 2.0, "512 should be 2.0")
    testing.expect(t, wayland.wl_fixed_to_double(128) == 0.5, "128 should be 0.5")
    testing.expect(t, wayland.wl_fixed_to_double(-256) == -1.0, "-256 should be -1.0")
    testing.expect(t, wayland.wl_fixed_to_double(0) == 0.0, "0 should be 0.0")
}

@(test)
test_wl_double_to_fixed :: proc(t: ^testing.T) {
    testing.expect(t, wayland.wl_double_to_fixed(1.0) == 256, "1.0 should be 256")
    testing.expect(t, wayland.wl_double_to_fixed(2.0) == 512, "2.0 should be 512")
    testing.expect(t, wayland.wl_double_to_fixed(0.5) == 128, "0.5 should be 128")
    testing.expect(t, wayland.wl_double_to_fixed(-1.0) == -256, "-1.0 should be -256")
    testing.expect(t, wayland.wl_double_to_fixed(0.0) == 0, "0.0 should be 0")
}

@(test)
test_button_codes :: proc(t: ^testing.T) {
    // Verify button codes are correct Linux input event codes
    testing.expect(t, wayland.BTN_LEFT == 0x110, "BTN_LEFT should be 0x110")
    testing.expect(t, wayland.BTN_RIGHT == 0x111, "BTN_RIGHT should be 0x111")
    testing.expect(t, wayland.BTN_MIDDLE == 0x112, "BTN_MIDDLE should be 0x112")
}

@(test)
test_seat_capabilities :: proc(t: ^testing.T) {
    // Verify capability flags
    testing.expect(t, u32(wayland.Wl_Seat_Capability.POINTER) == 1, "POINTER should be 1")
    testing.expect(t, u32(wayland.Wl_Seat_Capability.KEYBOARD) == 2, "KEYBOARD should be 2")
    testing.expect(t, u32(wayland.Wl_Seat_Capability.TOUCH) == 4, "TOUCH should be 4")
}

@(test)
test_pointer_button_state :: proc(t: ^testing.T) {
    testing.expect(t, u32(wayland.Wl_Pointer_Button_State.RELEASED) == 0, "RELEASED should be 0")
    testing.expect(t, u32(wayland.Wl_Pointer_Button_State.PRESSED) == 1, "PRESSED should be 1")
}

@(test)
test_keyboard_key_state :: proc(t: ^testing.T) {
    testing.expect(t, u32(wayland.Wl_Keyboard_Key_State.RELEASED) == 0, "RELEASED should be 0")
    testing.expect(t, u32(wayland.Wl_Keyboard_Key_State.PRESSED) == 1, "PRESSED should be 1")
}

@(test)
test_keyboard_keymap_format :: proc(t: ^testing.T) {
    testing.expect(t, u32(wayland.Wl_Keyboard_Keymap_Format.NO_KEYMAP) == 0, "NO_KEYMAP should be 0")
    testing.expect(t, u32(wayland.Wl_Keyboard_Keymap_Format.XKB_V1) == 1, "XKB_V1 should be 1")
}
