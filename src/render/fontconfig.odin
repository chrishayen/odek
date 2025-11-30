package render

import "core:c"

// Fontconfig types
FcConfig :: struct {}
FcPattern :: struct {}

// FcResult enum
FcResult :: enum c.int {
    Match = 0,
    NoMatch = 1,
    TypeMismatch = 2,
    NoId = 3,
    OutOfMemory = 4,
}

// FcMatchKind enum
FcMatchKind :: enum c.int {
    Pattern = 0,
    Font = 1,
    Scan = 2,
}

// FcBool type
FcBool :: c.int
FcTrue :: 1
FcFalse :: 0

// Property names
FC_SIZE :: "size"
FC_DPI :: "dpi"
FC_FAMILY :: "family"

// Foreign bindings to libfontconfig
foreign import fontconfig "system:fontconfig"

@(default_calling_convention = "c")
foreign fontconfig {
    // Initialize fontconfig
    FcInit :: proc() -> FcBool ---

    // Pattern functions
    FcPatternCreate :: proc() -> ^FcPattern ---
    FcPatternDestroy :: proc(p: ^FcPattern) ---

    // Config substitution
    FcConfigSubstitute :: proc(
        config: ^FcConfig,
        p: ^FcPattern,
        kind: FcMatchKind,
    ) -> FcBool ---

    // Apply default values
    FcDefaultSubstitute :: proc(pattern: ^FcPattern) ---

    // Font matching
    FcFontMatch :: proc(
        config: ^FcConfig,
        p: ^FcPattern,
        result: ^FcResult,
    ) -> ^FcPattern ---

    // Get double value from pattern
    FcPatternGetDouble :: proc(
        p: ^FcPattern,
        object: cstring,
        n: c.int,
        d: ^f64,
    ) -> FcResult ---

    // Get string value from pattern
    FcPatternGetString :: proc(
        p: ^FcPattern,
        object: cstring,
        n: c.int,
        s: ^cstring,
    ) -> FcResult ---
}

// Get system default font size in points
// Returns the size in points, or 0 if it couldn't be determined
fc_get_default_size :: proc() -> f64 {
    if FcInit() == FcFalse {
        return 0
    }

    pattern := FcPatternCreate()
    if pattern == nil {
        return 0
    }
    defer FcPatternDestroy(pattern)

    // Apply config substitutions and defaults
    FcConfigSubstitute(nil, pattern, .Pattern)
    FcDefaultSubstitute(pattern)

    // Match the font
    result: FcResult
    match := FcFontMatch(nil, pattern, &result)
    if match == nil || result != .Match {
        return 0
    }
    defer FcPatternDestroy(match)

    // Get the size
    size: f64
    if FcPatternGetDouble(match, FC_SIZE, 0, &size) != .Match {
        return 0
    }

    return size
}

// Get system DPI setting
// Returns the DPI, or 0 if it couldn't be determined
fc_get_dpi :: proc() -> f64 {
    if FcInit() == FcFalse {
        return 0
    }

    pattern := FcPatternCreate()
    if pattern == nil {
        return 0
    }
    defer FcPatternDestroy(pattern)

    FcConfigSubstitute(nil, pattern, .Pattern)
    FcDefaultSubstitute(pattern)

    result: FcResult
    match := FcFontMatch(nil, pattern, &result)
    if match == nil || result != .Match {
        return 0
    }
    defer FcPatternDestroy(match)

    dpi: f64
    if FcPatternGetDouble(match, FC_DPI, 0, &dpi) != .Match {
        return 0
    }

    return dpi
}

// Get system default font size in pixels
// Combines font size (points) and DPI to calculate pixel size
// Returns pixel size, or fallback value if system settings unavailable
fc_get_default_pixel_size :: proc(fallback: u32 = 14) -> u32 {
    size_pt := fc_get_default_size()
    if size_pt <= 0 {
        return fallback
    }

    dpi := fc_get_dpi()
    if dpi <= 0 {
        dpi = 96  // Standard DPI fallback
    }

    // Convert points to pixels: pixels = points * dpi / 72
    pixels := size_pt * dpi / 72.0
    if pixels < 1 {
        return fallback
    }

    return u32(pixels + 0.5)  // Round to nearest
}
