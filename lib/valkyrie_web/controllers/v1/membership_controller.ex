defmodule ValkyrieWeb.V1.MembershipController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Organizations

  def update_role(conn, params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "membership.write") do
      organization_id = Map.get(params, "organization_id")
      user_id = parse_int(Map.get(params, "user_id"), -1)
      role = Map.get(params, "role")

      cond do
        not is_binary(organization_id) or user_id <= 0 or not is_binary(role) ->
          render_error(
            conn,
            :bad_request,
            "invalid_request",
            "organization_id, user_id and role are required"
          )

        organization_id != scope.organization_id ->
          render_error(
            conn,
            :forbidden,
            "forbidden",
            "cross-organization role updates are forbidden"
          )

        true ->
          case Organizations.update_membership_role(organization_id, user_id, role) do
            {:ok, membership} ->
              json(conn, %{
                organization_id: membership.organization_id,
                user_id: membership.user_id,
                role: membership.role
              })

            {:error, :not_found} ->
              render_error(conn, :not_found, "not_found", "membership not found")

            {:error, changeset} ->
              render_error(conn, :bad_request, "invalid_role", inspect(changeset.errors))
          end
      end
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
    end
  end
end
