# Requirement: "a builder for binary Debian-format packages"

Builds a .deb archive from a manifest and a staged file tree. A .deb is an ar archive containing debian-binary, control.tar.gz, and data.tar.gz.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file's contents as bytes
      - returns error when the file cannot be opened
      # filesystem
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns all file paths under root in depth-first order
      - returns error when root does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, creating or truncating
      - returns error when the path is not writable
      # filesystem
  std.archive
    std.archive.tar_create
      @ (entries: list[tar_entry]) -> bytes
      + builds an uncompressed tar archive from a list of named entries with mode and content
      # archiving
    std.archive.ar_create
      @ (entries: list[ar_entry]) -> bytes
      + builds a Unix ar archive from named entries in the given order
      # archiving
  std.compress
    std.compress.gzip
      @ (data: bytes) -> bytes
      + compresses data with gzip
      # compression
  std.hash
    std.hash.md5_hex
      @ (data: bytes) -> string
      + returns the lowercase hex md5 of data
      # hashing

deb
  deb.new_manifest
    @ (name: string, version: string, arch: string) -> deb_manifest
    + creates a manifest with the required package fields
    # construction
  deb.set_field
    @ (m: deb_manifest, key: string, value: string) -> deb_manifest
    + sets an optional control field such as Maintainer, Depends, or Description
    # construction
  deb.validate_manifest
    @ (m: deb_manifest) -> result[void, string]
    + checks that name, version, and architecture are non-empty and well-formed
    - returns error when version does not match a permitted version format
    # validation
  deb.render_control
    @ (m: deb_manifest, installed_size: i64) -> string
    + renders a control file body from the manifest including the computed installed size
    # control
  deb.collect_payload
    @ (staging_root: string) -> result[list[tar_entry], string]
    + walks the staging directory and converts each file into a tar entry rooted at "./"
    - returns error when the staging root does not exist
    # payload
    -> std.fs.walk
    -> std.fs.read_all
  deb.build_md5sums
    @ (payload: list[tar_entry]) -> string
    + renders an md5sums file mapping each payload path to its md5 digest
    # control
    -> std.hash.md5_hex
  deb.build_control_tar
    @ (control: string, md5sums: string) -> bytes
    + builds the gzipped control tar containing the control and md5sums files
    # control
    -> std.archive.tar_create
    -> std.compress.gzip
  deb.build_data_tar
    @ (payload: list[tar_entry]) -> bytes
    + builds the gzipped data tar containing the payload tree
    # payload
    -> std.archive.tar_create
    -> std.compress.gzip
  deb.build
    @ (m: deb_manifest, staging_root: string) -> result[bytes, string]
    + validates the manifest, collects the payload, and returns a complete .deb ar archive
    - returns error when validation or payload collection fails
    # pipeline
    -> deb.validate_manifest
    -> deb.collect_payload
    -> deb.render_control
    -> deb.build_md5sums
    -> deb.build_control_tar
    -> deb.build_data_tar
    -> std.archive.ar_create
