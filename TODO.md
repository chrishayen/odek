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

## Remaining

### Phase 4: Widget System
- [ ] Base Widget struct with vtable
- [ ] Widget tree (parent/children)
- [ ] Event dispatch and hit testing
- [ ] Focus management
- [ ] Simple box layout (padding, margin)
- [ ] Damage tracking

### Phase 5: MVP Widgets
- [ ] Label - text content, font/size/color, alignment
- [ ] Button - normal/hover/pressed states, on_click callback
- [ ] Text Input - cursor, single-line editing, on_change/on_submit
- [ ] Container - vertical/horizontal box layout with spacing

## Test Status

- 31 tests passing
- Core types: rect math, color conversion
- Rendering: fill rect, clipping
- Input: fixed-point conversion, button codes, state enums
- Text: renderer init, font loading, glyph caching, text measurement, text drawing

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
  main.odin       - Demo application
tests/
  core_test.odin
  render_test.odin
  input_test.odin
  text_test.odin
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
