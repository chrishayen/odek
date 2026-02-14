defmodule ValkyrieWeb.V1.ChatController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Chat

  def create_thread(conn, params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "chat.write"),
         {:ok, thread} <- Chat.create_thread(scope, params) do
      conn
      |> put_status(:created)
      |> json(%{
        id: thread.id,
        project_id: thread.project_id,
        story_id: thread.story_id,
        title: thread.title
      })
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "project/story not found")

      {:error, changeset} ->
        render_error(conn, :bad_request, "chat_error", inspect(changeset.errors))
    end
  end

  def list_threads(conn, params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "chat.read"),
         project_id when is_binary(project_id) <- params["project_id"] do
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
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
      _ -> render_error(conn, :bad_request, "invalid_request", "project_id is required")
    end
  end

  def create_message(conn, %{"thread_id" => thread_id} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "chat.write"),
         {:ok, message, _thread} <- Chat.create_message(scope, thread_id, params) do
      conn
      |> put_status(:created)
      |> json(%{id: message.id, thread_id: message.thread_id})
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "thread not found")

      {:error, changeset} ->
        render_error(conn, :bad_request, "chat_error", inspect(changeset.errors))
    end
  end

  def list_messages(conn, %{"thread_id" => thread_id} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "chat.read") do
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
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
    end
  end
end
