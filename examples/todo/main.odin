package todo

import "../../src/app"
import "../../src/widgets"
import "core:fmt"

main :: proc() {
    fmt.println("Odek - Todo App")

    a := app.create("Todo App", 400, 500)
    if a == nil {
        fmt.eprintln("Failed to create application")
        return
    }
    defer app.destroy(a)

    // // Set minimum window size to prevent layout overflow
    // app.set_min_size(a, 400, 250)

    // Enable debug borders to visualize widget bounds
    widgets.debug_borders_set(true)

    init_state()
    load_todos()
    build_ui(a)
    setup_bindings()
    sync_ui()

    fmt.println("Todo app started")
    app.run(a)
    fmt.println("Application closed")
}
