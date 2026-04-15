# Requirement: "a library for reading, writing, and querying street-map data"

Street-map data is modeled as nodes, ways, and relations. The library parses an XML-based street-map file, exposes lookups and filters, and can serialize it back.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the file contents as a string
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes contents to a file
      - returns error when the path cannot be written
      # filesystem
  std.xml
    std.xml.parse
      fn (raw: string) -> result[xml_node, string]
      + returns a parsed XML tree
      - returns error on malformed XML
      # xml
    std.xml.serialize
      fn (node: xml_node) -> string
      + serializes an XML tree back to a string
      # xml

streetmap
  streetmap.parse_node
    fn (elem: xml_node) -> result[map_node, string]
    + returns a map node with id, lat, lon, and tags
    - returns error when id, lat, or lon is missing
    # parsing
  streetmap.parse_way
    fn (elem: xml_node) -> result[map_way, string]
    + returns a way with id, ordered node references, and tags
    - returns error when no node references are present
    # parsing
  streetmap.parse_relation
    fn (elem: xml_node) -> result[map_relation, string]
    + returns a relation with id, members, and tags
    - returns error when a member is missing a type or ref
    # parsing
  streetmap.read
    fn (path: string) -> result[street_map, string]
    + reads a street-map file and assembles all features
    - returns error when parsing fails for any element
    # reading
    -> std.fs.read_all
    -> std.xml.parse
  streetmap.node_by_id
    fn (map: street_map, id: i64) -> optional[map_node]
    + returns the node with the given id
    - returns none when no such node exists
    # lookup
  streetmap.way_by_id
    fn (map: street_map, id: i64) -> optional[map_way]
    + returns the way with the given id
    - returns none when no such way exists
    # lookup
  streetmap.filter_by_tag
    fn (map: street_map, key: string, value: string) -> list[map_feature]
    + returns every feature whose tags contain the key=value pair
    # query
  streetmap.bounding_box
    fn (map: street_map) -> optional[bounds]
    + returns min/max lat and lon across all nodes
    - returns none when the map has no nodes
    # geometry
  streetmap.way_length_meters
    fn (map: street_map, way: map_way) -> result[f64, string]
    + returns the great-circle length of a way in meters
    - returns error when any node reference is missing from the map
    # geometry
  streetmap.write
    fn (path: string, map: street_map) -> result[void, string]
    + serializes the map to XML and writes it to disk
    - returns error when the path cannot be written
    # writing
    -> std.xml.serialize
    -> std.fs.write_all
