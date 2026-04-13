# Requirement: "a repository template describing a default project layout"

Returns the canonical file list an empty project template would contain.

std: (all units exist)

repo_template
  repo_template.default_files
    @ () -> list[string]
    + returns the standard set of template file paths
    ? the list is hardcoded; no parameters
    # template
