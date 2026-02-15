defmodule ValkyrieWeb.WorkspacePlaceholderLive do
  use ValkyrieWeb, :live_view

  import ValkyrieWeb.WorkspaceComponents

  @impl true
  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope}>
      <.workspace_shell current_scope={@current_scope} active_nav={@active_nav}>
        <:actions>
          <.link
            navigate={~p"/projects"}
            class="workspace-ghost-button"
            id="placeholder-back-to-projects"
          >
            <.icon name="hero-arrow-left" class="size-4" /> Projects
          </.link>
        </:actions>

        <section id="workspace-placeholder-page" class="workspace-placeholder-card">
          <p class="workspace-placeholder-label">Coming soon</p>
          <h1 class="workspace-placeholder-title">{@title}</h1>
          <p class="workspace-placeholder-copy">{@copy}</p>
        </section>
      </.workspace_shell>
    </Layouts.app>
    """
  end

  @impl true
  def mount(_params, _session, socket) do
    {active_nav, title, copy} = content_for(socket.assigns.live_action)

    {:ok,
     socket
     |> assign(:active_nav, active_nav)
     |> assign(:title, title)
     |> assign(:copy, copy)}
  end

  defp content_for(:rules) do
    {
      :rules,
      "Rules",
      "Organization and project rules will live here. For now this route is a navigation placeholder."
    }
  end

  defp content_for(:skills) do
    {
      :skills,
      "Skills",
      "Skill libraries and reusable profile mappings will live here. This route is a placeholder for the UI shell."
    }
  end

  defp content_for(:agents) do
    {
      :agents,
      "Agents",
      "Agent-level presets and defaults will live here. This route is currently a placeholder."
    }
  end

  defp content_for(:settings) do
    {
      :settings,
      "Workspace Settings",
      "Workspace-level preferences and integrations will live here. This route is currently a placeholder."
    }
  end
end
