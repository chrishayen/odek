# Requirement: "a recipe manager with shopping lists and meal planning"

Parses structured recipes, aggregates ingredients into shopping lists, and assigns recipes to a weekly plan.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the full contents of a file as text
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_all
      fn (path: string, content: string) -> result[void, string]
      + writes content to path, replacing any existing file
      # filesystem
  std.json
    std.json.encode
      fn (value: json_value) -> string
      + serializes a json value to a canonical string
      # serialization
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses text into a json value
      - returns error on malformed input
      # serialization

recipes
  recipes.parse_recipe
    fn (text: string) -> result[recipe, string]
    + parses a recipe with title, ingredients (name, quantity, unit), and steps
    - returns error when the title line is missing
    - returns error when an ingredient line is malformed
    # recipe_parsing
  recipes.load_recipe
    fn (path: string) -> result[recipe, string]
    + reads and parses a recipe file
    # recipe_loading
    -> std.fs.read_all
    -> recipes.parse_recipe
  recipes.scale
    fn (r: recipe, factor: f64) -> recipe
    + multiplies every ingredient quantity by factor
    ? step text is not rewritten; only quantities change
    # scaling
  recipes.new_pantry
    fn () -> pantry
    + returns an empty pantry with no ingredients on hand
    # construction
  recipes.set_on_hand
    fn (p: pantry, ingredient: string, quantity: f64, unit: string) -> pantry
    + records or updates the quantity of an ingredient on hand
    # pantry
  recipes.shopping_list
    fn (selected: list[recipe], p: pantry) -> list[shopping_item]
    + sums ingredient quantities by (name, unit) across recipes
    + subtracts pantry quantities and omits items already fully covered
    ? ingredients with mismatched units are reported as separate lines
    # shopping_list
  recipes.new_week_plan
    fn () -> week_plan
    + returns an empty plan with seven day slots
    # construction
  recipes.assign_day
    fn (plan: week_plan, day: i32, r: recipe) -> result[week_plan, string]
    + assigns r to day (0 = Monday .. 6 = Sunday)
    - returns error when day is outside 0..6
    # planning
  recipes.plan_shopping_list
    fn (plan: week_plan, p: pantry) -> list[shopping_item]
    + aggregates every assigned recipe in the plan into one shopping list
    # planning
    -> recipes.shopping_list
  recipes.save_plan
    fn (plan: week_plan, path: string) -> result[void, string]
    + serializes the plan as JSON and writes it to path
    # persistence
    -> std.json.encode
    -> std.fs.write_all
  recipes.load_plan
    fn (path: string) -> result[week_plan, string]
    + reads and deserializes a plan from disk
    - returns error when the file does not contain a valid plan
    # persistence
    -> std.fs.read_all
    -> std.json.parse
