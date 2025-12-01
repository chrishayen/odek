package todo

import "../../src/app"
import "../../src/core"
import "../../src/widgets"
import "core:fmt"

// Bindings layer - connects UI to business logic

// Widget references
Todo_UI :: struct {
    app:            ^app.App,
    root:           ^widgets.Container,
    input:          ^widgets.Text_Input,
    add_button:     ^widgets.Button,
    todo_list:      ^widgets.Container,
    filter_buttons: [Filter]^widgets.Button,
    clear_btn:      ^widgets.Button,
    count_label:    ^widgets.Label,
}

g_ui: Todo_UI

// Wire up all callbacks
setup_bindings :: proc() {
    // Add button
    g_ui.add_button.on_click = on_add_click

    // Text input submit
    g_ui.input.on_submit = on_input_submit

    // Filter buttons
    for f in Filter {
        btn := g_ui.filter_buttons[f]
        btn.user_data = rawptr(uintptr(f))
        btn.on_click = on_filter_click
    }

    // Clear completed
    g_ui.clear_btn.on_click = on_clear_click
}

// Callback: Add button clicked
on_add_click :: proc(btn: ^widgets.Button) {
    text := widgets.text_input_get_text(g_ui.input)
    if len(text) == 0 {
        return
    }

    add_todo(text)
    widgets.text_input_set_text(g_ui.input, "")
    save_todos()
    sync_ui()
}

// Callback: Enter pressed in input
on_input_submit :: proc(ti: ^widgets.Text_Input) {
    text := widgets.text_input_get_text(ti)
    if len(text) == 0 {
        return
    }

    add_todo(text)
    widgets.text_input_set_text(ti, "")
    save_todos()
    sync_ui()
}

// Callback: Filter button clicked
on_filter_click :: proc(btn: ^widgets.Button) {
    f := Filter(uintptr(btn.user_data))
    set_filter(f)
    sync_ui()
}

// Callback: Clear completed clicked
on_clear_click :: proc(btn: ^widgets.Button) {
    clear_completed()
    save_todos()
    sync_ui()
}

// Callback: Todo checkbox toggled
on_todo_toggle :: proc(cb: ^widgets.Checkbox) {
    id := int(uintptr(cb.user_data))
    toggle_todo(id)
    save_todos()
    sync_ui()
}

// Callback: Delete button clicked
on_delete_click :: proc(btn: ^widgets.Button) {
    id := int(uintptr(btn.user_data))
    delete_todo(id)
    save_todos()
    sync_ui()
}

// Sync UI from state - rebuild todo list and update labels
sync_ui :: proc() {
    // Clear existing todo list children
    clear_todo_list_children()

    // Get visible todos based on filter
    visible := get_visible_todos()

    // Rebuild todo rows
    for todo in visible {
        create_todo_row(todo)
    }

    // Update count label
    active_count := get_active_count()
    count_text := fmt.tprintf("%d item%s left", active_count, active_count == 1 ? "" : "s")
    widgets.label_set_text(g_ui.count_label, count_text)

    // Update filter button styles
    update_filter_button_styles()

    // Request redraw - layout will be done in draw callback with correct window size
    app.request_redraw(g_ui.app)
}

// Clear children from todo list container
clear_todo_list_children :: proc() {
    // Remove all children in reverse order
    for len(g_ui.todo_list.children) > 0 {
        child := g_ui.todo_list.children[len(g_ui.todo_list.children) - 1]
        widgets.widget_remove_child(g_ui.todo_list, child)
        widgets.widget_destroy(child)
    }
}

// Create a row widget for a todo item
create_todo_row :: proc(todo: Todo_Item) {
    row := app.create_container(.Row)
    row.min_size = core.Size{0, 36}
    row.spacing = 10
    row.align_items = .Center
    row.background = core.color_hex(0x1f2b4d)
    row.padding = widgets.edges_symmetric(10, 5)
    widgets.widget_add_child(g_ui.todo_list, row)

    // Checkbox
    cb := app.create_checkbox(g_ui.app)
    widgets.checkbox_set_checked(cb, todo.completed)
    cb.user_data = rawptr(uintptr(todo.id))
    cb.on_change = on_todo_toggle
    widgets.widget_add_child(row, cb)

    // Label with strikethrough if completed
    label := app.create_label(g_ui.app, todo.text)
    label.flex = 1
    label.wrap = false
    if todo.completed {
        label.color = core.color_hex(0x666666)
        widgets.label_set_strikethrough(label, true)
    } else {
        label.color = core.COLOR_WHITE
    }
    widgets.widget_add_child(row, label)

    // Delete button
    del_btn := app.create_button(g_ui.app, "X")
    del_btn.min_size = core.Size{28, 28}
    del_btn.padding = widgets.edges_all(5)
    del_btn.user_data = rawptr(uintptr(todo.id))
    del_btn.on_click = on_delete_click
    widgets.button_set_colors(
        del_btn,
        core.color_hex(0x8b0000),
        core.color_hex(0xa00000),
        core.color_hex(0x700000),
    )
    widgets.widget_add_child(row, del_btn)
}

// Update filter button appearance based on current filter
update_filter_button_styles :: proc() {
    theme := widgets.theme_get()

    for f in Filter {
        btn := g_ui.filter_buttons[f]
        if f == g_state.current_filter {
            // Active filter - accent colors
            widgets.button_set_colors(btn, theme.accent, theme.accent_hover, theme.accent_pressed)
        } else {
            // Inactive filter - subdued colors
            widgets.button_set_colors(
                btn,
                core.color_hex(0x404040),
                core.color_hex(0x505050),
                core.color_hex(0x303030),
            )
        }
    }
}
