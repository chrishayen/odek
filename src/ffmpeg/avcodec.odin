package ffmpeg

// FFmpeg libavcodec bindings
// Codec and decoder types/functions

foreign import avcodec "system:avcodec"

// Codec IDs - subset of commonly used video codecs
AVCodecID :: enum i32 {
	AV_CODEC_ID_NONE = 0,
	AV_CODEC_ID_MPEG1VIDEO = 1,
	AV_CODEC_ID_MPEG2VIDEO = 2,
	AV_CODEC_ID_H261 = 3,
	AV_CODEC_ID_H263 = 4,
	AV_CODEC_ID_MPEG4 = 12,
	AV_CODEC_ID_MSMPEG4V1 = 15,
	AV_CODEC_ID_MSMPEG4V2 = 16,
	AV_CODEC_ID_MSMPEG4V3 = 17,
	AV_CODEC_ID_WMV1 = 18,
	AV_CODEC_ID_WMV2 = 19,
	AV_CODEC_ID_H264 = 27,
	AV_CODEC_ID_VP8 = 139,
	AV_CODEC_ID_VP9 = 167,
	AV_CODEC_ID_HEVC = 173,
	AV_CODEC_ID_AV1 = 226,
}

// AVCodec - codec descriptor (opaque struct)
AVCodec :: struct {
	name:             cstring,
	long_name:        cstring,
	type_:            AVMediaType,
	id:               AVCodecID,
	capabilities:     i32,
	// More fields exist but we don't need them
	_padding:         [512]u8,
}

// AVCodecParameters - codec parameters from stream
AVCodecParameters :: struct {
	codec_type:       AVMediaType,
	codec_id:         AVCodecID,
	codec_tag:        u32,
	extradata:        [^]u8,
	extradata_size:   i32,
	format:           i32,
	bit_rate:         i64,
	bits_per_coded_sample: i32,
	bits_per_raw_sample: i32,
	profile:          i32,
	level:            i32,
	width:            i32,
	height:           i32,
	sample_aspect_ratio: AVRational,
	field_order:      i32,
	color_range:      i32,
	color_primaries:  i32,
	color_trc:        i32,
	color_space:      i32,
	chroma_location:  i32,
	video_delay:      i32,
	// Audio fields follow but we only need video
	_padding:         [256]u8,
}

// AVCodecContext - main codec context
AVCodecContext :: struct {
	av_class:         ^AVClass,
	log_level_offset: i32,
	codec_type:       AVMediaType,
	codec:            ^AVCodec,
	codec_id:         AVCodecID,
	codec_tag:        u32,
	priv_data:        rawptr,
	internal:         rawptr,
	opaque:           rawptr,
	bit_rate:         i64,
	bit_rate_tolerance: i32,
	global_quality:   i32,
	compression_level: i32,
	flags:            i32,
	flags2:           i32,
	extradata:        [^]u8,
	extradata_size:   i32,
	time_base:        AVRational,
	ticks_per_frame:  i32,
	delay:            i32,
	width:            i32,
	height:           i32,
	coded_width:      i32,
	coded_height:     i32,
	gop_size:         i32,
	pix_fmt:          AVPixelFormat,
	// Many more fields exist - pad to be safe
	_padding:         [2048]u8,
}

// AVPacket - compressed data packet
AVPacket :: struct {
	buf:              ^AVBufferRef,
	pts:              i64,
	dts:              i64,
	data:             [^]u8,
	size:             i32,
	stream_index:     i32,
	flags:            i32,
	side_data:        rawptr,
	side_data_elems:  i32,
	duration:         i64,
	pos:              i64,
	// Additional fields
	_padding:         [64]u8,
}

// Packet flags
AV_PKT_FLAG_KEY :: 0x0001
AV_PKT_FLAG_CORRUPT :: 0x0002

@(default_calling_convention = "c")
@(link_prefix = "avcodec_")
foreign avcodec {
	// Codec context
	alloc_context3 :: proc(codec: ^AVCodec) -> ^AVCodecContext ---
	@(link_name = "avcodec_free_context")
	codec_free_context :: proc(avctx: ^^AVCodecContext) ---
	open2 :: proc(avctx: ^AVCodecContext, codec: ^AVCodec, options: ^^AVDictionary) -> i32 ---
	close :: proc(avctx: ^AVCodecContext) -> i32 ---

	// Codec parameters
	parameters_to_context :: proc(codec: ^AVCodecContext, par: ^AVCodecParameters) -> i32 ---

	// Decoding
	send_packet :: proc(avctx: ^AVCodecContext, avpkt: ^AVPacket) -> i32 ---
	receive_frame :: proc(avctx: ^AVCodecContext, frame: ^AVFrame) -> i32 ---

	// Codec finding
	find_decoder :: proc(id: AVCodecID) -> ^AVCodec ---
	find_decoder_by_name :: proc(name: cstring) -> ^AVCodec ---

	// Flush
	flush_buffers :: proc(avctx: ^AVCodecContext) ---
}

@(default_calling_convention = "c")
@(link_prefix = "av_")
foreign avcodec {
	// Packet allocation
	packet_alloc :: proc() -> ^AVPacket ---
	packet_free :: proc(pkt: ^^AVPacket) ---
	packet_unref :: proc(pkt: ^AVPacket) ---
}
