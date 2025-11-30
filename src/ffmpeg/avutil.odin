package ffmpeg

// FFmpeg libavutil bindings
// Basic types and error handling

foreign import avutil "system:avutil"

// Pixel formats - subset of commonly used ones
AVPixelFormat :: enum i32 {
	AV_PIX_FMT_NONE = -1,
	AV_PIX_FMT_YUV420P = 0,
	AV_PIX_FMT_YUYV422 = 1,
	AV_PIX_FMT_RGB24 = 2,
	AV_PIX_FMT_BGR24 = 3,
	AV_PIX_FMT_YUV422P = 4,
	AV_PIX_FMT_YUV444P = 5,
	AV_PIX_FMT_YUV410P = 6,
	AV_PIX_FMT_YUV411P = 7,
	AV_PIX_FMT_GRAY8 = 8,
	AV_PIX_FMT_MONOWHITE = 9,
	AV_PIX_FMT_MONOBLACK = 10,
	AV_PIX_FMT_PAL8 = 11,
	AV_PIX_FMT_YUVJ420P = 12,
	AV_PIX_FMT_YUVJ422P = 13,
	AV_PIX_FMT_YUVJ444P = 14,
	AV_PIX_FMT_NV12 = 23,
	AV_PIX_FMT_NV21 = 24,
	AV_PIX_FMT_ARGB = 25,
	AV_PIX_FMT_RGBA = 26,
	AV_PIX_FMT_ABGR = 27,
	AV_PIX_FMT_BGRA = 28,
	AV_PIX_FMT_RGB0 = 295,
	AV_PIX_FMT_BGR0 = 296,
	AV_PIX_FMT_0RGB = 297,
	AV_PIX_FMT_0BGR = 298,
}

// Media types
AVMediaType :: enum i32 {
	AVMEDIA_TYPE_UNKNOWN = -1,
	AVMEDIA_TYPE_VIDEO = 0,
	AVMEDIA_TYPE_AUDIO = 1,
	AVMEDIA_TYPE_DATA = 2,
	AVMEDIA_TYPE_SUBTITLE = 3,
	AVMEDIA_TYPE_ATTACHMENT = 4,
	AVMEDIA_TYPE_NB = 5,
}

// Rational number
AVRational :: struct {
	num: i32,
	den: i32,
}

// AVFrame - video/audio frame
// This is an opaque struct - we only need pointers to it
AVFrame :: struct {
	data:            [8]rawptr,
	linesize:        [8]i32,
	extended_data:   ^rawptr,
	width:           i32,
	height:          i32,
	nb_samples:      i32,
	format:          i32,
	key_frame:       i32,
	pict_type:       i32,
	sample_aspect_ratio: AVRational,
	pts:             i64,
	pkt_dts:         i64,
	time_base:       AVRational,
	// Additional fields exist but we only need the above
	_padding:        [256]u8,
}

// AVBuffer reference - opaque
AVBufferRef :: struct {}

// AVDictionary - opaque key/value storage
AVDictionary :: struct {}

// AVClass - used for logging/options - opaque
AVClass :: struct {}

@(default_calling_convention = "c")
@(link_prefix = "av_")
foreign avutil {
	// Frame allocation
	frame_alloc :: proc() -> ^AVFrame ---
	frame_free :: proc(frame: ^^AVFrame) ---
	frame_unref :: proc(frame: ^AVFrame) ---
	frame_get_buffer :: proc(frame: ^AVFrame, align: i32) -> i32 ---

	// Memory
	malloc :: proc(size: uint) -> rawptr ---
	mallocz :: proc(size: uint) -> rawptr ---
	free :: proc(ptr: rawptr) ---
	freep :: proc(ptr: ^rawptr) ---

	// Error handling
	strerror :: proc(errnum: i32, errbuf: [^]u8, errbuf_size: uint) -> [^]u8 ---

	// Logging
	log_set_level :: proc(level: i32) ---
}

// Error codes (negative POSIX error codes)
AVERROR_EOF :: -541478725  // FFERRTAG('E','O','F',' ')
AVERROR_EAGAIN :: -11

// Log levels
AV_LOG_QUIET :: -8
AV_LOG_PANIC :: 0
AV_LOG_FATAL :: 8
AV_LOG_ERROR :: 16
AV_LOG_WARNING :: 24
AV_LOG_INFO :: 32
AV_LOG_VERBOSE :: 40
AV_LOG_DEBUG :: 48

// Helper to get error string
av_err2str :: proc(errnum: i32) -> string {
	@(static) buf: [64]u8
	strerror(errnum, &buf[0], 64)
	len := 0
	for i := 0; i < 64 && buf[i] != 0; i += 1 {
		len = i + 1
	}
	return string(buf[:len])
}
