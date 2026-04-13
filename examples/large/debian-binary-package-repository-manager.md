# Requirement: "a library for managing Debian-style binary package repositories"

Ingests package files into local mirrors, organizes them into named snapshots and publishes, and produces the signed metadata files expected by a Debian-style apt client.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full file contents
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes the full contents to the file
      # filesystem
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns all regular files under root
      # filesystem
  std.hash
    std.hash.sha256
      @ (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest
      # hashing
    std.hash.md5
      @ (data: bytes) -> bytes
      + returns the 16-byte MD5 digest
      # hashing
  std.compress
    std.compress.gzip
      @ (data: bytes) -> bytes
      + returns the gzip-compressed representation
      # compression

debrepo
  debrepo.parse_control
    @ (raw: string) -> result[map[string,string], string]
    + parses an RFC822-style control paragraph into a field-to-value map
    - returns error on malformed paragraph
    # parsing
  debrepo.read_package
    @ (path: string) -> result[package_meta, string]
    + opens a .deb file and extracts its control fields plus file size and hashes
    - returns error when the file is not a valid ar archive
    # ingestion
    -> std.fs.read_all
    -> std.hash.sha256
    -> std.hash.md5
  debrepo.new_mirror
    @ (name: string) -> mirror_state
    + creates an empty mirror with the given name
    # mirror
  debrepo.add_package
    @ (mirror: mirror_state, meta: package_meta) -> mirror_state
    + adds a package entry to the mirror
    # mirror
  debrepo.create_snapshot
    @ (mirror: mirror_state, snapshot_name: string) -> snapshot_state
    + freezes the current package list into a named snapshot
    # snapshot
  debrepo.merge_snapshots
    @ (a: snapshot_state, b: snapshot_state, name: string) -> snapshot_state
    + combines two snapshots, newer package versions winning on conflict
    # snapshot
  debrepo.diff_snapshots
    @ (a: snapshot_state, b: snapshot_state) -> list[snapshot_diff_entry]
    + returns added, removed, and upgraded package entries
    # snapshot
  debrepo.render_packages_index
    @ (snapshot: snapshot_state) -> string
    + renders the control-file-style Packages index for a snapshot
    # metadata
  debrepo.render_release
    @ (snapshot: snapshot_state, packages_index: string) -> string
    + renders the Release file containing hashes of each index
    # metadata
    -> std.hash.sha256
    -> std.hash.md5
  debrepo.publish
    @ (snapshot: snapshot_state, dist: string, out_root: string) -> result[void, string]
    + writes packages index and release files under out_root/dists/dist
    - returns error when out_root cannot be written
    # publish
    -> std.fs.write_all
    -> std.compress.gzip
