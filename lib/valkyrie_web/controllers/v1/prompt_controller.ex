defmodule ValkyrieWeb.V1.PromptController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Prompts

  def preview(conn, %{"project_id" => project_id} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.read") do
      prompt_key = Map.get(params, "prompt_key", "worker")
      user_prompt = Map.get(params, "user_prompt")
      prompt = Prompts.compose_prompt(scope, project_id, prompt_key, user_prompt)

      json(conn, %{project_id: project_id, prompt_key: prompt_key, prompt: prompt})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
    end
  end

  def list_keys(conn, _params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.read") do
      json(conn, %{keys: Prompts.list_prompt_keys(scope)})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
    end
  end

  def add_version(conn, %{"prompt_key" => prompt_key} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.write"),
         version when version > 0 <- parse_int(Map.get(params, "version"), -1),
         body when is_binary(body) <- Map.get(params, "body"),
         {:ok, result} <- Prompts.add_prompt_version(scope, prompt_key, version, body) do
      conn
      |> put_status(:created)
      |> json(%{prompt_key: result.prompt_key, version: result.version})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, changeset} ->
        render_error(conn, :bad_request, "prompt_error", inspect(changeset.errors))

      _ ->
        render_error(conn, :bad_request, "invalid_request", "version and body are required")
    end
  end

  def activate(conn, %{"prompt_key" => prompt_key} = params) do
    scope = current_scope!(conn)
    version = parse_int(Map.get(params, "version"), -1)

    with :ok <- authorize(conn, "projects.write"),
         true <- version > 0,
         {:ok, prompt} <- Prompts.activate_prompt_version(scope, prompt_key, version) do
      json(conn, %{prompt_key: prompt.prompt_key, active_version: prompt.active_version})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      false ->
        render_error(conn, :bad_request, "invalid_request", "version must be greater than 0")

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "prompt version not found")

      {:error, changeset} ->
        render_error(conn, :bad_request, "prompt_error", inspect(changeset.errors))
    end
  end

  def create_skill(conn, params) do
    create_profile_resource(conn, params, &Prompts.create_skill/2, "skill")
  end

  def create_rule(conn, params) do
    create_profile_resource(conn, params, &Prompts.create_rule/2, "rule")
  end

  def create_skill_profile(conn, params) do
    create_profile_resource(conn, params, &Prompts.create_skill_profile/2, "skill_profile")
  end

  def create_rule_profile(conn, params) do
    create_profile_resource(conn, params, &Prompts.create_rule_profile/2, "rule_profile")
  end

  def add_skill_profile_item(conn, %{"profile_id" => profile_id} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.write"),
         skill_id when is_binary(skill_id) <- Map.get(params, "skill_id"),
         {:ok, link} <- Prompts.add_skill_profile_item(scope, profile_id, skill_id) do
      json(conn, %{profile_id: link.profile_id, skill_id: link.skill_id, added: true})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "profile or skill not found")

      {:error, changeset} ->
        render_error(conn, :bad_request, "invalid_request", inspect(changeset.errors))

      _ ->
        render_error(conn, :bad_request, "invalid_request", "skill_id is required")
    end
  end

  def add_rule_profile_item(conn, %{"profile_id" => profile_id} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.write"),
         rule_id when is_binary(rule_id) <- Map.get(params, "rule_id"),
         {:ok, link} <- Prompts.add_rule_profile_item(scope, profile_id, rule_id) do
      json(conn, %{profile_id: link.profile_id, rule_id: link.rule_id, added: true})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "profile or rule not found")

      {:error, changeset} ->
        render_error(conn, :bad_request, "invalid_request", inspect(changeset.errors))

      _ ->
        render_error(conn, :bad_request, "invalid_request", "rule_id is required")
    end
  end

  def assign_skill_profile(conn, %{"project_id" => project_id} = params) do
    assign_profile(conn, project_id, params["profile_id"], &Prompts.assign_skill_profile/3)
  end

  def assign_rule_profile(conn, %{"project_id" => project_id} = params) do
    assign_profile(conn, project_id, params["profile_id"], &Prompts.assign_rule_profile/3)
  end

  def list_profile_kinds(conn, _params) do
    with :ok <- authorize(conn, "projects.read") do
      json(conn, %{kinds: Prompts.list_profile_kinds()})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
    end
  end

  defp create_profile_resource(conn, params, creator, label) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.write"),
         {:ok, resource} <- creator.(scope, params) do
      conn
      |> put_status(:created)
      |> json(%{id: resource.id, name: resource.name, kind: label})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, changeset} ->
        render_error(conn, :bad_request, "invalid_request", inspect(changeset.errors))
    end
  end

  defp assign_profile(conn, project_id, profile_id, assigner) when is_binary(profile_id) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.write"),
         {:ok, link} <- assigner.(scope, project_id, profile_id) do
      json(conn, %{project_id: link.project_id, profile_id: link.profile_id, assigned: true})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "project or profile not found")

      {:error, changeset} ->
        render_error(conn, :bad_request, "invalid_request", inspect(changeset.errors))
    end
  end

  defp assign_profile(conn, _project_id, _profile_id, _assigner) do
    render_error(conn, :bad_request, "invalid_request", "profile_id is required")
  end
end
