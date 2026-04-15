# Requirement: "a dominant color extractor for raster images"

Given raw RGB pixel data, return the most prominent colors by bucketing into a coarse color cube.

std: (all units exist)

dominant
  dominant.extract
    fn (pixels: bytes, count: i32) -> list[u32]
    + returns up to count colors packed as 0x00RRGGBB sorted by frequency
    + buckets each pixel into a coarse cube to tolerate noise
    - returns an empty list when pixels is empty
    ? pixels are interpreted as contiguous RGB triples
    # color_analysis
  dominant.histogram
    fn (pixels: bytes) -> map[u32, i32]
    + returns a bucket-to-pixel-count map
    - returns an empty map when pixels is empty
    # color_analysis
