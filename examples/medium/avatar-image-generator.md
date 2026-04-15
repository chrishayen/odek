# Requirement: "an avatar image generator"

A library that composes an avatar image from a seed by layering a background, body, and face parts. Image encoding happens through a generic std primitive.

std
  std.hash
    std.hash.fnv64
      fn (data: bytes) -> u64
      + returns a 64-bit FNV-1a hash
      # hashing
  std.image
    std.image.new_rgba
      fn (width: i32, height: i32) -> image
      + creates a blank RGBA image of the given dimensions
      # image
    std.image.draw_layer
      fn (canvas: image, layer: image, x: i32, y: i32) -> image
      + alpha-composites a layer onto the canvas at the given offset
      # image
    std.image.encode_png
      fn (img: image) -> bytes
      + returns PNG bytes for the given image
      # image

avatar
  avatar.seed_from_string
    fn (input: string) -> u64
    + returns a deterministic seed for a given input string
    ? different inputs map to distinguishable avatars with high probability
    # seeding
    -> std.hash.fnv64
  avatar.pick_parts
    fn (seed: u64) -> avatar_parts
    + selects background, body, and face part indices from the seed
    ? uses the seed to index fixed part lists so the same seed yields the same avatar
    # selection
  avatar.render
    fn (parts: avatar_parts, size: i32) -> image
    + composes the chosen parts into a square avatar of the given size
    # rendering
    -> std.image.new_rgba
    -> std.image.draw_layer
  avatar.generate_png
    fn (input: string, size: i32) -> bytes
    + end-to-end: returns PNG bytes for the avatar derived from the input string
    # pipeline
    -> std.image.encode_png
