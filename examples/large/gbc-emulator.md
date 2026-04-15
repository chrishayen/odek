# Requirement: "a game boy color emulator"

A color handheld console emulator. Adds double-speed CPU mode, bank switching, and color palette memory on top of the monochrome model. The host surface (web, desktop, whatever) is out of scope.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads an entire file into a byte buffer
      - returns error when the path does not exist
      # filesystem
  std.bits
    std.bits.get_bit
      fn (value: u8, index: u8) -> bool
      + returns true when the bit at index is set
      # bit_manipulation
    std.bits.set_bit
      fn (value: u8, index: u8, on: bool) -> u8
      + returns value with the bit at index set or cleared
      # bit_manipulation

gbc
  gbc.load_cartridge
    fn (rom: bytes) -> result[cartridge_state, string]
    + parses the header and detects the MBC type and cgb flag
    - returns error when the header checksum does not match
    # cartridge
  gbc.new_machine
    fn (cart: cartridge_state) -> machine_state
    + creates a machine with two work ram banks and eight video ram banks for color mode
    # construction
  gbc.read_byte
    fn (machine: machine_state, address: u16) -> u8
    + routes reads through the active rom bank, ram bank, and wram bank
    # memory_bus
  gbc.write_byte
    fn (machine: machine_state, address: u16, value: u8) -> machine_state
    + writes to the active bank and handles mbc register writes for bank switching
    # memory_bus
  gbc.step_cpu
    fn (machine: machine_state) -> tuple[machine_state, u32]
    + executes one instruction; cycles are reported in double-speed units when enabled
    # cpu
    -> std.bits.get_bit
    -> std.bits.set_bit
  gbc.step_ppu
    fn (machine: machine_state, cycles: u32) -> machine_state
    + advances the pixel processor using the active color palette memory
    + raises the vblank interrupt at line 144
    # graphics
  gbc.set_color_palette
    fn (machine: machine_state, index: u8, rgb555: u16) -> machine_state
    + writes a color to the background or sprite palette memory
    # palette
  gbc.run_frame
    fn (machine: machine_state) -> machine_state
    + steps the machine until one full frame has been produced
    # frame_loop
  gbc.get_framebuffer_rgba
    fn (machine: machine_state) -> bytes
    + returns the 160x144 framebuffer as packed RGBA8 bytes for the most recent frame
    # output
  gbc.set_buttons
    fn (machine: machine_state, buttons: u8) -> machine_state
    + updates the joypad register from a bitmask of pressed buttons
    # input
