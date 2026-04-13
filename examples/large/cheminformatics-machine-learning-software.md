# Requirement: "a cheminformatics and machine learning library"

A core molecule representation plus a few canonical cheminformatics operations and a minimal feature extraction step that feeds a learned model. The project layer is kept thin by pushing parsing, linear algebra, and hashing into std primitives.

std
  std.io
    std.io.read_all
      @ (path: string) -> result[bytes, string]
      + reads the whole file into memory
      - returns error when the path does not exist
      # io
  std.math
    std.math.dot
      @ (a: list[f64], b: list[f64]) -> f64
      + returns sum of pairwise products
      ? vectors must be the same length
      # linear_algebra
    std.math.sigmoid
      @ (x: f64) -> f64
      + returns 1 / (1 + exp(-x))
      # activation
  std.hash
    std.hash.fnv1a_32
      @ (data: bytes) -> u32
      + returns a 32-bit FNV-1a hash
      # hashing

chem
  chem.parse_smiles
    @ (smiles: string) -> result[molecule, string]
    + parses a SMILES string into an atom/bond graph
    - returns error on unbalanced ring closures
    - returns error on unknown element symbols
    # parsing
  chem.canonical_smiles
    @ (mol: molecule) -> string
    + returns a canonical SMILES representation with deterministic atom ordering
    # canonicalization
  chem.morgan_fingerprint
    @ (mol: molecule, radius: i32, bits: i32) -> list[bool]
    + returns a folded circular fingerprint of the given bit length
    ? each atom environment is hashed and folded into the bit vector
    # featurization
    -> std.hash.fnv1a_32
  chem.tanimoto
    @ (a: list[bool], b: list[bool]) -> f64
    + returns the Tanimoto similarity between two fingerprints
    - returns 0.0 when both fingerprints are all zero
    # similarity
  chem.mol_weight
    @ (mol: molecule) -> f64
    + returns the sum of standard atomic weights for all atoms
    # descriptors
  chem.logistic_model_new
    @ (weights: list[f64], bias: f64) -> logistic_model
    + constructs a logistic regression model from learned parameters
    # model_construction
  chem.logistic_model_predict
    @ (model: logistic_model, features: list[f64]) -> f64
    + returns a probability in [0, 1] for the input feature vector
    # prediction
    -> std.math.dot
    -> std.math.sigmoid
  chem.fingerprint_to_features
    @ (fp: list[bool]) -> list[f64]
    + returns a float vector where each bit maps to 1.0 or 0.0
    # feature_conversion
  chem.load_sdf
    @ (path: string) -> result[list[molecule], string]
    + reads and parses an SDF file into a list of molecules
    - returns error when the file is missing or malformed
    # dataset_loading
    -> std.io.read_all
