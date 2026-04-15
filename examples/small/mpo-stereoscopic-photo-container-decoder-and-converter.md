# Requirement: "a decoder and converter for multi-picture-object stereoscopic photo containers"

An MPO file is a concatenation of JPEG streams. The library splits the container and exposes each frame as raw JPEG bytes.

std: (all units exist)

mpo_decode
  mpo_decode.split_frames
    fn (container: bytes) -> result[list[bytes], string]
    + returns each embedded JPEG stream as a separate byte slice
    + the first slice is the primary image; subsequent slices are alternate views
    - returns error when the container does not begin with a JPEG SOI marker
    - returns error when a frame is truncated before its EOI marker
    # parsing
  mpo_decode.frame_count
    fn (container: bytes) -> result[i32, string]
    + returns the number of embedded JPEG streams
    - returns error when the container is not a valid MPO
    # inspection
  mpo_decode.primary_frame
    fn (container: bytes) -> result[bytes, string]
    + returns the first embedded JPEG stream
    - returns error when the container has no frames
    # access
  mpo_decode.side_by_side
    fn (left: bytes, right: bytes) -> bytes
    + concatenates two frames into one container with the left frame primary
    ? the output is a well-formed multi-frame container with exactly two entries
    # conversion
