defmodule ValkyrieWeb.Plugs.RequirePermission do
  @moduledoc """
  Enforces role permissions for API principals.
  """

  import Plug.Conn

  alias Valkyrie.Authorization

  def init(opts), do: opts

  def call(conn, opts) do
    permission = Keyword.fetch!(opts, :permission)
    role = conn.assigns[:current_scope] && conn.assigns.current_scope.role

    if is_binary(role) and Authorization.can?(role, permission) do
      conn
    else
      forbidden(conn)
    end
  end

  defp forbidden(conn) do
    body = Jason.encode!(%{error: %{code: "forbidden", message: "insufficient permissions"}})

    conn
    |> put_resp_content_type("application/json")
    |> send_resp(:forbidden, body)
    |> halt()
  end
end
