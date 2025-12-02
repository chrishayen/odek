package render

import "core:container/lru"
import "core:strings"

// Image cache entry (stored in LRU cache)
Image_Cache_Entry :: struct {
    image:     Image,
    thumbnail: Image,
}

// LRU image cache wrapper
Image_Cache :: struct {
    cache:      lru.Cache(string, Image_Cache_Entry),
    thumb_size: i32,
}

// Create image cache
image_cache_create :: proc(max_entries: int = 100, thumb_size: i32 = 256) -> ^Image_Cache {
    cache := new(Image_Cache)
    lru.init(&cache.cache, max_entries)
    cache.cache.on_remove = on_entry_remove
    cache.thumb_size = thumb_size
    return cache
}

// Cleanup callback when entries are evicted or removed
@(private)
on_entry_remove :: proc(key: string, value: Image_Cache_Entry, _: rawptr) {
    delete(key)
    img := value.image
    thumb := value.thumbnail
    image_destroy(&img)
    image_destroy(&thumb)
}

// Destroy cache and free all images
image_cache_destroy :: proc(cache: ^Image_Cache) {
    if cache == nil {
        return
    }
    lru.destroy(&cache.cache, true)
    free(cache)
}

// Get cached image and thumbnail by path
// Returns (full_image, thumbnail, found)
image_cache_get :: proc(cache: ^Image_Cache, path: string) -> (^Image, ^Image, bool) {
    if cache == nil {
        return nil, nil, false
    }
    entry, ok := lru.get_ptr(&cache.cache, path)
    if !ok {
        return nil, nil, false
    }
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
    if entry, ok := lru.get_ptr(&cache.cache, path); ok {
        return &entry.image, &entry.thumbnail, true
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

    // Clone path - the cache takes ownership
    owned_path := strings.clone(path)

    // Add to cache (may evict old entries via on_remove callback)
    lru.set(&cache.cache, owned_path, Image_Cache_Entry{image = img, thumbnail = thumb})

    // Get pointer to cached entry
    entry, _ := lru.get_ptr(&cache.cache, path)
    return &entry.image, &entry.thumbnail, true
}

// Clear all cached entries
image_cache_clear :: proc(cache: ^Image_Cache) {
    if cache == nil {
        return
    }
    lru.clear(&cache.cache, true)
}

// Get number of cached entries
image_cache_count :: proc(cache: ^Image_Cache) -> int {
    if cache == nil {
        return 0
    }
    return cache.cache.count
}
