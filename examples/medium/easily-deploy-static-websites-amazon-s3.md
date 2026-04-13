# Requirement: "a static website deployment library targeting an object store"

Walks a local directory, computes hashes, compares against remote object metadata, uploads changed files, and deletes orphaned remote objects.

std
  std.fs
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns the list of regular file paths under root, recursively
      - returns error when the root does not exist
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire file contents
      - returns error on missing or unreadable files
      # filesystem
  std.crypto
    std.crypto.md5
      @ (data: bytes) -> string
      + returns the md5 digest as a lowercase hex string
      # hashing
  std.mime
    std.mime.guess
      @ (path: string) -> string
      + returns a content type guess based on the file extension
      + returns "application/octet-stream" when unknown
      # mime

site_deploy
  site_deploy.plan
    @ (root: string, remote_index: map[string, string]) -> result[deploy_plan, string]
    + computes the set of files to upload and objects to delete
    + a file is uploaded when its local hash differs from the remote hash
    + a remote object is deleted when no local file corresponds to it
    - returns error when the root cannot be walked
    # planning
    -> std.fs.walk
    -> std.fs.read_all
    -> std.crypto.md5
  site_deploy.upload_file
    @ (local_path: string, remote_key: string, content_type: string) -> result[void, string]
    + uploads the file under the given key with the given content type
    - returns error when the read or upload fails
    # upload
    -> std.fs.read_all
    -> std.mime.guess
  site_deploy.execute
    @ (plan: deploy_plan) -> result[deploy_report, string]
    + performs every upload and delete in the plan and returns counts
    - returns error on the first failed operation
    # execution
