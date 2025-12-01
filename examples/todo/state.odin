package todo

// Pure data layer - no widget imports

Todo_Item :: struct {
    id:        int,
    text:      string,
    completed: bool,
}

Filter :: enum {
    All,
    Active,
    Completed,
}

Todo_State :: struct {
    todos:          [dynamic]Todo_Item,
    next_id:        int,
    current_filter: Filter,
}

g_state: Todo_State

// Initialize empty state
init_state :: proc() {
    g_state.todos = make([dynamic]Todo_Item)
    g_state.next_id = 1
    g_state.current_filter = .All
}

// Cleanup state
destroy_state :: proc() {
    for &todo in g_state.todos {
        delete(todo.text)
    }
    delete(g_state.todos)
}
