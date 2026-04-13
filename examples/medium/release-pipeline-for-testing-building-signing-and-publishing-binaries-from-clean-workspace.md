# Requirement: "a pipeline that tests, builds, signs, and publishes binaries from a clean workspace"

Each pipeline stage is a small project rune calling thin std utilities. The orchestrator runs the stages in sequence and stops on the first failure.

std
  std.process
    std.process.run
      @ (cmd: string, args: list[string], cwd: string) -> result[i32, string]
      + runs a command in a directory and returns its exit code
      - returns error when the binary is missing
      # process
  std.fs
    std.fs.make_temp_dir
      @ (prefix: string) -> result[string, string]
      + creates an empty temporary directory and returns its absolute path
      # filesystem
    std.fs.remove_tree
      @ (path: string) -> result[void, string]
      + recursively deletes a directory
      # filesystem
  std.crypto
    std.crypto.sign_detached
      @ (private_key: bytes, data: bytes) -> bytes
      + produces a detached signature over the given bytes
      # cryptography

release_pipeline
  release_pipeline.prepare_workspace
    @ (source_dir: string) -> result[string, string]
    + copies the source into a fresh temp directory and returns its path
    - returns error when the source directory is empty
    # workspace
    -> std.fs.make_temp_dir
  release_pipeline.run_tests
    @ (workspace: string, test_cmd: string) -> result[void, string]
    + runs the test command in the workspace
    - returns error when the test command exits non-zero
    # test
    -> std.process.run
  release_pipeline.build_binary
    @ (workspace: string, build_cmd: string, output_path: string) -> result[string, string]
    + runs the build command and returns the path of the produced artifact
    - returns error when the build command fails
    # build
    -> std.process.run
  release_pipeline.sign_artifact
    @ (artifact_path: string, private_key: bytes) -> result[string, string]
    + writes a detached signature file alongside the artifact and returns its path
    # signing
    -> std.crypto.sign_detached
  release_pipeline.publish_artifact
    @ (artifact_path: string, signature_path: string, destination_uri: string) -> result[void, string]
    + uploads the artifact and signature to a pluggable destination
    - returns error when the destination rejects the upload
    # publish
  release_pipeline.run_all
    @ (source_dir: string, test_cmd: string, build_cmd: string, private_key: bytes, destination_uri: string) -> result[void, string]
    + runs prepare, test, build, sign, and publish in order
    + cleans up the temporary workspace on both success and failure
    - stops and returns the first failing stage's error
    # orchestration
    -> std.fs.remove_tree
