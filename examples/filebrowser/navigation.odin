package filebrowser

import "../../src/app"
import "../../src/widgets"
import "core:fmt"
import "core:os"
import "core:path/filepath"
import "core:strings"

navigate_to_directory :: proc(dir_path: string, add_to_history: bool = true) {
	if add_to_history && len(g_browser.current_directory) > 0 {
		append(&g_browser.directory_history, strings.clone(g_browser.current_directory))
	}

	if len(g_browser.current_directory) > 0 {
		delete(g_browser.current_directory)
	}
	g_browser.current_directory = strings.clone(dir_path)

	widgets.image_grid_clear(g_browser.image_grid)
	app.clear_image_loads(g_browser.app)

	handle, err := os.open(dir_path)
	if err != nil {
		fmt.eprintln("Failed to open directory:", dir_path)
		app.request_redraw(g_browser.app)
		return
	}
	defer os.close(handle)

	entries, read_err := os.read_dir(handle, -1)
	if read_err != nil {
		fmt.eprintln("Failed to read directory")
		app.request_redraw(g_browser.app)
		return
	}
	defer delete(entries)

	counts := load_entries(dir_path, entries)

	fmt.printf(
		"Loaded %d folders, %d images, %d videos, %d files from %s\n",
		counts.folders,
		counts.images,
		counts.videos,
		counts.files,
		dir_path,
	)

	if g_browser.header_label != nil {
		widgets.label_set_text(g_browser.header_label, g_browser.current_directory)
	}

	app.request_redraw(g_browser.app)
}

Entry_Counts :: struct {
	folders, images, videos, files: int,
}

load_entries :: proc(dir_path: string, entries: []os.File_Info) -> Entry_Counts {
	counts: Entry_Counts

	// First pass: folders
	for entry in entries {
		if !entry.is_dir || strings.has_prefix(entry.name, ".") {
			continue
		}

		full_path := filepath.join({dir_path, entry.name})
		name_clone := strings.clone(entry.name)
		path_clone := strings.clone(full_path)
		delete(full_path)

		widgets.image_grid_add_folder(g_browser.image_grid, name_clone, path_clone)
		counts.folders += 1
	}

	// Second pass: files
	for entry in entries {
		if entry.is_dir || strings.has_prefix(entry.name, ".") {
			continue
		}

		full_path := filepath.join({dir_path, entry.name})
		name_clone := strings.clone(entry.name)
		path_clone := strings.clone(full_path)

		ext := strings.to_lower(filepath.ext(entry.name))
		defer delete(ext)

		switch ext {
		case ".png", ".jpg", ".jpeg":
			app.queue_image_load(g_browser.app, g_browser.image_grid, path_clone, name_clone)
			counts.images += 1
		case ".mp4", ".mkv", ".avi", ".webm", ".mov":
			widgets.image_grid_add_video(g_browser.image_grid, name_clone, path_clone)
			counts.videos += 1
		case:
			widgets.image_grid_add_file(g_browser.image_grid, name_clone, path_clone)
			counts.files += 1
		}
		delete(full_path)
	}

	return counts
}

navigate_back :: proc() {
	if len(g_browser.directory_history) == 0 {
		return
	}
	prev := pop(&g_browser.directory_history)
	navigate_to_directory(prev, add_to_history = false)
	delete(prev)
}
