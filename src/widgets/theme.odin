package widgets

import "../core"

// Theme defines the color palette for widgets
Theme :: struct {
    // Backgrounds
    bg_primary:     core.Color,  // Main background (windows, containers)
    bg_secondary:   core.Color,  // Secondary background (cards, panels)
    bg_tertiary:    core.Color,  // Tertiary background (headers, sidebars)

    // Interactive states
    bg_hover:       core.Color,  // Hover state
    bg_pressed:     core.Color,  // Pressed/active state
    bg_selected:    core.Color,  // Selected item

    // Accent colors
    accent:         core.Color,  // Primary accent (buttons, links)
    accent_hover:   core.Color,  // Accent hover
    accent_pressed: core.Color,  // Accent pressed

    // Text colors
    text_primary:   core.Color,  // Main text
    text_secondary: core.Color,  // Muted text
    text_disabled:  core.Color,  // Disabled text
    text_on_accent: core.Color,  // Text on accent background

    // Borders and dividers
    border:         core.Color,  // Default border
    border_focus:   core.Color,  // Focus ring
    divider:        core.Color,  // Divider lines

    // Input fields
    input_bg:       core.Color,  // Input background
    input_border:   core.Color,  // Input border
    input_focus:    core.Color,  // Input focus border

    // Scrollbar
    scrollbar_track: core.Color,
    scrollbar_thumb: core.Color,
    scrollbar_hover: core.Color,

    // Semantic colors
    error:          core.Color,
    warning:        core.Color,
    success:        core.Color,
    info:           core.Color,
}

// Helper to create Color from hex at compile time
@(private)
hex_color :: #force_inline proc "contextless" (hex: u32) -> core.Color {
    return core.Color{
        r = u8((hex >> 16) & 0xFF),
        g = u8((hex >> 8) & 0xFF),
        b = u8(hex & 0xFF),
        a = 255,
    }
}

// Dark theme - default
DARK_THEME := Theme{
    // Backgrounds
    bg_primary     = hex_color(0x1E1E1E),
    bg_secondary   = hex_color(0x2D2D2D),
    bg_tertiary    = hex_color(0x333333),

    // Interactive states
    bg_hover       = hex_color(0x3D3D3D),
    bg_pressed     = hex_color(0x4D4D4D),
    bg_selected    = hex_color(0x094771),

    // Accent colors (blue)
    accent         = hex_color(0x4A90D9),
    accent_hover   = hex_color(0x5BA0E9),
    accent_pressed = hex_color(0x3A80C9),

    // Text colors
    text_primary   = hex_color(0xE0E0E0),
    text_secondary = hex_color(0x909090),
    text_disabled  = hex_color(0x606060),
    text_on_accent = hex_color(0xFFFFFF),

    // Borders and dividers
    border         = hex_color(0x444444),
    border_focus   = hex_color(0x4A90D9),
    divider        = hex_color(0x3D3D3D),

    // Input fields
    input_bg       = hex_color(0x252526),
    input_border   = hex_color(0x3C3C3C),
    input_focus    = hex_color(0x4A90D9),

    // Scrollbar
    scrollbar_track = hex_color(0x2D2D2D),
    scrollbar_thumb = hex_color(0x555555),
    scrollbar_hover = hex_color(0x666666),

    // Semantic colors
    error          = hex_color(0xF44336),
    warning        = hex_color(0xFF9800),
    success        = hex_color(0x4CAF50),
    info           = hex_color(0x2196F3),
}

// Light theme
LIGHT_THEME := Theme{
    // Backgrounds
    bg_primary     = hex_color(0xFFFFFF),
    bg_secondary   = hex_color(0xF5F5F5),
    bg_tertiary    = hex_color(0xE8E8E8),

    // Interactive states
    bg_hover       = hex_color(0xE0E0E0),
    bg_pressed     = hex_color(0xD0D0D0),
    bg_selected    = hex_color(0xCCE8FF),

    // Accent colors (blue)
    accent         = hex_color(0x0078D4),
    accent_hover   = hex_color(0x1084D8),
    accent_pressed = hex_color(0x006CBD),

    // Text colors
    text_primary   = hex_color(0x1E1E1E),
    text_secondary = hex_color(0x606060),
    text_disabled  = hex_color(0xA0A0A0),
    text_on_accent = hex_color(0xFFFFFF),

    // Borders and dividers
    border         = hex_color(0xD0D0D0),
    border_focus   = hex_color(0x0078D4),
    divider        = hex_color(0xE0E0E0),

    // Input fields
    input_bg       = hex_color(0xFFFFFF),
    input_border   = hex_color(0xCCCCCC),
    input_focus    = hex_color(0x0078D4),

    // Scrollbar
    scrollbar_track = hex_color(0xF0F0F0),
    scrollbar_thumb = hex_color(0xC0C0C0),
    scrollbar_hover = hex_color(0xA0A0A0),

    // Semantic colors
    error          = hex_color(0xD32F2F),
    warning        = hex_color(0xF57C00),
    success        = hex_color(0x388E3C),
    info           = hex_color(0x1976D2),
}

// Global theme - can be changed at runtime
current_theme: ^Theme = &DARK_THEME

// Set the current theme
theme_set :: proc(theme: ^Theme) {
    current_theme = theme
}

// Get the current theme
theme_get :: proc() -> ^Theme {
    return current_theme
}

// Use dark theme
theme_set_dark :: proc() {
    current_theme = &DARK_THEME
}

// Use light theme
theme_set_light :: proc() {
    current_theme = &LIGHT_THEME
}
