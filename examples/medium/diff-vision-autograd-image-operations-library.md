# Requirement: "a differentiable computer vision library with autograd-friendly image operations"

Image operations on a generic tensor type, each preserving a gradient path. The host's autograd engine supplies the tensor primitive.

std: (all units exist)

diff_vision
  diff_vision.from_image
    fn (pixels: bytes, width: i32, height: i32, channels: i32) -> tensor
    + returns a tensor of shape [channels, height, width] with values in [0, 1]
    # ingestion
  diff_vision.to_image
    fn (t: tensor) -> bytes
    + returns raw pixels clamped to [0, 255]
    # export
  diff_vision.rgb_to_grayscale
    fn (t: tensor) -> tensor
    + returns a 1-channel tensor using differentiable luma weights
    ? gradients flow back to each input channel
    # color_conversion
  diff_vision.warp_affine
    fn (t: tensor, matrix: tensor) -> tensor
    + returns t sampled under the given 2x3 affine matrix with bilinear interpolation
    ? sampling is differentiable through the matrix and input
    # geometric_transform
  diff_vision.gaussian_blur
    fn (t: tensor, sigma: f32) -> tensor
    + returns the tensor convolved with a separable Gaussian kernel
    # filtering
  diff_vision.sobel_edges
    fn (t: tensor) -> tensor
    + returns the gradient magnitude using differentiable Sobel operators
    # edge_detection
  diff_vision.ssim_loss
    fn (pred: tensor, target: tensor) -> tensor
    + returns 1 - SSIM as a scalar loss tensor
    - returns error when shapes differ
    # loss
  diff_vision.homography
    fn (src_points: list[tensor], dst_points: list[tensor]) -> tensor
    + returns the 3x3 homography matrix solving the correspondences
    - returns error when fewer than 4 points are provided
    # geometry
  diff_vision.normalize
    fn (t: tensor, mean: list[f32], std: list[f32]) -> tensor
    + returns (t - mean) / std per channel
    # preprocessing
