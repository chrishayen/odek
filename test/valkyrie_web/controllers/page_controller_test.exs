defmodule ValkyrieWeb.PageControllerTest do
  use ValkyrieWeb.ConnCase

  test "GET /landing", %{conn: conn} do
    conn = get(conn, ~p"/landing")
    assert html_response(conn, 200) =~ "Peace of mind from prototype to production"
  end
end
