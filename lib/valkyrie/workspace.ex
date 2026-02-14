defmodule Valkyrie.Workspace do
  @moduledoc """
  Projects, features, stories, pagination, and soft-delete behavior.
  """

  import Ecto.Query, warn: false

  alias Valkyrie.Accounts.Scope
  alias Valkyrie.Organizations
  alias Valkyrie.Repo
  alias Valkyrie.Workspace.{Feature, Project, Story}

  def canonical_story_states, do: Story.states()

  def create_project(%Scope{} = scope, attrs) do
    with {:ok, organization_id} <- resolve_organization(scope, attrs),
         :ok <- ensure_membership(organization_id, scope.user.id) do
      attrs
      |> Map.new(fn {k, v} -> {to_string(k), v} end)
      |> Map.put("organization_id", organization_id)
      |> then(&(%Project{} |> Project.changeset(&1) |> Repo.insert()))
    end
  end

  def get_project(%Scope{} = scope, project_id) do
    Repo.one(
      from p in Project,
        where:
          p.id == ^project_id and p.organization_id == ^scope.organization_id and
            is_nil(p.deleted_at)
    )
  end

  def list_projects(%Scope{} = scope) do
    Repo.all(
      from p in Project,
        where: p.organization_id == ^scope.organization_id and is_nil(p.deleted_at),
        order_by: [asc: p.inserted_at, asc: p.id]
    )
  end

  def update_project(%Scope{} = scope, project_id, attrs) do
    case get_project(scope, project_id) do
      nil -> {:error, :not_found}
      project -> project |> Project.update_changeset(attrs) |> Repo.update()
    end
  end

  def create_feature(%Scope{} = scope, attrs) do
    project_id = attr(attrs, "project_id")

    with {:ok, %Project{}} <- fetch_project(scope.organization_id, project_id) do
      attrs
      |> Map.new(fn {k, v} -> {to_string(k), v} end)
      |> then(&(%Feature{} |> Feature.changeset(&1) |> Repo.insert()))
    end
  end

  def create_story(%Scope{} = scope, attrs) do
    project_id = attr(attrs, "project_id")
    feature_id = attr(attrs, "feature_id")

    with {:ok, %Project{}} <- fetch_project(scope.organization_id, project_id),
         :ok <- validate_feature_project(project_id, feature_id) do
      params =
        attrs
        |> Map.new(fn {k, v} -> {to_string(k), v} end)
        |> Map.put("feature_id", normalize_optional(feature_id))

      %Story{}
      |> Story.changeset(params)
      |> Repo.insert()
    end
  end

  def get_story(%Scope{} = scope, story_id) do
    Repo.one(
      from s in Story,
        join: p in Project,
        on: p.id == s.project_id,
        where:
          s.id == ^story_id and p.organization_id == ^scope.organization_id and
            is_nil(s.deleted_at) and
            is_nil(p.deleted_at),
        preload: [:project]
    )
  end

  def patch_story_details(%Scope{} = scope, story_id, attrs) do
    case get_story(scope, story_id) do
      nil ->
        {:error, :not_found}

      story ->
        feature_id = attr(attrs, "feature_id")

        with :ok <- validate_feature_project(story.project_id, feature_id) do
          params =
            attrs
            |> Map.new(fn {k, v} -> {to_string(k), v} end)
            |> maybe_set_optional_feature_id(feature_id)

          story
          |> Story.update_details_changeset(params)
          |> Repo.update()
        end
    end
  end

  def update_story_state(%Scope{} = scope, story_id, state) do
    case get_story(scope, story_id) do
      nil -> {:error, :not_found}
      story -> story |> Story.state_changeset(%{state: state}) |> Repo.update()
    end
  end

  def soft_delete_story(%Scope{} = scope, story_id) do
    case get_story(scope, story_id) do
      nil ->
        {:error, :not_found}

      story ->
        story
        |> Story.changeset(%{deleted_at: DateTime.utc_now()})
        |> Repo.update()
    end
  end

  def list_stories(%Scope{} = scope, project_id, cursor, limit) do
    query =
      from s in Story,
        join: p in Project,
        on: p.id == s.project_id,
        where:
          p.organization_id == ^scope.organization_id and is_nil(s.deleted_at) and
            is_nil(p.deleted_at),
        where: ^is_nil(project_id) or s.project_id == ^project_id,
        order_by: [asc: s.inserted_at, asc: s.id]

    query = apply_cursor(query, scope.organization_id, cursor)
    fetch_cursor_page(query, limit)
  end

  def poll_ready_stories(%Scope{} = scope, project_id, cursor, limit) do
    query =
      from s in Story,
        join: p in Project,
        on: p.id == s.project_id,
        where:
          p.organization_id == ^scope.organization_id and is_nil(s.deleted_at) and
            is_nil(p.deleted_at) and
            s.state == "ready" and s.project_id == ^project_id,
        order_by: [asc: s.inserted_at, asc: s.id]

    query
    |> apply_cursor(scope.organization_id, cursor)
    |> fetch_cursor_page(limit)
  end

  defp apply_cursor(query, _organization_id, nil), do: query
  defp apply_cursor(query, _organization_id, ""), do: query

  defp apply_cursor(query, organization_id, cursor) do
    cursor_story =
      Repo.one(
        from s in Story,
          join: p in Project,
          on: p.id == s.project_id,
          where:
            s.id == ^cursor and p.organization_id == ^organization_id and is_nil(s.deleted_at) and
              is_nil(p.deleted_at),
          select: %{id: s.id, inserted_at: s.inserted_at}
      )

    case cursor_story do
      nil ->
        query

      cs ->
        from s in query,
          where:
            s.inserted_at > ^cs.inserted_at or
              (s.inserted_at == ^cs.inserted_at and s.id > ^cs.id)
    end
  end

  defp fetch_cursor_page(query, limit) do
    limit = normalize_limit(limit, 20)
    rows = Repo.all(from s in query, limit: ^(limit + 1))
    {page, extra} = Enum.split(rows, limit)

    next_cursor =
      if extra == [] do
        nil
      else
        page |> List.last() |> Map.get(:id)
      end

    %{stories: page, next_cursor: next_cursor}
  end

  defp normalize_limit(limit, _fallback) when is_integer(limit) and limit > 0, do: limit
  defp normalize_limit(_limit, fallback), do: fallback

  defp fetch_project(org_id, project_id) do
    case Repo.one(
           from p in Project,
             where: p.id == ^project_id and p.organization_id == ^org_id and is_nil(p.deleted_at)
         ) do
      nil -> {:error, :not_found}
      project -> {:ok, project}
    end
  end

  defp validate_feature_project(_project_id, nil), do: :ok
  defp validate_feature_project(_project_id, ""), do: :ok

  defp validate_feature_project(project_id, feature_id) do
    case Repo.one(
           from f in Feature,
             where: f.id == ^feature_id and f.project_id == ^project_id and is_nil(f.deleted_at)
         ) do
      nil -> {:error, :not_found}
      _feature -> :ok
    end
  end

  defp resolve_organization(%Scope{organization_id: nil}, attrs) do
    case attr(attrs, "organization_id") do
      nil -> {:error, :invalid_request}
      organization_id -> {:ok, organization_id}
    end
  end

  defp resolve_organization(%Scope{organization_id: org_id}, attrs) do
    case attr(attrs, "organization_id") do
      nil -> {:ok, org_id}
      ^org_id -> {:ok, org_id}
      _other -> {:error, :forbidden}
    end
  end

  defp ensure_membership(organization_id, user_id) do
    if Organizations.member?(organization_id, user_id), do: :ok, else: {:error, :forbidden}
  end

  defp attr(attrs, key) do
    Map.get(attrs, key) || Map.get(attrs, String.to_atom(key))
  end

  defp normalize_optional(nil), do: nil
  defp normalize_optional(""), do: nil
  defp normalize_optional(value), do: value

  defp maybe_set_optional_feature_id(params, feature_id) do
    if is_nil(feature_id) do
      params
    else
      Map.put(params, "feature_id", normalize_optional(feature_id))
    end
  end
end
