package render

import "../ffmpeg"
import "core:fmt"
import "core:strings"
import "core:mem"

// Video decoder for thumbnail playback
Video_Decoder :: struct {
	format_ctx:       ^ffmpeg.AVFormatContext,
	codec_ctx:        ^ffmpeg.AVCodecContext,
	sws_ctx:          ^ffmpeg.SwsContext,
	video_stream_idx: i32,
	frame:            ^ffmpeg.AVFrame,
	packet:           ^ffmpeg.AVPacket,
	width:            i32,
	height:           i32,
	// Scaled intermediate dimensions
	thumb_width:      i32,
	thumb_height:     i32,
	// Final cropped square size
	crop_size:        i32,
	// RGBA buffers
	scaled_buffer:    []u8,  // Intermediate scaled buffer
	rgba_buffer:      []u8,  // Final cropped square buffer
	// Frame timing
	frame_duration:   f64,  // seconds per frame
	accumulated_time: f64,  // time since last frame
	// State
	is_open:          bool,
	path:             string,
}

// Open a video file for decoding
video_decoder_open :: proc(path: string, thumb_size: i32 = 128) -> ^Video_Decoder {
	decoder := new(Video_Decoder)
	decoder.path = strings.clone(path)

	// Convert path to cstring
	cpath := strings.clone_to_cstring(path)
	defer delete(cpath)

	// Open input file
	ret := ffmpeg.open_input(&decoder.format_ctx, cpath, nil, nil)
	if ret < 0 {
		fmt.eprintln("Failed to open video:", path, ffmpeg.av_err2str(ret))
		video_decoder_close(decoder)
		return nil
	}

	// Find stream info
	ret = ffmpeg.find_stream_info(decoder.format_ctx, nil)
	if ret < 0 {
		fmt.eprintln("Failed to find stream info:", ffmpeg.av_err2str(ret))
		video_decoder_close(decoder)
		return nil
	}

	// Find video stream
	stream_idx, video_codec, ok := ffmpeg.find_video_stream(decoder.format_ctx)
	if !ok {
		fmt.eprintln("No video stream found")
		video_decoder_close(decoder)
		return nil
	}
	decoder.video_stream_idx = stream_idx

	// Get stream
	stream := decoder.format_ctx.streams[stream_idx]

	// Allocate codec context
	decoder.codec_ctx = ffmpeg.alloc_context3(video_codec)
	if decoder.codec_ctx == nil {
		fmt.eprintln("Failed to allocate codec context")
		video_decoder_close(decoder)
		return nil
	}

	// Copy codec parameters
	ret = ffmpeg.parameters_to_context(decoder.codec_ctx, stream.codecpar)
	if ret < 0 {
		fmt.eprintln("Failed to copy codec parameters:", ffmpeg.av_err2str(ret))
		video_decoder_close(decoder)
		return nil
	}

	// Open codec
	ret = ffmpeg.open2(decoder.codec_ctx, video_codec, nil)
	if ret < 0 {
		fmt.eprintln("Failed to open codec:", ffmpeg.av_err2str(ret))
		video_decoder_close(decoder)
		return nil
	}

	// Store codec dimensions (may differ from actual frame due to rotation)
	decoder.width = decoder.codec_ctx.width
	decoder.height = decoder.codec_ctx.height

	// Note: sws_ctx will be created on first frame when we know actual dimensions

	// Allocate frame and packet
	decoder.frame = ffmpeg.frame_alloc()
	decoder.packet = ffmpeg.packet_alloc()
	if decoder.frame == nil || decoder.packet == nil {
		fmt.eprintln("Failed to allocate frame/packet")
		video_decoder_close(decoder)
		return nil
	}

	// Store final crop size (square) - buffers allocated on first frame
	decoder.crop_size = thumb_size

	// Calculate frame duration from stream
	if stream.avg_frame_rate.den > 0 && stream.avg_frame_rate.num > 0 {
		decoder.frame_duration = f64(stream.avg_frame_rate.den) / f64(stream.avg_frame_rate.num)
	} else {
		decoder.frame_duration = 1.0 / 30.0  // Default 30 fps
	}

	decoder.is_open = true
	return decoder
}

// Read the next frame and return RGBA pixels
// Returns nil if no frame available
video_decoder_read_frame :: proc(decoder: ^Video_Decoder) -> (pixels: []u8, width, height: i32) {
	if decoder == nil || !decoder.is_open {
		return nil, 0, 0
	}

	for {
		// Try to receive a decoded frame
		ret := ffmpeg.receive_frame(decoder.codec_ctx, decoder.frame)
		if ret == 0 {
			// Got a frame - convert to RGBA and crop
			convert_frame_to_rgba(decoder)
			return decoder.rgba_buffer, decoder.crop_size, decoder.crop_size
		}

		if ret != ffmpeg.AVERROR_EAGAIN {
			// Error or EOF
			if ret == ffmpeg.AVERROR_EOF {
				// Loop back to start
				video_decoder_seek_start(decoder)
				continue
			}
			return nil, 0, 0
		}

		// Need more data - read a packet
		ret = ffmpeg.read_frame(decoder.format_ctx, decoder.packet)
		if ret < 0 {
			if ret == ffmpeg.AVERROR_EOF {
				// End of file - loop
				video_decoder_seek_start(decoder)
				continue
			}
			return nil, 0, 0
		}

		// Only process video packets
		if decoder.packet.stream_index != decoder.video_stream_idx {
			ffmpeg.packet_unref(decoder.packet)
			continue
		}

		// Send packet to decoder
		ret = ffmpeg.send_packet(decoder.codec_ctx, decoder.packet)
		ffmpeg.packet_unref(decoder.packet)

		if ret < 0 {
			fmt.eprintln("Error sending packet:", ffmpeg.av_err2str(ret))
			return nil, 0, 0
		}
	}
}

// Check if enough time has passed to decode next frame
// Returns true if we should decode
video_decoder_update :: proc(decoder: ^Video_Decoder, delta_time: f64) -> bool {
	if decoder == nil || !decoder.is_open {
		return false
	}

	decoder.accumulated_time += delta_time

	// Target ~12 FPS for thumbnails
	target_frame_time := 1.0 / 12.0

	if decoder.accumulated_time >= target_frame_time {
		decoder.accumulated_time = 0
		return true
	}
	return false
}

// Seek to start of video
video_decoder_seek_start :: proc(decoder: ^Video_Decoder) {
	if decoder == nil || !decoder.is_open {
		return
	}

	ffmpeg.seek_frame(
		decoder.format_ctx,
		decoder.video_stream_idx,
		0,
		ffmpeg.AVSEEK_FLAG_BACKWARD,
	)
	ffmpeg.flush_buffers(decoder.codec_ctx)
}

// Close the decoder and free resources
video_decoder_close :: proc(decoder: ^Video_Decoder) {
	if decoder == nil {
		return
	}

	if decoder.sws_ctx != nil {
		ffmpeg.freeContext(decoder.sws_ctx)
	}

	if decoder.frame != nil {
		ffmpeg.frame_free(&decoder.frame)
	}

	if decoder.packet != nil {
		ffmpeg.packet_free(&decoder.packet)
	}

	if decoder.codec_ctx != nil {
		ffmpeg.codec_free_context(&decoder.codec_ctx)
	}

	if decoder.format_ctx != nil {
		ffmpeg.close_input(&decoder.format_ctx)
	}

	if len(decoder.scaled_buffer) > 0 {
		delete(decoder.scaled_buffer)
	}

	if len(decoder.rgba_buffer) > 0 {
		delete(decoder.rgba_buffer)
	}

	if len(decoder.path) > 0 {
		delete(decoder.path)
	}

	decoder.is_open = false
	free(decoder)
}

// Initialize scaler and buffers on first frame (uses actual frame dimensions)
@(private)
init_scaler_for_frame :: proc(decoder: ^Video_Decoder) -> bool {
	// Use actual frame dimensions (may differ from codec due to rotation)
	decoder.width = decoder.frame.width
	decoder.height = decoder.frame.height

	if decoder.width <= 0 || decoder.height <= 0 {
		return false
	}

	// Calculate thumbnail dimensions for zoom-to-fill (crop to square)
	aspect := f32(decoder.width) / f32(decoder.height)

	if decoder.width > decoder.height {
		// Landscape: scale height to crop_size, width will be larger
		decoder.thumb_height = decoder.crop_size
		decoder.thumb_width = i32(f32(decoder.crop_size) * aspect)
	} else {
		// Portrait: scale width to crop_size, height will be larger
		decoder.thumb_width = decoder.crop_size
		decoder.thumb_height = i32(f32(decoder.crop_size) / aspect)
	}

	// Ensure even dimensions for scaling
	if decoder.thumb_width % 2 != 0 {
		decoder.thumb_width += 1
	}
	if decoder.thumb_height % 2 != 0 {
		decoder.thumb_height += 1
	}

	// Create scaler context with actual frame dimensions
	decoder.sws_ctx = ffmpeg.create_rgba_scaler(
		decoder.width,
		decoder.height,
		decoder.codec_ctx.pix_fmt,
		decoder.thumb_width,
		decoder.thumb_height,
	)
	if decoder.sws_ctx == nil {
		fmt.eprintln("Failed to create swscale context")
		return false
	}

	// Allocate scaled buffer (intermediate, may be larger than final)
	scaled_buffer_size := int(decoder.thumb_width * decoder.thumb_height * 4)
	decoder.scaled_buffer = make([]u8, scaled_buffer_size)

	// Allocate final RGBA buffer (cropped square)
	final_buffer_size := int(decoder.crop_size * decoder.crop_size * 4)
	decoder.rgba_buffer = make([]u8, final_buffer_size)

	return true
}

// Convert decoded frame to RGBA with center crop
@(private)
convert_frame_to_rgba :: proc(decoder: ^Video_Decoder) {
	// Initialize scaler on first frame
	if decoder.sws_ctx == nil {
		if !init_scaler_for_frame(decoder) {
			ffmpeg.frame_unref(decoder.frame)
			return
		}
	}

	// Setup destination pointers for intermediate scaled buffer
	dst_data: [8]rawptr
	dst_stride: [8]i32
	dst_data[0] = raw_data(decoder.scaled_buffer)
	dst_stride[0] = decoder.thumb_width * 4

	// Copy source data pointers and strides from frame
	src_data: [8]rawptr
	src_stride: [8]i32
	for i in 0..<8 {
		src_data[i] = decoder.frame.data[i]
		src_stride[i] = decoder.frame.linesize[i]
	}

	// Scale and convert to intermediate buffer
	ffmpeg.scale(
		decoder.sws_ctx,
		&src_data[0],
		&src_stride[0],
		0,
		decoder.height,
		&dst_data[0],
		&dst_stride[0],
	)

	// Unref frame for next decode
	ffmpeg.frame_unref(decoder.frame)

	// Crop center from scaled buffer to final rgba_buffer
	crop_x := (decoder.thumb_width - decoder.crop_size) / 2
	crop_y := (decoder.thumb_height - decoder.crop_size) / 2

	src_stride_bytes := int(decoder.thumb_width * 4)
	dst_stride_bytes := int(decoder.crop_size * 4)

	for row in 0..<decoder.crop_size {
		src_row := int(crop_y + row)
		src_offset := src_row * src_stride_bytes + int(crop_x * 4)
		dst_offset := int(row) * dst_stride_bytes

		// Copy row
		for col in 0..<dst_stride_bytes {
			decoder.rgba_buffer[dst_offset + col] = decoder.scaled_buffer[src_offset + col]
		}
	}
}

// Check if path is a video file
is_video_file :: proc(path: string) -> bool {
	lower := strings.to_lower(path)
	defer delete(lower)

	extensions := []string{".mp4", ".mkv", ".avi", ".webm", ".mov", ".wmv", ".flv", ".m4v"}
	for ext in extensions {
		if strings.has_suffix(lower, ext) {
			return true
		}
	}
	return false
}
