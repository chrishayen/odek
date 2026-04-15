# Requirement: "a chemical toolbox that reads and writes a variety of chemical data formats and exposes a common molecule model"

A cheminformatics library centered on an in-memory molecule graph. Format drivers convert between the model and textual representations; queries operate over the graph.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file contents
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to disk
      - returns error on permission failure
      # filesystem

chem_toolbox
  chem_toolbox.new_molecule
    fn () -> molecule
    + returns an empty molecule with no atoms or bonds
    # construction
  chem_toolbox.add_atom
    fn (mol: molecule, element: string, charge: i32) -> tuple[molecule, atom_id]
    + adds an atom and returns the updated molecule plus the new atom id
    # mutation
  chem_toolbox.add_bond
    fn (mol: molecule, a: atom_id, b: atom_id, order: i32) -> result[molecule, string]
    + adds a bond between two atoms with the given bond order
    - returns error when either atom id is unknown
    - returns error when order is outside 1..3
    # mutation
  chem_toolbox.parse_smiles
    fn (input: string) -> result[molecule, string]
    + returns the molecule represented by a SMILES string
    - returns error on malformed input
    # format_driver
  chem_toolbox.render_smiles
    fn (mol: molecule) -> string
    + returns a canonical SMILES string for the molecule
    # format_driver
  chem_toolbox.parse_molfile
    fn (raw: string) -> result[molecule, string]
    + returns the molecule represented by an MDL molfile
    - returns error on malformed input
    # format_driver
  chem_toolbox.render_molfile
    fn (mol: molecule) -> string
    + returns an MDL molfile for the molecule
    # format_driver
  chem_toolbox.load_file
    fn (path: string, format: string) -> result[molecule, string]
    + loads and parses a molecule from disk, dispatching on format
    - returns error when the format is unknown
    # io
    -> std.fs.read_all
  chem_toolbox.save_file
    fn (path: string, mol: molecule, format: string) -> result[void, string]
    + writes the molecule to disk in the given format
    - returns error when the format is unknown
    # io
    -> std.fs.write_all
  chem_toolbox.molecular_formula
    fn (mol: molecule) -> string
    + returns the Hill-ordered molecular formula
    # analysis
  chem_toolbox.molecular_weight
    fn (mol: molecule) -> f64
    + returns the average molecular weight in Daltons
    # analysis
  chem_toolbox.atom_count
    fn (mol: molecule, element: string) -> i32
    + returns the number of atoms of the given element
    # analysis
  chem_toolbox.substructure_match
    fn (mol: molecule, pattern: molecule) -> list[list[atom_id]]
    + returns every mapping of the pattern's atoms onto atoms in mol
    + returns empty list when no match exists
    # search
  chem_toolbox.canonicalize
    fn (mol: molecule) -> molecule
    + returns the molecule with atoms reordered into canonical form
    ? canonical order is used by render_smiles to produce deterministic output
    # normalization
