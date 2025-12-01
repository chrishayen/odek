package todo

import "../../src/app"
import "../../src/core"
import "../../src/widgets"

// UI layer - builds widget hierarchy (no business logic)

// Build the entire UI tree
build_ui :: proc(a: ^app.App) {
    g_ui.app = a

    root := app.get_root(a)
    g_ui.root = root

    root.padding = widgets.edges_all(15)
    root.background = core.color_hex(0x1a1a2e)
    root.spacing = 15
    root.align_items = .Stretch

    build_header(a, root)
    build_input_row(a, root)
    build_todo_list(a, root)
    build_footer(a, root)

    // Layout is handled in the draw callback where root.rect has valid window dimensions
}

// Build header with title
build_header :: proc(a: ^app.App, root: ^widgets.Container) {
    header := app.create_container(.Row)
    header.min_size = core.Size{0, 40}
    header.align_items = .Center
    header.justify_content = .Center
    widgets.widget_add_child(root, header)

    title := app.create_label(a, "Todo App")
    title.color = core.COLOR_WHITE
    widgets.widget_add_child(header, title)
}

// Build input row with text input and add button
build_input_row :: proc(a: ^app.App, root: ^widgets.Container) {
    input_row := app.create_container(.Row)
    input_row.spacing = 10
    input_row.align_items = .Center
    widgets.widget_add_child(root, input_row)

    g_ui.input = app.create_text_input(a)
    g_ui.input.flex = 1
    g_ui.input.min_size = core.Size{0, 36}
    widgets.text_input_set_placeholder(g_ui.input, "What needs to be done?")
    widgets.widget_add_child(input_row, g_ui.input)

    add_btn := app.create_button(a, "+")
    add_btn.min_size = core.Size{40, 36}
    g_ui.add_button = add_btn
    widgets.widget_add_child(input_row, add_btn)
}

// Build scrollable todo list container
build_todo_list :: proc(a: ^app.App, root: ^widgets.Container) {
    scroll := app.create_scroll_container(.Vertical)
    scroll.flex = 1
    scroll.background = core.color_hex(0x16213e)
    widgets.widget_add_child(root, scroll)

    g_ui.todo_list = app.create_container(.Column)
    g_ui.todo_list.spacing = 5
    g_ui.todo_list.padding = widgets.edges_all(10)
    g_ui.todo_list.align_items = .Stretch
    widgets.scroll_container_set_content(scroll, g_ui.todo_list)
}

// Build footer with filter buttons and count
build_footer :: proc(a: ^app.App, root: ^widgets.Container) {
    footer := app.create_container(.Column)
    footer.spacing = 10
    footer.align_items = .Stretch
    widgets.widget_add_child(root, footer)

    // Filter buttons row
    filter_row := app.create_container(.Row)
    filter_row.spacing = 8
    filter_row.align_items = .Center
    filter_row.justify_content = .Center
    widgets.widget_add_child(footer, filter_row)

    // Create filter buttons
    filter_names := [Filter]string{
        .All = "All",
        .Active = "Active",
        .Completed = "Completed",
    }

    for f in Filter {
        btn := app.create_button(a, filter_names[f])
        btn.min_size = core.Size{70, 28}
        btn.padding = widgets.edges_symmetric(10, 5)
        g_ui.filter_buttons[f] = btn
        widgets.widget_add_child(filter_row, btn)
    }

    // Clear completed button
    g_ui.clear_btn = app.create_button(a, "Clear Done")
    g_ui.clear_btn.min_size = core.Size{80, 28}
    g_ui.clear_btn.padding = widgets.edges_symmetric(10, 5)
    widgets.button_set_colors(
        g_ui.clear_btn,
        core.color_hex(0x404040),
        core.color_hex(0x505050),
        core.color_hex(0x303030),
    )
    widgets.widget_add_child(filter_row, g_ui.clear_btn)

    // Count label row
    count_row := app.create_container(.Row)
    count_row.align_items = .Center
    count_row.justify_content = .Center
    widgets.widget_add_child(footer, count_row)

    g_ui.count_label = app.create_label(a, "0 items left")
    g_ui.count_label.color = core.color_hex(0x888888)
    widgets.widget_add_child(count_row, g_ui.count_label)
}
