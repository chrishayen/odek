# Requirement: "a library for parsing benchmark results and producing a multi-axis visualization dataset"

Parse text-format benchmark output and shape it for plotting against multiple dimensions.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns file contents
      - returns error when unreadable
      # filesystem

benchviz
  benchviz.parse_results
    @ (raw: string) -> result[list[bench_sample], string]
    + returns one sample per bench line with name, iterations, ns_per_op, bytes_per_op, allocs_per_op
    - returns error when a line is malformed
    # parsing
  benchviz.load_file
    @ (path: string) -> result[list[bench_sample], string]
    + reads and parses a file of benchmark output
    - returns error when the file is unreadable
    # loading
    -> std.fs.read_all
  benchviz.group_by_name
    @ (samples: list[bench_sample]) -> map[string, list[bench_sample]]
    + groups samples by benchmark name
    # grouping
  benchviz.to_4d_points
    @ (samples: list[bench_sample]) -> list[bench_point]
    + returns one point per sample with axes (ns_per_op, bytes_per_op, allocs_per_op, iterations)
    # projection
  benchviz.summarize
    @ (samples: list[bench_sample]) -> bench_summary
    + returns mean and stddev for each numeric axis
    # summary
