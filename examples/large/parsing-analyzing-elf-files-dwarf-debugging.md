# Requirement: "an ELF file and DWARF debug information parser"

The project separates the container (ELF) from the debug payload (DWARF). std carries little/big-endian binary reads and LEB128 decoding.

std
  std.binary
    std.binary.read_u16_le
      @ (data: bytes, offset: i32) -> u16
      + reads a little-endian u16 at offset
      # binary
    std.binary.read_u32_le
      @ (data: bytes, offset: i32) -> u32
      + reads a little-endian u32 at offset
      # binary
    std.binary.read_u64_le
      @ (data: bytes, offset: i32) -> u64
      + reads a little-endian u64 at offset
      # binary
    std.binary.read_uleb128
      @ (data: bytes, offset: i32) -> result[tuple[u64, i32], string]
      + reads an unsigned LEB128 value, returning (value, bytes_consumed)
      - returns error when the stream ends mid-value
      # binary
    std.binary.read_sleb128
      @ (data: bytes, offset: i32) -> result[tuple[i64, i32], string]
      + reads a signed LEB128 value, returning (value, bytes_consumed)
      # binary
  std.strings
    std.strings.cstring_at
      @ (data: bytes, offset: i32) -> string
      + returns the NUL-terminated string starting at offset
      # strings

elfdwarf
  elfdwarf.parse_elf
    @ (data: bytes) -> result[elf_file, string]
    + validates the ELF magic and header, dispatching by class and endianness
    - returns error on bad magic or truncated header
    # parsing
    -> std.binary.read_u16_le
    -> std.binary.read_u32_le
    -> std.binary.read_u64_le
  elfdwarf.sections
    @ (f: elf_file) -> list[elf_section]
    + returns every section with name, type, address, offset, and size
    # query
  elfdwarf.section_by_name
    @ (f: elf_file, name: string) -> optional[elf_section]
    + returns the first section whose name matches
    # query
    -> std.strings.cstring_at
  elfdwarf.segments
    @ (f: elf_file) -> list[elf_segment]
    + returns every program header with virtual address, file offset, and size
    # query
  elfdwarf.symbols
    @ (f: elf_file) -> list[elf_symbol]
    + returns every symbol from .symtab and .dynsym with name, value, and binding
    # query
    -> std.strings.cstring_at
  elfdwarf.parse_dwarf
    @ (f: elf_file) -> result[dwarf_info, string]
    + parses .debug_info, .debug_abbrev, .debug_str, and .debug_line
    - returns error when a required DWARF section is missing
    # parsing
    -> std.binary.read_uleb128
    -> std.binary.read_sleb128
  elfdwarf.compilation_units
    @ (d: dwarf_info) -> list[compilation_unit]
    + returns every CU with its DIE tree and version
    # query
  elfdwarf.die_children
    @ (die: debug_info_entry) -> list[debug_info_entry]
    + returns the immediate children of a DIE
    # query
  elfdwarf.die_attribute
    @ (die: debug_info_entry, name: string) -> optional[dwarf_attr_value]
    + returns the attribute value by DW_AT_ name
    # query
  elfdwarf.line_program
    @ (cu: compilation_unit) -> list[line_entry]
    + returns the resolved (address, file, line, column) rows from the line number program
    # query
  elfdwarf.address_to_line
    @ (cu: compilation_unit, address: u64) -> optional[line_entry]
    + returns the line entry covering the given runtime address
    # query
  elfdwarf.function_at_address
    @ (d: dwarf_info, address: u64) -> optional[debug_info_entry]
    + returns the DW_TAG_subprogram DIE whose ranges contain address
    # query
