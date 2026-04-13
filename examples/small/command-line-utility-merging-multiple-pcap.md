# Requirement: "a library for merging packet capture files into a single file ordered by timestamp"

Reads records from multiple capture files and emits them in ascending timestamp order.

std
  std.fs
    std.fs.open_read
      @ (path: string) -> result[file_handle, string]
      + opens a file for sequential reads
      - returns error when the file cannot be opened
      # filesystem
    std.fs.open_write
      @ (path: string) -> result[file_handle, string]
      + opens a file for sequential writes, truncating any existing content
      # filesystem

pcap_merge
  pcap_merge.read_header
    @ (f: file_handle) -> result[pcap_header, string]
    + parses the global capture header
    - returns error on a bad magic number
    # parsing
  pcap_merge.read_next_record
    @ (f: file_handle) -> result[optional[pcap_record], string]
    + returns the next record or none at end-of-file
    - returns error on a truncated record
    # parsing
  pcap_merge.merge
    @ (inputs: list[string], output: string) -> result[i64, string]
    + writes every record from every input file into output in ascending timestamp order
    + returns the total number of records written
    - returns error when input headers disagree on link-layer type
    # merging
    -> std.fs.open_read
    -> std.fs.open_write
    -> pcap_merge.read_header
    -> pcap_merge.read_next_record
