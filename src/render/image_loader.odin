package render

import "base:builtin"
import "core:sync"
import "core:thread"
import "core:strings"
import "core:sys/linux"

// Result of an async image load
Load_Result :: struct {
    path:       string,      // Path that was loaded (owned)
    grid_index: i32,         // Index in grid to update
    image:      ^Image,      // Loaded image (nil if failed)
    thumbnail:  ^Image,      // Loaded thumbnail (nil if failed)
    success:    bool,
}

// Pending load request
Load_Request :: struct {
    path:       string,      // Owned copy
    grid_index: i32,
}

NUM_WORKERS :: 8

// Async image loader with multiple worker threads
Image_Loader :: struct {
    thumb_size:     i32,

    // Pending requests (main writes, workers read)
    pending:        [dynamic]Load_Request,
    pending_lock:   sync.Mutex,

    // Completed results (workers write, main reads)
    completed:      [dynamic]Load_Result,
    completed_lock: sync.Mutex,

    // Track pending loads for polling
    pending_count:  i32,

    // Worker thread control
    active:         bool,
    workers:        [NUM_WORKERS]^thread.Thread,

    // Event notification (for waking main thread)
    notify_fd:      linux.Fd,
}

// Create eventfd (not in Odin stdlib, call syscall directly)
// EFD_NONBLOCK = 0x800, EFD_CLOEXEC = 0x80000
eventfd :: proc "contextless" (initval: u32, flags: u32) -> linux.Fd {
    ret := linux.syscall(linux.SYS_eventfd2, uintptr(initval), uintptr(flags))
    if ret < 0 {
        return linux.Fd(-1)
    }
    return linux.Fd(ret)
}

EFD_NONBLOCK :: 0x800
EFD_CLOEXEC :: 0x80000

// Create image loader
image_loader_create :: proc(thumb_size: i32 = 256) -> ^Image_Loader {
    loader := new(Image_Loader)
    loader.thumb_size = thumb_size
    loader.active = true

    // Create eventfd for notification
    loader.notify_fd = eventfd(0, EFD_NONBLOCK | EFD_CLOEXEC)

    // Start worker threads
    for i in 0..<NUM_WORKERS {
        loader.workers[i] = thread.create_and_start_with_poly_data(loader, worker_proc)
    }

    return loader
}

// Destroy image loader
image_loader_destroy :: proc(loader: ^Image_Loader) {
    if loader == nil {
        return
    }

    // Signal workers to stop
    loader.active = false

    // Wait for all workers to finish
    for i in 0..<NUM_WORKERS {
        if loader.workers[i] != nil {
            thread.join(loader.workers[i])
            thread.destroy(loader.workers[i])
        }
    }

    // Close eventfd
    if loader.notify_fd >= 0 {
        linux.close(loader.notify_fd)
    }

    // Clean up pending requests
    sync.mutex_lock(&loader.pending_lock)
    for &req in loader.pending {
        delete(req.path)
    }
    delete(loader.pending)
    sync.mutex_unlock(&loader.pending_lock)

    // Clean up completed results
    sync.mutex_lock(&loader.completed_lock)
    for &result in loader.completed {
        if len(result.path) > 0 {
            delete(result.path)
        }
    }
    delete(loader.completed)
    sync.mutex_unlock(&loader.completed_lock)

    free(loader)
}

// Queue an image for async loading
image_loader_queue :: proc(loader: ^Image_Loader, path: string, grid_index: i32) {
    if loader == nil || !loader.active {
        return
    }

    // Track pending load
    sync.atomic_add(&loader.pending_count, 1)

    // Add to pending queue
    req := Load_Request{
        path = strings.clone(path),
        grid_index = grid_index,
    }

    sync.mutex_lock(&loader.pending_lock)
    append(&loader.pending, req)
    sync.mutex_unlock(&loader.pending_lock)
}

// Worker thread procedure
worker_proc :: proc(loader: ^Image_Loader) {
    for loader.active {
        // Try to get a task
        req: Load_Request
        has_task := false

        sync.mutex_lock(&loader.pending_lock)
        if len(loader.pending) > 0 {
            req = loader.pending[0]
            ordered_remove(&loader.pending, 0)
            has_task = true
        }
        sync.mutex_unlock(&loader.pending_lock)

        if !has_task {
            // No work, sleep briefly
            thread.yield()
            continue
        }

        // Process the task
        path := req.path
        grid_index := req.grid_index

        // Load the image
        img, ok := image_load(path)
        if !ok {
            // Failed
            result := Load_Result{
                path = path,
                grid_index = grid_index,
                success = false,
            }
            sync.mutex_lock(&loader.completed_lock)
            append(&loader.completed, result)
            sync.mutex_unlock(&loader.completed_lock)
            sync.atomic_sub(&loader.pending_count, 1)
            notify_completion(loader)
            continue
        }

        // Create thumbnail
        thumb, thumb_ok := image_create_thumbnail(&img, loader.thumb_size, loader.thumb_size)
        if !thumb_ok {
            image_destroy(&img)
            result := Load_Result{
                path = path,
                grid_index = grid_index,
                success = false,
            }
            sync.mutex_lock(&loader.completed_lock)
            append(&loader.completed, result)
            sync.mutex_unlock(&loader.completed_lock)
            sync.atomic_sub(&loader.pending_count, 1)
            notify_completion(loader)
            continue
        }

        // Create success result
        img_ptr := new(Image)
        img_ptr^ = img

        thumb_ptr := new(Image)
        thumb_ptr^ = thumb

        result := Load_Result{
            path = path,
            grid_index = grid_index,
            image = img_ptr,
            thumbnail = thumb_ptr,
            success = true,
        }

        sync.mutex_lock(&loader.completed_lock)
        append(&loader.completed, result)
        sync.mutex_unlock(&loader.completed_lock)
        sync.atomic_sub(&loader.pending_count, 1)
        notify_completion(loader)
    }
}

// Get completed loads (call from main thread)
image_loader_get_completed :: proc(loader: ^Image_Loader, allocator := context.allocator) -> []Load_Result {
    if loader == nil {
        return nil
    }

    sync.mutex_lock(&loader.completed_lock)
    defer sync.mutex_unlock(&loader.completed_lock)

    if len(loader.completed) == 0 {
        return nil
    }

    // Move results to caller
    results := make([]Load_Result, len(loader.completed), allocator)
    builtin.copy(results, loader.completed[:])
    builtin.clear(&loader.completed)

    return results
}

// Check if there are pending loads (for progressive loading)
image_loader_has_pending :: proc(loader: ^Image_Loader) -> bool {
    if loader == nil {
        return false
    }
    return sync.atomic_load(&loader.pending_count) > 0
}

// Get the notification FD for polling
image_loader_get_fd :: proc(loader: ^Image_Loader) -> linux.Fd {
    if loader == nil {
        return linux.Fd(-1)
    }
    return loader.notify_fd
}

// Acknowledge notifications (drain the eventfd)
image_loader_acknowledge :: proc(loader: ^Image_Loader) {
    if loader == nil || loader.notify_fd < 0 {
        return
    }
    val: [8]u8  // u64
    linux.read(loader.notify_fd, val[:])
}

// Internal: notify main thread that work is complete
@(private)
notify_completion :: proc(loader: ^Image_Loader) {
    if loader.notify_fd < 0 {
        return
    }
    val: u64 = 1
    buf := transmute([8]u8)val
    linux.write(loader.notify_fd, buf[:])
}

// Get pending load count (for debugging)
image_loader_pending_count :: proc(loader: ^Image_Loader) -> i32 {
    if loader == nil {
        return 0
    }
    return sync.atomic_load(&loader.pending_count)
}

// Clear pending loads (call when navigating away)
image_loader_clear :: proc(loader: ^Image_Loader) {
    if loader == nil {
        return
    }

    // Clear pending queue
    sync.mutex_lock(&loader.pending_lock)
    for &req in loader.pending {
        delete(req.path)
        sync.atomic_sub(&loader.pending_count, 1)
    }
    builtin.clear(&loader.pending)
    sync.mutex_unlock(&loader.pending_lock)

    // Clear completed results
    sync.mutex_lock(&loader.completed_lock)
    for &result in loader.completed {
        if len(result.path) > 0 {
            delete(result.path)
        }
        if result.image != nil {
            image_destroy(result.image)
            free(result.image)
        }
        if result.thumbnail != nil {
            image_destroy(result.thumbnail)
            free(result.thumbnail)
        }
    }
    builtin.clear(&loader.completed)
    sync.mutex_unlock(&loader.completed_lock)
}
