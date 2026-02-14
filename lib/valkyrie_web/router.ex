defmodule ValkyrieWeb.Router do
  use ValkyrieWeb, :router

  import ValkyrieWeb.UserAuth

  pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_live_flash
    plug :put_root_layout, html: {ValkyrieWeb.Layouts, :root}
    plug :protect_from_forgery
    plug :put_secure_browser_headers
    plug :fetch_current_scope_for_user
    plug :require_password_change
  end

  pipeline :api do
    plug :accepts, ["json"]
  end

  pipeline :api_with_session do
    plug :accepts, ["json"]
    plug :fetch_session
    plug :fetch_current_scope_for_user
  end

  pipeline :api_key_auth do
    plug ValkyrieWeb.Plugs.RequireAPIKey
  end

  scope "/", ValkyrieWeb do
    pipe_through :browser

    get "/", PageController, :home
  end

  scope "/", ValkyrieWeb do
    pipe_through :api

    get "/healthz", HealthController, :show
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through :api

    get "/resources", ResourcesController, :index
    post "/organizations", BootstrapController, :create_organization
    post "/users", BootstrapController, :create_user
    post "/memberships", BootstrapController, :create_membership
  end

  scope "/v1/auth", ValkyrieWeb.V1 do
    pipe_through :api_with_session

    post "/login", AuthController, :login
    post "/change-password", AuthController, :change_password
    post "/sessions", AuthController, :create_session
    delete "/sessions/:session_id", AuthController, :revoke_session
    post "/sessions/:session_id/switch-org", AuthController, :switch_org
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through :api_with_session

    get "/frontend/projects", AuthController, :frontend_projects
    post "/api-keys", APIKeyController, :create
    get "/api-keys", APIKeyController, :index
    post "/api-keys/:key_id/revoke", APIKeyController, :revoke
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api, :api_key_auth]

    get "/roles", RolesController, :index
    post "/memberships/role", MembershipController, :update_role

    post "/projects", ProjectController, :create
    get "/projects/:project_id", ProjectController, :show
    patch "/projects/:project_id", ProjectController, :update

    post "/features", FeatureController, :create

    get "/stories/poll", WorkflowController, :poll
    post "/stories", StoryController, :create
    get "/stories", StoryController, :index
    get "/stories/:story_id", StoryController, :show
    patch "/stories/:story_id", StoryController, :update
    delete "/stories/:story_id", StoryController, :delete

    post "/stories/:story_id/claim", WorkflowController, :claim
    post "/stories/:story_id/state", WorkflowController, :update_state
    post "/stories/:story_id/comments", WorkflowController, :create_comment
    get "/stories/:story_id/comments", WorkflowController, :list_comments
    get "/stories/:story_id/history", WorkflowController, :list_history

    post "/projects/:project_id/prompt-preview", PromptController, :preview
    get "/system-prompts/keys", PromptController, :list_keys
    post "/system-prompts/:prompt_key/versions", PromptController, :add_version
    post "/system-prompts/:prompt_key/activate", PromptController, :activate

    post "/skills", PromptController, :create_skill
    post "/rules", PromptController, :create_rule
    post "/skill-profiles", PromptController, :create_skill_profile
    post "/rule-profiles", PromptController, :create_rule_profile
    post "/skill-profiles/:profile_id/items", PromptController, :add_skill_profile_item
    post "/rule-profiles/:profile_id/items", PromptController, :add_rule_profile_item
    post "/projects/:project_id/skill-profiles", PromptController, :assign_skill_profile
    post "/projects/:project_id/rule-profiles", PromptController, :assign_rule_profile
    get "/profiles/kinds", PromptController, :list_profile_kinds

    get "/events/stream", EventsController, :stream
    post "/chat/threads", ChatController, :create_thread
    get "/chat/threads", ChatController, :list_threads
    post "/chat/threads/:thread_id/messages", ChatController, :create_message
    get "/chat/threads/:thread_id/messages", ChatController, :list_messages

    get "/projects/:project_id/context", ContextController, :show
  end

  # Other scopes may use custom stacks.
  # scope "/api", ValkyrieWeb do
  #   pipe_through :api
  # end

  # Enable LiveDashboard and Swoosh mailbox preview in development
  if Application.compile_env(:valkyrie, :dev_routes) do
    # If you want to use the LiveDashboard in production, you should put
    # it behind authentication and allow only admins to access it.
    # If your application does not have an admins-only section yet,
    # you can use Plug.BasicAuth to set up some basic authentication
    # as long as you are also using SSL (which you should anyway).
    import Phoenix.LiveDashboard.Router

    scope "/dev" do
      pipe_through :browser

      live_dashboard "/dashboard", metrics: ValkyrieWeb.Telemetry
      forward "/mailbox", Plug.Swoosh.MailboxPreview
    end
  end

  ## Authentication routes

  scope "/", ValkyrieWeb do
    pipe_through [:browser, :require_authenticated_user]

    live_session :require_authenticated_user,
      on_mount: [{ValkyrieWeb.UserAuth, :require_authenticated}] do
      live "/users/settings", UserLive.Settings, :edit
      live "/users/settings/confirm-email/:token", UserLive.Settings, :confirm_email
    end

    post "/users/update-password", UserSessionController, :update_password
  end

  scope "/", ValkyrieWeb do
    pipe_through [:browser]

    live_session :current_user,
      on_mount: [{ValkyrieWeb.UserAuth, :mount_current_scope}] do
      live "/users/log-in", UserLive.Login, :new
    end

    post "/users/log-in", UserSessionController, :create
    delete "/users/log-out", UserSessionController, :delete
  end
end
