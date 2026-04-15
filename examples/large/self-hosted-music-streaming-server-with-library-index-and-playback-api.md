# Requirement: "a self-hosted music streaming server library with a library index and playback API"

Scans a directory for audio files, reads tags, builds an index, and serves byte-range reads for playback. Playlists live in-memory.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + yields regular file paths under root
      - returns error when root is not a directory
      # filesystem
    std.fs.file_size
      fn (path: string) -> result[i64, string]
      + returns the file size in bytes
      # filesystem
    std.fs.read_range
      fn (path: string, offset: i64, length: i64) -> result[bytes, string]
      + reads length bytes starting at offset
      - returns error when the range is outside the file
      # filesystem
  std.audio
    std.audio.read_tags
      fn (path: string) -> result[audio_tags, string]
      + reads title, artist, album, duration, and bitrate from the file
      - returns error on unsupported container
      # audio

music_server
  music_server.new
    fn (library_root: string) -> server_state
    + initializes a server rooted at library_root with an empty index
    # construction
  music_server.scan
    fn (state: server_state) -> result[server_state, string]
    + walks the library root, reads tags, and populates the index
    - returns error when the root does not exist
    # indexing
    -> std.fs.walk
    -> std.audio.read_tags
  music_server.search_tracks
    fn (state: server_state, query: string) -> list[track]
    + returns tracks whose title, artist, or album contains query (case-insensitive)
    - returns empty list when nothing matches
    # search
  music_server.track_by_id
    fn (state: server_state, id: track_id) -> optional[track]
    + returns the track record for the given id
    - returns none when id is unknown
    # query
  music_server.albums_by_artist
    fn (state: server_state, artist: string) -> list[string]
    + returns unique album names for the artist
    # query
  music_server.stream_range
    fn (state: server_state, id: track_id, offset: i64, length: i64) -> result[bytes, string]
    + returns the requested byte range of the audio file for HTTP range serving
    - returns error when id is unknown
    - returns error when the range is invalid
    # streaming
    -> std.fs.read_range
    -> std.fs.file_size
  music_server.create_playlist
    fn (state: server_state, name: string) -> result[tuple[server_state, playlist_id], string]
    + creates an empty playlist and returns its id
    - returns error when name is empty
    # playlists
  music_server.add_to_playlist
    fn (state: server_state, id: playlist_id, track: track_id) -> result[server_state, string]
    + appends a track to the playlist
    - returns error when either id is unknown
    # playlists
  music_server.remove_from_playlist
    fn (state: server_state, id: playlist_id, position: i32) -> result[server_state, string]
    + removes the track at position
    - returns error when position is out of range
    # playlists
  music_server.playlist_tracks
    fn (state: server_state, id: playlist_id) -> result[list[track], string]
    + returns the ordered tracks in the playlist
    - returns error when id is unknown
    # playlists
