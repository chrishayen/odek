defmodule ValkyrieWeb.Plugs.RequirePermissionTest do
  use ValkyrieWeb.ConnCase, async: true

  alias Valkyrie.Accounts.Scope
  alias ValkyrieWeb.Plugs.RequirePermission

  test "allows requests when role has permission", %{conn: conn} do
    conn =
      conn
      |> Plug.Conn.assign(:current_scope, %Scope{role: "viewer"})
      |> RequirePermission.call(permission: "projects.read")

    refute conn.halted
  end

  test "halts with forbidden when role lacks permission", %{conn: conn} do
    conn =
      conn
      |> Plug.Conn.assign(:current_scope, %Scope{role: "viewer"})
      |> RequirePermission.call(permission: "projects.write")

    assert conn.halted
    assert conn.status == 403
    assert %{"error" => %{"code" => "forbidden"}} = Jason.decode!(conn.resp_body)
  end
end
