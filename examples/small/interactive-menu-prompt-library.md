# Requirement: "a menu prompt library for interactive applications"

Builds a menu of labeled choices and resolves user input to a selection.

std: (all units exist)

menu
  menu.new
    fn (title: string) -> menu_state
    + creates an empty menu with the given title
    # construction
  menu.add_option
    fn (state: menu_state, key: string, label: string) -> result[menu_state, string]
    + appends a selectable option keyed by a short identifier
    - returns error when key already exists
    # registration
  menu.render
    fn (state: menu_state) -> string
    + returns the menu title followed by numbered options
    # rendering
  menu.select
    fn (state: menu_state, input: string) -> result[string, string]
    + returns the option key matching the entered number or key
    - returns error when input does not match any option
    # selection
