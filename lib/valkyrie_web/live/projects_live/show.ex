defmodule ValkyrieWeb.ProjectsLive.Show do
  use ValkyrieWeb, :live_view

  import ValkyrieWeb.WorkspaceComponents

  alias Valkyrie.Workspace
  alias ValkyrieWeb.WorkspaceScope

  @impl true
  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope}>
      <.workspace_shell current_scope={@current_scope} active_nav={:projects}>
        <:actions>
          <.link navigate={~p"/projects"} class="workspace-ghost-button" id="back-to-projects-button">
            <.icon name="hero-arrow-left" class="size-4" /> Projects
          </.link>
        </:actions>

        <section id="project-detail-page" class="project-detail-card">
          <p class="project-detail-label">Project</p>
          <h1 class="project-detail-title">{@project.name}</h1>
          <p class="project-detail-copy">
            <%= if blank_text?(@project.definition_of_done) do %>
              No description yet. Add stories and project context next.
            <% else %>
              {@project.definition_of_done}
            <% end %>
          </p>
        </section>
      </.workspace_shell>
    </Layouts.app>
    """
  end

  @impl true
  def mount(%{"project_id" => project_id}, _session, socket) do
    workspace_scope = WorkspaceScope.resolve(socket.assigns.current_scope)

    case workspace_scope && Workspace.get_project(workspace_scope, project_id) do
      nil ->
        {:ok,
         socket
         |> put_flash(:error, "Project not found.")
         |> push_navigate(to: ~p"/projects")}

      project ->
        {:ok, assign(socket, workspace_scope: workspace_scope, project: project)}
    end
  end

  defp blank_text?(value) when value in [nil, ""], do: true
  defp blank_text?(value), do: String.trim(value) == ""
end
