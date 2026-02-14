defmodule ValkyrieWeb.V1.ContextController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Context

  def show(conn, %{"project_id" => project_id}) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.read"),
         {:ok, context} <- Context.build_project_context(scope, project_id) do
      json(conn, %{
        project: project_json(context.project),
        stories: Enum.map(context.stories, &story_json/1),
        comments: Enum.map(context.comments, &comment_json/1),
        history: Enum.map(context.history, &history_json/1),
        threads: Enum.map(context.threads, &thread_json/1),
        messages: Enum.map(context.messages, &message_json/1)
      })
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
      {:error, :not_found} -> render_error(conn, :not_found, "not_found", "project not found")
    end
  end

  defp project_json(project) do
    %{
      id: project.id,
      organization_id: project.organization_id,
      name: project.name,
      definition_of_done: project.definition_of_done,
      created_at: project.inserted_at,
      updated_at: project.updated_at
    }
  end

  defp story_json(story) do
    %{
      id: story.id,
      project_id: story.project_id,
      feature_id: story.feature_id,
      name: story.name,
      description: story.description,
      state: story.state,
      created_at: story.inserted_at,
      updated_at: story.updated_at
    }
  end

  defp comment_json(comment) do
    %{
      id: comment.id,
      story_id: comment.story_id,
      author_user_id: comment.author_user_id,
      body: comment.body,
      created_at: comment.inserted_at
    }
  end

  defp history_json(event) do
    %{
      id: event.id,
      story_id: event.story_id,
      actor_user_id: event.actor_user_id,
      event_type: event.event_type,
      payload: event.payload,
      created_at: event.inserted_at
    }
  end

  defp thread_json(thread) do
    %{
      id: thread.id,
      organization_id: thread.organization_id,
      project_id: thread.project_id,
      story_id: thread.story_id,
      title: thread.title,
      created_by_user_id: thread.created_by_user_id,
      created_at: thread.inserted_at
    }
  end

  defp message_json(message) do
    %{
      id: message.id,
      thread_id: message.thread_id,
      sender_kind: message.sender_kind,
      sender_user_id: message.sender_user_id,
      body: message.body,
      metadata: message.metadata,
      created_at: message.inserted_at
    }
  end
end
