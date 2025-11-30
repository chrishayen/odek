package ffmpeg

// FFmpeg libswscale bindings
// Pixel format conversion

foreign import swscale "system:swscale"

// SwsContext - scaling/conversion context (opaque)
SwsContext :: struct {}

// Scaling algorithms
SWS_FAST_BILINEAR :: 1
SWS_BILINEAR :: 2
SWS_BICUBIC :: 4
SWS_X :: 8
SWS_POINT :: 0x10
SWS_AREA :: 0x20
SWS_BICUBLIN :: 0x40
SWS_GAUSS :: 0x80
SWS_SINC :: 0x100
SWS_LANCZOS :: 0x200
SWS_SPLINE :: 0x400

@(default_calling_convention = "c")
@(link_prefix = "sws_")
foreign swscale {
	// Context creation/destruction
	getContext :: proc(
		srcW: i32,
		srcH: i32,
		srcFormat: AVPixelFormat,
		dstW: i32,
		dstH: i32,
		dstFormat: AVPixelFormat,
		flags: i32,
		srcFilter: rawptr,  // SwsFilter*
		dstFilter: rawptr,  // SwsFilter*
		param: [^]f64,
	) -> ^SwsContext ---

	freeContext :: proc(swsContext: ^SwsContext) ---

	// Scaling
	scale :: proc(
		c: ^SwsContext,
		srcSlice: [^]rawptr,
		srcStride: [^]i32,
		srcSliceY: i32,
		srcSliceH: i32,
		dst: [^]rawptr,
		dstStride: [^]i32,
	) -> i32 ---
}

// Helper to create a context for converting to RGBA
create_rgba_converter :: proc(width, height: i32, src_fmt: AVPixelFormat) -> ^SwsContext {
	return getContext(
		width, height, src_fmt,
		width, height, .AV_PIX_FMT_RGBA,
		SWS_BILINEAR,
		nil, nil, nil,
	)
}

// Helper to create a context for converting and scaling to RGBA
create_rgba_scaler :: proc(src_w, src_h: i32, src_fmt: AVPixelFormat, dst_w, dst_h: i32) -> ^SwsContext {
	return getContext(
		src_w, src_h, src_fmt,
		dst_w, dst_h, .AV_PIX_FMT_RGBA,
		SWS_BILINEAR,
		nil, nil, nil,
	)
}
