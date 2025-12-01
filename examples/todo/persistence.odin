package todo

import "core:encoding/json"
import "core:os"
import "core:strings"
import "core:fmt"

// File I/O layer - no widget imports

// JSON-serializable structure
Todo_Data :: struct {
    id:        int,
    text:      string,
    completed: bool,
}

Save_Data :: struct {
    todos:   []Todo_Data,
    next_id: int,
}

// Get the path to save file
get_save_path :: proc() -> string {
    home := os.get_env("HOME")
    if home == "" {
        return "todos.json"
    }

    // Create data directory if needed
    data_dir := strings.concatenate({home, "/.local/share/odek-todo"})
    defer delete(data_dir)

    os.make_directory(data_dir)

    return strings.concatenate({home, "/.local/share/odek-todo/todos.json"})
}

// Save todos to JSON file
save_todos :: proc() {
    path := get_save_path()
    defer delete(path)

    // Build save data
    todos_data := make([]Todo_Data, len(g_state.todos))
    defer delete(todos_data)

    for todo, i in g_state.todos {
        todos_data[i] = Todo_Data{
            id = todo.id,
            text = todo.text,
            completed = todo.completed,
        }
    }

    save_data := Save_Data{
        todos = todos_data,
        next_id = g_state.next_id,
    }

    // Marshal to JSON
    json_bytes, err := json.marshal(save_data, {pretty = true})
    if err != nil {
        fmt.eprintln("Failed to marshal todos:", err)
        return
    }
    defer delete(json_bytes)

    // Write to file
    ok := os.write_entire_file(path, json_bytes)
    if !ok {
        fmt.eprintln("Failed to write todos file")
    }
}

// Load todos from JSON file
load_todos :: proc() {
    path := get_save_path()
    defer delete(path)

    // Read file
    data, ok := os.read_entire_file(path)
    if !ok {
        // File doesn't exist yet - that's fine
        return
    }
    defer delete(data)

    // Parse JSON
    save_data: Save_Data
    err := json.unmarshal(data, &save_data)
    if err != nil {
        fmt.eprintln("Failed to parse todos file:", err)
        return
    }
    defer {
        // Clean up the parsed data's allocated strings
        for &todo in save_data.todos {
            delete(todo.text)
        }
        delete(save_data.todos)
    }

    // Populate state
    clear(&g_state.todos)
    g_state.next_id = save_data.next_id

    for todo in save_data.todos {
        append(&g_state.todos, Todo_Item{
            id = todo.id,
            text = strings.clone(todo.text),
            completed = todo.completed,
        })
    }
}
