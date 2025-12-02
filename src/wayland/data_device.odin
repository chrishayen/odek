package wayland

import "core:c"

// Opaque types
Wl_Data_Device_Manager :: struct {}
Wl_Data_Device :: struct {}
Wl_Data_Source :: struct {}
Wl_Data_Offer :: struct {}

// ============================================================================
// Types arrays for messages with object references
// Each position corresponds to an argument in the signature:
// - Object/new_id args get &interface pointer
// - Other args (u, i, s, f, a, h) get nil
// ============================================================================

// For data_device_manager.create_data_source (signature "n")
@(private)
data_device_manager_create_data_source_types := [1]rawptr{
    &wl_data_source_interface,  // n: new_id<wl_data_source>
}

// For data_device_manager.get_data_device (signature "no")
@(private)
data_device_manager_get_data_device_types := [2]rawptr{
    &wl_data_device_interface,  // n: new_id<wl_data_device>
    &wl_seat_interface,         // o: object<wl_seat>
}

// For data_device.start_drag (signature "?oo?ou")
@(private)
data_device_start_drag_types := [4]rawptr{
    &wl_data_source_interface,  // ?o: object<wl_data_source>, allow-null
    &wl_surface_interface,      // o: object<wl_surface>
    &wl_surface_interface,      // ?o: object<wl_surface>, allow-null
    nil,                        // u: uint
}

// For data_device.set_selection (signature "?ou")
@(private)
data_device_set_selection_types := [2]rawptr{
    &wl_data_source_interface,  // ?o: object<wl_data_source>, allow-null
    nil,                        // u: uint
}

// For data_device.data_offer event (signature "n")
@(private)
data_device_data_offer_types := [1]rawptr{
    &wl_data_offer_interface,   // n: new_id<wl_data_offer>
}

// For data_device.enter event (signature "uoff?o")
@(private)
data_device_enter_types := [5]rawptr{
    nil,                        // u: uint (serial)
    &wl_surface_interface,      // o: object<wl_surface>
    nil,                        // f: fixed (x)
    nil,                        // f: fixed (y)
    &wl_data_offer_interface,   // ?o: object<wl_data_offer>, allow-null
}

// For data_device.selection event (signature "?o")
@(private)
data_device_selection_types := [1]rawptr{
    &wl_data_offer_interface,   // ?o: object<wl_data_offer>, allow-null
}

// ============================================================================
// Message definitions
// ============================================================================

// Data device manager requests (2 requests)
@(private)
data_device_manager_requests := [2]Wl_Message{
    {name = "create_data_source", signature = "n", types = &data_device_manager_create_data_source_types[0]},
    {name = "get_data_device", signature = "no", types = &data_device_manager_get_data_device_types[0]},
}

// Data source requests (3 requests)
@(private)
data_source_requests := [3]Wl_Message{
    {name = "offer", signature = "s", types = nil},
    {name = "destroy", signature = "", types = nil},
    {name = "set_actions", signature = "u", types = nil},
}

// Data source events (6 events)
@(private)
data_source_events := [6]Wl_Message{
    {name = "target", signature = "?s", types = nil},
    {name = "send", signature = "sh", types = nil},
    {name = "cancelled", signature = "", types = nil},
    {name = "dnd_drop_performed", signature = "", types = nil},
    {name = "dnd_finished", signature = "", types = nil},
    {name = "action", signature = "u", types = nil},
}

// Data offer requests (5 requests)
@(private)
data_offer_requests := [5]Wl_Message{
    {name = "accept", signature = "u?s", types = nil},
    {name = "receive", signature = "sh", types = nil},
    {name = "destroy", signature = "", types = nil},
    {name = "finish", signature = "", types = nil},
    {name = "set_actions", signature = "uu", types = nil},
}

// Data offer events (3 events)
@(private)
data_offer_events := [3]Wl_Message{
    {name = "offer", signature = "s", types = nil},
    {name = "source_actions", signature = "u", types = nil},
    {name = "action", signature = "u", types = nil},
}

// Data device requests (3 requests)
@(private)
data_device_requests := [3]Wl_Message{
    {name = "start_drag", signature = "?oo?ou", types = &data_device_start_drag_types[0]},
    {name = "set_selection", signature = "?ou", types = &data_device_set_selection_types[0]},
    {name = "release", signature = "", types = nil},
}

// Data device events (6 events)
@(private)
data_device_events := [6]Wl_Message{
    {name = "data_offer", signature = "n", types = &data_device_data_offer_types[0]},
    {name = "enter", signature = "uoff?o", types = &data_device_enter_types[0]},
    {name = "leave", signature = "", types = nil},
    {name = "motion", signature = "uff", types = nil},
    {name = "drop", signature = "", types = nil},
    {name = "selection", signature = "?o", types = &data_device_selection_types[0]},
}

// ============================================================================
// Interface definitions
// ============================================================================

wl_data_device_manager_interface := Wl_Interface{
    name = "wl_data_device_manager",
    version = 3,
    method_count = 2,
    methods = &data_device_manager_requests[0],
    event_count = 0,
    events = nil,
}

wl_data_source_interface := Wl_Interface{
    name = "wl_data_source",
    version = 3,
    method_count = 3,
    methods = &data_source_requests[0],
    event_count = 6,
    events = &data_source_events[0],
}

wl_data_offer_interface := Wl_Interface{
    name = "wl_data_offer",
    version = 3,
    method_count = 5,
    methods = &data_offer_requests[0],
    event_count = 3,
    events = &data_offer_events[0],
}

wl_data_device_interface := Wl_Interface{
    name = "wl_data_device",
    version = 3,
    method_count = 3,
    methods = &data_device_requests[0],
    event_count = 6,
    events = &data_device_events[0],
}

// ============================================================================
// Listeners
// ============================================================================

Wl_Data_Source_Listener :: struct {
    target:             proc "c" (data: rawptr, source: ^Wl_Data_Source, mime_type: cstring),
    send:               proc "c" (data: rawptr, source: ^Wl_Data_Source, mime_type: cstring, fd: i32),
    cancelled:          proc "c" (data: rawptr, source: ^Wl_Data_Source),
    dnd_drop_performed: proc "c" (data: rawptr, source: ^Wl_Data_Source),
    dnd_finished:       proc "c" (data: rawptr, source: ^Wl_Data_Source),
    action:             proc "c" (data: rawptr, source: ^Wl_Data_Source, dnd_action: u32),
}

Wl_Data_Offer_Listener :: struct {
    offer:          proc "c" (data: rawptr, offer: ^Wl_Data_Offer, mime_type: cstring),
    source_actions: proc "c" (data: rawptr, offer: ^Wl_Data_Offer, source_actions: u32),
    action:         proc "c" (data: rawptr, offer: ^Wl_Data_Offer, dnd_action: u32),
}

Wl_Data_Device_Listener :: struct {
    data_offer: proc "c" (data: rawptr, device: ^Wl_Data_Device, offer: ^Wl_Data_Offer),
    enter:      proc "c" (data: rawptr, device: ^Wl_Data_Device, serial: u32, surface: ^Wl_Surface, x: i32, y: i32, offer: ^Wl_Data_Offer),
    leave:      proc "c" (data: rawptr, device: ^Wl_Data_Device),
    motion:     proc "c" (data: rawptr, device: ^Wl_Data_Device, time: u32, x: i32, y: i32),
    drop:       proc "c" (data: rawptr, device: ^Wl_Data_Device),
    selection:  proc "c" (data: rawptr, device: ^Wl_Data_Device, offer: ^Wl_Data_Offer),
}

// ============================================================================
// Data Device Manager operations
// ============================================================================

// Create a data source (opcode 0)
data_device_manager_create_data_source :: proc(manager: ^Wl_Data_Device_Manager) -> ^Wl_Data_Source {
    return cast(^Wl_Data_Source)wl_proxy_marshal_flags(
        manager, 0, &wl_data_source_interface, wl_proxy_get_version(manager), 0)
}

// Get data device for a seat (opcode 1)
data_device_manager_get_data_device :: proc(manager: ^Wl_Data_Device_Manager, seat: ^Wl_Seat) -> ^Wl_Data_Device {
    args: [2]Wl_Argument
    args[0].o = nil  // new_id placeholder
    args[1].o = seat
    return cast(^Wl_Data_Device)wl_proxy_marshal_array_flags(
        manager, 1, &wl_data_device_interface, wl_proxy_get_version(manager), 0, &args[0])
}

// Destroy data device manager (just destroy proxy, no protocol request)
data_device_manager_destroy :: proc(manager: ^Wl_Data_Device_Manager) {
    wl_proxy_destroy(manager)
}

// ============================================================================
// Data Source operations
// ============================================================================

data_source_add_listener :: proc(source: ^Wl_Data_Source, listener: ^Wl_Data_Source_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(source, listener, data)
}

// Offer a MIME type (opcode 0)
data_source_offer :: proc(source: ^Wl_Data_Source, mime_type: cstring) {
    wl_proxy_marshal_flags(source, 0, nil, wl_proxy_get_version(source), 0, mime_type)
}

// Destroy data source (opcode 1)
data_source_destroy :: proc(source: ^Wl_Data_Source) {
    wl_proxy_marshal_flags(source, 1, nil, wl_proxy_get_version(source), WL_MARSHAL_FLAG_DESTROY)
}

// ============================================================================
// Data Offer operations
// ============================================================================

data_offer_add_listener :: proc(offer: ^Wl_Data_Offer, listener: ^Wl_Data_Offer_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(offer, listener, data)
}

// Request data for a MIME type (opcode 1)
data_offer_receive :: proc(offer: ^Wl_Data_Offer, mime_type: cstring, fd: i32) {
    wl_proxy_marshal_flags(offer, 1, nil, wl_proxy_get_version(offer), 0, mime_type, fd)
}

// Destroy data offer (opcode 2)
data_offer_destroy :: proc(offer: ^Wl_Data_Offer) {
    wl_proxy_marshal_flags(offer, 2, nil, wl_proxy_get_version(offer), WL_MARSHAL_FLAG_DESTROY)
}

// ============================================================================
// Data Device operations
// ============================================================================

data_device_add_listener :: proc(device: ^Wl_Data_Device, listener: ^Wl_Data_Device_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(device, listener, data)
}

// Set clipboard selection (opcode 1)
data_device_set_selection :: proc(device: ^Wl_Data_Device, source: ^Wl_Data_Source, serial: u32) {
    wl_proxy_marshal_flags(device, 1, nil, wl_proxy_get_version(device), 0, source, serial)
}

// Release data device (opcode 2, version 2+)
data_device_release :: proc(device: ^Wl_Data_Device) {
    if wl_proxy_get_version(device) >= 2 {
        wl_proxy_marshal_flags(device, 2, nil, wl_proxy_get_version(device), WL_MARSHAL_FLAG_DESTROY)
    } else {
        wl_proxy_destroy(device)
    }
}
