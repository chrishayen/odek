defmodule ValkyrieWeb.Plugs.RequireAPIKey do
  @moduledoc """
  Authenticates API requests via bearer API key.
  """

  import Plug.Conn

  alias Valkyrie.APIKeys

  def init(opts), do: opts

  def call(conn, _opts) do
    with [header] <- get_req_header(conn, "authorization"),
         {:ok, token} <- bearer_token(header),
         {:ok, scope} <- APIKeys.authenticate_api_key(token) do
      assign(conn, :current_scope, scope)
    else
      _ -> unauthorized(conn)
    end
  end

  defp bearer_token(header) do
    case String.split(header, " ", parts: 2) do
      [scheme, token] when scheme in ["Bearer", "bearer"] and byte_size(token) > 0 ->
        {:ok, token}

      _ ->
        {:error, :invalid_header}
    end
  end

  defp unauthorized(conn) do
    body = Jason.encode!(%{error: %{code: "unauthorized", message: "invalid api key"}})

    conn
    |> put_resp_content_type("application/json")
    |> send_resp(:unauthorized, body)
    |> halt()
  end
end
