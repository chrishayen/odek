# Requirement: "a source code line counter with complexity and COCOMO estimates"

A per-file analyzer that classifies each line as code, comment, or blank, counts simple branch points for complexity, and applies the Basic COCOMO formula over totals.

std: (all units exist)

scc
  scc.language_for_extension
    @ (extension: string) -> optional[language_rules]
    + returns comment syntax and branch keywords for a known extension
    - returns none for unknown extensions
    # detection
  scc.count_lines
    @ (source: string, rules: language_rules) -> line_counts
    + returns counts of code, comment, and blank lines
    + treats a line that starts with a line-comment marker as comment
    + treats lines inside a block comment as comment until the closer
    # counting
  scc.complexity
    @ (source: string, rules: language_rules) -> i32
    + returns the count of branch keywords across code lines
    + returns 0 when there are no branch keywords
    # analysis
  scc.cocomo_basic
    @ (total_code_lines: i64) -> cocomo_estimate
    + returns person-months as 2.4 * (kloc ^ 1.05) and schedule-months as 2.5 * (pm ^ 0.38)
    + returns a zero estimate when total_code_lines is zero
    # estimation
  scc.summarize
    @ (files: list[file_counts]) -> summary
    + returns totals across all files grouped by language
    + includes code, comment, blank, complexity, and file counts
    # aggregation
