# Requirement: "a library that renders a bar chart as ASCII text"

Given labels and numeric values, produce a multi-line string where each row is a horizontal bar made of a fill character.

std: (all units exist)

asciichart
  asciichart.render
    @ (labels: list[string], values: list[f64], width: i32) -> result[string, string]
    + returns one line per label with a bar scaled so the maximum value fills `width` columns
    + right-pads labels to a uniform column
    - returns error when labels and values have different lengths
    - returns error when width is less than 1
    ? negative values are treated as zero
    # rendering
  asciichart.render_with_axis
    @ (labels: list[string], values: list[f64], width: i32) -> result[string, string]
    + like render but appends the numeric value at the end of each row
    - returns error when labels and values have different lengths
    # rendering
