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

  pipeline :api_principal_session_no_org do
    plug :accepts, ["json"]
    plug :fetch_session
    plug :fetch_current_scope_for_user

    plug ValkyrieWeb.Plugs.LoadAPIPrincipal,
      allow_api_key: false,
      allow_session: true,
      require_org_context: false
  end

  pipeline :api_principal_session do
    plug :accepts, ["json"]
    plug :fetch_session
    plug :fetch_current_scope_for_user

    plug ValkyrieWeb.Plugs.LoadAPIPrincipal,
      allow_api_key: false,
      allow_session: true,
      require_org_context: true
  end

  pipeline :api_principal_api_key do
    plug :accepts, ["json"]

    plug ValkyrieWeb.Plugs.LoadAPIPrincipal,
      allow_api_key: true,
      allow_session: false,
      require_org_context: false
  end

  pipeline :perm_membership_write do
    plug ValkyrieWeb.Plugs.RequirePermission, permission: "membership.write"
  end

  pipeline :perm_keys_write do
    plug ValkyrieWeb.Plugs.RequirePermission, permission: "keys.write"
  end

  pipeline :perm_projects_read do
    plug ValkyrieWeb.Plugs.RequirePermission, permission: "projects.read"
  end

  pipeline :perm_projects_write do
    plug ValkyrieWeb.Plugs.RequirePermission, permission: "projects.write"
  end

  pipeline :perm_stories_read do
    plug ValkyrieWeb.Plugs.RequirePermission, permission: "stories.read"
  end

  pipeline :perm_stories_write do
    plug ValkyrieWeb.Plugs.RequirePermission, permission: "stories.write"
  end

  pipeline :perm_chat_read do
    plug ValkyrieWeb.Plugs.RequirePermission, permission: "chat.read"
  end

  pipeline :perm_chat_write do
    plug ValkyrieWeb.Plugs.RequirePermission, permission: "chat.write"
  end

  scope "/", ValkyrieWeb do
    pipe_through :browser

    get "/landing", PageController, :home
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
    post "/sessions", AuthController, :create_session
  end

  scope "/v1/auth", ValkyrieWeb.V1 do
    pipe_through :api_principal_session_no_org

    post "/change-password", AuthController, :change_password
    delete "/sessions/:session_id", AuthController, :revoke_session
    post "/sessions/:session_id/switch-org", AuthController, :switch_org
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api_principal_session, :perm_projects_read]

    get "/frontend/projects", AuthController, :frontend_projects
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api_principal_session, :perm_keys_write]

    post "/api-keys", APIKeyController, :create
    get "/api-keys", APIKeyController, :index
    post "/api-keys/:key_id/revoke", APIKeyController, :revoke
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through :api_principal_api_key

    get "/roles", RolesController, :index
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api_principal_api_key, :perm_membership_write]

    post "/memberships/role", MembershipController, :update_role
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api_principal_api_key, :perm_projects_write]

    post "/projects", ProjectController, :create
    patch "/projects/:project_id", ProjectController, :update

    post "/features", FeatureController, :create
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api_principal_api_key, :perm_projects_read]

    get "/projects/:project_id", ProjectController, :show
    post "/projects/:project_id/prompt-preview", PromptController, :preview
    get "/system-prompts/keys", PromptController, :list_keys
    get "/profiles/kinds", PromptController, :list_profile_kinds
    get "/projects/:project_id/context", ContextController, :show
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api_principal_api_key, :perm_stories_write]

    post "/stories", StoryController, :create
    patch "/stories/:story_id", StoryController, :update
    delete "/stories/:story_id", StoryController, :delete

    post "/stories/:story_id/claim", WorkflowController, :claim
    post "/stories/:story_id/state", WorkflowController, :update_state
    post "/stories/:story_id/comments", WorkflowController, :create_comment
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api_principal_api_key, :perm_stories_read]

    get "/stories/poll", WorkflowController, :poll
    get "/stories", StoryController, :index
    get "/stories/:story_id", StoryController, :show
    get "/stories/:story_id/comments", WorkflowController, :list_comments
    get "/stories/:story_id/history", WorkflowController, :list_history
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api_principal_api_key, :perm_projects_write]

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
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api_principal_api_key, :perm_chat_read]

    get "/events/stream", EventsController, :stream
    get "/chat/threads", ChatController, :list_threads
    get "/chat/threads/:thread_id/messages", ChatController, :list_messages
  end

  scope "/v1", ValkyrieWeb.V1 do
    pipe_through [:api_principal_api_key, :perm_chat_write]

    post "/chat/threads", ChatController, :create_thread
    post "/chat/threads/:thread_id/messages", ChatController, :create_message
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
      live "/", ProjectsLive.Index, :index
      live "/projects", ProjectsLive.Index, :index
      live "/projects/:project_id", ProjectsLive.Show, :show
      live "/rules", WorkspacePlaceholderLive, :rules
      live "/skills", WorkspacePlaceholderLive, :skills
      live "/agents", WorkspacePlaceholderLive, :agents
      live "/settings", WorkspacePlaceholderLive, :settings
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
