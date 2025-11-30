package wayland

import "core:c"

// Fractional scale protocol - compositor-driven fractional scaling
// https://wayland.app/protocols/fractional-scale-v1

// Opaque types
Wp_Fractional_Scale_Manager_V1 :: struct {}
Wp_Fractional_Scale_V1 :: struct {}

// Message definitions for wp_fractional_scale_manager_v1
// Requests: destroy (0), get_fractional_scale (1)
@(private)
wp_fractional_scale_manager_v1_requests := [2]Wl_Message{
	{name = "destroy", signature = "", types = nil},
	{name = "get_fractional_scale", signature = "no", types = &wp_fractional_scale_v1_types[0]},
}

@(private)
wp_fractional_scale_v1_types := [2]rawptr{
	&wp_fractional_scale_v1_interface,
	&wl_surface_interface,
}

// Message definitions for wp_fractional_scale_v1
// Requests: destroy (0)
// Events: preferred_scale (0)
@(private)
wp_fractional_scale_v1_requests := [1]Wl_Message{
	{name = "destroy", signature = "", types = nil},
}

@(private)
wp_fractional_scale_v1_events := [1]Wl_Message{
	{name = "preferred_scale", signature = "u", types = nil}, // scale as u32, denominator is 120
}

// Interface definitions
wp_fractional_scale_manager_v1_interface := Wl_Interface{
	name = "wp_fractional_scale_manager_v1",
	version = 1,
	method_count = 2,
	methods = &wp_fractional_scale_manager_v1_requests[0],
	event_count = 0,
	events = nil,
}

wp_fractional_scale_v1_interface := Wl_Interface{
	name = "wp_fractional_scale_v1",
	version = 1,
	method_count = 1,
	methods = &wp_fractional_scale_v1_requests[0],
	event_count = 1,
	events = &wp_fractional_scale_v1_events[0],
}

// Listener for wp_fractional_scale_v1 events
Wp_Fractional_Scale_V1_Listener :: struct {
	// preferred_scale: compositor's preferred scale factor
	// scale is numerator with denominator 120 (e.g., 120 = 1.0, 180 = 1.5, 240 = 2.0)
	preferred_scale: proc "c" (data: rawptr, fractional_scale: ^Wp_Fractional_Scale_V1, scale: u32),
}

// wp_fractional_scale_manager_v1 operations

fractional_scale_manager_destroy :: proc(manager: ^Wp_Fractional_Scale_Manager_V1) {
	wl_proxy_marshal_flags(manager, 0, nil, wl_proxy_get_version(manager), WL_MARSHAL_FLAG_DESTROY)
}

fractional_scale_manager_get_fractional_scale :: proc(
	manager: ^Wp_Fractional_Scale_Manager_V1,
	surface: ^Wl_Surface,
) -> ^Wp_Fractional_Scale_V1 {
	args: [2]Wl_Argument
	args[0].o = nil // new_id placeholder
	args[1].o = surface
	return cast(^Wp_Fractional_Scale_V1)wl_proxy_marshal_array_flags(
		manager, 1, &wp_fractional_scale_v1_interface, wl_proxy_get_version(manager), 0, &args[0])
}

// wp_fractional_scale_v1 operations

fractional_scale_add_listener :: proc(
	fractional_scale: ^Wp_Fractional_Scale_V1,
	listener: ^Wp_Fractional_Scale_V1_Listener,
	data: rawptr,
) -> c.int {
	return wl_proxy_add_listener(fractional_scale, listener, data)
}

fractional_scale_destroy :: proc(fractional_scale: ^Wp_Fractional_Scale_V1) {
	wl_proxy_marshal_flags(fractional_scale, 0, nil, wl_proxy_get_version(fractional_scale), WL_MARSHAL_FLAG_DESTROY)
}

// Helper to convert scale factor from protocol format to f64
// Protocol uses numerator with denominator 120
fractional_scale_to_f64 :: proc(scale: u32) -> f64 {
	return f64(scale) / 120.0
}
