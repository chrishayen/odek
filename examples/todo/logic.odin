package todo

import "core:strings"

// Business logic layer - no widget imports

// Add a new todo item
add_todo :: proc(text: string) {
    if len(text) == 0 {
        return
    }

    todo := Todo_Item{
        id = g_state.next_id,
        text = strings.clone(text),
        completed = false,
    }
    append(&g_state.todos, todo)
    g_state.next_id += 1
}

// Toggle a todo's completed state
toggle_todo :: proc(id: int) {
    for &todo in g_state.todos {
        if todo.id == id {
            todo.completed = !todo.completed
            return
        }
    }
}

// Delete a todo by id
delete_todo :: proc(id: int) {
    for i := 0; i < len(g_state.todos); i += 1 {
        if g_state.todos[i].id == id {
            delete(g_state.todos[i].text)
            ordered_remove(&g_state.todos, i)
            return
        }
    }
}

// Clear all completed todos
clear_completed :: proc() {
    i := 0
    for i < len(g_state.todos) {
        if g_state.todos[i].completed {
            delete(g_state.todos[i].text)
            ordered_remove(&g_state.todos, i)
        } else {
            i += 1
        }
    }
}

// Set the current filter
set_filter :: proc(f: Filter) {
    g_state.current_filter = f
}

// Get visible todos based on current filter
get_visible_todos :: proc() -> []Todo_Item {
    #partial switch g_state.current_filter {
    case .Active:
        // Count active todos
        count := 0
        for todo in g_state.todos {
            if !todo.completed {
                count += 1
            }
        }
        if count == 0 {
            return {}
        }

        // Build result slice - using static buffer for simplicity
        @(static) result: [256]Todo_Item
        idx := 0
        for todo in g_state.todos {
            if !todo.completed && idx < len(result) {
                result[idx] = todo
                idx += 1
            }
        }
        return result[:idx]

    case .Completed:
        count := 0
        for todo in g_state.todos {
            if todo.completed {
                count += 1
            }
        }
        if count == 0 {
            return {}
        }

        @(static) result: [256]Todo_Item
        idx := 0
        for todo in g_state.todos {
            if todo.completed && idx < len(result) {
                result[idx] = todo
                idx += 1
            }
        }
        return result[:idx]
    }

    // All - return all todos
    return g_state.todos[:]
}

// Count active (incomplete) todos
get_active_count :: proc() -> int {
    count := 0
    for todo in g_state.todos {
        if !todo.completed {
            count += 1
        }
    }
    return count
}

// Count completed todos
get_completed_count :: proc() -> int {
    count := 0
    for todo in g_state.todos {
        if todo.completed {
            count += 1
        }
    }
    return count
}
