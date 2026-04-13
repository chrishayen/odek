# Requirement: "a library for publishing vector GIS data to a spatial database and map server"

Reads shapefiles, writes them to a spatial SQL store, and registers a layer on a map server.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the file contents
      - returns error when the path cannot be opened
      # filesystem
  std.sql
    std.sql.connect
      @ (connection_string: string) -> result[db_handle, string]
      + returns an open database connection
      - returns error when the connection string is malformed or unreachable
      # database
    std.sql.execute
      @ (db: db_handle, query: string, params: list[string]) -> result[i64, string]
      + returns the number of rows affected
      - returns error when the query is rejected
      # database
  std.http
    std.http.post
      @ (url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + returns the server response
      - returns error when the request cannot be delivered
      # http

gis_publisher
  gis_publisher.decode_shapefile
    @ (data: bytes) -> result[list[feature], string]
    + returns the features parsed from a shapefile byte stream
    - returns error on malformed shapefile
    # parsing
    -> std.fs.read_all
  gis_publisher.create_spatial_table
    @ (db: db_handle, table: string, srid: i32) -> result[void, string]
    + creates a table with a geometry column of the given srid
    - returns error when the table already exists
    # database
    -> std.sql.execute
  gis_publisher.insert_features
    @ (db: db_handle, table: string, features: list[feature]) -> result[i64, string]
    + returns the number of rows inserted
    - returns error when any feature geometry is invalid
    # database
    -> std.sql.execute
  gis_publisher.register_layer
    @ (server_url: string, workspace: string, layer_name: string, table: string) -> result[void, string]
    + publishes a layer definition on the map server pointing at the given table
    - returns error when the server rejects the request
    # publication
    -> std.http.post
  gis_publisher.publish
    @ (shapefile_path: string, db: db_handle, table: string, server_url: string, workspace: string) -> result[void, string]
    + end-to-end: decode, create table, insert features, register layer
    - returns error at the first failing step
    # pipeline
