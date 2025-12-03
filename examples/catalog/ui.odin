package catalog

import "../../src/app"
import "../../src/widgets"
import "../../src/core"

// Build the catalog UI
build_ui :: proc(a: ^app.App) {
    root := app.get_root(a)
    root.spacing = 0
    root.padding = widgets.edges_all(0)
    root.align_items = .Stretch

    // Header with theme selector
    header := app.create_container(.Row)
    header.padding = widgets.edges_symmetric(20, 12)
    header.spacing = 10
    header.align_items = .Center
    header.background = widgets.theme_get().bg_tertiary
    widgets.widget_add_child(root, header)

    title := app.create_label(a, "Component Catalog")
    widgets.widget_add_child(header, title)

    // Spacer
    spacer := app.create_container(.Row)
    spacer.flex = 1
    widgets.widget_add_child(header, spacer)

    // Theme dropdown
    theme_dropdown := app.create_dropdown(a)
    widgets.dropdown_add_option(theme_dropdown, "Dark")
    widgets.dropdown_add_option(theme_dropdown, "Light")
    theme_dropdown.min_size = core.Size{100, 0}
    theme_dropdown.on_change = proc(d: ^widgets.Dropdown, index: i32) {
        if index == 0 {
            widgets.theme_set_dark()
        } else {
            widgets.theme_set_light()
        }
    }
    widgets.widget_add_child(header, theme_dropdown)

    // Main content in scroll container
    scroll := app.create_scroll_container(.Vertical)
    scroll.flex = 1
    widgets.widget_add_child(root, scroll)

    content := app.create_container(.Column)
    content.padding = widgets.edges_all(20)
    content.spacing = 24
    content.align_items = .Stretch
    widgets.scroll_container_set_content(scroll, content)

    // Buttons section
    add_section(a, content, "Buttons", proc(a: ^app.App, container: ^widgets.Container) {
        row := app.create_container(.Row)
        row.spacing = 10
        widgets.widget_add_child(container, row)

        btn1 := app.create_button(a, "Normal")
        widgets.widget_add_child(row, btn1)

        btn2 := app.create_button(a, "Disabled")
        btn2.enabled = false
        widgets.widget_add_child(row, btn2)
    })

    // Text Input section
    add_section(a, content, "Text Input", proc(a: ^app.App, container: ^widgets.Container) {
        input := app.create_text_input(a)
        input.placeholder = "Enter some text..."
        input.min_size = core.Size{250, 0}
        widgets.widget_add_child(container, input)
    })

    // Checkbox section
    add_section(a, content, "Checkbox", proc(a: ^app.App, container: ^widgets.Container) {
        row1 := app.create_container(.Row)
        row1.spacing = 8
        row1.align_items = .Center
        widgets.widget_add_child(container, row1)

        cb1 := app.create_checkbox(a)
        cb1.checked = true
        widgets.widget_add_child(row1, cb1)
        lbl1 := app.create_label(a, "Enabled checkbox")
        widgets.widget_add_child(row1, lbl1)

        row2 := app.create_container(.Row)
        row2.spacing = 8
        row2.align_items = .Center
        widgets.widget_add_child(container, row2)

        cb2 := app.create_checkbox(a)
        widgets.widget_add_child(row2, cb2)
        lbl2 := app.create_label(a, "Unchecked checkbox")
        widgets.widget_add_child(row2, lbl2)
    })

    // Labels section
    add_section(a, content, "Labels", proc(a: ^app.App, container: ^widgets.Container) {
        lbl1 := app.create_label(a, "Regular text")
        widgets.widget_add_child(container, lbl1)

        lbl2 := app.create_label(a, "Bold text")
        widgets.label_set_bold(lbl2, true)
        widgets.widget_add_child(container, lbl2)

        lbl3 := app.create_label(a, "Strikethrough text")
        lbl3.strikethrough = true
        widgets.widget_add_child(container, lbl3)

        lbl4 := app.create_label(a, "Secondary color text")
        lbl4.color = widgets.theme_get().text_secondary
        widgets.widget_add_child(container, lbl4)
    })

    // Dropdown section
    add_section(a, content, "Dropdown", proc(a: ^app.App, container: ^widgets.Container) {
        dd := app.create_dropdown(a)
        widgets.dropdown_add_option(dd, "Option 1")
        widgets.dropdown_add_option(dd, "Option 2")
        widgets.dropdown_add_option(dd, "Option 3")
        dd.min_size = core.Size{150, 0}
        widgets.widget_add_child(container, dd)
    })

    // Toggle Group section
    add_section(a, content, "Toggle Group", proc(a: ^app.App, container: ^widgets.Container) {
        options := []string{"Day", "Week", "Month"}
        tg := app.create_toggle_group(a, options)
        widgets.widget_add_child(container, tg)
    })
}

// Helper to create a section with title and content
add_section :: proc(a: ^app.App, parent: ^widgets.Container, title: string, build_content: proc(a: ^app.App, container: ^widgets.Container)) {
    section := app.create_container(.Column)
    section.spacing = 10
    widgets.widget_add_child(parent, section)

    // Section title
    title_label := app.create_label(a, title)
    title_label.color = widgets.theme_get().text_secondary
    widgets.widget_add_child(section, title_label)

    // Divider
    divider := app.create_container(.Row)
    divider.background = widgets.theme_get().divider
    divider.min_size = core.Size{0, 1}
    widgets.widget_add_child(section, divider)

    // Content container
    content := app.create_container(.Column)
    content.spacing = 8
    content.padding = widgets.edges_symmetric(0, 8)
    widgets.widget_add_child(section, content)

    build_content(a, content)
}
