# Requirement: "a toolkit for biological sequence computation"

FASTA parsing plus the everyday operations a biology user reaches for on DNA strings.

std: (all units exist)

bio
  bio.parse_fasta
    @ (text: string) -> result[list[sequence_record], string]
    + parses FASTA text into records with id, description, and sequence
    - returns error when a sequence block has no header line
    # parsing
  bio.write_fasta
    @ (records: list[sequence_record], line_width: i32) -> string
    + serializes records with sequence lines wrapped to the given width
    # serialization
  bio.complement
    @ (dna: string) -> string
    + returns the per-base complement (A<->T, C<->G, case preserved)
    ? non-ACGT characters pass through unchanged
    # sequence_ops
  bio.reverse_complement
    @ (dna: string) -> string
    + returns the reverse complement of the DNA string
    # sequence_ops
    -> bio.complement
  bio.transcribe
    @ (dna: string) -> string
    + returns the RNA transcript (T -> U)
    # sequence_ops
  bio.translate
    @ (rna: string) -> string
    + translates RNA codons to a single-letter amino-acid string using the standard table
    + stops at the first stop codon
    - returns "" for input shorter than 3 nucleotides
    # translation
  bio.gc_content
    @ (dna: string) -> f64
    + returns the fraction of G and C bases in [0.0, 1.0]
    - returns 0.0 for an empty string
    # statistics
