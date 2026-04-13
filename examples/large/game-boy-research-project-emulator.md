# Requirement: "a game boy emulator"

A full 8-bit handheld console emulator with CPU, memory bus, PPU, timer, and cartridge loading. Rendering and input are exposed as pure state transitions; the caller drives the frame loop.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads an entire file into a byte buffer
      - returns error when the path does not exist
      # filesystem
  std.bits
    std.bits.get_bit
      @ (value: u8, index: u8) -> bool
      + returns true when the bit at index is set
      ? index is 0..7 where 0 is the least significant
      # bit_manipulation
    std.bits.set_bit
      @ (value: u8, index: u8, on: bool) -> u8
      + returns value with the bit at index set or cleared
      # bit_manipulation

gameboy
  gameboy.load_cartridge
    @ (rom: bytes) -> result[cartridge_state, string]
    + parses the cartridge header and returns a cartridge state
    - returns error when the header checksum does not match
    - returns error when the rom size is smaller than the header region
    # cartridge
  gameboy.new_machine
    @ (cart: cartridge_state) -> machine_state
    + returns a fresh machine with CPU registers at boot values and work RAM cleared
    ? does not execute a boot rom; starts at the post-boot state
    # construction
  gameboy.read_byte
    @ (machine: machine_state, address: u16) -> u8
    + routes the read through cartridge, work ram, video ram, oam, and io regions
    # memory_bus
  gameboy.write_byte
    @ (machine: machine_state, address: u16, value: u8) -> machine_state
    + routes writes to the correct region and triggers any side effects (dma, timer reset)
    # memory_bus
  gameboy.step_cpu
    @ (machine: machine_state) -> tuple[machine_state, u32]
    + decodes and executes a single instruction at the program counter
    + returns the number of machine cycles consumed
    # cpu
    -> std.bits.get_bit
    -> std.bits.set_bit
  gameboy.step_timer
    @ (machine: machine_state, cycles: u32) -> machine_state
    + advances the divider and TIMA counters; raises the timer interrupt on overflow
    # timer
  gameboy.step_ppu
    @ (machine: machine_state, cycles: u32) -> machine_state
    + advances the pixel processor through oam scan, drawing, hblank, and vblank modes
    + raises the vblank interrupt when line 144 is reached
    # graphics
  gameboy.handle_interrupts
    @ (machine: machine_state) -> machine_state
    + services pending interrupts in priority order when the master enable is set
    # interrupts
  gameboy.run_frame
    @ (machine: machine_state) -> machine_state
    + steps the CPU, timer, and PPU cooperatively until one full 154-line frame has been produced
    # frame_loop
  gameboy.get_framebuffer
    @ (machine: machine_state) -> list[u8]
    + returns the 160x144 indexed pixel buffer produced by the most recent frame
    # output
  gameboy.set_buttons
    @ (machine: machine_state, buttons: u8) -> machine_state
    + updates the joypad register from a bitmask of pressed buttons
    # input
