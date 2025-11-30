package wayland

import "core:c"

// Viewporter protocol - surface cropping and scaling
// https://wayland.app/protocols/viewporter

// Opaque types
Wp_Viewporter :: struct {}
Wp_Viewport :: struct {}

// Message definitions for wp_viewporter
// Requests: destroy (0), get_viewport (1)
@(private)
wp_viewporter_requests := [2]Wl_Message{
	{name = "destroy", signature = "", types = nil},
	{name = "get_viewport", signature = "no", types = &wp_viewport_types[0]},
}

@(private)
wp_viewport_types := [2]rawptr{
	&wp_viewport_interface,
	&wl_surface_interface,
}

// Message definitions for wp_viewport
// Requests: destroy (0), set_source (1), set_destination (2)
@(private)
wp_viewport_requests := [3]Wl_Message{
	{name = "destroy", signature = "", types = nil},
	{name = "set_source", signature = "ffff", types = nil}, // wl_fixed_t x4
	{name = "set_destination", signature = "ii", types = nil},
}

// Interface definitions
wp_viewporter_interface := Wl_Interface{
	name = "wp_viewporter",
	version = 1,
	method_count = 2,
	methods = &wp_viewporter_requests[0],
	event_count = 0,
	events = nil,
}

wp_viewport_interface := Wl_Interface{
	name = "wp_viewport",
	version = 1,
	method_count = 3,
	methods = &wp_viewport_requests[0],
	event_count = 0,
	events = nil,
}

// wp_viewporter operations

viewporter_destroy :: proc(viewporter: ^Wp_Viewporter) {
	wl_proxy_marshal_flags(viewporter, 0, nil, wl_proxy_get_version(viewporter), WL_MARSHAL_FLAG_DESTROY)
}

viewporter_get_viewport :: proc(viewporter: ^Wp_Viewporter, surface: ^Wl_Surface) -> ^Wp_Viewport {
	args: [2]Wl_Argument
	args[0].o = nil // new_id placeholder
	args[1].o = surface
	return cast(^Wp_Viewport)wl_proxy_marshal_array_flags(
		viewporter, 1, &wp_viewport_interface, wl_proxy_get_version(viewporter), 0, &args[0])
}

// wp_viewport operations

viewport_destroy :: proc(viewport: ^Wp_Viewport) {
	wl_proxy_marshal_flags(viewport, 0, nil, wl_proxy_get_version(viewport), WL_MARSHAL_FLAG_DESTROY)
}

// Set source rectangle (crop region) in buffer coordinates
// Pass -1.0 for all values to unset
viewport_set_source :: proc(viewport: ^Wp_Viewport, x, y, width, height: f64) {
	wl_proxy_marshal_flags(
		viewport, 1, nil, wl_proxy_get_version(viewport), 0,
		wl_double_to_fixed(x),
		wl_double_to_fixed(y),
		wl_double_to_fixed(width),
		wl_double_to_fixed(height))
}

// Set destination size in surface coordinates (logical size)
// Pass -1 for both to unset
viewport_set_destination :: proc(viewport: ^Wp_Viewport, width, height: i32) {
	wl_proxy_marshal_flags(viewport, 2, nil, wl_proxy_get_version(viewport), 0, width, height)
}
