# Requirement: "a cache-friendly 2d grid library with pathfinding, observers, and import/export"

A tile grid backed by a flat array of u32 cells, with A* pathfinding, tile-change observers, and binary import/export.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file fully into memory
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to a file
      # filesystem

tilegrid
  tilegrid.new
    fn (width: u32, height: u32) -> grid_state
    + creates a grid with all cells initialized to zero
    ? cells are stored in a single contiguous buffer in row-major order
    # construction
  tilegrid.get
    fn (grid: grid_state, x: u32, y: u32) -> result[u32, string]
    + returns the cell value at (x, y)
    - returns error when the coordinates are out of range
    # access
  tilegrid.set
    fn (grid: grid_state, x: u32, y: u32, value: u32) -> result[grid_state, string]
    + returns a grid with the cell at (x, y) set to value and notifies observers
    - returns error when the coordinates are out of range
    # mutation
  tilegrid.add_observer
    fn (grid: grid_state, observer: fn(u32, u32, u32, u32) -> void) -> grid_state
    + registers a callback fired on every set, called with (x, y, old, new)
    # observers
  tilegrid.neighbors4
    fn (grid: grid_state, x: u32, y: u32) -> list[tuple[u32, u32]]
    + returns the in-bounds orthogonal neighbors of (x, y)
    # topology
  tilegrid.find_path
    fn (grid: grid_state, start: tuple[u32, u32], goal: tuple[u32, u32], passable: fn(u32) -> bool) -> optional[list[tuple[u32, u32]]]
    + returns an A* path from start to goal using passable to test each cell
    - returns none when no path exists
    - returns none when start or goal is impassable
    # pathfinding
  tilegrid.export_binary
    fn (grid: grid_state, path: string) -> result[void, string]
    + writes the grid to a file as width, height, and flat cell data
    # io
    -> std.fs.write_all
  tilegrid.import_binary
    fn (path: string) -> result[grid_state, string]
    + loads a grid previously written by export_binary
    - returns error on truncated or malformed input
    # io
    -> std.fs.read_all
