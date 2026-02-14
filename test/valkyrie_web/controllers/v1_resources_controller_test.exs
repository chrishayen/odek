defmodule ValkyrieWeb.V1.ResourcesControllerTest do
  use ValkyrieWeb.ConnCase, async: true

  test "GET /healthz returns ok", %{conn: conn} do
    conn = get(conn, "/healthz")

    assert %{"status" => "ok"} = json_response(conn, 200)
  end

  test "GET /v1/resources excludes unsupported resources", %{conn: conn} do
    conn = get(conn, "/v1/resources")
    %{"resources" => resources} = json_response(conn, 200)

    assert "projects" in resources
    refute "agents" in resources
    refute "memory_entries" in resources
    refute "context_entries" in resources
  end
end
