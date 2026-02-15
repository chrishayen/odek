defmodule ValkyrieWeb.Plugs.LoadAPIPrincipalTest do
  use ValkyrieWeb.ConnCase, async: true

  import Valkyrie.V1Fixtures

  alias Valkyrie.APIKeys
  alias Valkyrie.Accounts.Scope
  alias ValkyrieWeb.Plugs.LoadAPIPrincipal
  alias ValkyrieWeb.UserAuth

  test "loads api-key principal from bearer token", %{conn: conn} do
    %{user: user, organization: organization} = user_with_membership_fixture("member")

    {:ok, _key, raw_key} = APIKeys.create_api_key(Scope.for_user(user), organization.id)

    conn =
      conn
      |> put_req_header("authorization", "Bearer #{raw_key}")
      |> LoadAPIPrincipal.call(
        allow_api_key: true,
        allow_session: false,
        require_org_context: false
      )

    assert conn.assigns.current_scope.auth_mode == :api_key
    assert conn.assigns.current_scope.organization_id == organization.id
    assert conn.assigns.current_scope.user.id == user.id
  end

  test "falls back to session principal with org context", %{conn: conn} do
    %{user: user, organization: organization} = user_with_membership_fixture("member")

    conn =
      conn
      |> log_in_user(user)
      |> Plug.Conn.put_session(:active_organization_id, organization.id)
      |> UserAuth.fetch_current_scope_for_user([])
      |> LoadAPIPrincipal.call(
        allow_api_key: false,
        allow_session: true,
        require_org_context: true
      )

    assert conn.assigns.current_scope.auth_mode == :session
    assert conn.assigns.current_scope.organization_id == organization.id
    assert conn.assigns.current_scope.role == "member"
  end

  test "bearer token takes precedence over session when both are present", %{conn: conn} do
    %{user: user, organization: organization} = user_with_membership_fixture("member")

    {:ok, _key, raw_key} = APIKeys.create_api_key(Scope.for_user(user), organization.id)

    conn =
      conn
      |> log_in_user(user)
      |> Plug.Conn.put_session(:active_organization_id, organization.id)
      |> UserAuth.fetch_current_scope_for_user([])
      |> put_req_header("authorization", "Bearer #{raw_key}")
      |> LoadAPIPrincipal.call(
        allow_api_key: true,
        allow_session: true,
        require_org_context: true
      )

    assert conn.assigns.current_scope.auth_mode == :api_key
  end

  test "halts with unauthorized when no auth can be resolved", %{conn: conn} do
    conn =
      conn
      |> LoadAPIPrincipal.call(
        allow_api_key: true,
        allow_session: true,
        require_org_context: true
      )

    assert conn.halted
    assert conn.status == 401
    assert %{"error" => %{"code" => "unauthorized"}} = Jason.decode!(conn.resp_body)
  end
end
