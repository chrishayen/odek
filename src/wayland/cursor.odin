package wayland

import "core:c"

// ============================================================================
// libwayland-cursor bindings
// ============================================================================

// Opaque cursor theme handle
Wl_Cursor_Theme :: struct {}

// Cursor (may have multiple images for animation)
Wl_Cursor :: struct {
    image_count: c.uint,
    images:      [^]^Wl_Cursor_Image,
    name:        cstring,
}

// Single cursor image
Wl_Cursor_Image :: struct {
    width:     u32,
    height:    u32,
    hotspot_x: u32,
    hotspot_y: u32,
    delay:     u32,  // Animation delay in ms
}

foreign import wayland_cursor "system:wayland-cursor"

@(default_calling_convention = "c")
foreign wayland_cursor {
    // Load a cursor theme
    // name: theme name (nil for default), size: cursor size in pixels
    wl_cursor_theme_load :: proc(name: cstring, size: c.int, shm: ^Wl_Shm) -> ^Wl_Cursor_Theme ---

    // Destroy a cursor theme
    wl_cursor_theme_destroy :: proc(theme: ^Wl_Cursor_Theme) ---

    // Get a cursor by name from the theme
    // Returns nil if cursor not found
    wl_cursor_theme_get_cursor :: proc(theme: ^Wl_Cursor_Theme, name: cstring) -> ^Wl_Cursor ---

    // Get the wl_buffer for a cursor image
    wl_cursor_image_get_buffer :: proc(image: ^Wl_Cursor_Image) -> ^Wl_Buffer ---
}

// Standard cursor names (X cursor names)
CURSOR_LEFT_PTR :: "left_ptr"       // Default arrow
CURSOR_HAND1 :: "hand1"             // Pointing hand (clickable)
CURSOR_HAND2 :: "hand2"             // Pointing hand alternative
CURSOR_POINTER :: "pointer"         // Another name for pointing hand
CURSOR_WATCH :: "watch"             // Busy/wait
CURSOR_XTERM :: "xterm"             // Text cursor (I-beam)
CURSOR_TEXT :: "text"               // Text cursor alternative
CURSOR_CROSSHAIR :: "crosshair"     // Crosshair
CURSOR_GRABBING :: "grabbing"       // Grabbing/dragging
CURSOR_FLEUR :: "fleur"             // Move cursor
CURSOR_TOP_SIDE :: "top_side"       // Resize N
CURSOR_BOTTOM_SIDE :: "bottom_side" // Resize S
CURSOR_LEFT_SIDE :: "left_side"     // Resize W
CURSOR_RIGHT_SIDE :: "right_side"   // Resize E
CURSOR_TOP_LEFT_CORNER :: "top_left_corner"     // Resize NW
CURSOR_TOP_RIGHT_CORNER :: "top_right_corner"   // Resize NE
CURSOR_BOTTOM_LEFT_CORNER :: "bottom_left_corner"   // Resize SW
CURSOR_BOTTOM_RIGHT_CORNER :: "bottom_right_corner" // Resize SE

// Cursor type enum for easy switching
Cursor_Type :: enum {
    Arrow,      // Default arrow
    Hand,       // Pointing hand for clickable items
    Text,       // I-beam for text
    Wait,       // Busy/loading
    Crosshair,  // Crosshair
    Move,       // Move/drag
    Grab,       // Grabbing
    ResizeN,
    ResizeS,
    ResizeE,
    ResizeW,
    ResizeNE,
    ResizeNW,
    ResizeSE,
    ResizeSW,
}

// Get cursor name for a cursor type
cursor_type_to_name :: proc(cursor_type: Cursor_Type) -> cstring {
    switch cursor_type {
    case .Arrow:
        return CURSOR_LEFT_PTR
    case .Hand:
        return CURSOR_HAND1
    case .Text:
        return CURSOR_XTERM
    case .Wait:
        return CURSOR_WATCH
    case .Crosshair:
        return CURSOR_CROSSHAIR
    case .Move:
        return CURSOR_FLEUR
    case .Grab:
        return CURSOR_GRABBING
    case .ResizeN:
        return CURSOR_TOP_SIDE
    case .ResizeS:
        return CURSOR_BOTTOM_SIDE
    case .ResizeE:
        return CURSOR_RIGHT_SIDE
    case .ResizeW:
        return CURSOR_LEFT_SIDE
    case .ResizeNE:
        return CURSOR_TOP_RIGHT_CORNER
    case .ResizeNW:
        return CURSOR_TOP_LEFT_CORNER
    case .ResizeSE:
        return CURSOR_BOTTOM_RIGHT_CORNER
    case .ResizeSW:
        return CURSOR_BOTTOM_LEFT_CORNER
    }
    return CURSOR_LEFT_PTR
}
