# Odek

A pure Odin UI toolkit for building native Wayland desktop applications.

## Features

- Native Wayland support with XDG Shell window management
- HiDPI and fractional scaling support
- Flexbox-inspired layout system
- Software rendering with double buffering
- Async image loading with worker threads
- FreeType text rendering with glyph caching
- Fontconfig for automatic system font discovery
- Focus management and keyboard navigation
- Theme system with dark and light themes
- Comprehensive event system (pointer, keyboard, scroll)
- Clipboard support (copy/paste)
- Video thumbnail decoding via FFmpeg

## Requirements

- Linux with Wayland compositor
- Odin compiler
- FreeType library (`libfreetype`)
- Fontconfig library (`libfontconfig`)
- Wayland client libraries (`libwayland-client`, `libwayland-cursor`)
- FFmpeg libraries (`libavformat`, `libavcodec`, `libswscale`, `libavutil`) - optional, for video support

## Building

```bash
odin build src -out:odek
```

## Quick Start

The high-level `app` package handles all boilerplate automatically:

```odin
package main

import "app"
import "widgets"

// Global state for callbacks (Odin has no closures)
g_label: ^widgets.Label
g_input: ^widgets.Text_Input

main :: proc() {
    a := app.create("My App", 400, 300)
    if a == nil {
        return
    }
    defer app.destroy(a)

    // Create widgets (auto-uses app's font)
    g_label = app.label(a, "Enter your name:")
    g_input = app.text_input(a)
    widgets.text_input_set_placeholder(g_input, "Type here...")

    btn := app.button(a, "Submit")
    btn.on_click = proc(b: ^widgets.Button) {
        text := widgets.text_input_get_text(g_input)
        widgets.label_set_text(g_label, text)
    }

    // Layout
    app.column(a, {g_label, g_input, btn})

    // Run (all events handled automatically)
    app.run(a)
}
```

The `app` package automatically:
- Initializes Wayland and creates the window
- Sets up text rendering and loads a system font
- Wires up all event callbacks (draw, pointer, keyboard, scroll)
- Manages focus and hit testing

Interact with widgets directly using `widgets.*` functions or field assignment.

## Advanced Usage

For more control, use the low-level API directly:

```odin
package main

import "core"
import "widgets"
import "render"
import "core:fmt"

// Global state - required because callbacks can't capture local variables
App_State :: struct {
    window: ^core.Window,
    root: ^widgets.Container,
    text_renderer: render.Text_Renderer,
    font: render.Font,
    focus_manager: widgets.Focus_Manager,
    hit_state: widgets.Hit_Test_State,
    label: ^widgets.Label,
    input: ^widgets.Text_Input,
}

g_state: App_State

main :: proc() {
    // Initialize application
    app := core.init()
    if app == nil {
        fmt.eprintln("Failed to initialize Wayland")
        return
    }
    defer core.shutdown(app)

    // Create window
    g_state.window = core.create_window(app, "My App", 400, 300)
    if g_state.window == nil {
        fmt.eprintln("Failed to create window")
        return
    }

    // Initialize text rendering
    ok: bool
    g_state.text_renderer, ok = render.text_renderer_init()
    if !ok {
        fmt.eprintln("Failed to initialize text renderer")
        return
    }
    defer render.text_renderer_destroy(&g_state.text_renderer)

    g_state.font, ok = render.font_load(&g_state.text_renderer, "/usr/share/fonts/TTF/DejaVuSans.ttf", 14)
    if !ok {
        fmt.eprintln("Failed to load font")
        return
    }
    defer render.font_destroy(&g_state.font)

    // Build UI
    build_ui()

    // Set up callbacks
    g_state.window.on_draw = draw_callback
    g_state.window.on_pointer_motion = pointer_motion_callback
    g_state.window.on_pointer_button = pointer_button_callback
    g_state.window.on_key = key_callback
    g_state.window.on_close = close_callback

    // Run event loop
    core.run(app)
}

build_ui :: proc() {
    // Root container with theme background
    g_state.root = widgets.container_create(.Column)
    g_state.root.padding = widgets.edges_all(20)
    g_state.root.spacing = 10
    g_state.root.background = widgets.theme_get().bg_primary

    // Label
    g_state.label = widgets.label_create("Enter your name:", &g_state.font)
    widgets.widget_add_child(g_state.root, g_state.label)

    // Text input
    g_state.input = widgets.text_input_create(&g_state.font)
    g_state.input.min_size = core.Size{0, 32}
    widgets.text_input_set_placeholder(g_state.input, "Type here...")
    widgets.widget_add_child(g_state.root, g_state.input)

    // Button
    button := widgets.button_create("Submit", &g_state.font)
    button.min_size = core.Size{100, 36}
    button.on_click = proc(b: ^widgets.Button) {
        text := widgets.text_input_get_text(g_state.input)
        widgets.label_set_text(g_state.label, fmt.tprintf("Hello, %s!", text))
    }
    widgets.widget_add_child(g_state.root, button)

    // Initialize focus manager
    g_state.focus_manager = widgets.focus_manager_init(g_state.root)
}

draw_callback :: proc(win: ^core.Window, pixels: [^]u32, w, h, stride: i32) {
    ctx := render.context_create_scaled(
        pixels, w, h, stride,
        win.logical_width, win.logical_height, win.scale)
    render.clear(&ctx, widgets.theme_get().bg_primary)

    g_state.root.rect = core.Rect{0, 0, win.logical_width, win.logical_height}
    widgets.widget_layout(g_state.root)
    widgets.widget_draw(g_state.root, &ctx)
}

pointer_motion_callback :: proc(win: ^core.Window, x, y: f64) {
    widgets.update_hover(&g_state.hit_state, g_state.root, i32(x), i32(y))
    core.window_request_redraw(win)
}

pointer_button_callback :: proc(win: ^core.Window, button: u32, pressed: bool) {
    x, y := core.get_pointer_pos(win.app)
    event := core.event_pointer_button(core.Mouse_Button(button), pressed, i32(x), i32(y), 0)
    widgets.dispatch_pointer_event(&g_state.hit_state, g_state.root, &event)
    core.window_request_redraw(win)
}

key_callback :: proc(win: ^core.Window, keycode: u32, pressed: bool, utf8: string) {
    if !pressed {
        return
    }

    event := core.Event{
        type = .Key_Press,
        keycode = keycode,
    }

    // Handle Tab for focus navigation
    if widgets.focus_handle_tab(&g_state.focus_manager, &event) {
        core.window_request_redraw(win)
        return
    }

    // Send to focused widget
    focused := widgets.focus_get(&g_state.focus_manager)
    if focused != nil {
        widgets.widget_handle_event(focused, &event)
        core.window_request_redraw(win)
    }
}

close_callback :: proc(win: ^core.Window) {
    win.app.running = false
}
```

## Architecture

```
src/
├── app/            High-level API (recommended)
│   └── app.odin        GTK-style convenience API
├── core/           Application & event management
│   ├── app.odin        Window creation, event loop, Wayland setup
│   ├── event.odin      Event types and constructors
│   ├── keys.odin       Keyboard constants (keycodes and keysyms)
│   └── types.odin      Basic types (Rect, Point, Size, Color)
├── widgets/        UI components
│   ├── widget.odin     Base widget with VTable polymorphism
│   ├── container.odin  Flexbox-lite layout container
│   ├── button.odin     Interactive button widget
│   ├── checkbox.odin   Toggle checkbox widget
│   ├── dropdown.odin   Dropdown select widget
│   ├── toggle_group.odin  Toggle button group
│   ├── label.odin      Text display with wrapping
│   ├── text_input.odin Single-line text input
│   ├── scroll_container.odin  Scrollable container
│   ├── scroll.odin     Reusable scroll state
│   ├── image_grid.odin Scrollable image grid
│   ├── theme.odin      Color theme system
│   ├── focus.odin      Focus management
│   └── hit_test.odin   Hit testing and event dispatch
├── render/         Graphics rendering
│   ├── buffer.odin     Draw context and primitives
│   ├── text.odin       FreeType text rendering
│   ├── fontconfig.odin System font discovery
│   ├── image.odin      Image loading (PNG/JPEG)
│   ├── image_loader.odin  Async loading with threads
│   ├── image_cache.odin   LRU image cache
│   └── video.odin      FFmpeg video decoding
├── ffmpeg/         FFmpeg bindings (optional)
└── wayland/        Wayland protocol bindings
```

## Widgets

### Container

Flexbox-inspired layout container supporting row and column directions.

```odin
container := widgets.container_create(.Column)
container.padding = widgets.edges_all(10)
container.spacing = 5
container.background = widgets.theme_get().bg_secondary
container.align_items = .Stretch

widgets.widget_add_child(container, child1)
widgets.widget_add_child(container, child2)
```

**Properties:**
| Property | Type | Description |
|----------|------|-------------|
| `direction` | `Direction` | `.Row` or `.Column` |
| `align_items` | `Align` | `.Start`, `.Center`, `.End`, `.Stretch` |
| `spacing` | `i32` | Gap between children (pixels) |
| `background` | `Color` | Fill color |
| `padding` | `Edges` | Internal spacing |

### Button

Interactive button with hover and press states. Uses theme colors by default.

```odin
button := widgets.button_create("Submit", &font)
button.min_size = core.Size{120, 40}
button.corner_radius = 4

// Click handler (must use global state, not local captures)
button.on_click = proc(b: ^widgets.Button) {
    // Handle click using g_state
}

widgets.widget_add_child(parent, button)
```

**Properties:**
| Property | Type | Description |
|----------|------|-------------|
| `text` | `string` | Button label |
| `font` | `^Font` | Text font |
| `bg_normal` | `Color` | Default background (from theme) |
| `bg_hover` | `Color` | Hover background (from theme) |
| `bg_pressed` | `Color` | Pressed background (from theme) |
| `text_color` | `Color` | Label color (from theme) |
| `corner_radius` | `i32` | Border radius |
| `on_click` | `proc` | Click callback |

### Label

Text display with optional word wrapping.

```odin
label := widgets.label_create("Hello World", &font)
widgets.label_set_align(label, .Center)
label.wrap = true

widgets.widget_add_child(parent, label)
```

**Properties:**
| Property | Type | Description |
|----------|------|-------------|
| `text` | `string` | Display text |
| `font` | `^Font` | Text font |
| `color` | `Color` | Text color (from theme) |
| `h_align` | `Align` | `.Start`, `.Center`, `.End` |
| `wrap` | `bool` | Enable word wrapping |

### Text Input

Single-line text input with cursor and editing.

```odin
input := widgets.text_input_create(&font)
input.min_size = core.Size{200, 32}
widgets.text_input_set_placeholder(input, "Enter text...")

// Callbacks
input.on_change = proc(ti: ^widgets.Text_Input) {
    // Text changed
}
input.on_submit = proc(ti: ^widgets.Text_Input) {
    // Enter pressed
}

widgets.widget_add_child(parent, input)

// Get/set text
text := widgets.text_input_get_text(input)
widgets.text_input_set_text(input, "New text")
```

**Properties:**
| Property | Type | Description |
|----------|------|-------------|
| `font` | `^Font` | Text font |
| `placeholder` | `string` | Placeholder text |
| `bg_color` | `Color` | Background (from theme) |
| `text_color` | `Color` | Text color (from theme) |
| `border_color` | `Color` | Border color (from theme) |
| `focus_color` | `Color` | Focus border (from theme) |
| `on_change` | `proc` | Text change callback |
| `on_submit` | `proc` | Enter key callback |

### Checkbox

Toggle checkbox widget with checked/unchecked state.

```odin
checkbox := widgets.checkbox_create()
checkbox.checked = true

checkbox.on_change = proc(cb: ^widgets.Checkbox) {
    if cb.checked {
        // Handle checked state
    }
}

widgets.widget_add_child(parent, checkbox)
```

**Properties:**
| Property | Type | Description |
|----------|------|-------------|
| `checked` | `bool` | Current toggle state |
| `box_size` | `i32` | Size of the checkbox box |
| `corner_radius` | `i32` | Border radius |
| `box_color` | `Color` | Box background (from theme) |
| `check_color` | `Color` | Checkmark color (from theme) |
| `border_color` | `Color` | Border color (from theme) |
| `on_change` | `proc` | State change callback |

### Dropdown

Dropdown select widget for choosing from a list of options.

```odin
dropdown := widgets.dropdown_create(&font)
widgets.dropdown_add_option(dropdown, "Option 1")
widgets.dropdown_add_option(dropdown, "Option 2")
widgets.dropdown_add_option(dropdown, "Option 3")
dropdown.selected_index = 0

dropdown.on_change = proc(d: ^widgets.Dropdown, index: i32) {
    // Handle selection change
}

widgets.widget_add_child(parent, dropdown)
```

**Properties:**
| Property | Type | Description |
|----------|------|-------------|
| `options` | `[dynamic]string` | List of options |
| `selected_index` | `i32` | Currently selected option index |
| `is_open` | `bool` | Whether dropdown is expanded |
| `font` | `^Font` | Text font |
| `corner_radius` | `i32` | Border radius |
| `item_height` | `i32` | Height of each option item |
| `on_change` | `proc` | Selection change callback |

### Toggle Group

Horizontal button group where only one option can be selected at a time.

```odin
options := []string{"Day", "Week", "Month"}
toggle := widgets.toggle_group_create(options, &font)

toggle.on_change = proc(g: ^widgets.Toggle_Group) {
    index := widgets.toggle_group_get_selected(g)
    text := widgets.toggle_group_get_selected_text(g)
}

widgets.widget_add_child(parent, toggle)
```

**Properties:**
| Property | Type | Description |
|----------|------|-------------|
| `options` | `[dynamic]string` | List of options (cloned on create) |
| `selected_index` | `int` | Currently selected option index |
| `font` | `^Font` | Text font |
| `corner_radius` | `i32` | Border radius |
| `bg_normal` | `Color` | Unselected background (from theme) |
| `bg_hover` | `Color` | Hover background (from theme) |
| `bg_selected` | `Color` | Selected background (from theme) |
| `text_color` | `Color` | Text color (from theme) |
| `on_change` | `proc` | Selection change callback |

**Functions:**
| Function | Description |
|----------|-------------|
| `toggle_group_get_selected(g)` | Get selected index |
| `toggle_group_get_selected_text(g)` | Get selected option text |
| `toggle_group_set_selected(g, index)` | Set selected index |

**Keyboard:** Left/Right arrows change selection when focused.

### Scroll Container

Wraps content in a scrollable container.

```odin
scroll := widgets.scroll_container_create(.Vertical)

content := widgets.container_create(.Column)
// Add many children to content...

widgets.scroll_container_set_content(scroll, content)
widgets.widget_add_child(parent, scroll)
```

**Properties:**
| Property | Type | Description |
|----------|------|-------------|
| `direction` | `Scroll_Direction` | `.Vertical`, `.Horizontal`, `.Both` |
| `scrollbar_width` | `i32` | Scrollbar width |
| `track_color` | `Color` | Scrollbar track (from theme) |
| `thumb_color` | `Color` | Scrollbar thumb (from theme) |

### Image Grid

Scrollable grid for displaying images with selection support.

```odin
grid := widgets.image_grid_create()
grid.cell_width = 150
grid.cell_height = 150
grid.spacing = 10
grid.font = &font

widgets.image_grid_add_folder(grid, "Pictures", "/home/user/Pictures")
widgets.image_grid_add_item(grid, image, thumbnail, "/path/to/image.png")

grid.on_click = proc(g: ^widgets.Image_Grid, idx: i32, item: ^widgets.Grid_Item) {
    // Handle selection
}

widgets.widget_add_child(parent, grid)
```

## Theme System

Odek includes a theme system with dark and light themes:

```odin
// Get current theme
theme := widgets.theme_get()

// Use theme colors
container.background = theme.bg_primary
label.color = theme.text_primary

// Switch themes
widgets.theme_set_dark()   // Default
widgets.theme_set_light()

// Or set custom theme
widgets.theme_set(&my_custom_theme)
```

**Theme colors:**
| Color | Description |
|-------|-------------|
| `bg_primary` | Main background |
| `bg_secondary` | Secondary background |
| `bg_hover` | Hover state |
| `bg_pressed` | Pressed state |
| `accent` | Primary accent (buttons) |
| `accent_hover` | Accent hover |
| `text_primary` | Main text |
| `text_secondary` | Muted text |
| `input_bg` | Input background |
| `input_border` | Input border |
| `scrollbar_track` | Scrollbar track |
| `scrollbar_thumb` | Scrollbar thumb |

## Keyboard Constants

Use named constants instead of magic numbers:

```odin
import "core"

// Keycodes (evdev scancodes)
if event.keycode == u32(core.Keycode.Enter) { ... }
if event.keycode == u32(core.Keycode.Space) { ... }
if event.keycode == u32(core.Keycode.Tab) { ... }
if event.keycode == u32(core.Keycode.Backspace) { ... }

// Keysyms (XKB)
if event.keysym == u32(core.Keysym.Return) { ... }
if event.keysym == u32(core.Keysym.Tab) { ... }
```

## Layout System

Odek uses a flexbox-inspired layout system:

```odin
// Fixed size child
child1 := widgets.container_create(.Row)
child1.min_size = core.Size{100, 50}

// Flexible child (takes remaining space)
child2 := widgets.container_create(.Row)
child2.flex = 1.0

// Proportional flex
child3 := widgets.container_create(.Row)
child3.flex = 2.0  // Gets 2x the space of child2
```

### Edge Insets

```odin
padding := widgets.edges_all(10)               // All edges
padding := widgets.edges_symmetric(20, 10)     // Horizontal, vertical
padding := widgets.Edges{top = 5, right = 10, bottom = 5, left = 10}
```

## Event Handling

### Window Callbacks

```odin
window.on_draw = proc(win: ^Window, pixels: [^]u32, w, h, stride: i32) {
    // Render frame using global state
}

window.on_pointer_motion = proc(win: ^Window, x, y: f64) {
    widgets.update_hover(&g_state.hit_state, g_state.root, i32(x), i32(y))
    core.window_request_redraw(win)
}

window.on_pointer_button = proc(win: ^Window, button: u32, pressed: bool) {
    // Create event and dispatch
}

window.on_key = proc(win: ^Window, key: u32, pressed: bool, utf8: string) {
    // Handle keyboard input
}

window.on_scroll = proc(win: ^Window, delta: i32, axis: u32) {
    // Handle scroll wheel
}
```

### Focus Management

```odin
focus_manager := widgets.focus_manager_init(root)

// In key handler
event := core.Event{type = .Key_Press, keycode = keycode}
if widgets.focus_handle_tab(&focus_manager, &event) {
    // Tab navigation handled
    return
}

// Send to focused widget
focused := widgets.focus_get(&focus_manager)
if focused != nil {
    widgets.widget_handle_event(focused, &event)
}
```

## Rendering

### Draw Context

```odin
ctx := render.context_create_scaled(
    pixels, phys_width, phys_height, stride,
    logical_width, logical_height, scale)

render.clear(&ctx, background_color)
render.fill_rect(&ctx, rect, color)
render.fill_rounded_rect(&ctx, rect, radius, color)
render.draw_text_top(&ctx, &font, "text", x, y, color)
render.draw_image(&ctx, &image, x, y)
```

### Text Rendering

```odin
text_renderer, ok := render.text_renderer_init()
defer render.text_renderer_destroy(&text_renderer)

font, ok := render.font_load(&text_renderer, "/path/to/font.ttf", 14)
defer render.font_destroy(&font)

// Measure text
width := render.text_measure(&font, "Hello")

// Handle HiDPI
render.font_set_scale(&font, window.scale)
```

## Examples

The `examples/` directory contains complete sample applications:

- **todo** - A simple todo list app with persistence
- **catalog** - Component showcase demonstrating all widgets
- **filebrowser** - Image browser with directory navigation and video thumbnail preview

Run an example:

```bash
odin build examples/todo -out:todo && ./todo
odin build examples/catalog -out:catalog && ./catalog
odin build examples/filebrowser -out:filebrowser && ./filebrowser
```

## Limitations

- **Wayland-only**: No X11 support. Requires a Wayland compositor.
- **Linux-only**: Uses Linux-specific APIs (eventfd, etc.)
- **Software rendering**: No GPU acceleration. Performance may be limited with many widgets.
- **No closures**: Odin doesn't support closures, so callbacks must use global state.

## Testing

Run the test suite:

```bash
odin test tests
```

## License

MIT
