package core

// Linux input event keycodes (from linux/input-event-codes.h)
// These match evdev scancodes used by Wayland
Keycode :: enum u32 {
    Escape       = 1,
    Num_1        = 2,
    Num_2        = 3,
    Num_3        = 4,
    Num_4        = 5,
    Num_5        = 6,
    Num_6        = 7,
    Num_7        = 8,
    Num_8        = 9,
    Num_9        = 10,
    Num_0        = 11,
    Minus        = 12,
    Equal        = 13,
    Backspace    = 14,
    Tab          = 15,
    Q            = 16,
    W            = 17,
    E            = 18,
    R            = 19,
    T            = 20,
    Y            = 21,
    U            = 22,
    I            = 23,
    O            = 24,
    P            = 25,
    Left_Bracket = 26,
    Right_Bracket = 27,
    Enter        = 28,
    Left_Ctrl    = 29,
    A            = 30,
    S            = 31,
    D            = 32,
    F            = 33,
    G            = 34,
    H            = 35,
    J            = 36,
    K            = 37,
    L            = 38,
    Semicolon    = 39,
    Apostrophe   = 40,
    Grave        = 41,
    Left_Shift   = 42,
    Backslash    = 43,
    Z            = 44,
    X            = 45,
    C            = 46,
    V            = 47,
    B            = 48,
    N            = 49,
    M            = 50,
    Comma        = 51,
    Period       = 52,
    Slash        = 53,
    Right_Shift  = 54,
    KP_Asterisk  = 55,
    Left_Alt     = 56,
    Space        = 57,
    Caps_Lock    = 58,
    F1           = 59,
    F2           = 60,
    F3           = 61,
    F4           = 62,
    F5           = 63,
    F6           = 64,
    F7           = 65,
    F8           = 66,
    F9           = 67,
    F10          = 68,
    Num_Lock     = 69,
    Scroll_Lock  = 70,
    KP_7         = 71,
    KP_8         = 72,
    KP_9         = 73,
    KP_Minus     = 74,
    KP_4         = 75,
    KP_5         = 76,
    KP_6         = 77,
    KP_Plus      = 78,
    KP_1         = 79,
    KP_2         = 80,
    KP_3         = 81,
    KP_0         = 82,
    KP_Period    = 83,
    F11          = 87,
    F12          = 88,
    KP_Enter     = 96,
    Right_Ctrl   = 97,
    KP_Slash     = 98,
    Right_Alt    = 100,
    Home         = 102,
    Up           = 103,
    Page_Up      = 104,
    Left         = 105,
    Right        = 106,
    End          = 107,
    Down         = 108,
    Page_Down    = 109,
    Insert       = 110,
    Delete       = 111,
    Left_Meta    = 125,
    Right_Meta   = 126,
}

// XKB keysyms for common keys
// These are used by the XKB keyboard handling
Keysym :: enum u32 {
    // TTY function keys
    BackSpace   = 0xFF08,
    Tab         = 0xFF09,
    Return      = 0xFF0D,
    Escape      = 0xFF1B,
    Delete      = 0xFFFF,

    // Cursor control
    Home        = 0xFF50,
    Left        = 0xFF51,
    Up          = 0xFF52,
    Right       = 0xFF53,
    Down        = 0xFF54,
    Page_Up     = 0xFF55,
    Page_Down   = 0xFF56,
    End         = 0xFF57,

    // Misc
    Insert      = 0xFF63,

    // Keypad
    KP_Enter    = 0xFF8D,
    KP_Home     = 0xFF95,
    KP_Left     = 0xFF96,
    KP_Up       = 0xFF97,
    KP_Right    = 0xFF98,
    KP_Down     = 0xFF99,
    KP_Page_Up  = 0xFF9A,
    KP_Page_Down = 0xFF9B,
    KP_End      = 0xFF9C,
    KP_Delete   = 0xFF9F,

    // Function keys
    F1          = 0xFFBE,
    F2          = 0xFFBF,
    F3          = 0xFFC0,
    F4          = 0xFFC1,
    F5          = 0xFFC2,
    F6          = 0xFFC3,
    F7          = 0xFFC4,
    F8          = 0xFFC5,
    F9          = 0xFFC6,
    F10         = 0xFFC7,
    F11         = 0xFFC8,
    F12         = 0xFFC9,

    // Modifiers
    Shift_L     = 0xFFE1,
    Shift_R     = 0xFFE2,
    Control_L   = 0xFFE3,
    Control_R   = 0xFFE4,
    Caps_Lock   = 0xFFE5,
    Alt_L       = 0xFFE9,
    Alt_R       = 0xFFEA,
    Super_L     = 0xFFEB,
    Super_R     = 0xFFEC,

    // Latin 1 (ASCII subset)
    Space       = 0x0020,
}

// Helper to check if a keycode matches
keycode_is :: proc(event_keycode: u32, key: Keycode) -> bool {
    return event_keycode == u32(key)
}

// Helper to check if a keysym matches
keysym_is :: proc(event_keysym: u32, key: Keysym) -> bool {
    return event_keysym == u32(key)
}

// Check if either keycode or keysym matches (useful for key handling)
key_matches :: proc(event_keycode, event_keysym: u32, keycode: Keycode, keysym: Keysym) -> bool {
    return event_keycode == u32(keycode) || event_keysym == u32(keysym)
}
