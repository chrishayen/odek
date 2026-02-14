defmodule ValkyrieWeb.AuthRoutesTest do
  use ValkyrieWeb.ConnCase, async: true

  describe "disabled auth routes" do
    test "returns 404 for registration route", %{conn: conn} do
      conn = get(conn, "/users/register")
      assert html_response(conn, 404)
    end

    test "returns 404 for magic-link token login route", %{conn: conn} do
      conn = get(conn, "/users/log-in/some-token")
      assert html_response(conn, 404)
    end
  end
end
