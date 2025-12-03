package filebrowser

import "../../src/app"
import "core:fmt"

main :: proc() {
	fmt.println("Odek - Image Browser Demo")

	a := app.create("Odek Image Browser", 800, 600, "com.odek.filebrowser")
	if a == nil {
		fmt.eprintln("Failed to create application")
		return
	}
	defer app.destroy(a)

	g_browser.app = a

	build_ui(a)

	a.on_draw_overlay = draw_preview_overlay
	a.on_key = handle_key

	navigate_to_directory("/home/chris", add_to_history = false)

	fmt.println("Window created, entering event loop...")
	app.run(a)
	fmt.println("Application closed")
}
