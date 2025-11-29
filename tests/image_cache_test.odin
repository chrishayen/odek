package tests

import "../src/render"
import "core:testing"

@(test)
test_image_cache_create_destroy :: proc(t: ^testing.T) {
    cache := render.image_cache_create(50, 128)
    testing.expect(t, cache != nil, "Cache should be created")
    testing.expect(t, render.image_cache_count(cache) == 0, "Cache should be empty initially")

    render.image_cache_destroy(cache)
    // Should complete without crash
    testing.expect(t, true, "Cache destroy should not crash")
}

@(test)
test_image_cache_destroy_nil :: proc(t: ^testing.T) {
    // Should not crash on nil
    render.image_cache_destroy(nil)
    testing.expect(t, true, "Destroying nil cache should not crash")
}

@(test)
test_image_cache_get_nonexistent :: proc(t: ^testing.T) {
    cache := render.image_cache_create()
    defer render.image_cache_destroy(cache)

    img, thumb, ok := render.image_cache_get(cache, "nonexistent.png")
    testing.expect(t, !ok, "Getting nonexistent image should fail")
    testing.expect(t, img == nil, "Image should be nil")
    testing.expect(t, thumb == nil, "Thumbnail should be nil")
}

@(test)
test_image_cache_get_nil_cache :: proc(t: ^testing.T) {
    img, thumb, ok := render.image_cache_get(nil, "test.png")
    testing.expect(t, !ok, "Getting from nil cache should fail")
    testing.expect(t, img == nil, "Image should be nil")
    testing.expect(t, thumb == nil, "Thumbnail should be nil")
}

@(test)
test_image_cache_load_nonexistent :: proc(t: ^testing.T) {
    cache := render.image_cache_create()
    defer render.image_cache_destroy(cache)

    img, thumb, ok := render.image_cache_load(cache, "nonexistent_file_12345.png")
    testing.expect(t, !ok, "Loading nonexistent file should fail")
    testing.expect(t, img == nil, "Image should be nil on failure")
    testing.expect(t, thumb == nil, "Thumbnail should be nil on failure")
    testing.expect(t, render.image_cache_count(cache) == 0, "Cache should still be empty")
}

@(test)
test_image_cache_clear :: proc(t: ^testing.T) {
    cache := render.image_cache_create()
    defer render.image_cache_destroy(cache)

    // Clear empty cache should not crash
    render.image_cache_clear(cache)
    testing.expect(t, render.image_cache_count(cache) == 0, "Cache should be empty after clear")
}

@(test)
test_image_cache_clear_nil :: proc(t: ^testing.T) {
    // Should not crash
    render.image_cache_clear(nil)
    testing.expect(t, true, "Clearing nil cache should not crash")
}

@(test)
test_image_cache_count_nil :: proc(t: ^testing.T) {
    count := render.image_cache_count(nil)
    testing.expect(t, count == 0, "Count of nil cache should be 0")
}

@(test)
test_image_cache_get_thumbnail_nonexistent :: proc(t: ^testing.T) {
    cache := render.image_cache_create()
    defer render.image_cache_destroy(cache)

    thumb, ok := render.image_cache_get_thumbnail(cache, "nonexistent.png")
    testing.expect(t, !ok, "Getting thumbnail of nonexistent file should fail")
    testing.expect(t, thumb == nil, "Thumbnail should be nil")
}

@(test)
test_image_cache_evict_lru_empty :: proc(t: ^testing.T) {
    cache := render.image_cache_create()
    defer render.image_cache_destroy(cache)

    // Evicting from empty cache should not crash
    render.image_cache_evict_lru(cache)
    testing.expect(t, render.image_cache_count(cache) == 0, "Cache should still be empty")
}

@(test)
test_image_cache_evict_lru_nil :: proc(t: ^testing.T) {
    // Should not crash
    render.image_cache_evict_lru(nil)
    testing.expect(t, true, "Evicting from nil cache should not crash")
}
