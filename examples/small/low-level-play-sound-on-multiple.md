# Requirement: "a low-level audio playback library that pushes PCM samples to an output device"

Manages a device handle and a sample queue. The host supplies the actual audio backend via a driver identifier.

std: (all units exist)

audio
  audio.open
    @ (sample_rate_hz: i32, channels: i32, driver: string) -> result[audio_device, string]
    + opens a playback device at the given sample rate and channel count
    - returns error when the driver is unknown or parameters are unsupported
    # device_lifecycle
  audio.write_samples
    @ (device: audio_device, samples: list[f32]) -> result[audio_device, string]
    + appends interleaved PCM samples to the device's pending queue
    - returns error when the device is closed
    # playback
  audio.pending_frames
    @ (device: audio_device) -> i32
    + returns the number of frames still queued for playback
    # query
  audio.close
    @ (device: audio_device) -> result[void, string]
    + drains pending samples and releases the device
    - returns error when the device is already closed
    # device_lifecycle
