# Odek MVP Progress

## Completed

### Phase 1: Wayland Foundation
- [x] Wayland client bindings (wl_display, wl_registry, wl_compositor, wl_shm)
- [x] XDG shell (xdg_wm_base, xdg_surface, xdg_toplevel)
- [x] SHM buffer management with memfd_create + mmap
- [x] Window resize handling (deferred resize when buffer busy)
- [x] Frame callback event loop
- [x] Basic software rendering primitives (rect, rounded rect, lines)

### Phase 2: Input Handling
- [x] wl_seat binding and capability handling
- [x] Pointer events (enter, leave, motion, button)
- [x] Keyboard events (keymap, key, modifiers)
- [x] Full XKB integration via libxkbcommon
- [x] Key-to-UTF8 conversion
- [x] Modifier tracking (Shift, Ctrl, Alt, Super)
- [x] Input state query functions

### Phase 3: FreeType Text Rendering
- [x] FFI bindings to libfreetype
- [x] Font loading (system fonts or bundled TTF)
- [x] Glyph rendering with antialiasing
- [x] Glyph caching for performance
- [x] Text measurement for layout

### Phase 4: Widget System
- [x] Base Widget struct with vtable
- [x] Widget tree (parent/children)
- [x] Event dispatch and hit testing
- [x] Focus management
- [x] Flexbox-lite layout (padding, margin, direction, alignment, flex grow)
- [x] Immediate damage tracking (dirty flag propagation)

## Remaining

### Phase 5: MVP Widgets
- [ ] Label - text content, font/size/color, alignment
- [ ] Button - normal/hover/pressed states, on_click callback
- [ ] Text Input - cursor, single-line editing, on_change/on_submit

### Future: Cursor Support
- [ ] Implement libwayland-cursor for proper cursor theming
- [ ] Or implement wp_cursor_shape_manager_v1 protocol bindings

## Known Issues

- **Cursor shows as text cursor**: Requires libwayland-cursor implementation (cursor theme loading)
- **Demo sidebar items don't navigate**: These are static text, not Button widgets (Button is Phase 5)

## Test Status

- 57 tests passing
- Core types: rect math, color conversion
- Rendering: fill rect, clipping
- Input: fixed-point conversion, button codes, state enums
- Text: renderer init, font loading, glyph caching, text measurement, text drawing
- Widgets: creation/destruction, parent/child, hit testing, focus, container layout, alignment

## Files

```
src/
  core/
    app.odin      - Application lifecycle, window management, input handling
    event.odin    - Event types
    types.odin    - Rect, Color, Point
  wayland/
    client.odin   - Core Wayland FFI bindings
    shm.odin      - SHM buffer management
    xdg_shell.odin - XDG shell protocol
    seat.odin     - Input device bindings
    xkb.odin      - libxkbcommon bindings
  render/
    buffer.odin   - Software rendering primitives
    freetype.odin - FreeType FFI bindings
    text.odin     - Text rendering, font loading, glyph caching
  widgets/
    widget.odin   - Base Widget struct, vtable, core operations
    container.odin - Container with flexbox-lite layout
    hit_test.odin - Hit testing and event dispatch
    focus.odin    - Focus management
  main.odin       - Demo application
tests/
  core_test.odin
  render_test.odin
  input_test.odin
  text_test.odin
  widget_test.odin
```

## Dependencies

- libwayland-client.so (system)
- libxkbcommon.so (system)
- libfreetype.so (system)

## Build

```bash
odin build src -out:odek
odin test tests
```
