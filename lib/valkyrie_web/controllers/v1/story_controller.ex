defmodule ValkyrieWeb.V1.StoryController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Workflow
  alias Valkyrie.Workspace

  def create(conn, params) do
    scope = current_scope!(conn)

    case Workspace.create_story(scope, params) do
      {:ok, story} ->
        conn
        |> put_status(:created)
        |> json(%{
          id: story.id,
          project_id: story.project_id,
          feature_id: story.feature_id,
          name: story.name,
          description: story.description,
          state: story.state
        })

      {:error, :not_found} ->
        not_found(conn, "project/feature not found")

      {:error, changeset} ->
        validation_error(conn, "story_error", changeset)
    end
  end

  def show(conn, %{"story_id" => story_id}) do
    scope = current_scope!(conn)

    case Workspace.get_story(scope, story_id) do
      nil ->
        not_found(conn, "story not found")

      story ->
        json(conn, %{
          id: story.id,
          project_id: story.project_id,
          feature_id: story.feature_id,
          name: story.name,
          description: story.description,
          state: story.state,
          created_at: story.inserted_at,
          updated_at: story.updated_at
        })
    end
  end

  def update(conn, %{"story_id" => story_id} = params) do
    scope = current_scope!(conn)

    case Workspace.get_story(scope, story_id) do
      nil ->
        not_found(conn, "story not found")

      story_before ->
        case Workspace.patch_story_details(scope, story_id, params) do
          {:ok, story_after} ->
            changed_fields = changed_fields(story_before, story_after)
            _ = Workflow.record_field_update(scope, story_id, changed_fields)

            json(conn, %{
              id: story_after.id,
              project_id: story_after.project_id,
              feature_id: story_after.feature_id,
              name: story_after.name,
              description: story_after.description,
              updated: true
            })

          {:error, :not_found} ->
            not_found(conn, "story or feature not found")

          {:error, changeset} ->
            validation_error(conn, "story_error", changeset)
        end
    end
  end

  def index(conn, params) do
    scope = current_scope!(conn)
    project_id = Map.get(params, "project_id")
    cursor = Map.get(params, "cursor")
    limit = parse_limit(params, 20)

    result = Workspace.list_stories(scope, project_id, cursor, limit)

    response = %{
      stories: Enum.map(result.stories, &story_to_json/1)
    }

    response =
      if result.next_cursor do
        Map.put(response, :next_cursor, result.next_cursor)
      else
        response
      end

    json(conn, response)
  end

  def delete(conn, %{"story_id" => story_id}) do
    scope = current_scope!(conn)

    case Workspace.soft_delete_story(scope, story_id) do
      {:ok, story} ->
        json(conn, %{id: story.id, deleted: true})

      {:error, :not_found} ->
        not_found(conn, "story not found")

      {:error, changeset} ->
        validation_error(conn, "delete_error", changeset)
    end
  end

  defp story_to_json(story) do
    %{
      id: story.id,
      project_id: story.project_id,
      feature_id: story.feature_id,
      name: story.name,
      description: story.description,
      state: story.state,
      created_at: story.inserted_at
    }
  end

  defp changed_fields(before_story, after_story) do
    [
      {:name, before_story.name, after_story.name},
      {:description, before_story.description, after_story.description},
      {:feature_id, before_story.feature_id, after_story.feature_id}
    ]
    |> Enum.reduce(%{}, fn {field, old_value, new_value}, acc ->
      if old_value == new_value do
        acc
      else
        Map.put(acc, Atom.to_string(field), %{old: old_value, new: new_value})
      end
    end)
  end
end
