package ffmpeg

// FFmpeg libavformat bindings
// Container format handling

foreign import avformat "system:avformat"

// AVInputFormat - demuxer descriptor (opaque)
AVInputFormat :: struct {
	name:             cstring,
	long_name:        cstring,
	// More fields exist
	_padding:         [512]u8,
}

// AVStream - single stream in a container (FFmpeg 6.x layout)
AVStream :: struct {
	av_class:         ^AVClass,
	index:            i32,
	id:               i32,
	codecpar:         ^AVCodecParameters,
	priv_data:        rawptr,
	time_base:        AVRational,
	start_time:       i64,
	duration:         i64,
	nb_frames:        i64,
	disposition:      i32,
	discard:          i32,
	sample_aspect_ratio: AVRational,
	metadata:         ^AVDictionary,
	avg_frame_rate:   AVRational,
	attached_pic:     AVPacket,
	side_data:        rawptr,
	nb_side_data:     i32,
	event_flags:      i32,
	r_frame_rate:     AVRational,
	pts_wrap_bits:    i32,
	// More fields
	_padding:         [256]u8,
}

// AVIOContext - I/O context (opaque)
AVIOContext :: struct {}

// AVFormatContext - main format context
AVFormatContext :: struct {
	av_class:         ^AVClass,
	iformat:          ^AVInputFormat,
	oformat:          rawptr,  // AVOutputFormat*
	priv_data:        rawptr,
	pb:               ^AVIOContext,
	ctx_flags:        i32,
	nb_streams:       u32,
	streams:          [^]^AVStream,
	url:              cstring,
	// Deprecated filename field removed in newer FFmpeg
	start_time:       i64,
	duration:         i64,
	bit_rate:         i64,
	packet_size:      u32,
	max_delay:        i32,
	flags:            i32,
	probesize:        i64,
	max_analyze_duration: i64,
	key:              [^]u8,
	keylen:           i32,
	nb_programs:      u32,
	programs:         rawptr,
	video_codec_id:   AVCodecID,
	audio_codec_id:   AVCodecID,
	subtitle_codec_id: AVCodecID,
	max_index_size:   u32,
	max_picture_buffer: u32,
	nb_chapters:      u32,
	chapters:         rawptr,
	metadata:         ^AVDictionary,
	// More fields exist
	_padding:         [1024]u8,
}

// Format context flags
AVFMT_FLAG_GENPTS :: 0x0001
AVFMT_FLAG_IGNIDX :: 0x0002
AVFMT_FLAG_NONBLOCK :: 0x0004
AVFMT_FLAG_IGNDTS :: 0x0008
AVFMT_FLAG_NOFILLIN :: 0x0010
AVFMT_FLAG_NOPARSE :: 0x0020
AVFMT_FLAG_NOBUFFER :: 0x0040
AVFMT_FLAG_CUSTOM_IO :: 0x0080
AVFMT_FLAG_DISCARD_CORRUPT :: 0x0100
AVFMT_FLAG_FLUSH_PACKETS :: 0x0200
AVFMT_FLAG_BITEXACT :: 0x0400
AVFMT_FLAG_SORT_DTS :: 0x10000
AVFMT_FLAG_FAST_SEEK :: 0x80000
AVFMT_FLAG_SHORTEST :: 0x100000
AVFMT_FLAG_AUTO_BSF :: 0x200000

// Seek flags
AVSEEK_FLAG_BACKWARD :: 1
AVSEEK_FLAG_BYTE :: 2
AVSEEK_FLAG_ANY :: 4
AVSEEK_FLAG_FRAME :: 8

@(default_calling_convention = "c")
@(link_prefix = "avformat_")
foreign avformat {
	// Format context
	alloc_context :: proc() -> ^AVFormatContext ---
	@(link_name = "avformat_free_context")
	format_free_context :: proc(s: ^AVFormatContext) ---

	// Open/close
	open_input :: proc(ps: ^^AVFormatContext, url: cstring, fmt: ^AVInputFormat, options: ^^AVDictionary) -> i32 ---
	close_input :: proc(s: ^^AVFormatContext) ---

	// Stream info
	find_stream_info :: proc(ic: ^AVFormatContext, options: ^^AVDictionary) -> i32 ---
}

@(default_calling_convention = "c")
@(link_prefix = "av_")
foreign avformat {
	// Reading
	read_frame :: proc(s: ^AVFormatContext, pkt: ^AVPacket) -> i32 ---

	// Find best stream
	find_best_stream :: proc(
		ic: ^AVFormatContext,
		type_: AVMediaType,
		wanted_stream_nb: i32,
		related_stream: i32,
		decoder_ret: ^^AVCodec,
		flags: i32,
	) -> i32 ---

	// Seeking
	seek_frame :: proc(s: ^AVFormatContext, stream_index: i32, timestamp: i64, flags: i32) -> i32 ---
	seek_file :: proc(s: ^AVFormatContext, stream_index: i32, min_ts: i64, ts: i64, max_ts: i64, flags: i32) -> i32 ---

	// Utility
	dump_format :: proc(ic: ^AVFormatContext, index: i32, url: cstring, is_output: i32) ---
}

// Helper to find best video stream
find_video_stream :: proc(fmt_ctx: ^AVFormatContext) -> (stream_idx: i32, codec: ^AVCodec, ok: bool) {
	codec_ptr: ^AVCodec
	idx := find_best_stream(fmt_ctx, .AVMEDIA_TYPE_VIDEO, -1, -1, &codec_ptr, 0)
	if idx < 0 {
		return -1, nil, false
	}
	return idx, codec_ptr, true
}
