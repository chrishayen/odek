# Requirement: "a static site generator with plugin hooks"

Like a plain static site generator, but every stage of the pipeline goes through an ordered plugin list so callers can inject behavior.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the full contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating parent directories
      # filesystem
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns every regular file beneath root
      - returns error when root does not exist
      # filesystem

site_builder
  site_builder.new_document
    fn (source_path: string, raw: bytes) -> document
    + constructs a document with empty metadata and raw content
    # model
  site_builder.register_plugin
    fn (registry: plugin_registry, name: string, hook: plugin_hook) -> plugin_registry
    + appends the plugin to the ordered list
    ? order of registration is preserved; earlier plugins run first
    # registry
  site_builder.run_load_hooks
    fn (registry: plugin_registry, doc: document) -> result[document, string]
    + runs every load hook in order, each receiving the previous output
    - returns error when any hook fails
    # pipeline_load
  site_builder.run_transform_hooks
    fn (registry: plugin_registry, doc: document) -> result[document, string]
    + runs every transform hook in order
    - returns error when any hook fails
    # pipeline_transform
  site_builder.run_render_hooks
    fn (registry: plugin_registry, doc: document) -> result[string, string]
    + runs render hooks; the last one must produce final HTML
    - returns error when no render hook produced output
    # pipeline_render
  site_builder.run_write_hooks
    fn (registry: plugin_registry, doc: document, html: string, out_root: string) -> result[string, string]
    + lets write hooks decide the destination path, defaulting to a mirrored path
    + writes the HTML to the chosen path and returns it
    # pipeline_write
    -> std.fs.write_all
  site_builder.build_one
    fn (registry: plugin_registry, source_path: string, out_root: string) -> result[string, string]
    + runs load, transform, render, and write hooks for one source file
    - returns error when any stage fails
    # build_one
    -> std.fs.read_all
  site_builder.build_all
    fn (registry: plugin_registry, source_root: string, out_root: string) -> result[i32, string]
    + walks source_root and builds every file, returning the count written
    - returns error when the source root is missing
    # build_all
    -> std.fs.walk
