package filebrowser

import "../../src/app"
import "../../src/core"
import "../../src/render"
import "../../src/widgets"
import "core:path/filepath"

open_preview :: proc() {
	if g_browser.image_grid == nil {
		return
	}

	item, ok := widgets.image_grid_get_selected(g_browser.image_grid)
	if !ok || item.type != .Image || item.image == nil {
		return
	}

	g_browser.preview_active = true
	g_browser.preview_image = item.image
	g_browser.preview_path = item.path
}

close_preview :: proc() {
	g_browser.preview_active = false
	g_browser.preview_image = nil
	g_browser.preview_path = ""
}

navigate_preview :: proc(direction: i32) {
	if g_browser.image_grid == nil || direction == 0 {
		return
	}

	current := g_browser.image_grid.selected_idx
	count := widgets.image_grid_count(g_browser.image_grid)
	if count == 0 {
		return
	}

	new_idx := current
	for {
		new_idx += direction
		if new_idx < 0 || new_idx >= count {
			return
		}

		item := &g_browser.image_grid.items[new_idx]
		if item.type == .Image && item.image != nil {
			widgets.image_grid_set_selected(g_browser.image_grid, new_idx)
			g_browser.preview_image = item.image
			g_browser.preview_path = item.path
			return
		}
	}
}

draw_preview_overlay :: proc(a: ^app.App, ctx: ^render.Draw_Context) {
	if !g_browser.preview_active || g_browser.preview_image == nil {
		return
	}

	overlay_color := core.color_rgba(0, 0, 0, 200)
	render.fill_rect(ctx, core.Rect{0, 0, ctx.logical_width, ctx.logical_height}, overlay_color)

	padding: i32 = 40
	max_w := ctx.logical_width - padding * 2
	max_h := ctx.logical_height - padding * 2

	img_w := g_browser.preview_image.width
	img_h := g_browser.preview_image.height

	scale_w := f32(max_w) / f32(img_w)
	scale_h := f32(max_h) / f32(img_h)
	scale := min(scale_w, scale_h, 1.0)

	display_w := i32(f32(img_w) * scale)
	display_h := i32(f32(img_h) * scale)

	x := (ctx.logical_width - display_w) / 2
	y := (ctx.logical_height - display_h) / 2

	dest_rect := core.Rect{x, y, display_w, display_h}
	render.draw_image_scaled(ctx, g_browser.preview_image, dest_rect)

	font := app.get_font(a)
	if font == nil || len(g_browser.preview_path) == 0 {
		return
	}

	filename := filepath.base(g_browser.preview_path)
	text_color := core.color_hex(0xCCCCCC)
	text_y := ctx.logical_height - 30
	text_width := i32(f64(render.text_measure(font, filename)) / ctx.scale)
	text_x := (ctx.logical_width - text_width) / 2
	render.draw_text_top(ctx, font, filename, text_x, text_y, text_color)
}
