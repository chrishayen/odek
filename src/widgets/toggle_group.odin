package widgets

import "../core"
import "../render"

// Toggle button group - horizontal row of buttons where only one can be selected
Toggle_Group :: struct {
	using base: Widget,

	// Options (owned, cloned on creation)
	options:        [dynamic]string,
	selected_index: int,

	// Font
	font: ^render.Font,

	// Colors
	bg_normal:     core.Color,
	bg_hover:      core.Color,
	bg_selected:   core.Color,
	text_color:    core.Color,
	border_color:  core.Color,

	// Visual properties
	corner_radius: i32,
	button_padding: Edges,

	// Interaction state
	hovered_index: int,
	pressed_index: int,

	// Callback
	on_change: proc(group: ^Toggle_Group),
}

// Shared vtable
toggle_group_vtable := Widget_VTable{
	draw         = toggle_group_draw,
	handle_event = toggle_group_handle_event,
	layout       = toggle_group_layout,
	destroy      = toggle_group_destroy,
	measure      = toggle_group_measure,
}

// Create a new toggle group with the given options
toggle_group_create :: proc(options: []string, font: ^render.Font = nil) -> ^Toggle_Group {
	g := new(Toggle_Group)
	g.vtable = &toggle_group_vtable
	g.visible = true
	g.enabled = true
	g.dirty = true
	g.focusable = true

	// Clone options into owned dynamic array
	g.options = make([dynamic]string, len(options))
	for i := 0; i < len(options); i += 1 {
		g.options[i] = options[i]
	}
	g.selected_index = 0
	g.hovered_index = -1
	g.pressed_index = -1
	g.font = font

	// Colors from theme
	theme := theme_get()
	g.bg_normal = theme.bg_secondary
	g.bg_hover = theme.bg_hover
	g.bg_selected = theme.bg_pressed
	g.text_color = theme.text_primary
	g.border_color = theme.border

	g.corner_radius = 4
	g.button_padding = edges_symmetric(12, 6)
	g.padding = edges_all(0)

	return g
}

// Set selected index
toggle_group_set_selected :: proc(g: ^Toggle_Group, index: int) {
	if index < 0 || index >= len(g.options) {
		return
	}
	if g.selected_index == index {
		return
	}
	g.selected_index = index
	widget_mark_dirty(g)
}

// Get selected index
toggle_group_get_selected :: proc(g: ^Toggle_Group) -> int {
	return g.selected_index
}

// Get selected option text
toggle_group_get_selected_text :: proc(g: ^Toggle_Group) -> string {
	if g.selected_index >= 0 && g.selected_index < len(g.options) {
		return g.options[g.selected_index]
	}
	return ""
}

// Set on_change callback
toggle_group_set_on_change :: proc(g: ^Toggle_Group, callback: proc(group: ^Toggle_Group)) {
	g.on_change = callback
}

// Calculate button widths - returns total width and individual button width
@(private)
toggle_group_button_width :: proc(g: ^Toggle_Group) -> i32 {
	max_text_width: i32 = 0
	if g.font != nil {
		for opt in g.options {
			w := render.text_measure_logical(g.font, opt)
			max_text_width = max(max_text_width, w)
		}
	}
	return max_text_width + g.button_padding.left + g.button_padding.right
}

// Get button index at position
@(private)
toggle_group_index_at :: proc(g: ^Toggle_Group, x, y: i32) -> int {
	abs_rect := widget_get_absolute_rect(g)

	if y < abs_rect.y || y >= abs_rect.y + abs_rect.height {
		return -1
	}

	button_width := toggle_group_button_width(g)
	rel_x := x - abs_rect.x

	if rel_x < 0 {
		return -1
	}

	index := int(rel_x / button_width)
	if index >= len(g.options) {
		return -1
	}

	return index
}

// Draw toggle group
toggle_group_draw :: proc(w: ^Widget, ctx: ^render.Draw_Context) {
	g := cast(^Toggle_Group)w
	abs_rect := widget_get_absolute_rect(w)

	button_width := toggle_group_button_width(g)
	button_height := abs_rect.height

	// Draw each button
	for i := 0; i < len(g.options); i += 1 {
		opt := g.options[i]
		button_x := abs_rect.x + i32(i) * button_width
		button_rect := core.Rect{button_x, abs_rect.y, button_width, button_height}

		// Determine background color
		bg_color: core.Color
		is_selected := i == g.selected_index
		is_hovered := i == g.hovered_index
		is_pressed := i == g.pressed_index

		if is_selected {
			bg_color = g.bg_selected
		} else if is_pressed {
			bg_color = g.bg_selected
		} else if is_hovered {
			bg_color = g.bg_hover
		} else {
			bg_color = g.bg_normal
		}

		// Draw button background with appropriate corners
		if i == 0 && len(g.options) == 1 {
			// Single button - all corners rounded
			render.fill_rounded_rect(ctx, button_rect, g.corner_radius, bg_color)
		} else if i == 0 {
			// First button - left corners rounded
			render.fill_rounded_rect_corners(ctx, button_rect, g.corner_radius, true, false, false, true, bg_color)
		} else if i == len(g.options) - 1 {
			// Last button - right corners rounded
			render.fill_rounded_rect_corners(ctx, button_rect, g.corner_radius, false, true, true, false, bg_color)
		} else {
			// Middle button - no rounded corners
			render.fill_rect(ctx, button_rect, bg_color)
		}

		// Draw text centered
		if g.font != nil {
			text_width := render.text_measure_logical(g.font, opt)
			text_x := button_x + (button_width - text_width) / 2
			text_y := abs_rect.y + (button_height - render.font_get_logical_line_height(g.font)) / 2
			render.draw_text_top(ctx, g.font, opt, text_x, text_y, g.text_color)
		}

		// Draw separator between buttons (except after last)
		if i < len(g.options) - 1 {
			sep_x := button_x + button_width
			render.draw_vline(ctx, sep_x, abs_rect.y + 2, abs_rect.y + button_height - 2, g.border_color)
		}
	}

	// Draw outer border
	total_width := button_width * i32(len(g.options))
	outer_rect := core.Rect{abs_rect.x, abs_rect.y, total_width, button_height}
	render.draw_rounded_rect(ctx, outer_rect, g.corner_radius, g.border_color)
}

// Handle toggle group events
toggle_group_handle_event :: proc(w: ^Widget, event: ^core.Event) -> bool {
	g := cast(^Toggle_Group)w

	if !w.enabled {
		return false
	}

	#partial switch event.type {
	case .Pointer_Enter:
		g.hovered_index = toggle_group_index_at(g, event.pointer_x, event.pointer_y)
		widget_mark_dirty(g)
		return true

	case .Pointer_Leave:
		g.hovered_index = -1
		g.pressed_index = -1
		widget_mark_dirty(g)
		return true

	case .Pointer_Motion:
		new_hovered := toggle_group_index_at(g, event.pointer_x, event.pointer_y)
		if new_hovered != g.hovered_index {
			g.hovered_index = new_hovered
			widget_mark_dirty(g)
		}
		return true

	case .Pointer_Button_Press:
		if event.button == .Left {
			index := toggle_group_index_at(g, event.pointer_x, event.pointer_y)
			if index >= 0 {
				g.pressed_index = index
				widget_mark_dirty(g)
				return true
			}
		}

	case .Pointer_Button_Release:
		if event.button == .Left && g.pressed_index >= 0 {
			index := toggle_group_index_at(g, event.pointer_x, event.pointer_y)
			if index == g.pressed_index && index != g.selected_index {
				g.selected_index = index
				widget_mark_dirty(g)
				if g.on_change != nil {
					g.on_change(g)
				}
			}
			g.pressed_index = -1
			widget_mark_dirty(g)
			return true
		}

	case .Key_Press:
		if w.focused {
			// Left/Right arrow keys to change selection
			if event.keysym == u32(core.Keysym.Left) && g.selected_index > 0 {
				g.selected_index -= 1
				widget_mark_dirty(g)
				if g.on_change != nil {
					g.on_change(g)
				}
				return true
			}
			if event.keysym == u32(core.Keysym.Right) && g.selected_index < len(g.options) - 1 {
				g.selected_index += 1
				widget_mark_dirty(g)
				if g.on_change != nil {
					g.on_change(g)
				}
				return true
			}
		}
	}

	return false
}

// Toggle group layout (no children)
toggle_group_layout :: proc(w: ^Widget) {
	// No children to layout
}

// Measure toggle group preferred size
toggle_group_measure :: proc(w: ^Widget, available_width: i32) -> core.Size {
	g := cast(^Toggle_Group)w

	button_width := toggle_group_button_width(g)
	total_width := button_width * i32(len(g.options))

	text_height: i32 = 16
	if g.font != nil {
		text_height = render.font_get_logical_line_height(g.font)
	}
	height := text_height + g.button_padding.top + g.button_padding.bottom

	return core.Size{
		width = max(total_width + w.padding.left + w.padding.right, w.min_size.width),
		height = max(height + w.padding.top + w.padding.bottom, w.min_size.height),
	}
}

// Destroy toggle group
toggle_group_destroy :: proc(w: ^Widget) {
	g := cast(^Toggle_Group)w
	delete(g.options)
}
