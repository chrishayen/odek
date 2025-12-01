package filebrowser

import "../../src/app"
import "../../src/core"
import "../../src/widgets"
import "core:fmt"

build_ui :: proc(a: ^app.App) {
	root := app.get_root(a)
	g_browser.root = root

	root.padding = widgets.edges_all(10)
	root.background = core.color_hex(0x2D2D2D)
	root.spacing = 10
	root.align_items = .Stretch

	build_header(a, root)
	build_content(a, root)

	widgets.widget_layout(root)
}

build_header :: proc(a: ^app.App, root: ^widgets.Container) {
	g_browser.header = app.create_container(.Row)
	g_browser.header.min_size = core.Size{0, 50}
	g_browser.header.background = core.color_hex(0x404040)
	g_browser.header.padding = widgets.edges_symmetric(15, 10)
	g_browser.header.spacing = 10
	g_browser.header.align_items = .Center
	widgets.widget_add_child(root, g_browser.header)

	if app.get_font(a) == nil {
		return
	}

	g_browser.back_button = app.create_button(a, "<")
	g_browser.back_button.min_size = core.Size{36, 30}
	g_browser.back_button.on_click = back_button_callback
	widgets.button_set_colors(
		g_browser.back_button,
		core.color_hex(0x505050),
		core.color_hex(0x606060),
		core.color_hex(0x404040),
	)
	widgets.widget_add_child(g_browser.header, g_browser.back_button)

	g_browser.header_label = app.create_label(a, "Odek Image Browser")
	widgets.label_set_color(g_browser.header_label, core.COLOR_WHITE)
	widgets.widget_add_child(g_browser.header, g_browser.header_label)
}

build_content :: proc(a: ^app.App, root: ^widgets.Container) {
	content := app.create_container(.Row)
	content.flex = 1
	content.spacing = 10
	content.align_items = .Stretch
	widgets.widget_add_child(root, content)

	build_sidebar(a, content)
	build_main_area(a, content)
}

build_sidebar :: proc(a: ^app.App, parent: ^widgets.Container) {
	g_browser.sidebar = app.create_container(.Column)
	g_browser.sidebar.min_size = core.Size{150, 0}
	g_browser.sidebar.background = core.color_hex(0x383838)
	g_browser.sidebar.padding = widgets.edges_all(10)
	g_browser.sidebar.spacing = 8
	g_browser.sidebar.align_items = .Stretch
	widgets.widget_add_child(parent, g_browser.sidebar)

	if app.get_font(a) == nil {
		return
	}

	Bookmark :: struct {
		name: string,
		path: cstring,
	}
	bookmarks := []Bookmark {
		{"Home", "/home/chris"},
		{"Pictures", "/home/chris/Pictures"},
		{"Downloads", "/home/chris/Downloads"},
		{"Code", "/home/chris/Code"},
	}

	for bm in bookmarks {
		btn := app.create_button(a, bm.name)
		btn.user_data = rawptr(bm.path)
		btn.on_click = bookmark_callback
		btn.min_size = core.Size{0, 32}
		widgets.button_set_colors(
			btn,
			core.color_hex(0x444444),
			core.color_hex(0x555555),
			core.color_hex(0x333333),
		)
		widgets.widget_add_child(g_browser.sidebar, btn)
	}
}

build_main_area :: proc(a: ^app.App, parent: ^widgets.Container) {
	g_browser.main_area = app.create_container(.Column)
	g_browser.main_area.flex = 1
	g_browser.main_area.background = core.color_hex(0x333333)
	g_browser.main_area.padding = widgets.edges_all(20)
	g_browser.main_area.spacing = 15
	g_browser.main_area.align_items = .Stretch
	widgets.widget_add_child(parent, g_browser.main_area)

	g_browser.image_grid = app.create_image_grid(a)
	g_browser.image_grid.flex = 1
	g_browser.image_grid.cell_width = 150
	g_browser.image_grid.cell_height = 150
	g_browser.image_grid.spacing = 10
	g_browser.image_grid.padding = widgets.edges_all(10)
	g_browser.image_grid.on_click = image_grid_click_callback
	g_browser.image_grid.on_folder_click = folder_click_callback
	widgets.widget_add_child(g_browser.main_area, g_browser.image_grid)
}

back_button_callback :: proc(button: ^widgets.Button) {
	navigate_back()
}

bookmark_callback :: proc(button: ^widgets.Button) {
	path := cast(cstring)button.user_data
	if path != nil {
		navigate_to_directory(string(path))
	}
}

folder_click_callback :: proc(grid: ^widgets.Image_Grid, path: string) {
	fmt.printf("Navigating to folder: %s\n", path)
	navigate_to_directory(path)
}

image_grid_click_callback :: proc(grid: ^widgets.Image_Grid, index: i32, item: ^widgets.Grid_Item) {
	fmt.printf("Image clicked: index=%d, path=%s\n", index, item.path)
}
