# Requirement: "an offline-first Markdown editor library"

Document model with load, edit, render-to-HTML, and save operations. No UI; callers drive the model.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads the full file as UTF-8
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes contents to path, overwriting
      - returns error when the parent directory is missing
      # filesystem

md_editor
  md_editor.new_document
    fn () -> doc_state
    + returns an empty document
    # construction
  md_editor.open
    fn (path: string) -> result[doc_state, string]
    + loads a Markdown file from disk
    - returns error when the file does not exist
    # io
    -> std.fs.read_all
  md_editor.save
    fn (state: doc_state, path: string) -> result[doc_state, string]
    + writes the current contents to path and marks the document clean
    - returns error when the parent directory is missing
    # io
    -> std.fs.write_all
  md_editor.replace_range
    fn (state: doc_state, start: i32, end: i32, text: string) -> doc_state
    + replaces the byte range [start, end) with text and marks the document dirty
    # editing
  md_editor.insert_at
    fn (state: doc_state, offset: i32, text: string) -> doc_state
    + inserts text at offset and marks the document dirty
    # editing
  md_editor.undo
    fn (state: doc_state) -> doc_state
    + reverts the most recent edit
    + is a no-op when the undo stack is empty
    # history
  md_editor.redo
    fn (state: doc_state) -> doc_state
    + re-applies the most recently undone edit
    + is a no-op when the redo stack is empty
    # history
  md_editor.render_html
    fn (state: doc_state) -> string
    + returns the document rendered to HTML
    # rendering
  md_editor.is_dirty
    fn (state: doc_state) -> bool
    + returns true when there are unsaved changes
    # status
