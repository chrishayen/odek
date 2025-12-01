package filebrowser

import "../../src/app"
import "../../src/widgets"

handle_key :: proc(a: ^app.App, key: u32, pressed: bool, utf8: string) -> bool {
	if !pressed {
		return false
	}

	switch key {
	case 1: // Escape
		if g_browser.preview_active {
			close_preview()
			return true
		}
	case 57: // Spacebar
		if g_browser.preview_active {
			close_preview()
		} else {
			open_preview()
		}
		return true
	case 14: // Backspace
		navigate_back()
		return true
	case 28: // Enter
		return handle_enter()
	case 103, 105, 106, 108: // Arrow keys
		return handle_arrow(key)
	}

	return false
}

handle_enter :: proc() -> bool {
	if g_browser.image_grid == nil {
		return false
	}

	item, ok := widgets.image_grid_get_selected(g_browser.image_grid)
	if !ok {
		return false
	}

	switch item.type {
	case .Folder:
		navigate_to_directory(item.path)
	case .Image:
		open_preview()
	case .File, .Video:
		// No action
	}
	return true
}

handle_arrow :: proc(key: u32) -> bool {
	if g_browser.preview_active {
		direction: i32 = key == 105 ? -1 : (key == 106 ? 1 : 0)
		navigate_preview(direction)
		return true
	}

	navigate_grid(key)
	return true
}

navigate_grid :: proc(key: u32) {
	if g_browser.image_grid == nil {
		return
	}

	cols := widgets.image_grid_get_columns(g_browser.image_grid)
	count := widgets.image_grid_count(g_browser.image_grid)
	current := g_browser.image_grid.selected_idx

	if current < 0 && count > 0 {
		widgets.image_grid_set_selected(g_browser.image_grid, 0)
		return
	}

	new_idx := current
	switch key {
	case 105: // Left
		if current > 0 {
			new_idx = current - 1
		}
	case 106: // Right
		if current < count - 1 {
			new_idx = current + 1
		}
	case 103: // Up
		if current >= cols {
			new_idx = current - cols
		}
	case 108: // Down
		if current + cols < count {
			new_idx = current + cols
		}
	}

	if new_idx != current && new_idx >= 0 && new_idx < count {
		widgets.image_grid_set_selected(g_browser.image_grid, new_idx)
		widgets.image_grid_ensure_visible(g_browser.image_grid, new_idx)
	}
}
