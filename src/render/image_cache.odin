package render

import "base:builtin"
import "core:time"

// Image cache entry
Image_Cache_Entry :: struct {
    image:      Image,          // Full image
    thumbnail:  Image,          // Pre-generated thumbnail
    path:       string,         // File path (owned copy)
    last_used:  time.Time,      // For LRU eviction
}

// LRU image cache
Image_Cache :: struct {
    entries:        map[string]^Image_Cache_Entry,
    max_entries:    int,        // Maximum cached images
    thumb_size:     i32,        // Thumbnail max dimension
}

// Create image cache
image_cache_create :: proc(max_entries: int = 100, thumb_size: i32 = 256) -> ^Image_Cache {
    cache := new(Image_Cache)
    cache.entries = make(map[string]^Image_Cache_Entry)
    cache.max_entries = max_entries
    cache.thumb_size = thumb_size
    return cache
}

// Destroy cache and free all images
image_cache_destroy :: proc(cache: ^Image_Cache) {
    if cache == nil {
        return
    }

    for _, entry in cache.entries {
        image_destroy(&entry.image)
        image_destroy(&entry.thumbnail)
        delete(entry.path)
        free(entry)
    }

    delete(cache.entries)
    free(cache)
}

// Get cached image and thumbnail by path
// Returns (full_image, thumbnail, found)
image_cache_get :: proc(cache: ^Image_Cache, path: string) -> (^Image, ^Image, bool) {
    if cache == nil {
        return nil, nil, false
    }

    entry, ok := cache.entries[path]
    if !ok {
        return nil, nil, false
    }

    // Update last used time
    entry.last_used = time.now()

    return &entry.image, &entry.thumbnail, true
}

// Get just the thumbnail (loads if needed)
image_cache_get_thumbnail :: proc(cache: ^Image_Cache, path: string) -> (^Image, bool) {
    _, thumb, ok := image_cache_load(cache, path)
    return thumb, ok
}

// Load image into cache if not already cached
// Returns (full_image, thumbnail, success)
image_cache_load :: proc(cache: ^Image_Cache, path: string) -> (^Image, ^Image, bool) {
    if cache == nil {
        return nil, nil, false
    }

    // Check if already cached
    if entry, ok := cache.entries[path]; ok {
        entry.last_used = time.now()
        return &entry.image, &entry.thumbnail, true
    }

    // Evict if at capacity
    if len(cache.entries) >= cache.max_entries {
        image_cache_evict_lru(cache)
    }

    // Load image from disk
    img, ok := image_load(path)
    if !ok {
        return nil, nil, false
    }

    // Create thumbnail
    thumb, thumb_ok := image_create_thumbnail(&img, cache.thumb_size, cache.thumb_size)
    if !thumb_ok {
        image_destroy(&img)
        return nil, nil, false
    }

    // Create entry
    entry := new(Image_Cache_Entry)
    entry.image = img
    entry.thumbnail = thumb
    entry.path = clone_string(path)
    entry.last_used = time.now()

    cache.entries[entry.path] = entry

    return &entry.image, &entry.thumbnail, true
}

// Evict least recently used entry
image_cache_evict_lru :: proc(cache: ^Image_Cache) {
    if cache == nil || len(cache.entries) == 0 {
        return
    }

    // Find oldest entry
    oldest_path: string
    oldest_time := time.Time{}
    first := true

    for path, entry in cache.entries {
        if first || time.diff(entry.last_used, oldest_time) > 0 {
            oldest_path = path
            oldest_time = entry.last_used
            first = false
        }
    }

    // Remove oldest entry
    if oldest_path != "" {
        if entry, ok := cache.entries[oldest_path]; ok {
            image_destroy(&entry.image)
            image_destroy(&entry.thumbnail)
            delete(entry.path)
            free(entry)
            delete_key(&cache.entries, oldest_path)
        }
    }
}

// Clear all cached entries
image_cache_clear :: proc(cache: ^Image_Cache) {
    if cache == nil {
        return
    }

    for _, entry in cache.entries {
        image_destroy(&entry.image)
        image_destroy(&entry.thumbnail)
        delete(entry.path)
        free(entry)
    }

    builtin.clear(&cache.entries)
}

// Get number of cached entries
image_cache_count :: proc(cache: ^Image_Cache) -> int {
    if cache == nil {
        return 0
    }
    return len(cache.entries)
}

// Helper to clone a string
@(private)
clone_string :: proc(s: string) -> string {
    if len(s) == 0 {
        return ""
    }
    buf := make([]u8, len(s))
    copy(buf, transmute([]u8)s)
    return string(buf)
}
