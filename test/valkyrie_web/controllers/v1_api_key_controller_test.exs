defmodule ValkyrieWeb.V1.APIKeyControllerTest do
  use ValkyrieWeb.ConnCase, async: true

  import Valkyrie.V1Fixtures

  test "session user can create/list/revoke API keys and revoked key cannot authenticate", %{
    conn: conn
  } do
    %{user: user, organization: organization} = user_with_membership_fixture("member")

    conn =
      conn
      |> log_in_user(user)
      |> post("/v1/api-keys", %{"organization_id" => organization.id})

    %{"id" => key_id, "api_key" => raw_key} = json_response(conn, 201)
    assert is_binary(raw_key)

    conn =
      conn
      |> recycle()
      |> log_in_user(user)
      |> get("/v1/api-keys")

    %{"api_keys" => [first | _]} = json_response(conn, 200)
    assert first["id"] == key_id
    refute Map.has_key?(first, "api_key")

    auth_conn =
      build_conn()
      |> put_req_header("authorization", "Bearer #{raw_key}")
      |> get("/v1/roles")

    assert %{"roles" => roles} = json_response(auth_conn, 200)
    assert Enum.any?(roles, &(&1["role"] == "member"))

    revoke_conn =
      build_conn()
      |> log_in_user(user)
      |> post("/v1/api-keys/#{key_id}/revoke")

    assert %{"revoked" => true} = json_response(revoke_conn, 200)

    denied_conn =
      build_conn()
      |> put_req_header("authorization", "Bearer #{raw_key}")
      |> get("/v1/roles")

    assert %{"error" => %{"code" => "unauthorized"}} = json_response(denied_conn, 401)
  end
end
