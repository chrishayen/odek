defmodule ValkyrieWeb.V1.BootstrapControllerTest do
  use ValkyrieWeb.ConnCase, async: true

  alias Valkyrie.Organizations

  test "create_user provisions organization membership", %{conn: conn} do
    email = "bootstrap#{System.unique_integer([:positive])}@example.com"

    conn =
      post(conn, "/v1/users", %{
        "email" => email,
        "password" => "bootstrap secret",
        "name" => "Bootstrap User"
      })

    assert %{
             "id" => user_id,
             "email" => ^email,
             "name" => "Bootstrap User",
             "organization_id" => organization_id,
             "role" => "owner"
           } = json_response(conn, 201)

    assert is_integer(user_id)
    assert is_binary(organization_id)
    assert %{role: "owner"} = Organizations.get_membership(organization_id, user_id)
  end
end
