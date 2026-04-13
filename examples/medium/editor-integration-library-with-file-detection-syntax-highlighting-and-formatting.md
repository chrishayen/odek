# Requirement: "an editor integration library for file-type detection, syntax highlighting, and source formatting"

Detects a file's language from its name or contents, tokenizes it for highlighting, and invokes an external formatter.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the file contents as a string
      - returns error when the file cannot be read
      # filesystem
  std.proc
    std.proc.run_command
      @ (program: string, args: list[string], stdin: string) -> result[command_result, string]
      + runs an external program and returns its stdout, stderr, and exit code
      - returns error when the program cannot be launched
      # process

editor_lang
  editor_lang.detect_from_path
    @ (path: string) -> optional[string]
    + returns the language id when the file extension matches a known mapping
    - returns none when no mapping matches
    # detection
  editor_lang.detect_from_content
    @ (text: string) -> optional[string]
    + inspects shebang and first-line hints and returns a language id when recognized
    # detection
  editor_lang.tokenize
    @ (language: string, text: string) -> result[list[highlight_span], string]
    + returns a list of (start, end, token_kind) spans suitable for highlighting
    - returns error when the language is not supported
    # highlighting
  editor_lang.render_ansi
    @ (text: string, spans: list[highlight_span]) -> string
    + emits the text with ANSI color escapes for each span
    # highlighting
  editor_lang.format_file
    @ (path: string, formatter: string) -> result[string, string]
    + runs the external formatter on the file contents and returns the formatted text
    - returns error when the formatter exits non-zero
    # formatting
    -> std.fs.read_all
    -> std.proc.run_command
