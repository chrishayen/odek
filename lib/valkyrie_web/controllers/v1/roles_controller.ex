defmodule ValkyrieWeb.V1.RolesController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Organizations

  def index(conn, _params) do
    roles =
      Organizations.list_roles()
      |> Enum.sort()
      |> Enum.map(&%{role: &1})

    json(conn, %{roles: roles})
  end
end
