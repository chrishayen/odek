# Projects LiveView BDD

## Scope
Authenticated workspace project shell, project listing/search, create flow, and placeholder navigation routes.

## Scenarios
- Given an authenticated user, when the user visits `/projects` or `/`, then the projects list UI is rendered.
- Given an unauthenticated visitor, when the visitor requests `/projects` or `/`, then access is redirected to log in.
- Given a user has projects in workspace scope, when the page loads, then projects are shown sorted alphabetically by name.
- Given the user types into project search, when the query changes, then the list updates immediately with case-insensitive name matching.
- Given an empty workspace, when no query is active, then a helpful empty-state message and create-project call-to-action are shown.
- Given the user triggers create-project actions from the header, empty state, or grid new card, when activated, then a modal form opens for name and optional description.
- Given a successful create action, when valid name and optional description are submitted, then the user is redirected to `/projects/:project_id`.
- Given navigation placeholders, when the user opens `/rules`, `/skills`, `/agents`, or `/settings`, then a placeholder page renders inside the authenticated workspace shell.

## Invariants
- Routes are defined inside `:require_authenticated_user` pipeline and `live_session :require_authenticated_user`.
- LiveView templates start with `<Layouts.app flash={@flash} current_scope={@current_scope}>`.
- Project collections are rendered via LiveView streams.
