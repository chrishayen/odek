package todo

import "../../src/app"
import "core:fmt"

main :: proc() {
    fmt.println("Odek - Todo App")

    a := app.create("Todo App", 400, 500, "com.odek.todo")
    if a == nil {
        fmt.eprintln("Failed to create application")
        return
    }
    defer app.destroy(a)

    // // Set minimum window size to prevent layout overflow
    // app.set_min_size(a, 400, 250)

    init_state()
    load_todos()
    build_ui(a)
    setup_bindings()
    sync_ui()

    fmt.println("Todo app started")
    app.run(a)
    fmt.println("Application closed")
}
