# Requirement: "a graph-based retrieval augmented generation framework that transforms documents into knowledge graphs"

Ingests documents, extracts entities and relationships, builds a graph, and answers queries by combining graph traversal with vector retrieval.

std
  std.text
    std.text.split
      fn (s: string, sep: string) -> list[string]
      + splits a string on a separator
      # text
    std.text.ascii_lower
      fn (s: string) -> string
      + lowercases ASCII letters
      # text
  std.math
    std.math.cosine_similarity
      fn (a: list[f32], b: list[f32]) -> f32
      + returns the cosine similarity of two equal-length vectors
      + returns 0.0 on either vector being all-zero
      # math
  std.hash
    std.hash.fnv1a_64
      fn (data: bytes) -> u64
      + returns the FNV-1a 64-bit hash
      # hashing

graphrag
  graphrag.graph_new
    fn () -> graph_state
    + returns an empty graph with no nodes or edges
    # construction
  graphrag.chunk_document
    fn (text: string, max_tokens: i32) -> list[string]
    + splits a document into overlapping chunks of roughly max_tokens words
    # ingestion
    -> std.text.split
  graphrag.extract_entities
    fn (chunk: string, entity_extractor_id: string) -> list[entity]
    + runs the named extractor on a chunk, returning surface form, type, and offset
    # extraction
    -> std.text.ascii_lower
  graphrag.extract_relations
    fn (chunk: string, entities: list[entity], relation_extractor_id: string) -> list[relation]
    + returns (subject, predicate, object) triples grounded in the chunk
    # extraction
  graphrag.add_entity
    fn (graph: graph_state, ent: entity) -> tuple[u64, graph_state]
    + merges the entity into the graph, returning a stable node id
    # graph_construction
    -> std.hash.fnv1a_64
  graphrag.add_relation
    fn (graph: graph_state, subj_id: u64, pred: string, obj_id: u64) -> graph_state
    + adds a directed edge between two nodes
    # graph_construction
  graphrag.embed_chunk
    fn (chunk: string, embedder_id: string) -> list[f32]
    + returns a dense vector embedding for the chunk
    # embeddings
  graphrag.index_embedding
    fn (graph: graph_state, node_id: u64, embedding: list[f32]) -> graph_state
    + attaches an embedding to a node for vector search
    # indexing
  graphrag.vector_search
    fn (graph: graph_state, query_embedding: list[f32], top_k: i32) -> list[tuple[u64, f32]]
    + returns the top-k nodes ranked by cosine similarity to the query
    # retrieval
    -> std.math.cosine_similarity
  graphrag.neighbors
    fn (graph: graph_state, node_id: u64, hops: i32) -> list[u64]
    + returns nodes reachable within the given number of hops
    # traversal
  graphrag.query
    fn (graph: graph_state, question: string, embedder_id: string, top_k: i32) -> list[tuple[u64, f32]]
    + embeds the question, finds top-k seed nodes, and expands to their neighborhood
    # query
    -> std.math.cosine_similarity
  graphrag.build_context
    fn (graph: graph_state, node_ids: list[u64]) -> string
    + concatenates entity descriptions and relations for the given nodes into a prompt-ready context string
    # context_assembly
