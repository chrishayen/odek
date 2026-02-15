defmodule ValkyrieWeb.WorkspaceComponents do
  @moduledoc """
  Reusable shell primitives for authenticated workspace pages.
  """

  use ValkyrieWeb, :html

  attr :current_scope, :map, required: true

  attr :active_nav, :atom,
    required: true,
    values: [:projects, :rules, :skills, :agents, :settings]

  slot :actions
  slot :inner_block, required: true

  def workspace_shell(assigns) do
    ~H"""
    <div id="workspace-shell" class="workspace-shell">
      <div class="workspace-scene" aria-hidden="true">
        <div class="workspace-dot-grid"></div>
        <div class="workspace-noise"></div>
      </div>

      <div class="workspace-frame">
        <header id="workspace-top-nav" class="workspace-top-nav">
          <div class="workspace-top-nav-left">
            <.link navigate={~p"/projects"} class="workspace-brand">
              <span class="workspace-brand-icon">
                <.icon name="hero-folder" class="size-4" />
              </span>
              <span class="workspace-brand-label">Valkyrie</span>
            </.link>

            <nav class="workspace-nav-links" aria-label="Workspace navigation">
              <.nav_item
                label="Workspace"
                icon="hero-squares-2x2"
                to={~p"/projects"}
                active?={@active_nav == :projects}
              />
              <.nav_item
                label="Rules"
                icon="hero-document-text"
                to={~p"/rules"}
                active?={@active_nav == :rules}
              />
              <.nav_item
                label="Skills"
                icon="hero-star"
                to={~p"/skills"}
                active?={@active_nav == :skills}
              />
              <.nav_item
                label="Agents"
                icon="hero-user"
                to={~p"/agents"}
                active?={@active_nav == :agents}
              />
            </nav>
          </div>

          <div class="workspace-nav-right">
            <div class="workspace-nav-actions">{render_slot(@actions)}</div>
            <.link navigate={~p"/settings"} class="workspace-icon-button" title="Settings">
              <.icon name="hero-cog-6-tooth" class="size-4" />
            </.link>
          </div>
        </header>

        <main class="workspace-main">
          {render_slot(@inner_block)}
        </main>
      </div>
    </div>
    """
  end

  attr :label, :string, required: true
  attr :icon, :string, required: true
  attr :to, :any, required: true
  attr :active?, :boolean, default: false

  defp nav_item(assigns) do
    ~H"""
    <.link
      navigate={@to}
      class={[
        "workspace-nav-item",
        @active? && "workspace-nav-item-active"
      ]}
    >
      <.icon name={@icon} class="size-4" />
      <span class="workspace-nav-item-label">{@label}</span>
    </.link>
    """
  end
end
