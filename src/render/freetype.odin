package render

import "core:c"

// FreeType opaque types
FT_Library :: struct {}
FT_Face :: struct {}

// FreeType error type
FT_Error :: c.int

// FT_Glyph_Format enum
FT_Glyph_Format :: enum c.int {
    NONE = 0,
    COMPOSITE = 0x636F6D70, // 'comp'
    BITMAP = 0x62697473,    // 'bits'
    OUTLINE = 0x6F75746C,   // 'outl'
    PLOTTER = 0x706C6F74,   // 'plot'
    SVG = 0x53564720,       // 'SVG '
}

// FT_Pixel_Mode enum
FT_Pixel_Mode :: enum c.uchar {
    NONE = 0,
    MONO = 1,
    GRAY = 2,
    GRAY2 = 3,
    GRAY4 = 4,
    LCD = 5,
    LCD_V = 6,
    BGRA = 7,
    MAX = 8,
}

// FT_Render_Mode enum
FT_Render_Mode :: enum c.uint {
    NORMAL = 0,
    LIGHT = 1,
    MONO = 2,
    LCD = 3,
    LCD_V = 4,
    SDF = 5,
    MAX = 6,
}

// FT_Kerning_Mode enum
FT_Kerning_Mode :: enum c.uint {
    DEFAULT = 0,
    UNFITTED = 1,
    UNSCALED = 2,
}

// Load flags
FT_LOAD_DEFAULT :: 0
FT_LOAD_NO_SCALE :: 1 << 0
FT_LOAD_NO_HINTING :: 1 << 1
FT_LOAD_RENDER :: 1 << 2
FT_LOAD_NO_BITMAP :: 1 << 3
FT_LOAD_VERTICAL_LAYOUT :: 1 << 4
FT_LOAD_FORCE_AUTOHINT :: 1 << 5
FT_LOAD_CROP_BITMAP :: 1 << 6
FT_LOAD_PEDANTIC :: 1 << 7
FT_LOAD_NO_RECURSE :: 1 << 10
FT_LOAD_IGNORE_TRANSFORM :: 1 << 11
FT_LOAD_MONOCHROME :: 1 << 12
FT_LOAD_LINEAR_DESIGN :: 1 << 13
FT_LOAD_NO_AUTOHINT :: 1 << 15
FT_LOAD_COLOR :: 1 << 20
FT_LOAD_COMPUTE_METRICS :: 1 << 21
FT_LOAD_BITMAP_METRICS_ONLY :: 1 << 22
FT_LOAD_NO_SVG :: 1 << 24

// FT_Vector - a 2D vector
FT_Vector :: struct {
    x: FT_Pos,
    y: FT_Pos,
}

// FT_Pos - a signed long for positions
FT_Pos :: c.long

// FT_Fixed - a signed long for fixed point (16.16)
FT_Fixed :: c.long

// FT_Bitmap - a bitmap
FT_Bitmap :: struct {
    rows:         c.uint,
    width:        c.uint,
    pitch:        c.int,
    buffer:       [^]u8,
    num_grays:    c.ushort,
    pixel_mode:   FT_Pixel_Mode,
    palette_mode: c.uchar,
    palette:      rawptr,
}

// FT_Glyph_Metrics - glyph metrics
FT_Glyph_Metrics :: struct {
    width:        FT_Pos,
    height:       FT_Pos,
    horiBearingX: FT_Pos,
    horiBearingY: FT_Pos,
    horiAdvance:  FT_Pos,
    vertBearingX: FT_Pos,
    vertBearingY: FT_Pos,
    vertAdvance:  FT_Pos,
}

// FT_GlyphSlotRec - glyph slot
FT_GlyphSlotRec :: struct {
    library:            ^FT_Library,
    face:               ^FT_FaceRec,
    next:               ^FT_GlyphSlotRec,
    glyph_index:        c.uint,
    generic:            FT_Generic,
    metrics:            FT_Glyph_Metrics,
    linearHoriAdvance:  FT_Fixed,
    linearVertAdvance:  FT_Fixed,
    advance:            FT_Vector,
    format:             FT_Glyph_Format,
    bitmap:             FT_Bitmap,
    bitmap_left:        c.int,
    bitmap_top:         c.int,
    outline:            FT_Outline,
    num_subglyphs:      c.uint,
    subglyphs:          rawptr,
    control_data:       rawptr,
    control_len:        c.long,
    lsb_delta:          FT_Pos,
    rsb_delta:          FT_Pos,
    other:              rawptr,
    internal:           rawptr,
}

// FT_Generic - client data
FT_Generic :: struct {
    data:      rawptr,
    finalizer: rawptr,
}

// FT_BBox - bounding box
FT_BBox :: struct {
    xMin: FT_Pos,
    yMin: FT_Pos,
    xMax: FT_Pos,
    yMax: FT_Pos,
}

// FT_Outline - outline
FT_Outline :: struct {
    n_contours: c.short,
    n_points:   c.short,
    points:     [^]FT_Vector,
    tags:       [^]c.char,
    contours:   [^]c.short,
    flags:      c.int,
}

// FT_Size_Metrics - size metrics
FT_Size_Metrics :: struct {
    x_ppem:      c.ushort,
    y_ppem:      c.ushort,
    x_scale:     FT_Fixed,
    y_scale:     FT_Fixed,
    ascender:    FT_Pos,
    descender:   FT_Pos,
    height:      FT_Pos,
    max_advance: FT_Pos,
}

// FT_SizeRec - size record
FT_SizeRec :: struct {
    face:     ^FT_FaceRec,
    generic:  FT_Generic,
    metrics:  FT_Size_Metrics,
    internal: rawptr,
}

// FT_CharMapRec - character map
FT_CharMapRec :: struct {
    face:        ^FT_FaceRec,
    encoding:    FT_Encoding,
    platform_id: c.ushort,
    encoding_id: c.ushort,
}

// FT_Encoding - encoding types
FT_Encoding :: enum c.int {
    NONE = 0,
    MS_SYMBOL = 0x73796D62,    // 'symb'
    UNICODE = 0x756E6963,       // 'unic'
    SJIS = 0x736A6973,          // 'sjis'
    PRC = 0x67622020,           // 'gb  '
    BIG5 = 0x62696735,          // 'big5'
    WANSUNG = 0x77616E73,       // 'wans'
    JOHAB = 0x6A6F6861,         // 'joha'
    ADOBE_STANDARD = 0x41444F42, // 'ADOB'
    ADOBE_EXPERT = 0x41444245,   // 'ADBE'
    ADOBE_CUSTOM = 0x41444243,   // 'ADBC'
    ADOBE_LATIN_1 = 0x6C617431,  // 'lat1'
    OLD_LATIN_2 = 0x6C617432,    // 'lat2'
    APPLE_ROMAN = 0x61726D6E,    // 'armn'
}

// FT_FaceRec - face record
FT_FaceRec :: struct {
    num_faces:           c.long,
    face_index:          c.long,
    face_flags:          c.long,
    style_flags:         c.long,
    num_glyphs:          c.long,
    family_name:         cstring,
    style_name:          cstring,
    num_fixed_sizes:     c.int,
    available_sizes:     rawptr,
    num_charmaps:        c.int,
    charmaps:            [^]^FT_CharMapRec,
    generic:             FT_Generic,
    bbox:                FT_BBox,
    units_per_EM:        c.ushort,
    ascender:            c.short,
    descender:           c.short,
    height:              c.short,
    max_advance_width:   c.short,
    max_advance_height:  c.short,
    underline_position:  c.short,
    underline_thickness: c.short,
    glyph:               ^FT_GlyphSlotRec,
    size:                ^FT_SizeRec,
    charmap:             ^FT_CharMapRec,
    driver:              rawptr,
    memory:              rawptr,
    stream:              rawptr,
    sizes_list_head:     rawptr,
    sizes_list_tail:     rawptr,
    autohint:            FT_Generic,
    extensions:          rawptr,
    internal:            rawptr,
}

// Foreign bindings to libfreetype
foreign import freetype "system:freetype"

@(default_calling_convention = "c")
foreign freetype {
    // Library functions
    FT_Init_FreeType :: proc(alibrary: ^^FT_Library) -> FT_Error ---
    FT_Done_FreeType :: proc(library: ^FT_Library) -> FT_Error ---

    // Face functions
    FT_New_Face :: proc(
        library: ^FT_Library,
        filepathname: cstring,
        face_index: c.long,
        aface: ^^FT_FaceRec,
    ) -> FT_Error ---

    FT_New_Memory_Face :: proc(
        library: ^FT_Library,
        file_base: [^]u8,
        file_size: c.long,
        face_index: c.long,
        aface: ^^FT_FaceRec,
    ) -> FT_Error ---

    FT_Done_Face :: proc(face: ^FT_FaceRec) -> FT_Error ---

    // Size functions
    FT_Set_Pixel_Sizes :: proc(
        face: ^FT_FaceRec,
        pixel_width: c.uint,
        pixel_height: c.uint,
    ) -> FT_Error ---

    FT_Set_Char_Size :: proc(
        face: ^FT_FaceRec,
        char_width: FT_Pos,
        char_height: FT_Pos,
        horz_resolution: c.uint,
        vert_resolution: c.uint,
    ) -> FT_Error ---

    // Glyph functions
    FT_Load_Char :: proc(
        face: ^FT_FaceRec,
        char_code: c.ulong,
        load_flags: c.int,
    ) -> FT_Error ---

    FT_Load_Glyph :: proc(
        face: ^FT_FaceRec,
        glyph_index: c.uint,
        load_flags: c.int,
    ) -> FT_Error ---

    FT_Render_Glyph :: proc(
        slot: ^FT_GlyphSlotRec,
        render_mode: FT_Render_Mode,
    ) -> FT_Error ---

    FT_Get_Char_Index :: proc(
        face: ^FT_FaceRec,
        charcode: c.ulong,
    ) -> c.uint ---

    // Kerning
    FT_Get_Kerning :: proc(
        face: ^FT_FaceRec,
        left_glyph: c.uint,
        right_glyph: c.uint,
        kern_mode: FT_Kerning_Mode,
        akerning: ^FT_Vector,
    ) -> FT_Error ---

    // Face flags check
    FT_HAS_KERNING :: proc(face: ^FT_FaceRec) -> c.int ---

    // Select charmap
    FT_Select_Charmap :: proc(
        face: ^FT_FaceRec,
        encoding: FT_Encoding,
    ) -> FT_Error ---
}

// Face flags (for checking face_flags field)
FT_FACE_FLAG_SCALABLE :: 1 << 0
FT_FACE_FLAG_FIXED_SIZES :: 1 << 1
FT_FACE_FLAG_FIXED_WIDTH :: 1 << 2
FT_FACE_FLAG_SFNT :: 1 << 3
FT_FACE_FLAG_HORIZONTAL :: 1 << 4
FT_FACE_FLAG_VERTICAL :: 1 << 5
FT_FACE_FLAG_KERNING :: 1 << 6
FT_FACE_FLAG_FAST_GLYPHS :: 1 << 7
FT_FACE_FLAG_MULTIPLE_MASTERS :: 1 << 8
FT_FACE_FLAG_GLYPH_NAMES :: 1 << 9
FT_FACE_FLAG_EXTERNAL_STREAM :: 1 << 10
FT_FACE_FLAG_HINTER :: 1 << 11
FT_FACE_FLAG_CID_KEYED :: 1 << 12
FT_FACE_FLAG_TRICKY :: 1 << 13
FT_FACE_FLAG_COLOR :: 1 << 14
FT_FACE_FLAG_VARIATION :: 1 << 15
FT_FACE_FLAG_SVG :: 1 << 16
FT_FACE_FLAG_SBIX :: 1 << 17
FT_FACE_FLAG_SBIX_OVERLAY :: 1 << 18

// Helper to check if face has kerning
face_has_kerning :: proc(face: ^FT_FaceRec) -> bool {
    return (face.face_flags & FT_FACE_FLAG_KERNING) != 0
}
