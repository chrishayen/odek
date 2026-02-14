defmodule ValkyrieWeb.V1.ResourcesController do
  use ValkyrieWeb.V1.BaseController

  @resources [
    "organizations",
    "users",
    "memberships",
    "projects",
    "features",
    "stories",
    "story_comments",
    "story_history_events",
    "skills",
    "rules",
    "skill_profiles",
    "rule_profiles",
    "system_prompts",
    "api_keys",
    "frontend_sessions",
    "chat_threads",
    "chat_messages",
    "story_claims"
  ]

  def index(conn, _params) do
    json(conn, %{resources: @resources})
  end
end
