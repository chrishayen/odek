defmodule ValkyrieWeb.V1.EventsController do
  use ValkyrieWeb.V1.BaseController

  import Plug.Conn

  alias Valkyrie.Chat

  def stream(conn, params) do
    scope = current_scope!(conn)
    project_id = Map.get(params, "project_id")

    with :ok <- authorize(conn, "chat.read"),
         false <- blank?(project_id) do
      :ok = Chat.subscribe(project_id)

      conn =
        conn
        |> put_resp_content_type("text/event-stream")
        |> put_resp_header("cache-control", "no-cache")
        |> put_resp_header("connection", "keep-alive")
        |> send_chunked(:ok)

      case chunk(
             conn,
             "event: ready\ndata: {\"organization_id\":\"#{scope.organization_id}\",\"project_id\":\"#{project_id}\"}\n\n"
           ) do
        {:ok, conn} -> stream_events(conn)
        {:error, _reason} -> conn
      end
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
      true -> render_error(conn, :bad_request, "invalid_request", "project_id is required")
    end
  end

  defp stream_events(conn) do
    receive do
      {:chat_message, event} ->
        payload = Jason.encode!(event)

        case chunk(conn, "event: chat_message\ndata: #{payload}\n\n") do
          {:ok, conn} -> stream_events(conn)
          {:error, _reason} -> conn
        end
    after
      15_000 ->
        case chunk(conn, ": keepalive\n\n") do
          {:ok, conn} -> stream_events(conn)
          {:error, _reason} -> conn
        end
    end
  end
end
