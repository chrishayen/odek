# Requirement: "a library for parsing and interpreting the results of computational chemistry packages"

Parses output files from chemistry computation tools into a normalized result structure. Format detection and extraction live in the project layer; file reading is a std primitive.

std
  std.fs
    std.fs.read_text
      @ (path: string) -> result[string, string]
      + returns the full text of the file
      - returns error when the file does not exist or is unreadable
      # filesystem

chem_parse
  chem_parse.detect_format
    @ (raw: string) -> result[string, string]
    + returns a short format tag (e.g. "gaussian", "orca", "nwchem") by inspecting header lines
    - returns error when no known format signature is found
    # format_detection
  chem_parse.parse
    @ (raw: string, format: string) -> result[chem_result, string]
    + returns a chem_result populated with whatever sections the format exposes
    - returns error when the format tag is unknown
    - returns error when the text is truncated mid-section
    # parsing
    -> std.fs.read_text
  chem_parse.scf_energies
    @ (res: chem_result) -> list[f64]
    + returns all SCF energies in the order they appear
    ? empty list when no SCF block was present
    # extraction
  chem_parse.final_geometry
    @ (res: chem_result) -> optional[list[tuple[string, f64, f64, f64]]]
    + returns the last optimized geometry as (element, x, y, z) tuples
    - returns none when no geometry block was parsed
    # extraction
  chem_parse.vibrational_modes
    @ (res: chem_result) -> list[tuple[f64, f64]]
    + returns (frequency, intensity) pairs for each vibrational mode
    # extraction
