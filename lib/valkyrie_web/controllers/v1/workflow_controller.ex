defmodule ValkyrieWeb.V1.WorkflowController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Workflow

  def poll(conn, params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "stories.read"),
         project_id when is_binary(project_id) <- Map.get(params, "project_id") do
      result =
        Workflow.poll_stories(
          scope,
          project_id,
          Map.get(params, "cursor"),
          parse_limit(params, 20)
        )

      response = %{
        stories: Enum.map(result.stories, &%{id: &1.id, project_id: &1.project_id, name: &1.name})
      }

      response =
        if result.next_cursor do
          Map.put(response, :next_cursor, result.next_cursor)
        else
          response
        end

      json(conn, response)
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
      _ -> render_error(conn, :bad_request, "invalid_request", "project_id is required")
    end
  end

  def claim(conn, %{"story_id" => story_id}) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "stories.write"),
         {:ok, claim} <- Workflow.claim_story(scope, story_id) do
      json(conn, %{story_id: claim.story_id, claim_id: claim.id, claimed: true})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "story not found")

      {:error, :already_claimed} ->
        render_error(conn, :conflict, "claim_conflict", "story already claimed")

      {:error, changeset} ->
        render_error(conn, :bad_request, "claim_error", inspect(changeset.errors))
    end
  end

  def update_state(conn, %{"story_id" => story_id} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "stories.write"),
         {:ok, story} <- Workflow.update_story_state(scope, story_id, params["state"]) do
      json(conn, %{story_id: story.id, state: story.state})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "story not found")

      {:error, changeset} ->
        render_error(conn, :bad_request, "invalid_state", inspect(changeset.errors))

      _ ->
        render_error(conn, :bad_request, "invalid_request", "state is required")
    end
  end

  def create_comment(conn, %{"story_id" => story_id} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "stories.write"),
         body when is_binary(body) <- params["body"],
         {:ok, comment} <- Workflow.add_comment(scope, story_id, body) do
      conn
      |> put_status(:created)
      |> json(%{id: comment.id, story_id: comment.story_id})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "story not found")

      {:error, changeset} ->
        render_error(conn, :bad_request, "comment_error", inspect(changeset.errors))

      _ ->
        render_error(conn, :bad_request, "invalid_request", "body is required")
    end
  end

  def list_comments(conn, %{"story_id" => story_id} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "stories.read") do
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
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
    end
  end

  def list_history(conn, %{"story_id" => story_id} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "stories.read") do
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
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
    end
  end
end
