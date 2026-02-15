defmodule ValkyrieWeb.V1.PromptController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Prompts

  def preview(conn, %{"project_id" => project_id} = params) do
    scope = current_scope!(conn)
    prompt_key = Map.get(params, "prompt_key", "worker")
    user_prompt = Map.get(params, "user_prompt")
    prompt = Prompts.compose_prompt(scope, project_id, prompt_key, user_prompt)

    json(conn, %{project_id: project_id, prompt_key: prompt_key, prompt: prompt})
  end

  def list_keys(conn, _params) do
    scope = current_scope!(conn)
    json(conn, %{keys: Prompts.list_prompt_keys(scope)})
  end

  def add_version(conn, %{"prompt_key" => prompt_key} = params) do
    scope = current_scope!(conn)
    version = parse_int(Map.get(params, "version"), -1)
    body = Map.get(params, "body")

    cond do
      version <= 0 or not is_binary(body) ->
        bad_request(conn, "invalid_request", "version and body are required")

      true ->
        case Prompts.add_prompt_version(scope, prompt_key, version, body) do
          {:ok, result} ->
            conn
            |> put_status(:created)
            |> json(%{prompt_key: result.prompt_key, version: result.version})

          {:error, changeset} ->
            validation_error(conn, "prompt_error", changeset)
        end
    end
  end

  def activate(conn, %{"prompt_key" => prompt_key} = params) do
    scope = current_scope!(conn)
    version = parse_int(Map.get(params, "version"), -1)

    cond do
      version <= 0 ->
        bad_request(conn, "invalid_request", "version must be greater than 0")

      true ->
        case Prompts.activate_prompt_version(scope, prompt_key, version) do
          {:ok, prompt} ->
            json(conn, %{prompt_key: prompt.prompt_key, active_version: prompt.active_version})

          {:error, :not_found} ->
            not_found(conn, "prompt version not found")

          {:error, changeset} ->
            validation_error(conn, "prompt_error", changeset)
        end
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

    case Map.get(params, "skill_id") do
      skill_id when is_binary(skill_id) ->
        case Prompts.add_skill_profile_item(scope, profile_id, skill_id) do
          {:ok, link} ->
            json(conn, %{profile_id: link.profile_id, skill_id: link.skill_id, added: true})

          {:error, :not_found} ->
            not_found(conn, "profile or skill not found")

          {:error, changeset} ->
            validation_error(conn, "invalid_request", changeset)
        end

      _ ->
        bad_request(conn, "invalid_request", "skill_id is required")
    end
  end

  def add_rule_profile_item(conn, %{"profile_id" => profile_id} = params) do
    scope = current_scope!(conn)

    case Map.get(params, "rule_id") do
      rule_id when is_binary(rule_id) ->
        case Prompts.add_rule_profile_item(scope, profile_id, rule_id) do
          {:ok, link} ->
            json(conn, %{profile_id: link.profile_id, rule_id: link.rule_id, added: true})

          {:error, :not_found} ->
            not_found(conn, "profile or rule not found")

          {:error, changeset} ->
            validation_error(conn, "invalid_request", changeset)
        end

      _ ->
        bad_request(conn, "invalid_request", "rule_id is required")
    end
  end

  def assign_skill_profile(conn, %{"project_id" => project_id} = params) do
    assign_profile(conn, project_id, params["profile_id"], &Prompts.assign_skill_profile/3)
  end

  def assign_rule_profile(conn, %{"project_id" => project_id} = params) do
    assign_profile(conn, project_id, params["profile_id"], &Prompts.assign_rule_profile/3)
  end

  def list_profile_kinds(conn, _params) do
    json(conn, %{kinds: Prompts.list_profile_kinds()})
  end

  defp create_profile_resource(conn, params, creator, label) do
    scope = current_scope!(conn)

    case creator.(scope, params) do
      {:ok, resource} ->
        conn
        |> put_status(:created)
        |> json(%{id: resource.id, name: resource.name, kind: label})

      {:error, changeset} ->
        validation_error(conn, "invalid_request", changeset)
    end
  end

  defp assign_profile(conn, project_id, profile_id, assigner) when is_binary(profile_id) do
    scope = current_scope!(conn)

    case assigner.(scope, project_id, profile_id) do
      {:ok, link} ->
        json(conn, %{project_id: link.project_id, profile_id: link.profile_id, assigned: true})

      {:error, :not_found} ->
        not_found(conn, "project or profile not found")

      {:error, changeset} ->
        validation_error(conn, "invalid_request", changeset)
    end
  end

  defp assign_profile(conn, _project_id, _profile_id, _assigner) do
    bad_request(conn, "invalid_request", "profile_id is required")
  end
end
