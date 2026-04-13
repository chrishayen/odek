# Requirement: "a version control library compatible with the git object model"

Object store, refs, index, commit graph, and pack-file transfer for a git-compatible repository.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the file contents
      - returns error when the path cannot be opened
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes the bytes to the file
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns the entries in the directory
      # filesystem
  std.crypto
    std.crypto.sha1
      @ (data: bytes) -> bytes
      + returns the 20-byte sha1 digest of the input
      # cryptography
  std.compression
    std.compression.zlib_deflate
      @ (data: bytes) -> bytes
      + returns the zlib-compressed form of the input
      # compression
    std.compression.zlib_inflate
      @ (data: bytes) -> result[bytes, string]
      + returns the decompressed bytes
      - returns error when the stream is corrupt
      # compression
  std.http
    std.http.get
      @ (url: string, headers: map[string, string]) -> result[http_response, string]
      + returns the http response
      # http
    std.http.post
      @ (url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + returns the http response
      # http

git
  git.init_repo
    @ (path: string) -> result[void, string]
    + creates the directory layout and an empty HEAD
    - returns error when the path already contains a repository
    # lifecycle
    -> std.fs.write_all
  git.hash_object
    @ (obj_type: string, content: bytes) -> string
    + returns the sha1 hex id for a git object with the given header and content
    # object_store
    -> std.crypto.sha1
  git.write_object
    @ (repo_path: string, obj_type: string, content: bytes) -> result[string, string]
    + stores the object and returns its id
    + is a no-op when the object already exists
    # object_store
    -> std.crypto.sha1
    -> std.compression.zlib_deflate
    -> std.fs.write_all
  git.read_object
    @ (repo_path: string, id: string) -> result[tuple[string, bytes], string]
    + returns the type and raw content
    - returns error when the object is not present
    # object_store
    -> std.fs.read_all
    -> std.compression.zlib_inflate
  git.build_tree
    @ (entries: list[tree_entry]) -> bytes
    + returns the canonical tree object content for the entries
    ? entries are sorted by name as git requires
    # objects
  git.build_commit
    @ (tree_id: string, parent_ids: list[string], author: string, message: string, timestamp_s: i64) -> bytes
    + returns the canonical commit object content
    # objects
  git.read_ref
    @ (repo_path: string, name: string) -> result[string, string]
    + returns the commit id the ref points to
    - returns error when the ref does not exist
    # refs
    -> std.fs.read_all
  git.update_ref
    @ (repo_path: string, name: string, commit_id: string) -> result[void, string]
    + writes the ref atomically
    # refs
    -> std.fs.write_all
  git.read_index
    @ (repo_path: string) -> result[list[index_entry], string]
    + returns the staged entries in the working index
    - returns error when the index is corrupt
    # index
    -> std.fs.read_all
  git.write_index
    @ (repo_path: string, entries: list[index_entry]) -> result[void, string]
    + writes the staged entries to the index file
    # index
    -> std.fs.write_all
  git.stage_file
    @ (repo_path: string, relative_path: string) -> result[void, string]
    + hashes and stores the file, then updates the index entry
    - returns error when the file is outside the repository
    # index
  git.commit
    @ (repo_path: string, author: string, message: string, timestamp_s: i64) -> result[string, string]
    + returns the id of the new commit, builds tree from current index, advances HEAD
    - returns error when the index is empty
    # workflow
  git.walk_commits
    @ (repo_path: string, start: string, limit: i32) -> result[list[commit_summary], string]
    + returns up to limit commits in topological order starting at the given id
    # history
  git.decode_pack
    @ (raw: bytes) -> result[list[pack_entry], string]
    + returns the objects contained in a pack file, resolving deltas
    - returns error when the pack header or checksum is invalid
    # pack
    -> std.compression.zlib_inflate
  git.encode_pack
    @ (objects: list[pack_entry]) -> bytes
    + returns a pack file containing the given objects
    # pack
    -> std.compression.zlib_deflate
  git.clone
    @ (remote_url: string, local_path: string) -> result[void, string]
    + fetches refs and objects via the smart http protocol and writes them locally
    - returns error when the remote cannot be reached
    # transport
    -> std.http.get
    -> std.http.post
  git.fetch
    @ (repo_path: string, remote_url: string) -> result[list[string], string]
    + returns the list of updated ref names
    # transport
    -> std.http.get
    -> std.http.post
