package filebrowser

import "../../src/app"
import "../../src/render"
import "../../src/widgets"

File_Browser :: struct {
	app:               ^app.App,

	// Widget references
	root:              ^widgets.Container,
	header:            ^widgets.Container,
	sidebar:           ^widgets.Container,
	main_area:         ^widgets.Container,
	back_button:       ^widgets.Button,
	header_label:      ^widgets.Label,
	image_grid:        ^widgets.Image_Grid,

	// Navigation
	current_directory: string,
	directory_history: [dynamic]string,

	// Image preview
	preview_active:    bool,
	preview_image:     ^render.Image,
	preview_path:      string,
}

g_browser: File_Browser
