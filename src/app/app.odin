package app

// High-level API for odek UI toolkit
// Provides GTK/Qt-style convenience - just create app, add widgets, run.

import "../core"
import "../render"
import "../widgets"
import "core:time"

// App encapsulates all state - no global variables needed by user
App :: struct {
	// Core Wayland state
	core_app:        ^core.App,
	window:          ^core.Window,

	// Text rendering (auto-initialized)
	text_renderer:   render.Text_Renderer,
	font:            render.Font,
	font_loaded:     bool,

	// Widget tree
	root:            ^widgets.Container,

	// State management (auto-handled)
	focus_manager:   widgets.Focus_Manager,
	hit_state:       widgets.Hit_Test_State,

	// Image loading (auto-initialized)
	image_cache:     ^render.Image_Cache,
	image_loader:    ^render.Image_Loader,

	// Delta time tracking for animations/video
	last_frame_time: f64,

	// Image grids to update with loaded images
	image_grids:     [dynamic]^widgets.Image_Grid,

	// User callbacks for customization
	on_draw_overlay: proc(a: ^App, ctx: ^render.Draw_Context),
	on_key:          proc(a: ^App, keycode: u32, pressed: bool, utf8: string) -> bool,
	on_scroll:       proc(a: ^App, delta: i32, axis: u32) -> bool,
}

// Common font paths to try
DEFAULT_FONT_PATHS :: []string {
	// Arch Linux / TTF
	"/usr/share/fonts/TTF/DejaVuSans.ttf",
	"/usr/share/fonts/TTF/LiberationSans-Regular.ttf",
	"/usr/share/fonts/noto/NotoSans-Regular.ttf",
	// Debian/Ubuntu
	"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
	"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
	"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",
	// Fedora
	"/usr/share/fonts/dejavu-sans-fonts/DejaVuSans.ttf",
	"/usr/share/fonts/liberation-sans/LiberationSans-Regular.ttf",
	// Generic
	"/usr/share/fonts/TTF/FreeSans.ttf",
	"/usr/share/fonts/truetype/freefont/FreeSans.ttf",
}

// Image cache configuration
IMAGE_CACHE_MIN_SIZE :: 50   // Minimum images to keep in cache
IMAGE_CACHE_MAX_SIZE :: 150  // Maximum images before eviction
IMAGE_LOADER_QUEUE_SIZE :: 256  // Max pending image load requests

// Global app pointer for callbacks (internal use only)
// This is necessary because Wayland callbacks can't capture state
@(private)
g_app: ^App

// Create a new application
create :: proc(title: string, width: i32 = 800, height: i32 = 600) -> ^App {
	// Only one App instance is supported due to Wayland callback limitations
	if g_app != nil {
		return nil
	}

	a := new(App)

	// Initialize Wayland
	a.core_app = core.init()
	if a.core_app == nil {
		free(a)
		return nil
	}

	// Create window
	a.window = core.create_window(a.core_app, title, width, height)
	if a.window == nil {
		core.shutdown(a.core_app)
		free(a)
		return nil
	}

	// Initialize text renderer
	ok: bool
	a.text_renderer, ok = render.text_renderer_init()
	if !ok {
		core.shutdown(a.core_app)
		free(a)
		return nil
	}

	// Auto-load font using system-configured size from fontconfig
	font_size := render.fc_get_default_pixel_size(14)
	a.font_loaded = false
	for path in DEFAULT_FONT_PATHS {
		a.font, ok = render.font_load(&a.text_renderer, path, font_size)
		if ok {
			a.font_loaded = true
			break
		}
	}

	// Create root container
	theme := widgets.theme_get()
	a.root = widgets.container_create(.Column)
	a.root.background = theme.bg_primary

	// Initialize focus manager
	a.focus_manager = widgets.focus_manager_init(a.root)

	// Initialize image cache and loader
	a.image_cache = render.image_cache_create(IMAGE_CACHE_MIN_SIZE, IMAGE_CACHE_MAX_SIZE)
	a.image_loader = render.image_loader_create(IMAGE_LOADER_QUEUE_SIZE)

	// Set up callbacks (store app pointer for callbacks)
	g_app = a

	// Set global state for widget destroy notifications
	widgets.hit_state_set_global(&a.hit_state)
	widgets.focus_manager_set_global(&a.focus_manager)

	// Register image loader notification FD for event-driven updates
	loader_fd := render.image_loader_get_fd(a.image_loader)
	core.app_add_poll_fd(a.core_app, loader_fd, _image_load_complete_callback)
	a.window.on_draw = _draw_callback
	a.window.on_close = _close_callback
	a.window.on_pointer_enter = _pointer_enter_callback
	a.window.on_pointer_leave = _pointer_leave_callback
	a.window.on_pointer_motion = _pointer_motion_callback
	a.window.on_pointer_button = _pointer_button_callback
	a.window.on_scroll = _scroll_callback
	a.window.on_key = _key_callback
	a.window.on_scale_changed = _scale_changed_callback

	return a
}

// Destroy the application
destroy :: proc(a: ^App) {
	if a == nil {
		return
	}

	if a.root != nil {
		widgets.widget_destroy(a.root)
	}

	if a.font_loaded {
		render.font_destroy(&a.font)
	}

	render.text_renderer_destroy(&a.text_renderer)

	// Clean up image loading
	if a.image_loader != nil {
		render.image_loader_destroy(a.image_loader)
	}
	if a.image_cache != nil {
		render.image_cache_destroy(a.image_cache)
	}
	delete(a.image_grids)

	if a.core_app != nil {
		core.shutdown(a.core_app)
	}

	if g_app == a {
		g_app = nil
	}

	free(a)
}

// Run the application event loop
run :: proc(a: ^App) {
	if a == nil || a.core_app == nil {
		return
	}
	core.run(a.core_app)
}

// Get the app's font (for custom widgets)
get_font :: proc(a: ^App) -> ^render.Font {
	if a == nil || !a.font_loaded {
		return nil
	}
	return &a.font
}

// ============================================================================
// Widget Factory Functions (add to root)
// ============================================================================

// Create a label with the app's font and add to root
label :: proc(a: ^App, text: string = "") -> ^widgets.Label {
	l := create_label(a, text)
	widgets.widget_add_child(a.root, l)
	return l
}

// Create a button with the app's font and add to root
button :: proc(a: ^App, text: string = "") -> ^widgets.Button {
	b := create_button(a, text)
	widgets.widget_add_child(a.root, b)
	return b
}

// Create a text input with the app's font and add to root
text_input :: proc(a: ^App) -> ^widgets.Text_Input {
	ti := create_text_input(a)
	widgets.widget_add_child(a.root, ti)
	return ti
}

// Create a container and add to root
container :: proc(a: ^App, direction: widgets.Direction = .Column) -> ^widgets.Container {
	c := create_container(direction)
	widgets.widget_add_child(a.root, c)
	return c
}

// Create a scroll container and add to root
scroll_container :: proc(
	a: ^App,
	direction: widgets.Scroll_Direction = .Vertical,
) -> ^widgets.Scroll_Container {
	sc := create_scroll_container(direction)
	widgets.widget_add_child(a.root, sc)
	return sc
}

// Create an image grid and add to root (auto-wired for async loading)
image_grid :: proc(a: ^App) -> ^widgets.Image_Grid {
	ig := create_image_grid(a)
	widgets.widget_add_child(a.root, ig)
	return ig
}

// ============================================================================
// Widget Creation Functions (don't add to root - for custom layouts)
// ============================================================================

// Create a label with the app's font (does not add to root)
create_label :: proc(a: ^App, text: string = "") -> ^widgets.Label {
	return widgets.label_create(text, get_font(a))
}

// Create a button with the app's font (does not add to root)
create_button :: proc(a: ^App, text: string = "") -> ^widgets.Button {
	return widgets.button_create(text, get_font(a))
}

// Create a text input with the app's font (does not add to root)
create_text_input :: proc(a: ^App) -> ^widgets.Text_Input {
	ti := widgets.text_input_create(get_font(a))
	ti.min_size = core.Size{0, 32}
	return ti
}

// Create a checkbox (does not add to root)
create_checkbox :: proc(a: ^App) -> ^widgets.Checkbox {
	return widgets.checkbox_create()
}

// Create a container (does not add to root)
create_container :: proc(direction: widgets.Direction = .Column) -> ^widgets.Container {
	return widgets.container_create(direction)
}

// Create a scroll container (does not add to root)
create_scroll_container :: proc(
	direction: widgets.Scroll_Direction = .Vertical,
) -> ^widgets.Scroll_Container {
	return widgets.scroll_container_create(direction)
}

// Create an image grid (does not add to root, auto-wired for async loading)
create_image_grid :: proc(a: ^App) -> ^widgets.Image_Grid {
	ig := widgets.image_grid_create()
	if a.font_loaded {
		ig.font = &a.font
	}
	// Track grid for image load updates
	append(&a.image_grids, ig)
	return ig
}

// Unregister an image grid from async loading updates (call before destroying the grid)
unregister_image_grid :: proc(a: ^App, grid: ^widgets.Image_Grid) {
	for i := 0; i < len(a.image_grids); i += 1 {
		if a.image_grids[i] == grid {
			unordered_remove(&a.image_grids, i)
			return
		}
	}
}

// ============================================================================
// Image Loading API
// ============================================================================

// Queue an image for async loading into a grid
queue_image_load :: proc(a: ^App, grid: ^widgets.Image_Grid, path: string, name: string) -> i32 {
	idx := widgets.image_grid_add_placeholder(grid, name, path)
	render.image_loader_queue(a.image_loader, path, idx)
	return idx
}

// Clear pending image loads (call when navigating away)
clear_image_loads :: proc(a: ^App) {
	if a.image_loader != nil {
		render.image_loader_clear(a.image_loader)
	}
}

// Get the image cache
get_image_cache :: proc(a: ^App) -> ^render.Image_Cache {
	return a.image_cache
}

// Get the image loader
get_image_loader :: proc(a: ^App) -> ^render.Image_Loader {
	return a.image_loader
}

// Get the window
get_window :: proc(a: ^App) -> ^core.Window {
	return a.window
}

// Get the core app
get_core_app :: proc(a: ^App) -> ^core.App {
	return a.core_app
}

// Request a redraw
request_redraw :: proc(a: ^App) {
	if a != nil && a.window != nil {
		core.window_request_redraw(a.window)
	}
}

// Get the root container for custom layouts
get_root :: proc(a: ^App) -> ^widgets.Container {
	return a.root
}

// Set minimum window size (hint to compositor)
set_min_size :: proc(a: ^App, width, height: i32) {
	if a != nil && a.window != nil {
		core.window_set_min_size(a.window, width, height)
	}
}

// Get current pointer position
get_pointer_pos :: proc(a: ^App) -> (x, y: f64) {
	if a == nil || a.core_app == nil {
		return 0, 0
	}
	return core.get_pointer_pos(a.core_app)
}

// Get the hit test state
get_hit_state :: proc(a: ^App) -> ^widgets.Hit_Test_State {
	return &a.hit_state
}

// Get the focus manager
get_focus_manager :: proc(a: ^App) -> ^widgets.Focus_Manager {
	return &a.focus_manager
}

// ============================================================================
// Layout Helpers
// ============================================================================

// Arrange children in a column layout
column :: proc(a: ^App, children: []^widgets.Widget, spacing: i32 = 10, padding: i32 = 20) {
	_set_layout(a, .Column, children, spacing, padding)
}

// Arrange children in a row layout
row :: proc(a: ^App, children: []^widgets.Widget, spacing: i32 = 10, padding: i32 = 20) {
	_set_layout(a, .Row, children, spacing, padding)
}

@(private)
_set_layout :: proc(a: ^App, direction: widgets.Direction, children: []^widgets.Widget, spacing: i32, padding: i32) {
	a.root.spacing = spacing
	a.root.padding = widgets.edges_all(padding)
	widgets.container_set_direction(a.root, direction)

	// Clear existing children and add new ones
	for len(a.root.children) > 0 {
		widgets.widget_remove_child(a.root, a.root.children[0])
	}

	for child in children {
		if child != nil {
			widgets.widget_add_child(a.root, child)
		}
	}

	// Re-init focus manager with new tree
	a.focus_manager = widgets.focus_manager_init(a.root)
	widgets.focus_manager_set_global(&a.focus_manager)
}

// ============================================================================
// Internal Callbacks
// ============================================================================

@(private)
_draw_callback :: proc(win: ^core.Window, pixels: [^]u32, w, h, stride: i32) {
	if g_app == nil {
		return
	}

	ctx := render.context_create_scaled(
		pixels,
		w,
		h,
		stride,
		win.logical_width,
		win.logical_height,
		win.scale,
	)

	// Calculate delta time for animations/video
	current_time := f64(time.now()._nsec) / 1_000_000_000.0
	delta_time: f64 = 0.0
	if g_app.last_frame_time > 0 {
		delta_time = current_time - g_app.last_frame_time
	}
	g_app.last_frame_time = current_time

	// Update video thumbnails in all tracked image grids
	if delta_time > 0 {
		for grid in g_app.image_grids {
			if widgets.image_grid_update_videos(grid, delta_time) {
				// Videos updated - request another frame
				core.window_request_redraw(win)
			}
		}
	}

	theme := widgets.theme_get()
	render.clear(&ctx, theme.bg_primary)

	g_app.root.rect = core.Rect{0, 0, win.logical_width, win.logical_height}
	widgets.widget_layout(g_app.root)
	widgets.widget_draw(g_app.root, &ctx)

	// Call user's overlay callback for custom drawing on top
	if g_app.on_draw_overlay != nil {
		g_app.on_draw_overlay(g_app, &ctx)
	}
}

@(private)
_close_callback :: proc(win: ^core.Window) {
	if g_app != nil && g_app.core_app != nil {
		g_app.core_app.running = false
	}
}

@(private)
_pointer_enter_callback :: proc(win: ^core.Window, x, y: f64) {
	// Nothing special needed
}

@(private)
_pointer_leave_callback :: proc(win: ^core.Window) {
	if g_app == nil {
		return
	}
	// Clear hover state
	widgets.update_hover(&g_app.hit_state, g_app.root, -1000, -1000)
	core.window_request_redraw(win)
}

@(private)
_pointer_motion_callback :: proc(win: ^core.Window, x, y: f64) {
	if g_app == nil {
		return
	}

	hovered := widgets.update_hover(&g_app.hit_state, g_app.root, i32(x), i32(y))

	// Update cursor based on hovered widget
	if hovered != nil {
		#partial switch hovered.cursor {
		case .Hand:
			core.set_cursor(g_app.core_app, .Hand)
		case .Text:
			core.set_cursor(g_app.core_app, .Text)
		case:
			core.set_cursor(g_app.core_app, .Arrow)
		}
	} else {
		core.set_cursor(g_app.core_app, .Arrow)
	}

	core.window_request_redraw(win)
}

@(private)
_pointer_button_callback :: proc(win: ^core.Window, button: u32, pressed: bool) {
	if g_app == nil {
		return
	}

	x, y := core.get_pointer_pos(g_app.core_app)
	event := core.event_pointer_button(core.Mouse_Button(button), pressed, i32(x), i32(y), 0)
	widgets.dispatch_pointer_event(&g_app.hit_state, g_app.root, &event)

	// Focus clicked widget if focusable
	if pressed && button == u32(core.Mouse_Button.Left) {
		hovered := g_app.hit_state.hovered
		if hovered != nil && hovered.focusable {
			widgets.focus_set(&g_app.focus_manager, hovered)
		}
	}

	core.window_request_redraw(win)
}

@(private)
_scroll_callback :: proc(win: ^core.Window, delta: i32, axis: u32) {
	if g_app == nil {
		return
	}

	// Let user handle scroll events first
	if g_app.on_scroll != nil {
		if g_app.on_scroll(g_app, delta, axis) {
			core.window_request_redraw(win)
			return
		}
	}

	x, y := core.get_pointer_pos(g_app.core_app)
	event := core.event_scroll(delta, axis, i32(x), i32(y))
	widgets.dispatch_pointer_event(&g_app.hit_state, g_app.root, &event)

	core.window_request_redraw(win)
}

@(private)
_key_callback :: proc(win: ^core.Window, keycode: u32, pressed: bool, utf8: string) {
	if g_app == nil {
		return
	}

	// Let user handle key events first
	if g_app.on_key != nil {
		if g_app.on_key(g_app, keycode, pressed, utf8) {
			core.window_request_redraw(win)
			return
		}
	}

	if !pressed {
		return
	}

	event := core.Event {
		type    = .Key_Press,
		keycode = keycode,
	}

	// Handle Tab for focus navigation
	if widgets.focus_handle_tab(&g_app.focus_manager, &event) {
		core.window_request_redraw(win)
		return
	}

	// Send to focused widget
	focused := widgets.focus_get(&g_app.focus_manager)
	if focused != nil {
		// Use UTF-8 from XKB (has correct case from modifier state)
		if len(utf8) > 0 {
			event.keysym = u32(utf8[0])
		}

		widgets.widget_handle_event(focused, &event)
		core.window_request_redraw(win)
	}
}

@(private)
_scale_changed_callback :: proc(win: ^core.Window, scale: f64) {
	if g_app == nil || !g_app.font_loaded {
		return
	}
	render.font_set_scale(&g_app.font, scale)
}

// Called when image loader has completed work (event-driven via eventfd)
@(private)
_image_load_complete_callback :: proc(app: ^core.App, user_data: rawptr) {
	if g_app == nil || g_app.image_loader == nil {
		return
	}

	render.image_loader_acknowledge(g_app.image_loader)

	// Process completed images
	completed := render.image_loader_get_completed(g_app.image_loader)
	if completed == nil || len(completed) == 0 {
		return
	}

	for result in completed {
		if result.success {
			// Update all tracked image grids with this result
			for grid in g_app.image_grids {
				widgets.image_grid_set_image(
					grid,
					result.grid_index,
					result.image,
					result.thumbnail,
				)
			}
		}
		delete(result.path)
	}
	delete(completed)

	// Request redraw
	if g_app.window != nil {
		core.window_request_redraw(g_app.window)
	}
}
