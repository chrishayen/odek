defmodule ValkyrieWeb.ProjectsLive.Index do
  use ValkyrieWeb, :live_view

  import ValkyrieWeb.WorkspaceComponents

  alias Ecto.Changeset
  alias Valkyrie.Workspace
  alias ValkyrieWeb.WorkspaceScope

  @project_form_types %{name: :string, description: :string}

  @impl true
  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope}>
      <.workspace_shell current_scope={@current_scope} active_nav={:projects}>
        <:actions>
          <button
            id="new-project-button"
            type="button"
            phx-click="toggle-create"
            class="workspace-primary-button"
          >
            <.icon name="hero-plus" class="size-4" />
            <span>New Project</span>
          </button>
        </:actions>

        <section id="projects-page">
          <header class="projects-header">
            <h1 class="projects-title">Projects</h1>
            <p class="projects-subtitle">
              {@total_projects_count} {pluralize_project(@total_projects_count)} in your workspace
            </p>
          </header>

          <div class="projects-tools">
            <div class="projects-tools-row">
              <.form
                for={@search_form}
                id="projects-search-form"
                phx-change="search"
                phx-submit="search"
                class="projects-search-shell"
              >
                <.icon name="hero-magnifying-glass" class="projects-search-icon size-4" />
                <.input
                  field={@search_form[:q]}
                  id="projects-search-input"
                  class="projects-search-input"
                  placeholder="Search projects..."
                  autocomplete="off"
                  phx-debounce="150"
                />
              </.form>

              <div class="projects-filter-tabs">
                <button type="button" class="projects-filter-tab projects-filter-tab-active">
                  All
                </button>
                <button type="button" class="projects-filter-tab" disabled>Active</button>
                <button type="button" class="projects-filter-tab" disabled>Archived</button>
              </div>
            </div>
          </div>

          <section
            :if={@show_create_form?}
            id="project-create-modal"
            class="project-create-modal"
            phx-window-keydown="toggle-create"
            phx-key="escape"
          >
            <% project_name = project_form_value(@project_form, :name) %>
            <% project_name_length = String.length(project_name) %>

            <button
              type="button"
              class="project-create-backdrop"
              phx-click="toggle-create"
              aria-label="Close create project modal"
            >
            </button>

            <div
              id="project-create-dialog"
              class="project-create-dialog"
              role="dialog"
              aria-modal="true"
              aria-labelledby="project-create-title"
            >
              <div class="project-create-panel">
                <div class="project-create-header">
                  <div class="project-create-header-left">
                    <h2 id="project-create-title" class="project-create-title">New project</h2>
                    <p class="project-create-subtitle">Add a project to your workspace</p>
                  </div>
                  <button
                    id="close-project-create"
                    type="button"
                    class="project-create-close"
                    phx-click="toggle-create"
                    aria-label="Close create project modal"
                  >
                    <.icon name="hero-x-mark" class="size-4" />
                  </button>
                </div>

                <.form
                  for={@project_form}
                  id="project-create-form"
                  phx-change="validate-project"
                  phx-submit="create-project"
                  class="project-create-form"
                >
                  <div class="project-create-field">
                    <div class="project-create-field-label">
                      <label class="project-create-label-text" for="project-name-input">
                        Project name
                      </label>
                      <span
                        :if={project_name_length > 0}
                        class={[
                          "project-create-char-count",
                          project_name_length > 140 && "project-create-char-count-warn"
                        ]}
                      >
                        {project_name_length}/160
                      </span>
                    </div>

                    <.input
                      field={@project_form[:name]}
                      id="project-name-input"
                      placeholder="Phoenix migration runner"
                      class="project-create-input"
                      maxlength="160"
                      required
                    />

                    <div class={[
                      "project-create-slug-row",
                      project_name_length > 0 && "project-create-slug-row-visible"
                    ]}>
                      <span class="project-create-slug-label">slug</span>
                      <span class="project-create-slug-value">{project_slug(project_name)}</span>
                    </div>
                  </div>

                  <div class="project-create-field">
                    <div class="project-create-field-label">
                      <label class="project-create-label-text" for="project-description-input">
                        Description
                      </label>
                      <span class="project-create-label-hint">optional</span>
                    </div>

                    <.input
                      field={@project_form[:description]}
                      id="project-description-input"
                      type="textarea"
                      placeholder="What this project is about and when it's done"
                      class="project-create-input project-create-textarea"
                      maxlength="4000"
                      rows="4"
                    />
                  </div>

                  <div class="project-create-actions">
                    <div class="project-create-shortcut-hint">
                      <span class="project-create-kbd">Ctrl</span>
                      <span>+</span>
                      <span class="project-create-kbd">Enter</span>
                    </div>

                    <div class="project-create-actions-right">
                      <button
                        id="cancel-project-create"
                        type="button"
                        class="project-create-cancel"
                        phx-click="toggle-create"
                      >
                        Cancel
                      </button>
                      <button
                        id="create-project-submit"
                        type="submit"
                        class="project-create-submit"
                        phx-disable-with="Creating..."
                        disabled={project_name == ""}
                      >
                        Create project
                      </button>
                    </div>
                  </div>
                </.form>
              </div>
            </div>
          </section>

          <div id="projects-grid" phx-update="stream" class="projects-grid">
            <div id="projects-empty" class="hidden only:block projects-empty">
              <h2 class="projects-empty-title">
                <%= if @search_query == "" do %>
                  No projects yet
                <% else %>
                  No matches for "{@search_query}"
                <% end %>
              </h2>
              <p class="projects-empty-copy">
                <%= if @search_query == "" do %>
                  Start by creating your first project. You can add details now and refine them later.
                <% else %>
                  Try a different name or clear your search to browse every project in this workspace.
                <% end %>
              </p>
              <div class="projects-empty-actions">
                <button
                  :if={@search_query == "" && !@show_create_form?}
                  id="empty-state-new-project"
                  type="button"
                  phx-click="toggle-create"
                  class="workspace-primary-button"
                >
                  <.icon name="hero-plus" class="size-4" /> Create New Project
                </button>
                <button
                  :if={@search_query != ""}
                  id="clear-project-search"
                  type="button"
                  class="workspace-ghost-button"
                  phx-click="clear-search"
                >
                  Clear search
                </button>
              </div>
            </div>

            <article :for={{dom_id, project} <- @streams.projects} id={dom_id} class="project-card">
              <div class="project-card-header">
                <div class="project-card-avatar" style={project_avatar_style(project.name)}>
                  {project_initials(project.name)}
                </div>
                <div class="project-card-identity">
                  <h2 class="project-card-title">{project.name}</h2>
                  <p class="project-card-slug">{project_slug(project.name)}</p>
                </div>
              </div>

              <p class={[
                "project-card-description",
                blank_text?(project.definition_of_done) && "project-card-description-empty"
              ]}>
                <%= if blank_text?(project.definition_of_done) do %>
                  No description yet.
                <% else %>
                  {project.definition_of_done}
                <% end %>
              </p>

              <div class="project-card-footer">
                <div class="project-card-stats">
                  <span class="project-card-stat">No stories yet</span>
                  <span class="project-card-stat">{format_relative_date(project.updated_at)}</span>
                </div>
                <.link
                  id={"project-link-#{project.id}"}
                  navigate={~p"/projects/#{project.id}"}
                  class="project-card-link"
                >
                  <.icon name="hero-arrow-right" class="size-4" />
                </.link>
              </div>
            </article>

            <button
              :if={@search_query == "" && !@show_create_form?}
              id="projects-grid-new-card"
              type="button"
              class="project-card project-card-new"
              phx-click="toggle-create"
            >
              <div class="project-card-new-inner">
                <span class="project-card-new-icon">
                  <.icon name="hero-plus" class="size-5" />
                </span>
                <span class="project-card-new-text">Create New Project</span>
              </div>
            </button>
          </div>
        </section>
      </.workspace_shell>
    </Layouts.app>
    """
  end

  @impl true
  def mount(_params, _session, socket) do
    workspace_scope = WorkspaceScope.resolve(socket.assigns.current_scope)
    projects = list_projects(workspace_scope)

    socket =
      socket
      |> assign(:workspace_scope, workspace_scope)
      |> assign(:all_projects, projects)
      |> assign(:search_query, "")
      |> assign(:total_projects_count, length(projects))
      |> assign(:show_create_form?, false)
      |> assign(:search_form, search_form(""))
      |> assign(:project_form, project_form(%{}))
      |> stream(:projects, projects, reset: true)

    {:ok, socket}
  end

  @impl true
  def handle_event("search", %{"search" => %{"q" => query}}, socket) do
    query = normalize_text(query) || ""
    filtered = filter_projects(socket.assigns.all_projects, query)

    {:noreply,
     socket
     |> assign(:search_query, query)
     |> assign(:search_form, search_form(query))
     |> stream(:projects, filtered, reset: true)}
  end

  def handle_event("clear-search", _params, socket) do
    {:noreply,
     socket
     |> assign(:search_query, "")
     |> assign(:search_form, search_form(""))
     |> stream(:projects, socket.assigns.all_projects, reset: true)}
  end

  def handle_event("toggle-create", _params, socket) do
    show_create_form? = !socket.assigns.show_create_form?

    {:noreply,
     socket
     |> assign(:show_create_form?, show_create_form?)
     |> assign(
       :project_form,
       if(show_create_form?, do: project_form(%{}), else: socket.assigns.project_form)
     )}
  end

  def handle_event("validate-project", %{"project" => project_params}, socket) do
    {:noreply,
     assign(
       socket,
       :project_form,
       project_params
       |> project_changeset()
       |> Map.put(:action, :validate)
       |> to_form(as: :project)
     )}
  end

  def handle_event("create-project", %{"project" => project_params}, socket) do
    case socket.assigns.workspace_scope do
      nil ->
        {:noreply,
         socket
         |> put_flash(:error, "You need an organization membership before creating projects.")
         |> assign(:show_create_form?, false)}

      workspace_scope ->
        form_changeset =
          project_params
          |> project_changeset()
          |> Map.put(:action, :insert)

        if form_changeset.valid? do
          attrs = %{
            "name" => Changeset.get_field(form_changeset, :name),
            "definition_of_done" => Changeset.get_field(form_changeset, :description) || ""
          }

          case Workspace.create_project(workspace_scope, attrs) do
            {:ok, project} ->
              {:noreply,
               socket
               |> put_flash(:info, "Project created.")
               |> push_navigate(to: ~p"/projects/#{project.id}")}

            {:error, %Ecto.Changeset{} = workspace_changeset} ->
              {:noreply,
               assign(
                 socket,
                 :project_form,
                 workspace_error_form(form_changeset, workspace_changeset)
               )}

            {:error, :invalid_request} ->
              {:noreply,
               put_flash(socket, :error, "An organization is required to create projects.")}

            {:error, :forbidden} ->
              {:noreply,
               put_flash(
                 socket,
                 :error,
                 "You are not allowed to create projects in this workspace."
               )}
          end
        else
          {:noreply, assign(socket, :project_form, to_form(form_changeset, as: :project))}
        end
    end
  end

  defp workspace_error_form(form_changeset, workspace_changeset) do
    merged =
      Enum.reduce(workspace_changeset.errors, form_changeset, fn
        {:name, {message, _opts}}, acc -> Changeset.add_error(acc, :name, message)
        _error, acc -> acc
      end)

    to_form(Map.put(merged, :action, :insert), as: :project)
  end

  defp list_projects(nil), do: []

  defp list_projects(workspace_scope) do
    Workspace.list_projects(workspace_scope)
  end

  defp filter_projects(projects, ""), do: projects

  defp filter_projects(projects, query) do
    query_downcase = String.downcase(query)

    Enum.filter(projects, fn project ->
      project.name
      |> String.downcase()
      |> String.contains?(query_downcase)
    end)
  end

  defp project_form(params) do
    params
    |> project_changeset()
    |> to_form(as: :project)
  end

  defp project_form_value(form, field) do
    form
    |> Access.get(field)
    |> case do
      nil -> ""
      field_state -> field_state.value || ""
    end
    |> to_string()
    |> String.trim()
  end

  defp project_changeset(params) do
    {%{}, @project_form_types}
    |> Changeset.cast(params, Map.keys(@project_form_types))
    |> Changeset.update_change(:name, &normalize_text/1)
    |> Changeset.update_change(:description, &normalize_text/1)
    |> Changeset.validate_required([:name])
    |> Changeset.validate_length(:name, max: 160)
    |> Changeset.validate_length(:description, max: 4000)
  end

  defp search_form(query), do: to_form(%{"q" => query}, as: :search)

  defp normalize_text(nil), do: nil
  defp normalize_text(value), do: String.trim(value)

  defp format_relative_date(nil), do: "just now"

  defp format_relative_date(datetime) do
    days = Date.diff(Date.utc_today(), DateTime.to_date(datetime))

    cond do
      days <= 0 -> "today"
      days == 1 -> "1 day ago"
      true -> "#{days} days ago"
    end
  end

  defp pluralize_project(1), do: "project"
  defp pluralize_project(_count), do: "projects"

  defp project_slug(name) do
    name
    |> String.downcase()
    |> String.replace(~r/[^a-z0-9]+/u, "-")
    |> String.trim("-")
    |> case do
      "" -> "project"
      slug -> slug
    end
  end

  defp project_initials(name) do
    name
    |> String.split(~r/\s+/, trim: true)
    |> Enum.take(2)
    |> Enum.map(&String.first/1)
    |> Enum.join()
    |> case do
      "" -> "PR"
      initials -> String.upcase(initials)
    end
  end

  defp project_avatar_style(name) do
    colors = [
      {"#06b6d4", "#0891b2"},
      {"#0ea5e9", "#0284c7"},
      {"#22c55e", "#16a34a"},
      {"#f97316", "#ea580c"}
    ]

    {start_color, end_color} = Enum.at(colors, :erlang.phash2(name, length(colors)))
    "background: linear-gradient(135deg, #{start_color} 0%, #{end_color} 100%);"
  end

  defp blank_text?(value) when value in [nil, ""], do: true
  defp blank_text?(value), do: String.trim(value) == ""
end
