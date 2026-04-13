# Requirement: "an opinionated zero-config static site generator using component templates and markdown-with-components"

Walks a source tree, renders markdown-with-components pages through a component template engine, writes static HTML to an output directory.

std
  std.fs
    std.fs.walk_files
      @ (root: string) -> result[list[string], string]
      + returns all file paths under root recursively
      - returns error when root does not exist
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when file is missing or unreadable
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + creates parent directories as needed and writes contents
      - returns error when the path is not writable
      # filesystem
  std.path
    std.path.replace_extension
      @ (path: string, new_ext: string) -> string
      + replaces the final extension on a path with new_ext
      + appends new_ext when the path has no extension
      # paths
    std.path.join
      @ (a: string, b: string) -> string
      + joins two path segments with a single separator
      # paths

ssg
  ssg.parse_component_template
    @ (source: string) -> result[component_template, string]
    + parses a component-template source string into a template tree
    - returns error on unclosed tags
    # template_parsing
  ssg.parse_markdown_with_components
    @ (source: string) -> result[page_ast, string]
    + parses markdown text with embedded component tags into a page AST
    - returns error when a component tag has no matching close
    # content_parsing
  ssg.render_page
    @ (page: page_ast, layout: component_template) -> result[string, string]
    + returns the rendered HTML for a page wrapped in the layout
    - returns error when the page references a component the layout cannot resolve
    # rendering
  ssg.build_site
    @ (src_dir: string, out_dir: string) -> result[i32, string]
    + walks src_dir, renders every content file, and writes HTML into out_dir mirroring the structure
    + returns the number of pages written
    - returns error when src_dir is missing
    # orchestration
    -> std.fs.walk_files
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.path.replace_extension
    -> std.path.join
