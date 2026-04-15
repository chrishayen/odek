# Requirement: "a Wilson-score-based average rating calculator"

Pure math. One function that returns the lower bound of a Wilson confidence interval given positive and negative vote counts.

std: (all units exist)

wilson_rating
  wilson_rating.lower_bound
    fn (positive: i64, negative: i64, confidence: f64) -> f64
    + returns the Wilson score lower bound for the given votes and confidence
    + returns 0.0 when positive and negative are both 0
    ? uses the standard normal inverse for the confidence level
    # ranking
