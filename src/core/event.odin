package core

// Event types for the widget system
Event_Type :: enum {
    None,

    // Window events
    Window_Close,
    Window_Resize,
    Window_Focus,
    Window_Unfocus,

    // Pointer events
    Pointer_Enter,
    Pointer_Leave,
    Pointer_Motion,
    Pointer_Button_Press,
    Pointer_Button_Release,

    // Keyboard events
    Key_Press,
    Key_Release,
    Text_Input,

    // Scroll events
    Scroll,

    // Frame/redraw events
    Frame,
}

// Mouse buttons
Mouse_Button :: enum u32 {
    Left = 0x110,   // BTN_LEFT
    Right = 0x111,  // BTN_RIGHT
    Middle = 0x112, // BTN_MIDDLE
}

// Keyboard modifier flags
Modifier_Flags :: bit_set[Modifier; u32]

Modifier :: enum {
    Shift,
    Ctrl,
    Alt,
    Super,
    Caps_Lock,
    Num_Lock,
}

// Unified event structure
Event :: struct {
    type: Event_Type,

    // Window events
    window_width: i32,
    window_height: i32,

    // Pointer events
    pointer_x: i32,
    pointer_y: i32,
    button: Mouse_Button,

    // Keyboard events
    keycode: u32,
    keysym: u32,
    modifiers: Modifier_Flags,
    text: [32]u8, // UTF-8 encoded text input
    text_len: int,

    // Scroll events
    scroll_delta: i32,   // Scroll amount (negative = up, positive = down)
    scroll_axis: u32,    // 0 = vertical, 1 = horizontal

    // Timing
    time: u32,
}

// Create specific event types
event_window_resize :: proc(width, height: i32) -> Event {
    return Event{
        type = .Window_Resize,
        window_width = width,
        window_height = height,
    }
}

event_window_close :: proc() -> Event {
    return Event{type = .Window_Close}
}

event_pointer_motion :: proc(x, y: i32, time: u32) -> Event {
    return Event{
        type = .Pointer_Motion,
        pointer_x = x,
        pointer_y = y,
        time = time,
    }
}

event_pointer_button :: proc(button: Mouse_Button, pressed: bool, x, y: i32, time: u32) -> Event {
    return Event{
        type = pressed ? .Pointer_Button_Press : .Pointer_Button_Release,
        button = button,
        pointer_x = x,
        pointer_y = y,
        time = time,
    }
}

event_key :: proc(keycode, keysym: u32, pressed: bool, modifiers: Modifier_Flags) -> Event {
    return Event{
        type = pressed ? .Key_Press : .Key_Release,
        keycode = keycode,
        keysym = keysym,
        modifiers = modifiers,
    }
}

event_text_input :: proc(text: []u8) -> Event {
    e := Event{type = .Text_Input}
    copy_len := min(len(text), len(e.text))
    for i in 0 ..< copy_len {
        e.text[i] = text[i]
    }
    e.text_len = copy_len
    return e
}

event_frame :: proc() -> Event {
    return Event{type = .Frame}
}

event_scroll :: proc(delta: i32, axis: u32 = 0, x: i32 = 0, y: i32 = 0) -> Event {
    return Event{
        type = .Scroll,
        scroll_delta = delta,
        scroll_axis = axis,
        pointer_x = x,
        pointer_y = y,
    }
}
