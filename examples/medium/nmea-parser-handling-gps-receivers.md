# Requirement: "an nmea sentence parser for gps receivers"

Validates checksums and decodes the common positional sentences into structured fixes.

std
  std.strings
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits a string by the given separator
      # strings
  std.parse
    std.parse.parse_f64
      @ (s: string) -> result[f64, string]
      + parses a 64-bit float from text
      - returns error when the string is not a valid float
      # parsing
    std.parse.parse_i32
      @ (s: string) -> result[i32, string]
      + parses a 32-bit signed integer from text
      - returns error when the string is not a valid integer
      # parsing

nmea
  nmea.verify_checksum
    @ (sentence: string) -> result[string, string]
    + returns the payload between '$' and '*' when the checksum matches
    - returns error when '$' or '*' is missing
    - returns error when the checksum byte does not match XOR of the payload
    # validation
  nmea.sentence_type
    @ (payload: string) -> string
    + returns the 5-character talker+type identifier
    - returns empty string when the payload has no comma
    # classification
    -> std.strings.split
  nmea.parse_lat_lon
    @ (raw: string, hemisphere: string) -> result[f64, string]
    + converts ddmm.mmmm or dddmm.mmmm plus N/S/E/W to decimal degrees
    - returns error when the format is unrecognized
    - returns error when the hemisphere letter is not one of N/S/E/W
    # coordinates
    -> std.parse.parse_f64
  nmea.parse_gga
    @ (payload: string) -> result[fix_record, string]
    + extracts time, latitude, longitude, fix quality, satellites, and altitude
    - returns error when required fields are missing
    # sentences
    -> std.strings.split
    -> nmea.parse_lat_lon
    -> std.parse.parse_i32
    -> std.parse.parse_f64
  nmea.parse_rmc
    @ (payload: string) -> result[fix_record, string]
    + extracts time, date, status, latitude, longitude, speed, and course
    - returns error when required fields are missing
    - returns error when status is not 'A' or 'V'
    # sentences
    -> std.strings.split
    -> nmea.parse_lat_lon
    -> std.parse.parse_f64
  nmea.parse_sentence
    @ (sentence: string) -> result[fix_record, string]
    + verifies checksum, dispatches on sentence type, and returns the fix
    - returns error on unsupported sentence types
    # dispatch
    -> nmea.verify_checksum
    -> nmea.sentence_type
    -> nmea.parse_gga
    -> nmea.parse_rmc
