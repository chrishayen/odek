# Requirement: "a library for parsing, filtering, transforming, and building streaming media manifests"

Supports two streaming manifest formats through a common neutral manifest model.

std: (all units exist)

manifest
  manifest.parse_playlist
    fn (raw: string) -> result[manifest_model, string]
    + parses a text-based playlist manifest into the neutral model
    - returns error on malformed directives
    # parsing
  manifest.parse_adaptive
    fn (raw: string) -> result[manifest_model, string]
    + parses an XML-based adaptive streaming manifest into the neutral model
    - returns error on invalid XML or missing required attributes
    # parsing
  manifest.parse
    fn (raw: string) -> result[manifest_model, string]
    + detects the format and dispatches to the appropriate parser
    - returns error when the format cannot be detected
    # parsing
    -> manifest.parse_playlist
    -> manifest.parse_adaptive
  manifest.filter_tracks
    fn (model: manifest_model, predicate: track_predicate) -> manifest_model
    + returns a manifest keeping only tracks matching the predicate
    # filtering
  manifest.map_segments
    fn (model: manifest_model, transform: segment_transform) -> manifest_model
    + applies a transform to every segment in every track
    # transformation
  manifest.set_base_url
    fn (model: manifest_model, base: string) -> manifest_model
    + rewrites all relative URIs to be absolute under the given base
    # transformation
  manifest.build_playlist
    fn (model: manifest_model) -> string
    + emits a text-based playlist from the neutral model
    # building
  manifest.build_adaptive
    fn (model: manifest_model) -> string
    + emits an XML-based adaptive streaming manifest from the neutral model
    # building
