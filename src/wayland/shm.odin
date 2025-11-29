package wayland

import "core:c"
import "core:c/libc"
import "core:mem"
import "core:sys/posix"

// Linux-specific syscalls for shared memory
foreign import libc_shm "system:c"

@(default_calling_convention = "c")
foreign libc_shm {
    memfd_create :: proc(name: cstring, flags: c.uint) -> c.int ---
    ftruncate :: proc(fd: c.int, length: c.long) -> c.int ---
}

MFD_CLOEXEC :: 0x0001
MFD_ALLOW_SEALING :: 0x0002

// SHM pool wrapper
Shm_Pool :: struct {
    wl_pool: ^Wl_Shm_Pool,
    fd: c.int,
    data: [^]u8,
    size: int,
    capacity: int,
}

// Double-buffered surface buffer
Buffer :: struct {
    wl_buffer: ^Wl_Buffer,
    data: [^]u32,
    width: i32,
    height: i32,
    stride: i32,
    busy: bool,
    listener: Wl_Buffer_Listener,
}

// Create a shared memory pool
shm_pool_create :: proc(shm: ^Wl_Shm, size: int) -> (^Shm_Pool, bool) {
    // Create anonymous file for shared memory
    fd := memfd_create("odek-shm", MFD_CLOEXEC)
    if fd < 0 {
        return nil, false
    }

    // Set the size
    if ftruncate(fd, c.long(size)) < 0 {
        posix.close(posix.FD(fd))
        return nil, false
    }

    // Map the memory
    data := posix.mmap(nil, uint(size), {.READ, .WRITE}, {.SHARED}, posix.FD(fd), 0)
    if data == posix.MAP_FAILED {
        posix.close(posix.FD(fd))
        return nil, false
    }

    // Create wl_shm_pool
    wl_pool := shm_create_pool(shm, fd, i32(size))
    if wl_pool == nil {
        posix.munmap(data, uint(size))
        posix.close(posix.FD(fd))
        return nil, false
    }

    pool := new(Shm_Pool)
    pool.wl_pool = wl_pool
    pool.fd = fd
    pool.data = cast([^]u8)data
    pool.size = size
    pool.capacity = size

    return pool, true
}

// Destroy a shared memory pool
shm_pool_destroy :: proc(pool: ^Shm_Pool) {
    if pool == nil {
        return
    }

    shm_pool_wl_destroy(pool.wl_pool)
    posix.munmap(pool.data, uint(pool.capacity))
    posix.close(posix.FD(pool.fd))
    free(pool)
}

// Create a buffer from the pool
buffer_create :: proc(pool: ^Shm_Pool, width, height: i32, format: Wl_Shm_Format) -> (^Buffer, bool) {
    stride := width * 4 // 4 bytes per pixel for ARGB8888
    size := stride * height

    if int(size) > pool.capacity {
        return nil, false
    }

    wl_buffer := shm_pool_create_buffer(pool.wl_pool, 0, width, height, stride, format)
    if wl_buffer == nil {
        return nil, false
    }

    buf := new(Buffer)
    buf.wl_buffer = wl_buffer
    buf.data = cast([^]u32)pool.data
    buf.width = width
    buf.height = height
    buf.stride = stride
    buf.busy = false

    // Set up release listener
    buf.listener = Wl_Buffer_Listener{
        release = buffer_release_callback,
    }
    buffer_add_listener(wl_buffer, &buf.listener, buf)

    return buf, true
}

// Buffer release callback - compositor is done with the buffer
buffer_release_callback :: proc "c" (data: rawptr, wl_buffer: ^Wl_Buffer) {
    buf := cast(^Buffer)data
    buf.busy = false
}

// Destroy a buffer
buffer_destroy_internal :: proc(buf: ^Buffer) {
    if buf == nil {
        return
    }
    buffer_destroy(buf.wl_buffer)
    free(buf)
}

// Low-level SHM operations
shm_create_pool :: proc(shm: ^Wl_Shm, fd: c.int, size: i32) -> ^Wl_Shm_Pool {
    return cast(^Wl_Shm_Pool)wl_proxy_marshal_flags(
        shm, 0, &wl_shm_pool_interface, wl_proxy_get_version(shm), 0, nil, fd, size)
}

shm_pool_create_buffer :: proc(
    pool: ^Wl_Shm_Pool,
    offset: i32,
    width, height: i32,
    stride: i32,
    format: Wl_Shm_Format,
) -> ^Wl_Buffer {
    return cast(^Wl_Buffer)wl_proxy_marshal_flags(
        pool, 0, &wl_buffer_interface, wl_proxy_get_version(pool), 0,
        nil, offset, width, height, stride, u32(format))
}

shm_pool_wl_destroy :: proc(pool: ^Wl_Shm_Pool) {
    // WL_MARSHAL_FLAG_DESTROY already destroys the proxy
    wl_proxy_marshal_flags(pool, 1, nil, wl_proxy_get_version(pool), WL_MARSHAL_FLAG_DESTROY)
}

shm_add_listener :: proc(shm: ^Wl_Shm, listener: ^Wl_Shm_Listener, data: rawptr) -> c.int {
    return wl_proxy_add_listener(shm, listener, data)
}
