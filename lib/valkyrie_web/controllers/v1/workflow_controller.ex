defmodule ValkyrieWeb.V1.WorkflowController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Workflow

  def poll(conn, params) do
    scope = current_scope!(conn)

    case Map.get(params, "project_id") do
      project_id when is_binary(project_id) ->
        result =
          Workflow.poll_stories(
            scope,
            project_id,
            Map.get(params, "cursor"),
            parse_limit(params, 20)
          )

        response = %{
          stories:
            Enum.map(result.stories, &%{id: &1.id, project_id: &1.project_id, name: &1.name})
        }

        response =
          if result.next_cursor do
            Map.put(response, :next_cursor, result.next_cursor)
          else
            response
          end

        json(conn, response)

      _ ->
        bad_request(conn, "invalid_request", "project_id is required")
    end
  end

  def claim(conn, %{"story_id" => story_id}) do
    scope = current_scope!(conn)

    case Workflow.claim_story(scope, story_id) do
      {:ok, claim} ->
        json(conn, %{story_id: claim.story_id, claim_id: claim.id, claimed: true})

      {:error, :not_found} ->
        not_found(conn, "story not found")

      {:error, :already_claimed} ->
        conflict(conn, "claim_conflict", "story already claimed")

      {:error, changeset} ->
        validation_error(conn, "claim_error", changeset)
    end
  end

  def update_state(conn, %{"story_id" => story_id} = params) do
    scope = current_scope!(conn)

    case Workflow.update_story_state(scope, story_id, params["state"]) do
      {:ok, story} ->
        json(conn, %{story_id: story.id, state: story.state})

      {:error, :not_found} ->
        not_found(conn, "story not found")

      {:error, changeset} ->
        validation_error(conn, "invalid_state", changeset)
    end
  end

  def create_comment(conn, %{"story_id" => story_id} = params) do
    scope = current_scope!(conn)
    body = params["body"]

    if is_binary(body) do
      case Workflow.add_comment(scope, story_id, body) do
        {:ok, comment} ->
          conn
          |> put_status(:created)
          |> json(%{id: comment.id, story_id: comment.story_id})

        {:error, :not_found} ->
          not_found(conn, "story not found")

        {:error, changeset} ->
          validation_error(conn, "comment_error", changeset)
      end
    else
      bad_request(conn, "invalid_request", "body is required")
    end
  end

  def list_comments(conn, %{"story_id" => story_id} = params) do
    scope = current_scope!(conn)
    comments = Workflow.list_comments(scope, story_id, parse_limit(params, 50))

    json(conn, %{
      comments:
        Enum.map(comments, fn comment ->
          %{
            id: comment.id,
            author_user_id: comment.author_user_id,
            body: comment.body,
            created_at: comment.inserted_at
          }
        end)
    })
  end

  def list_history(conn, %{"story_id" => story_id} = params) do
    scope = current_scope!(conn)
    history = Workflow.list_history(scope, story_id, parse_limit(params, 100))

    json(conn, %{
      history:
        Enum.map(history, fn event ->
          %{
            id: event.id,
            actor_user_id: event.actor_user_id,
            event_type: event.event_type,
            payload: event.payload,
            created_at: event.inserted_at
          }
        end)
    })
  end
end
