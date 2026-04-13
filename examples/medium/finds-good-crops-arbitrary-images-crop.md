# Requirement: "a library that finds good crops for arbitrary images and crop sizes"

Scores candidate crop rectangles against image saliency and returns the highest-scoring crop for the requested aspect ratio.

std
  std.image
    std.image.pixel_at
      @ (img: image, x: i32, y: i32) -> pixel
      + returns the RGB pixel at (x, y)
      # image
    std.image.dimensions
      @ (img: image) -> tuple[i32, i32]
      + returns width and height in pixels
      # image

smart_crop
  smart_crop.edge_strength
    @ (img: image, x: i32, y: i32) -> f64
    + returns a luminance gradient magnitude at the pixel
    + higher values indicate edges likely to be salient
    # analysis
    -> std.image.pixel_at
  smart_crop.skin_score
    @ (p: pixel) -> f64
    + returns a value in [0, 1] estimating how skin-tone-like the pixel is
    # analysis
  smart_crop.saturation_score
    @ (p: pixel) -> f64
    + returns saturation in [0, 1]
    # analysis
  smart_crop.build_saliency_map
    @ (img: image) -> saliency_map
    + returns a per-pixel weighted sum of edge, skin, and saturation scores
    # map_building
    -> std.image.dimensions
  smart_crop.candidate_rects
    @ (img_width: i32, img_height: i32, target_ratio: f64, step: i32) -> list[rect]
    + enumerates rectangles matching target_ratio at varying sizes and offsets
    + each successive candidate shifts by step pixels
    # candidate_generation
  smart_crop.score_rect
    @ (map: saliency_map, r: rect) -> f64
    + returns the sum of saliency values inside r, normalized for area
    + penalizes rectangles whose center is far from the saliency centroid
    # scoring
  smart_crop.find_best_crop
    @ (img: image, target_width: i32, target_height: i32) -> rect
    + returns the highest-scoring rectangle with the requested aspect ratio
    - returns the full image rect when no candidate fits
    # selection
