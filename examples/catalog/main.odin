package catalog

import "../../src/app"
import "../../src/widgets"

main :: proc() {
    a := app.create("Odek - Component Catalog", 600, 700)
    if a == nil {
        return
    }
    defer app.destroy(a)

    build_ui(a)

    app.run(a)
}
