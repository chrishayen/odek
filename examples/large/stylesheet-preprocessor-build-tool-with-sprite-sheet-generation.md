# Requirement: "a stylesheet preprocessor build tool that compiles sources and generates sprite sheets"

A small build pipeline: discover stylesheet sources, compile them through an external compiler, and produce sprite sheets from referenced image folders.

std
  std.fs
    std.fs.walk_dir
      fn (root: string) -> result[list[string], string]
      + returns every file path under root
      - returns error when root does not exist
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file's full content
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data atomically, replacing any existing file
      # filesystem
  std.path
    std.path.join
      fn (parts: list[string]) -> string
      + joins path segments using the platform separator
      # path
    std.path.ext
      fn (path: string) -> string
      + returns the extension including the leading dot
      # path
  std.image
    std.image.decode
      fn (data: bytes) -> result[image, string]
      + decodes PNG, JPEG, or GIF
      - returns error on unsupported format
      # image
    std.image.encode_png
      fn (img: image) -> bytes
      + encodes an image as PNG
      # image
    std.image.compose
      fn (canvas: image, sprite: image, x: i32, y: i32) -> image
      + draws sprite onto canvas at the given offset
      # image
    std.image.new_canvas
      fn (width: i32, height: i32) -> image
      + allocates a blank RGBA canvas
      # image
  std.process
    std.process.run
      fn (cmd: string, args: list[string], workdir: string) -> result[string, string]
      + runs an external command and returns stdout
      - returns error on non-zero exit
      # process

style_build
  style_build.discover_sources
    fn (root: string, extensions: list[string]) -> result[list[string], string]
    + returns files under root whose extension is in the allowed list
    - returns error when root does not exist
    # discovery
    -> std.fs.walk_dir
    -> std.path.ext
  style_build.compile_source
    fn (compiler: string, input_path: string, output_path: string) -> result[void, string]
    + invokes the compiler on one source and writes the result
    - returns error when the compiler exits non-zero
    # compilation
    -> std.process.run
    -> std.fs.write_all
  style_build.compile_all
    fn (compiler: string, sources: list[string], out_dir: string) -> result[i32, string]
    + compiles each source into out_dir and returns how many succeeded
    - returns error on the first compilation failure
    # compilation
    -> std.path.join
  style_build.load_sprite_set
    fn (dir: string) -> result[list[image], string]
    + reads every image in the directory into memory
    - returns error when any image cannot be decoded
    # sprite
    -> std.fs.walk_dir
    -> std.fs.read_all
    -> std.image.decode
  style_build.pack_sprites
    fn (images: list[image]) -> tuple[image, list[tuple[string,i32,i32]]]
    + tiles images left-to-right into a canvas and returns the canvas plus per-image offsets
    # sprite
    -> std.image.new_canvas
    -> std.image.compose
  style_build.write_sprite_sheet
    fn (canvas: image, out_path: string) -> result[void, string]
    + encodes the canvas as PNG and writes it to disk
    # sprite
    -> std.image.encode_png
    -> std.fs.write_all
  style_build.emit_sprite_map
    fn (offsets: list[tuple[string,i32,i32]]) -> string
    + produces stylesheet rules mapping each sprite name to its pixel offset in the sheet
    # sprite
  style_build.build
    fn (src_root: string, sprite_root: string, out_dir: string, compiler: string) -> result[void, string]
    + runs the full pipeline: discover, compile, pack sprites, write sheet and map
    - returns error on any step's failure
    # pipeline
