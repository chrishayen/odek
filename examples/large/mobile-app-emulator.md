# Requirement: "a high-level emulator for mobile phone applications"

Emulates an application binary by loading a package, translating guest CPU instructions, and stubbing the operating system's framework calls.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of the file
      # filesystem
  std.archive
    std.archive.extract_zip
      @ (archive: bytes) -> result[map[string, bytes], string]
      + extracts a zip archive into a map from entry name to entry bytes
      - returns error on a corrupt archive
      # archive
  std.gfx
    std.gfx.create_surface
      @ (width: i32, height: i32) -> surface_handle
      + allocates a pixel surface for rendering
      # graphics
    std.gfx.blit
      @ (dst: surface_handle, src: bytes, x: i32, y: i32, w: i32, h: i32) -> void
      + copies src pixels into the surface at (x, y)
      # graphics
    std.gfx.present
      @ (surface: surface_handle) -> void
      + displays the surface on the host window
      # graphics
  std.audio
    std.audio.open_output
      @ (sample_rate: i32, channels: i32) -> audio_handle
      + opens a PCM audio output stream
      # audio
    std.audio.queue_pcm
      @ (handle: audio_handle, samples: bytes) -> void
      + queues PCM samples for playback
      # audio

emulator
  emulator.load_bundle
    @ (path: string) -> result[app_bundle, string]
    + reads and extracts an application package into a bundle value
    - returns error when the bundle is missing required metadata
    # loading
    -> std.fs.read_all
    -> std.archive.extract_zip
  emulator.parse_macho
    @ (binary: bytes) -> result[macho_image, string]
    + parses a Mach-O executable into an image descriptor with segments and symbols
    - returns error on an unknown magic number
    # binary
  emulator.new_cpu
    @ () -> cpu_state
    + creates a fresh ARMv7 CPU state with zeroed registers
    # cpu
  emulator.map_image
    @ (cpu: cpu_state, image: macho_image) -> cpu_state
    + loads the image segments into the virtual address space
    # loading
  emulator.step
    @ (cpu: cpu_state) -> result[cpu_state, string]
    + decodes and executes the next instruction
    - returns error on an unimplemented opcode
    # cpu
  emulator.run_until_svc
    @ (cpu: cpu_state) -> result[cpu_state, i32]
    + executes until a supervisor call is issued, returning the call number
    # cpu
  emulator.register_stub
    @ (name: string, handler: fn(cpu: cpu_state) -> cpu_state) -> void
    + registers a host handler for a named framework function
    # os_stubs
  emulator.dispatch_stub
    @ (cpu: cpu_state, name: string) -> cpu_state
    + looks up and invokes a registered stub for name
    # os_stubs
  emulator.touch_event
    @ (cpu: cpu_state, x: f32, y: f32, phase: i32) -> cpu_state
    + injects a touch event into the app's pending event queue
    # input
  emulator.present_frame
    @ (cpu: cpu_state, surface: surface_handle) -> cpu_state
    + flushes the emulated framebuffer to the host surface
    # graphics
    -> std.gfx.blit
    -> std.gfx.present
  emulator.run_frame
    @ (cpu: cpu_state) -> result[cpu_state, string]
    + runs the guest for one frame's worth of execution
    # lifecycle
    -> emulator.run_until_svc
    -> emulator.dispatch_stub
