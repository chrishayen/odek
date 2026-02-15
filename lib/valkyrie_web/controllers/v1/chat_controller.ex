defmodule ValkyrieWeb.V1.ChatController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Chat

  def create_thread(conn, params) do
    scope = current_scope!(conn)

    case Chat.create_thread(scope, params) do
      {:ok, thread} ->
        conn
        |> put_status(:created)
        |> json(%{
          id: thread.id,
          project_id: thread.project_id,
          story_id: thread.story_id,
          title: thread.title
        })

      {:error, :not_found} ->
        not_found(conn, "project/story not found")

      {:error, changeset} ->
        validation_error(conn, "chat_error", changeset)
    end
  end

  def list_threads(conn, params) do
    scope = current_scope!(conn)

    case params["project_id"] do
      project_id when is_binary(project_id) ->
        threads = Chat.list_threads(scope, project_id, parse_limit(params, 50))

        json(conn, %{
          threads:
            Enum.map(threads, fn thread ->
              %{
                id: thread.id,
                project_id: thread.project_id,
                story_id: thread.story_id,
                title: thread.title,
                created_at: thread.inserted_at
              }
            end)
        })

      _ ->
        bad_request(conn, "invalid_request", "project_id is required")
    end
  end

  def create_message(conn, %{"thread_id" => thread_id} = params) do
    scope = current_scope!(conn)

    case Chat.create_message(scope, thread_id, params) do
      {:ok, message, _thread} ->
        conn
        |> put_status(:created)
        |> json(%{id: message.id, thread_id: message.thread_id})

      {:error, :not_found} ->
        not_found(conn, "thread not found")

      {:error, changeset} ->
        validation_error(conn, "chat_error", changeset)
    end
  end

  def list_messages(conn, %{"thread_id" => thread_id} = params) do
    scope = current_scope!(conn)
    messages = Chat.list_messages(scope, thread_id, parse_limit(params, 100))

    json(conn, %{
      messages:
        Enum.map(messages, fn message ->
          %{
            id: message.id,
            sender_kind: message.sender_kind,
            sender_user_id: message.sender_user_id,
            body: message.body,
            metadata: message.metadata,
            created_at: message.inserted_at
          }
        end)
    })
  end
end
