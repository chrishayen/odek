defmodule ValkyrieWeb.UserLive.Login do
  use ValkyrieWeb, :live_view

  @impl true
  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope}>
      <div class="login-scene" aria-hidden="true">
        <div class="dot-grid"></div>
        <div class="scan-lines"></div>
        <div class="noise"></div>
        <div class="vignette"></div>
      </div>

      <section id="login-page" class="login-container">
        <div class="login-left-panel">
          <div class="brand-icon">S</div>
          <h1 class="brand-name">Sessions</h1>
          <p class="brand-tagline">agent orchestration platform</p>
          <div class="brand-decoration">
            <span class="brand-dot"></span>
            <span class="brand-dot"></span>
            <span class="brand-dot"></span>
          </div>
        </div>

        <div class="login-right-panel">
          <div class="login-form-container">
            <h2 class="login-form-title">Sign In</h2>
            <p class="login-form-subtitle">
              <%= if @current_scope do %>
                You need to reauthenticate to perform sensitive actions on your account.
              <% else %>
                Welcome back. Enter your credentials to continue.
              <% end %>
            </p>

            <.form
              for={@form}
              id="login_form_password"
              action={~p"/users/log-in"}
              phx-submit="submit_password"
              phx-trigger-action={@trigger_submit}
            >
              <div class="form-section">
                <.input
                  readonly={!!@current_scope}
                  field={@form[:email]}
                  type="email"
                  label="Email"
                  autocomplete="email"
                  required
                  phx-mounted={JS.focus()}
                  class="form-input"
                />
              </div>

              <div class="form-section">
                <.input
                  field={@form[:password]}
                  type="password"
                  label="Password"
                  autocomplete="current-password"
                  required
                  class="form-input"
                />
              </div>

              <.button class="btn btn-primary btn--full">
                Sign In
              </.button>
            </.form>
          </div>
        </div>
      </section>
    </Layouts.app>
    """
  end

  @impl true
  def mount(_params, _session, socket) do
    email =
      Phoenix.Flash.get(socket.assigns.flash, :email) ||
        get_in(socket.assigns, [:current_scope, Access.key(:user), Access.key(:email)])

    form = to_form(%{"email" => email}, as: "user")

    {:ok, assign(socket, form: form, trigger_submit: false)}
  end

  @impl true
  def handle_event("submit_password", _params, socket) do
    {:noreply, assign(socket, :trigger_submit, true)}
  end
end
