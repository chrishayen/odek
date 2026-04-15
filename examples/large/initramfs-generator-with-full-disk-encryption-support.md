# Requirement: "an initramfs generator with full-disk encryption support"

Builds a cpio archive from a rootfs layout, pulls in the kernel modules needed at boot, and wires up an unlock step for encrypted block devices.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file fully into memory
      - returns error when unreadable
      # filesystem
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns the entries under a directory
      - returns error when not a directory
      # filesystem
    std.fs.stat
      fn (path: string) -> result[file_info, string]
      + returns mode, size, and type for a path
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to path
      # filesystem
  std.compress
    std.compress.zstd_compress
      fn (data: bytes, level: i32) -> bytes
      + returns zstd-compressed data at the given level
      # compression
  std.archive
    std.archive.cpio_newc_append
      fn (archive: bytes, name: string, mode: u32, content: bytes) -> bytes
      + appends an entry to a newc-format cpio archive
      # archive_format
    std.archive.cpio_newc_trailer
      fn (archive: bytes) -> bytes
      + writes the TRAILER!!! marker that closes a cpio stream
      # archive_format

initramfs
  initramfs.new
    fn () -> initramfs_spec
    + creates an empty spec with no files and no modules
    # construction
  initramfs.add_file
    fn (spec: initramfs_spec, path: string, mode: u32) -> result[initramfs_spec, string]
    + stages a host file to appear at the same path inside the image
    - returns error when the source is missing
    # staging
    -> std.fs.read_all
    -> std.fs.stat
  initramfs.add_directory
    fn (spec: initramfs_spec, path: string) -> result[initramfs_spec, string]
    + recursively stages a directory tree
    # staging
    -> std.fs.list_dir
    -> std.fs.stat
  initramfs.add_module
    fn (spec: initramfs_spec, module_name: string) -> result[initramfs_spec, string]
    + resolves a kernel module and its dependencies from the running kernel's modules.dep
    - returns error when the module is unknown
    # module_resolution
    -> std.fs.read_all
  initramfs.enable_encrypted_root
    fn (spec: initramfs_spec, device_uuid: string, key_source: i32) -> initramfs_spec
    + adds the cryptsetup binary, its libraries, and an init script that unlocks the device
    # crypto_support
  initramfs.set_init_script
    fn (spec: initramfs_spec, script: string) -> initramfs_spec
    + overrides the default /init script
    # configuration
  initramfs.build
    fn (spec: initramfs_spec) -> result[bytes, string]
    + emits a cpio archive containing all staged files
    - returns error when any staged file cannot be read
    # assembly
    -> std.archive.cpio_newc_append
    -> std.archive.cpio_newc_trailer
  initramfs.write
    fn (spec: initramfs_spec, path: string, compress: bool) -> result[void, string]
    + builds, optionally zstd-compresses, and writes to path
    # output
    -> std.compress.zstd_compress
    -> std.fs.write_all
